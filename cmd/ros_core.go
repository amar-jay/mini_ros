package main

import (
	"bufio"
	"encoding/json"
	"fmt"
	"net"
	"strconv"
	"strings"
	"sync"
)

type RosCore struct {
	mu          sync.RWMutex
	Subscribers map[string][]net.Conn // map of topic to subscribers (connections)
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

	fmt.Printf("roscore server listening on tcp://%s:%d/\n", host, port)

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

		var command, topic string
		if len(tokens) == 2 {
			command, topic = tokens[0], tokens[1]
		} else if len(tokens) == 1 {
			command = tokens[0]
		} else {
			conn.Write([]byte("Invalid command\n"))
			println("Invalid command", message)
			continue
		}
		//println(len(r.Subscribers[topic]), "Subscribers", topic, " ", command)

		switch command {
		case "SUBSCRIBE":
			r.Subscribe(topic, conn)
		case "UNSUBSCRIBE":
			r.Unsubscribe(topic, conn)
		case "PUBLISH":
			parts := strings.SplitAfterN(topic, " ", 2)
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

		case "LIST":
			//println("Number of Subscribers: ", len(r.Subscribers[topic]))
			topics := make([]string, 0, len(r.Subscribers))
			for t := range r.Subscribers {
				topics = append(topics, t)
			}
			st, err := json.Marshal(topics)
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
	/* THERE IS NO NEED TO SEND A MESSAGE TO THE CLIENT THAT THEY HAVE SUBSCRIBED SUCCESSFULLY
	msg, _ := json.Marshal(map[string]string{
		"message": "subscribed successfully",
	})

	conn.Write([]byte(topic + " " + string(msg) + "\n"))
	*/
}

func (r *RosCore) Unsubscribe(topic string, conn net.Conn) {
	r.mu.Lock()
	defer r.mu.Unlock()
	for i, c := range r.Subscribers[topic] {
		if c == conn {
			r.Subscribers[topic] = append(r.Subscribers[topic][:i], r.Subscribers[topic][i+1:]...)
			fmt.Println("Client", conn.RemoteAddr(), "unsubscribed from topic", topic)
			break
		}

	}

	// if empty delete
	if len(r.Subscribers[topic]) == 0 {
		delete(r.Subscribers, topic)
	}
	//conn.Write([]byte("Unsubscribed from " + topic + " successfully\n"))
}

func (r *RosCore) Publish(topic string, message []byte) {
	r.mu.RLock()
	defer r.mu.RUnlock()
	fmt.Printf("Published to topic %s (type %T) to %d subscribers\n", topic, message, len(r.Subscribers))
	//fmt.Printf("%s ---- %d %d\n", topic, len(r.Subscribers[topic]), len(r.Subscribers["/hello"]))
	for _, conn := range r.Subscribers[topic] {
		conn.Write([]byte(topic + " " + string(message) + "\n"))
	}
}
