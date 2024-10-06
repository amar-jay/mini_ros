package main

import (
	"bufio"
	"fmt"
	"net"
	"os"
	"strings"
	"time"
)

type Topic struct {
	Name     string
	Messages chan interface{}
}

func (t *Topic) Publish(msg interface{}) {
	t.Messages <- msg
}

func (t *Topic) Subscribe() chan interface{} {
	return t.Messages
}

// ConnectToServer connects to the TCP Roscore server
func DialServer(address string) net.Conn {
	conn, err := net.Dial("tcp", address)
	if err != nil {
		fmt.Println("Error connecting to server:", err)
		os.Exit(1)
	}
	return conn
}

// HandleServerResponse listens for messages from the server
func handleSubscribe(conn net.Conn, topic string) {
	reader := bufio.NewReader(conn)
	for {
		message, err := reader.ReadString('\n')
		if err != nil {
			fmt.Println("Server disconnected.")
			os.Exit(1)
		}
		message = strings.TrimSpace(message)
		m := strings.SplitN(message, " ", 2)
		if len(m) < 2 {
			fmt.Println("ERROR< Invalid message from server:", m, len(m))
		}

		_topic, _message := m[0], m[1]
		if topic == _topic {
			fmt.Println(_topic, "> ", _message)
		} else {
			fmt.Println("ERROR< Unkown response from server:", _topic, _message)
		}
	}
}

// RunPublisher publishes messages to the Roscore server.
func PublishTopic(conn net.Conn, topic string, message interface{}, timeInterval time.Duration) {
	for {
		fmt.Print(topic, "< ", message, "\n")

		// Send PUBLISH command to server
		fmt.Fprintf(conn, "PUBLISH %s %s\n", topic, message)
		time.Sleep(timeInterval)
	}
}

// RunPublisher publishes messages to the Roscore server.
func PublishOnceTopic(conn net.Conn, topic string, message interface{}) {
	fmt.Print(topic, "< ", message, "\n")

	// Send PUBLISH command to server
	fmt.Fprintf(conn, "PUBLISH %s %s\n", topic, message)
}

// RunSubscriber subscribes to a topic on the Roscore server.
func SubscribeTopic(conn net.Conn, topic string) {
	// Send SUBSCRIBE command to server
	fmt.Fprintf(conn, "SUBSCRIBE %s\n", topic)

	// Listen for messages from the server
	handleSubscribe(conn, topic)

	//select {} // Block the main thread forever
}
