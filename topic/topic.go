package topic

import (
	"bufio"
	"encoding/json"
	"fmt"
	"net"
	"os"
	"os/signal"
	"strings"
	"syscall"
	"time"
)

type Topic struct {
	Name    string
	Message interface{}
}

var topics = make([]string, 0)

func (t *Topic) Publish(msg interface{}) {
	//t.Messages <- msg
	return
}

func (t *Topic) Subscribe() {
	//return t.Messages
	return
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

func handleUnsubscribe(conn net.Conn, topic string) {
	fmt.Fprintf(conn, "UNSUBSCRIBE %s\n", topic)
	println("unsubscribed successfully", topic)
}

func handleSubscribeTopic(conn net.Conn, topic string) {
	reader := bufio.NewReader(conn)
	sig := make(chan os.Signal)
	signal.Notify(sig, os.Interrupt, syscall.SIGTERM)
	go func() {
		<-sig
		handleUnsubscribe(conn, topic)
		os.Exit(1)
	}()

	for {
		println("waiting for message...")
		message, err := reader.ReadString('\n')
		if err != nil {
			fmt.Println("Server disconnected.")
			return
		}
		message = strings.TrimSpace(message)
		m := strings.SplitAfterN(message, " ", 2)
		if len(m) < 2 {
			fmt.Println("ERROR< Invalid message from server:", m, len(m))
		}

		_topic, message := m[0], m[1]
		message = strings.TrimSpace(message)
		if len(message) == 0 {
			continue
		}

		var msg interface{}
		err = json.Unmarshal([]byte(message), &msg)
		if err != nil {
			fmt.Println("Unmarshal json error", err)
			continue
		}
		if strings.TrimSpace(_topic) == topic {
			fmt.Printf("%s > %s\n", _topic, message)
		}
	}
}

func handleSubscribe[T any](conn net.Conn, topic string, msg T, callback func()) {
	reader := bufio.NewReader(conn)
	sig := make(chan os.Signal)
	signal.Notify(sig, os.Interrupt, syscall.SIGTERM)
	go func() {
		<-sig
		handleUnsubscribe(conn, topic)
		os.Exit(1)
	}()

	for {
		println("waiting for message...")
		message, err := reader.ReadString('\n')
		if err != nil {
			fmt.Println("Server disconnected.")
			return
		}
		message = strings.TrimSpace(message)
		m := strings.SplitAfterN(message, " ", 2)
		if len(m) < 2 {
			fmt.Println("ERROR> Invalid message from server:", m, len(m))
		}

		_topic, message := m[0], m[1]
		message = strings.TrimSpace(message)
		if len(message) == 0 {
			continue
		}

		err = json.Unmarshal([]byte(message), &msg)
		if err != nil {
			fmt.Println("Unmarshal json error", err)
			continue
		}
		if strings.TrimSpace(_topic) == topic {
			fmt.Printf("%s> %s\n", _topic, message)
		}
		if callback != nil {
			callback()
		}
	}
}

func handleStatus(conn net.Conn, topic string) {
	reader := bufio.NewReader(conn)
	message, err := reader.ReadString('\n')
	if err != nil {
		fmt.Println("Server disconnected.")
		return
	}
	message = strings.TrimSpace(message)
	if len(message) == 0 {
		return
	}
	var msg struct {
		Subscribers map[string]int `json:"subscribers"`
	}

	err = json.Unmarshal([]byte(message), &msg)
	if err != nil {
		fmt.Println("Unmarshal json error", err)
		return
	}
	fmt.Printf("STATUS> %s:\t\t%d Subscribers\n", topic, msg.Subscribers[topic])
}

func handleList(conn net.Conn) {
	reader := bufio.NewReader(conn)
	message, err := reader.ReadString('\n')
	if err != nil {
		fmt.Println("Server disconnected.")
		return
	}
	message = strings.TrimSpace(message)
	if len(message) == 0 {
		return
	}

	err = json.Unmarshal([]byte(message), &topics)
	if err != nil {
		fmt.Println("Unmarshal json error: ", message, "\n", err)
		return
	}
	for _, topic := range topics {
		fmt.Println(topic)
	}
}

// RunPublisher publishes messages to the Roscore server.
func PublishTopic(conn net.Conn, topic string, message interface{}, timeInterval time.Duration) {
	for {
		msg, err := json.Marshal(Topic{
			Name:    topic,
			Message: message,
		})
		println("waiting for publish...", conn.RemoteAddr().String())
		if err != nil {
			fmt.Printf("invalid message type. unable to parse message")
		}
		fmt.Printf("%s < %v\n", topic, message)

		// Send PUBLISH command to server
		fmt.Fprintf(conn, "PUBLISH %s %s\n", topic, msg)
		time.Sleep(timeInterval)
	}
}

func Publish(conn net.Conn, topic string, message interface{}) {
	msg, err := json.Marshal(Topic{
		Name:    topic,
		Message: message,
	})
	println("waiting for publish...", conn.RemoteAddr().String())
	if err != nil {
		fmt.Printf("invalid message type. unable to parse message")
	}
	fmt.Printf("%s < %v\n", topic, message)

	// Send PUBLISH command to server
	fmt.Fprintf(conn, "PUBLISH %s %s\n", topic, msg)
}

// RunPublisher publishes messages to the Roscore server.
func PublishOnceTopic(conn net.Conn, topic string, message interface{}) {
	msg, err := json.Marshal(Topic{
		Name:    topic,
		Message: message,
	})
	if err != nil {
		fmt.Printf("invalid message type. unable to parse message")
	}
	fmt.Printf("%s < %v\n", topic, message)

	// Send PUBLISH command to server
	fmt.Fprintf(conn, "PUBLISH %s %s\n", topic, msg)
}

func SubscribeTopic(conn net.Conn, topic string) {
	fmt.Fprintf(conn, "SUBSCRIBE %s\n", topic)

	handleSubscribeTopic(conn, topic)

	//select {} // Block the main thread forever
}

func Subscribe[T any](conn net.Conn, topic string, msg T, callback func()) {
	fmt.Fprintf(conn, "SUBSCRIBE %s\n", topic)

	handleSubscribe(conn, topic, msg, callback)

	//select {} // Block the main thread forever
}
func List(conn net.Conn) {
	fmt.Fprintf(conn, "LIST\n")

	handleList(conn)

}
func SubscribeStatus(conn net.Conn, topic string) {
	fmt.Fprintf(conn, "STATUS %s\n", topic)

	handleStatus(conn, topic)

}
