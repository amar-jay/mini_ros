package main

import "github.com/amar-jay/mini_ros/node"

func main() {
	node := node.Init("simple_node")
	node.OnShutdown(func() {
		println("shutting down node")
	})

	node.Publish("/chatter", "Hello World")
}
