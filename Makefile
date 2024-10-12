# Define binary output paths
ROS_CORE_BIN := bin/mini_ros
SHELL_TYPE := $(HOME)/.zshrc

# Define source files
ROS_CORE_SRC := cmd/main.go

# Default build command
all: build

# Build the ros_core and topic binaries
build: $(ROS_CORE_BIN) export

$(ROS_CORE_BIN): $(ROS_CORE_SRC)
	@echo "Building ROS binary..."
	@go build -o $(ROS_CORE_BIN) $(ROS_CORE_SRC)

# @echo 'export PATH=$(PWD)/bin:$$PATH' >> $(SHELL_SCRIPT)
export:
	@echo "Exporting MINI_ROS_BIN to PATH in $(SHELL_TYPE)"
	echo 'export PATH=$(PWD)/bin:$$PATH' >> $(SHELL_TYPE)

# Run the ros_core binary
ros_core: $(ROS_CORE_BIN)
	@echo "Running ROS binary..."
	@$(ROS_CORE_BIN)

# Clean up build artifacts
clean:
	@echo "Cleaning up binaries..."
	@rm -f $(ROS_CORE_BIN)

.PHONY: all build ros_core topic clean

