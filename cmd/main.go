package main

import (
	"encoding/json"
	"log"
	"os"
	"time"

	"github.com/amar-jay/mini_ros/core"
	"github.com/amar-jay/mini_ros/topic"
	"github.com/urfave/cli/v2"
)

func main() {
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
						//Required: true,
					},

					&cli.StringFlag{
						Name:  "host",
						Value: "0.0.0.0",
						Usage: "ROS master host",
						//Required: true,
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
						//Required: true,
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
								Value:   "{\"message\": \"hello_world\"}",
								Usage:   "Message to send",
								//Required: true,
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
								topic.PublishOnceTopic(conn, cCtx.Args().Get(0), msg)
							} else {
								topic.PublishTopic(conn, cCtx.Args().Get(0), msg, 5*time.Second)
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
							topic.SubscribeTopic(conn, cCtx.Args().Get(0))
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
		/*
			Action: func(*cli.Context) error {
				fmt.Println("Invalid command. Use 'publish', 'subscribe', or 'status'")
				return nil
			},
		*/
	}

	if err := app.Run(os.Args); err != nil {
		log.Fatal(err)
	}
}
