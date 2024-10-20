<div id="header" align="center">
   <h1 id="badges">
            `mini_ros`
  </h1>
</div>

This is a lightweight reimplementation of core ROS concepts, focusing on __`roscore`__ and topic-based communication for subscribing and publishing. Built entirely in Go without external libraries, it mimics essential ROS behavior in a minimalistic way. The `roscore` server manages message exchanges between nodes, while topics enable asynchronous communication.

This aims to mimic the essential behavior of ROS in a minimalistic way, making it easier to understand the underlying mechanisms while maintaining flexibility and performance due to Go’s concurrency model.

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
- [x] message types
- [x] get topic metrics
- [x] better CLI
- [ ] Create more realistic topic `/cmd_vel` or `/raw_image`
- [x] ROS Node
- [x] ROS simple client library
- [ ] ROS service
- [ ] ROS launch file

### HOW TO USE

#### to build it,
1. *Linux*: simply run. **Note SHELL_TYPE** in Makefile may be different verify it if not working
```
make build
```

2. *MacOs*: change the __SHELL_TYPE__, to location of shell script (ie `~/.bashrc`) in the **Makefile**, then run `make build`
3. *Windows*: use WSL, and follow Linux setup

#### commands
You can check it for yourself using. 
```
mini_ros --help
```

