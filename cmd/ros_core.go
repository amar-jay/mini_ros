package main

import (
	"bufio"
	"encoding/json"
	"flag"
	"fmt"
	"net"
	"strconv"
	"strings"
	"sync"
)

type RosCore struct {
	mu          sync.RWMutex
	Subscribers map[string][]net.Conn
}

func NewRosCore() *RosCore {
	return &RosCore{
		mu:          sync.RWMutex{},
		Subscribers: make(map[string][]net.Conn),
	}
}

type Status struct {
	Subscribers map[string]int `json:"subscribers"`
}

func (r *RosCore) listen(host string, port int) {
	ln, err := net.Listen("tcp", host+":"+strconv.Itoa(port))
	if err != nil {
		fmt.Println("Error starting TCP server:", err)
		return
	}
	defer ln.Close()

	fmt.Println("Roscore Server listening on port", port)

	for {
		conn, err := ln.Accept()
		if err != nil {
			fmt.Println("Error accepting connection:", err)
			continue
		}

		go r.HandleConn(conn)
	}
}

func (r *RosCore) HandleConn(conn net.Conn) {
	defer conn.Close()
	reader := bufio.NewReader(conn)
	for {
		// Read the incoming message from the client
		message, err := reader.ReadString('\n')
		if err != nil {
			//fmt.Println("Client disconnected:", conn.RemoteAddr())
			break
		}

		message = strings.TrimSpace(message)
		tokens := strings.SplitN(message, " ", 2)

		if len(tokens) < 2 {
			conn.Write([]byte("Invalid command\n"))
			continue
		}

		command, topic := tokens[0], tokens[1]
		//println(len(r.Subscribers[topic]), "Subscribers", topic, " ", command)

		switch command {
		case "SUBSCRIBE":
			r.Subscribe(topic, conn)
		case "PUBLISH":
			parts := strings.SplitAfterN(topic, " ", 2)
			fmt.Println(parts)
			if len(parts) < 2 {
				conn.Write([]byte("Invalid publish format. Use: PUBLISH <topic> <message>\n"))
			} else {
				topic, message := parts[0], parts[1]
				message = strings.TrimSpace(message)
				topic = strings.TrimSpace(topic)
				/*
					var msg interface{}
					err = json.Unmarshal([]byte(message), &msg)
					if err != nil {
						panic(err)
					}
					fmt.Printf("%s> %v\n", topic, msg)
				*/
				r.Publish(topic, []byte(message))
			}
		case "STATUS":
			//println("Number of Subscribers: ", len(r.Subscribers[topic]))
			status := Status{Subscribers: map[string]int{topic: len(r.Subscribers[topic])}}
			for t, conns := range r.Subscribers {
				status.Subscribers[t] = len(conns)
			}
			st, err := json.Marshal(status)
			if err != nil {
				fmt.Println("Error marshalling status")
				return
			}

			conn.Write([]byte(string(st) + "\n"))
		default:
			conn.Write([]byte("Unknown command\n"))
		}
	}
}

func (r *RosCore) Subscribe(topic string, conn net.Conn) {
	r.mu.Lock()
	defer r.mu.Unlock()

	r.Subscribers[topic] = append(r.Subscribers[topic], conn)
	// since a topic can have multiple subscribers, we keep track of all the subscribers in a slice.
	fmt.Println("Client", conn.RemoteAddr(), "subscribed to topic", topic)
	//fmt.Printf("Subscribers: %d\n", len(r.Subscribers[topic]))
	msg, _ := json.Marshal(map[string]string{
		"message": "subscribed successfully",
	})

	conn.Write([]byte(topic + " " + string(msg) + "\n"))
}

func (r *RosCore) Publish(topic string, message []byte) {
	r.mu.RLock()
	defer r.mu.RUnlock()
	fmt.Printf("Published to topic %s message type %T\n", topic, message)
	//fmt.Printf("%s ---- %d %d\n", topic, len(r.Subscribers[topic]), len(r.Subscribers["/hello"]))
	for _, conn := range r.Subscribers[topic] {
		conn.Write([]byte(topic + " " + string(message) + "\n"))
	}
}

func main() {
	port := flag.Int("port", 11311, "ROS master port")
	host := flag.String("url", "", "ROS Master URL")
	flag.Parse()
	r := NewRosCore()
	r.listen(*host, *port)
}
