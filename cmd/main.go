package main

import (
	"encoding/json"
	"fmt"
	"log"
	"os"
	"time"

	"github.com/amar-jay/mini_ros/cmd/topic"
	"github.com/urfave/cli/v2"
)

func main() {
	app := &cli.App{
		Name:  "mini_ros",
		Usage: "a simple ROS implementation in Go for educational purposes",
		Commands: []*cli.Command{
			{
				Name:  "core",
				Usage: "start a ROS core server",
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

					r := NewRosCore()
					r.listen(host, port)
					return nil
				},
			},
			{
				Name:    "publish",
				Aliases: []string{"pub"},
				Usage:   "publish a ROS topic",

				Flags: []cli.Flag{
					&cli.StringFlag{
						Name:    "address",
						Aliases: []string{"add", "a"},
						Value:   "localhost:11311",
						Usage:   "ROS master host",
						//Required: true,
					},
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
					fmt.Println("topic name: ", cCtx.Args().First())

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
				Name:    "subscribe",
				Aliases: []string{"sub"},
				Usage:   "subscribe to a ROS topic",

				Flags: []cli.Flag{
					&cli.StringFlag{
						Name:    "address",
						Aliases: []string{"add", "a"},
						Value:   "localhost:11311",
						Usage:   "ROS master host",
						//Required: true,
					},
				},
				Action: func(cCtx *cli.Context) error {
					fmt.Println("topic name: ", cCtx.Args().First())
					if cCtx.NArg() == 0 {
						log.Fatal("Topic name is required")
					}
					conn := topic.DialServer(cCtx.String("address"))
					topic.SubscribeTopic(conn, cCtx.Args().Get(0))
					return nil
				},
			},

			{
				Name:    "status",
				Aliases: []string{"stats", "stat"},
				Usage:   "get stats of a ROS topic",

				Flags: []cli.Flag{
					&cli.StringFlag{
						Name:    "address",
						Aliases: []string{"add", "a"},
						Value:   "localhost:11311",
						Usage:   "ROS master host",
						//Required: true,
					},
				},
				Action: func(cCtx *cli.Context) error {
					fmt.Println("topic name: ", cCtx.Args().First())
					if cCtx.NArg() == 0 {
						log.Fatal("Topic name is required")
					}
					conn := topic.DialServer(cCtx.String("address"))
					topic.SubscribeStatus(conn, cCtx.Args().Get(0))
					return nil
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
