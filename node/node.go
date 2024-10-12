package node

import (
	"net"
	"os"
	"os/signal"
	"syscall"

	"github.com/amar-jay/mini_ros/cmd/topic"
)

type Node struct {
	Name       string
	onshutdown func()
	callback   func() // for subscribers
	conn       net.Conn
}

/** To create a new node */
func Init(name string) *Node {

	n := &Node{
		Name: name,
	}

	sig := make(chan os.Signal)
	signal.Notify(sig, os.Interrupt, syscall.SIGTERM)

	go func() {
		<-sig
		if n.onshutdown != nil {
			println("shutting down node: ", n.Name)
			n.onshutdown()
		}
	}()

	n.conn = topic.DialServer("localhost:11311")
	return n
}

func (n *Node) OnShutdown(f func()) {
	n.onshutdown = f
}

func (n *Node) Callback(f func()) {
	n.callback = f
}
func (p *Node) Publish(_topic string, msg interface{}) {
	println("publishing message: ", msg)
	topic.Publish(p.conn, _topic, msg)
}

func (s *Node) Subscribe(_topic string, msg interface{}) {
	topic.Subscribe(s.conn, _topic, msg, s.callback)
}

/*
 TODO: until there is a need to have different publishers and subscribers, there is no need for generics

type Publisher[T any] struct {
	Topic string
	Msg   T
	*Node
}

type Subscriber[T any] struct {
	Topic string
	Msg   T
	*Node
}

func (p *Publisher[T]) Publish(_topic string, msg T) {
	println("publishing message: ", msg)
	topic.Publish(p.conn, _topic, msg, p.rate)
}

func (s *Subscriber[T]) Subscribe(_topic string, msg T) {
	c := func() {
		println("subscribing to topic: ", _topic)
	}
	topic.Subscribe(s.conn, _topic, msg, c)
}
*/
