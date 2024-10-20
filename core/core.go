package core

import (
	"bufio"
	"encoding/json"
	"fmt"
	"net"
	"strconv"
	"strings"
	"sync"

	t "github.com/amar-jay/mini_ros/topic"
)

// hmmm! A synchronous map will be more useful here. However, this is just a simple example
type RosCore struct {
	mu          sync.RWMutex
	Subscribers map[string][]net.Conn // map of topic to subscribers (connections)
	Types       map[string]string     // ros topic types map
}

func NewRosCore() *RosCore {
	return &RosCore{
		mu:          sync.RWMutex{},
		Subscribers: make(map[string][]net.Conn),
		Types:       make(map[string]string),
	}
}

func (r *RosCore) Listen(host string, port int) {
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
			// most likely the client disconnected, so we close the connection, else user must reconnect
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

		_type := "unknown"
		switch command {
		case "SUBSCRIBE":
			r.Subscribe(topic, _type, conn)
		case "UNSUBSCRIBE":
			r.Unsubscribe(topic, conn)
		case "PUBLISH":
			r.Publish(topic, conn)
		case "STATUS":
			r.Status(topic, conn)
		case "LIST":
			r.List(conn)
		default:
			conn.Write([]byte("Unknown command\n"))
		}
	}
}

func (r *RosCore) Subscribe(topic string, _type string, conn net.Conn) {
	r.mu.Lock()
	defer r.mu.Unlock()

	r.Subscribers[topic] = append(r.Subscribers[topic], conn)
	r.Types[topic] = _type // TODO: implement type checking

	// since a topic can have multiple subscribers, we keep track of all the subscribers in a slice.
	fmt.Println("Client", conn.RemoteAddr(), "subscribed to topic", topic)
	// there is no need to send a message to the client that they have subscribed successfully
	/*
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
	delete(r.Types, topic)

	// if empty delete
	if len(r.Subscribers[topic]) == 0 {
		delete(r.Subscribers, topic)
	}
	// no need to send a message to the client that they have unsubscribed successfully
}

func (r *RosCore) Publish(topic_message string, conn net.Conn) {
	r.mu.RLock()
	defer r.mu.RUnlock()

	parts := strings.SplitAfterN(topic_message, " ", 2)
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
		fmt.Printf("Published to topic %s (type %T) to %d subscribers\n", topic, message, len(r.Subscribers))
		for _, conn := range r.Subscribers[topic] {
			conn.Write([]byte(topic + " " + string(message) + "\n"))
		}
	}
}

func (r *RosCore) Status(topic string, conn net.Conn) {
	status := t.Status{Subscribers: map[string]int{topic: len(r.Subscribers[topic])}, Type: r.Types[topic]}
	for t, conns := range r.Subscribers {
		status.Subscribers[t] = len(conns)
	}
	st, err := json.Marshal(status)
	if err != nil {
		fmt.Println("Error marshalling status")
		return
	}

	conn.Write([]byte(string(st) + "\n"))

}

func (r *RosCore) List(conn net.Conn) {
	//println("Number of Subscribers: ", len(r.Subscribers[topic]))
	topics := make([]t.Topic, 0, len(r.Subscribers))
	for _t := range r.Subscribers {
		topics = append(topics, t.Topic{Name: _t}) // that is to assume list does not need to know the type
	}
	st, err := json.Marshal(topics)
	if err != nil {
		fmt.Println("Error marshalling status")
		return
	}

	conn.Write([]byte(string(st) + "\n"))
}
