package main

import (
	"bufio"
	"flag"
	"fmt"
	"net"
	"strconv"
	"strings"
	"sync"
)

type RosCore struct {
	mu          sync.RWMutex
	subscribers map[string][]net.Conn
}

func NewRosCore() *RosCore {
	return &RosCore{
		mu:          sync.RWMutex{},
		subscribers: make(map[string][]net.Conn),
	}
}

var r = NewRosCore()

func startServer(host string, port int) {
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

		go handleConn(conn)
	}
}

func handleConn(conn net.Conn) {
	defer conn.Close()
	reader := bufio.NewReader(conn)
	for {
		// Read the incoming message from the client
		message, err := reader.ReadString('\n')
		if err != nil {
			fmt.Println("Client disconnected:", conn.RemoteAddr())
			break
		}

		message = strings.TrimSpace(message)
		tokens := strings.SplitN(message, " ", 2)

		if len(tokens) < 2 {
			conn.Write([]byte("Invalid command\n"))
			continue
		}

		command, topic := tokens[0], tokens[1]

		switch command {
		case "SUBSCRIBE":
			subscribe(topic, conn)
		case "PUBLISH":
			parts := strings.SplitN(topic, " ", 2)
			if len(parts) < 2 {
				conn.Write([]byte("Invalid publish format. Use: PUBLISH <topic> <message>\n"))
			} else {
				topic, message := parts[0], parts[1]
				publish(topic, message)
			}
		default:
			conn.Write([]byte("Unknown command\n"))
		}
	}
}

// Subscribe allows a client to subscribe to a specific topic.
func subscribe(topic string, conn net.Conn) {
	r.mu.Lock()
	defer r.mu.Unlock()

	r.subscribers[topic] = append(r.subscribers[topic], conn)
	// since a topic can have multiple subscribers, we keep track of all the subscribers in a slice.
	fmt.Println("Client", conn.RemoteAddr(), "subscribed to topic", topic)
	conn.Write([]byte(topic + " 'subscribed successfully'" + "\n"))
}

func publish(topic, message string) {
	r.mu.RLock()
	defer r.mu.RUnlock()
	fmt.Printf("Published to topic %s message type %T", topic, message)
	for _, conn := range r.subscribers[topic] {
		conn.Write([]byte(topic + " " + message + "\n"))
	}
}

func main() {
	port := flag.Int("port", 11311, "ROS master port")
	host := flag.String("url", "", "ROS Master URL")
	flag.Parse()
	startServer(*host, *port)
}
