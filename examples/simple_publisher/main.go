package main

import (
	"github.com/amar-jay/mini_ros/msgs"
	"github.com/amar-jay/mini_ros/node"
)

func main() {
	node := node.Init("simple_node")
	node.OnShutdown(func() {
		println("shutting down node")
	})

	msg := msgs.Quaternion{
		X: 0.1,
		Y: 0.2,
		Z: 0.3,
		W: 0.4,
	}
	//msg := "Hello World"
	node.Publish("/chatter", msg)
}
