package main

import (
	"fmt"
	"time"

	"github.com/amar-jay/mini_ros/msgs"
	"github.com/amar-jay/mini_ros/node"
)

func main() {
	node := node.Init("simple_node")
	node.OnShutdown(func() {
		println("shutting down node")
	})

	/*NOTE:
	     * does not work with &new(msgs.Quaternion) / var t *msgs.Quaternion / var t msgs.Quaternion.
			 * It works only with the expression msg := &msgs.Quaternion{}.
			 * Still trying to figure out why.
			 * This is the only way to make it work.
			 * Other types like string, int, float32, float64, bool, etc. work fine with their respective expressions.
			 * Tried other ways but does not have type safety nearly as good as this.
	*/

	t := &msgs.Quaternion{}
	node.Callback(func() {
		fmt.Printf("%v\n", t)
		println(time.Now().String(), "callback called")
	})
	node.Subscribe("/chatter", t)
}
