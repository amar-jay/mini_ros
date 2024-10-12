package main

import (
	"time"

	"github.com/amar-jay/mini_ros/msgs"
	"github.com/amar-jay/mini_ros/node"
)

func main() {
	node := node.Init("simple_node")
	node.OnShutdown(func() {
		println("shutting down node")
	})

	node.Callback(func() {
		println(time.Now().String(), "callback called")
	})
	var t msgs.String
	node.Subscribe("/chatter", t)
}
