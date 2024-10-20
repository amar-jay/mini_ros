package topic

import (
	"bufio"
	"encoding/json"
	"fmt"
	"net"
	"os"
	"os/signal"
	"reflect"
	"strings"
	"syscall"

	"github.com/amar-jay/mini_ros/msgs"
)

type Topic struct {
	Name    string `json:"name"`
	Type    reflect.Type
	Message interface{} `json:"message,omitempty"`
}

type Status struct {
	Subscribers map[string]int `json:"subscribers"`
	Type        string         `json:"type"`
}

var topics = make([]Topic, 0)

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

func handleSubscribe(conn net.Conn, topic string, msg msgs.ROS_MSG, callback func(topic string, message msgs.ROS_MSG)) {
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

		if callback != nil {
			callback(topic, msg)
		}

		if strings.TrimSpace(_topic) == topic {
			fmt.Printf("%s> %s\n", _topic, message)
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

	var msg Status

	err = json.Unmarshal([]byte(message), &msg)
	if err != nil {
		fmt.Println("Unmarshal json error", err)
		return
	}
	fmt.Printf("STATUS> %s:\t\t%d Subscribers\t %s Type", topic, msg.Subscribers[topic], msg.Type)
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
		fmt.Println(topic.Name)
	}
}

func Publish(conn net.Conn, topic string, message msgs.ROS_MSG) {
	msg, err := json.Marshal(Topic{
		Name:    topic,
		Type:    reflect.TypeOf(message),
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

func Subscribe(conn net.Conn, topic string, msg msgs.ROS_MSG, callback func(topic string, message msgs.ROS_MSG)) {
	fmt.Fprintf(conn, "SUBSCRIBE %s\n", topic)

	handleSubscribe(conn, topic, msg, callback)
}

func List(conn net.Conn) {
	fmt.Fprintf(conn, "LIST\n")

	handleList(conn)
}

func SubscribeStatus(conn net.Conn, topic string) {
	fmt.Fprintf(conn, "STATUS %s\n", topic)

	handleStatus(conn, topic)
}
