package main

import (
	"flag"
	"fmt"
	"os"
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
		SubscribeTopic(conn, *topic)
		println("Subscribed to", *topic)
	case "publish":
		message := &os.Args[3]
		if *message == "" {
			fmt.Println("Invalid message")
			return
		}
		conn := DialServer(*address)
		PublishOnceTopic(conn, *topic, *message)
		println("Published to", *topic)
	default:
		print("unkown method...")
	}
}
