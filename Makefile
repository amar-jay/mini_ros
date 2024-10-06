# Define binary output paths
ROS_CORE_BIN := bin/ros_core
TOPIC_BIN := bin/topic

# Define source files
ROS_CORE_SRC := cmd/ros_core.go
TOPIC_SRC := cmd/topic/*.go

# Default build command
all: build

# Build the ros_core and topic binaries
build: $(ROS_CORE_BIN) $(TOPIC_BIN)

$(ROS_CORE_BIN): $(ROS_CORE_SRC)
	@echo "Building Roscore binary..."
	@go build -o $(ROS_CORE_BIN) $(ROS_CORE_SRC)

$(TOPIC_BIN): $(TOPIC_SRC)
	@echo "Building Topic binary..."
	@go build -o $(TOPIC_BIN) $(TOPIC_SRC)

# Run the ros_core binary
ros_core: $(ROS_CORE_BIN)
	@echo "Running Roscore TCP server..."
	@$(ROS_CORE_BIN)

# Run the topic binary
topic: $(TOPIC_BIN)
	@echo "Running Topic binary..."
	@$(TOPIC_BIN) $(ARGS)

# Clean up build artifacts
clean:
	@echo "Cleaning up binaries..."
	@rm -f $(ROS_CORE_BIN) $(TOPIC_BIN)

.PHONY: all build ros_core topic clean

