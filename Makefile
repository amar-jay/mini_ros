# Define binary output paths
ROS_CORE_BIN := bin/mini_ros
SHELL_TYPE := $(HOME)/.zshrc

# Define source files
ROS_CORE_SRC := cmd/main.go

# Function to check if the directory is in PATH
define check_in_path
	if ! echo "$(PATH)" | grep -q "$(PWD)/bin"; then \
		echo "Exporting MINI_ROS_BIN to PATH in $(SHELL_TYPE)"; \
		echo 'export PATH=$(PWD)/bin:$$PATH' >> $(SHELL_TYPE); \
		zsh; \
	else \
		echo "$(PWD)/bin is already in PATH."; \
	fi
endef


# Default build command
all: build export

# Build the ros_core and topic binaries
build: $(ROS_CORE_BIN)

$(ROS_CORE_BIN): $(ROS_CORE_SRC)
	@echo "Building ROS binary..."
	@go build -o $(ROS_CORE_BIN) $(ROS_CORE_SRC)


# Run the ros_core binary
ros_core: $(ROS_CORE_BIN)
	@echo "Running ROS binary..."
	@$(ROS_CORE_BIN)

# Clean up build artifacts
clean:
	@echo "Cleaning up binaries..."
	@rm -f $(ROS_CORE_BIN)

export:
	@$(call check_in_path)

.PHONY: all build ros_core topic clean
