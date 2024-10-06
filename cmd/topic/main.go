package main

import (
	"encoding/json"
	"flag"
	"fmt"
	"os"
	"time"
)

func main() {
	// get first argument
	method := os.Args[1]
	topic := &os.Args[2]
	//topic := flag.String("topic", "", "Topic to publish/subscribe to")
	address := flag.String("address", "localhost:11311", "Server address")
	//message := flag.String("message", "", "Message to send")

	flag.Parse()

	if *topic == "" {
		fmt.Println("Invalid topic name")
		return
	}

	switch method {
	case "subscribe":
		conn := DialServer(*address)
		println("Subscribed to", *topic)
		SubscribeTopic(conn, *topic)
	case "publish":
		message := &os.Args[3]
		if *message == "" {
			fmt.Println("Invalid message")
			return
		}
		conn := DialServer(*address)

		var msg interface{}
		err := json.Unmarshal([]byte(*message), &msg)
		if err != nil {
			fmt.Println("Unable to unmarshal message")
			return
		}
		PublishTopic(conn, *topic, msg, 5*time.Second)
	case "status":
		conn := DialServer(*address)
		println("Subscribed to", *topic)
		SubscribeStatus(conn, *topic)
	default:
		print("unknown method...")
	}
}
