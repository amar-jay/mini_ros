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

	node.Callback(func(topic string, message msgs.ROS_MSG) {
		message, ok := message.(msgs.Quaternion)
		if !ok {
			println("failed to cast to Quaternion")
		}
		println(time.Now().String(), "callback called")
	})
	var t msgs.Quaternion
	node.Subscribe("/chatter", t)
}
