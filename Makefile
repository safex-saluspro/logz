ROOT_DIR := $(dir $(abspath $(lastword $(MAKEFILE_LIST))))
BINARY_NAME := $(ROOT_DIR)logz
CMD_DIR := $(ROOT_DIR)cmd
INSTALL_SCRIPT=$(ROOT_DIR)scripts/install.sh
ARGS :=

# Target for build
build:
	go build -ldflags "-s -w -X main.version=$(git describe --tags) -X main.commit=$(git rev-parse HEAD) -X main.date=$(date +%Y-%m-%d)" -trimpath -o $(BINARY_NAME) ${CMD_DIR} &&\
    upx $(BINARY_NAME) --force-overwrite --lzma --no-progress --no-color

# Target for development build (without compression)
build-dev:
	go build -ldflags "-s -w -X main.version=$(git describe --tags) -X main.commit=$(git rev-parse HEAD) -X main.date=$(date +%Y-%m-%d)" -trimpath -o $(BINARY_NAME) ${CMD_DIR}

# Target for installation
install:
	sh $(INSTALL_SCRIPT) install $(ARGS)

# Clean the generated binary
clean:
	rm -f $(BINARY_NAME)

# Help target
help:
	@echo "Available targets:"
	@echo "  make build      - Build the binary using install script"
	@echo "  make build-dev  - Build the binary without compression"
	@echo "  make install    - Install the binary and configure environment"
	@echo "  make clean      - Clean up build artifacts"
	@echo "  make help       - Display this help message"
	@echo ""
	@echo "Usage with arguments:"
	@echo "  make install ARGS='--custom-arg value' - Pass custom arguments to the install script"