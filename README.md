This is a lightweight reimplementation of core ROS concepts, focusing on `roscore` and topic-based communication for subscribing and publishing. Built entirely in Go without external libraries, it mimics essential ROS behavior in a minimalistic way. The `roscore` server manages message exchanges between nodes, while topics enable asynchronous communication.

This aims to mimic the essential behavior of ROS in a minimalistic way, making it easier to understand the underlying mechanisms while maintaining flexibility and performance due to Go’s concurrency model.

```
NAME:
   mini_ros - a simple ROS implementation in Go for educational purposes

USAGE:
   mini_ros [global options] command [command options]

COMMANDS:
   core                 start a ROS core server
   publish, pub         publish a ROS topic
   subscribe, sub       subscribe to a ROS topic
   status, stats, stat  get stats of a ROS topic
   help, h              Shows a list of commands or help for one command

GLOBAL OPTIONS:
   --help, -h  show help

\---


NAME:
   mini_ros core - start a ROS core server

USAGE:
   mini_ros core [command options]

OPTIONS:
   --port value  ROS master port (default: "11311")
   --host value  ROS master host (default: "0.0.0.0")
   --help, -h    show help


\---

NAME:
   mini_ros publish - publish a ROS topic

USAGE:
   mini_ros publish [command options]

OPTIONS:
   --address value, --add value, -a value  ROS master host (default: "localhost:11311")
   --message value, --msg value            Message to send (default: "{\"message\": \"hello_world\"}")
   --once, -o                              Publish message once (default: false)
   --help, -h                              show help


\---

NAME:
   mini_ros subscribe - subscribe to a ROS topic

USAGE:
   mini_ros subscribe [command options]

OPTIONS:
   --address value, --add value, -a value  ROS master host (default: "localhost:11311")
   --help, -h                              show help


NAME:
   mini_ros status - get stats of a ROS topic

USAGE:
   mini_ros status [command options]

OPTIONS:
   --address value, --add value, -a value  ROS master host (default: "localhost:11311")
   --help, -h                              show help

```

| commands  | purposes                                        |
| --------- | ----------------------------------------------- |
| core      | To start roscore server on master url as in ROS |
| subscribe | To subscribe to a topic                         |
| publish   | To publish a topic                              |
| status    | To get stats of a topic                         |

### TODO

- [x] ROS core
- [x] Publish topic
- [x] Subscribe to topic
- [x] get topic metrics
- [x] better CLI
- [ ] Create more realistic topic `/cmd_vel` or `/raw_image`
- [ ] ROS Node
- [ ] ROS simple client library
- [ ] ROS service
- [ ] ROS launch file
