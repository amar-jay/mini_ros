package main

import (
	"encoding/json"
	"log"
	"os"
	"reflect"
	"time"

	"github.com/amar-jay/mini_ros/core"
	"github.com/amar-jay/mini_ros/msgs"
	"github.com/amar-jay/mini_ros/topic"
	"github.com/urfave/cli/v2"
)

func main() {

	demoMsg := new(msgs.DemoMsg)
	demoMsg.Message = "Hello Mini ROS!"

	demoMsgBytes, _ := json.Marshal(demoMsg)

	app := &cli.App{
		Name:                 "mini_ros",
		EnableBashCompletion: true,
		Usage:                "a simple ROS implementation in Go for educational purposes",
		Commands: []*cli.Command{
			{
				Name:        "core",
				Usage:       "start a ROS core server",
				Subcommands: []*cli.Command{},
				Flags: []cli.Flag{
					&cli.StringFlag{
						Name:  "port",
						Value: "11311",
						Usage: "ROS master port",
					},

					&cli.StringFlag{
						Name:  "host",
						Value: "0.0.0.0",
						Usage: "ROS master host",
					},
				},
				Action: func(cCtx *cli.Context) error {
					host := cCtx.String("host")
					port := cCtx.Int("port")

					r := core.NewRosCore()
					r.Listen(host, port)
					return nil
				},
			},
			{
				Name:        "node",
				Usage:       "run node methods",
				Subcommands: []*cli.Command{},
			},
			{
				Name:  "topic",
				Usage: "run topic methods",

				Flags: []cli.Flag{

					&cli.StringFlag{
						Name:    "address",
						Aliases: []string{"add", "a"},
						Value:   "localhost:11311",
						Usage:   "ROS master host",
					},
				},
				Subcommands: []*cli.Command{
					{
						Name:     "publish",
						Category: "topic",
						Aliases:  []string{"pub"},
						Usage:    "publish a ROS topic",

						Flags: []cli.Flag{
							&cli.StringFlag{
								Name:    "message",
								Aliases: []string{"msg"},
								Value:   string(demoMsgBytes),
								Usage:   "Message to send",
							},
							&cli.BoolFlag{
								Name:    "once",
								Aliases: []string{"o"},
								Value:   false,
								Usage:   "Publish message once",
							},
						},
						Action: func(cCtx *cli.Context) error {

							if cCtx.NArg() == 0 {
								log.Fatal("Topic name is required")
							}
							message := cCtx.String("message")

							conn := topic.DialServer(cCtx.String("address"))

							var msg interface{}
							err := json.Unmarshal([]byte(message), &msg)
							if err != nil {
								log.Fatal("Unable to unmarshal message")
							}

							if cCtx.Bool("once") {
								topic.Publish(conn, cCtx.Args().Get(0), msg)
							} else {
								for {
									topic.Publish(conn, cCtx.Args().Get(0), msg)
									time.Sleep(5 * time.Second)
								}
							}

							return nil
						},
					},
					{
						Name:     "subscribe",
						Category: "topic",
						Aliases:  []string{"sub"},
						Usage:    "subscribe to a ROS topic",
						Action: func(cCtx *cli.Context) error {
							if cCtx.NArg() == 0 {
								log.Fatal("Topic name is required")
							}
							conn := topic.DialServer(cCtx.String("address"))
							msg := msgs.DemoMsg{}
							_topic := cCtx.Args().Get(0)
							callback := func() {
								log.Printf("%s>%s (type:%s)\n", _topic, msg.Message, reflect.TypeOf(msg))
							}

							topic.Subscribe(conn, _topic, &msg, callback)
							return nil
						},
					},
					{
						Name:     "status",
						Aliases:  []string{"stats", "stat"},
						Category: "topic",
						Usage:    "get stats of a ROS topic",
						Action: func(cCtx *cli.Context) error {
							if cCtx.NArg() == 0 {
								log.Fatal("Topic name is required")
							}
							conn := topic.DialServer(cCtx.String("address"))
							topic.SubscribeStatus(conn, cCtx.Args().Get(0))
							return nil
						},
					},
					{
						Name:     "list",
						Category: "topic",
						Usage:    "get list of all topics",
						Action: func(cCtx *cli.Context) error {
							conn := topic.DialServer(cCtx.String("address"))
							topic.List(conn)
							return nil
						},
					},
				},
			},
		},
	}

	if err := app.Run(os.Args); err != nil {
		log.Fatal(err)
	}
}
