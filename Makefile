
ROOT_DIR := $(dir $(abspath $(lastword $(MAKEFILE_LIST))))
BINARY_NAME := $(ROOT_DIR)logz
CMD_DIR := $(ROOT_DIR)cmd
SCRIPT=$(ROOT_DIR)scripts/install.sh

# Alvo para build
build:
	go build -ldflags "-s -w -X main.version=$(git describe --tags) -X main.commit=$(git rev-parse HEAD) -X main.date=$(date +%Y-%m-%d)" -trimpath -o $(BINARY_NAME) ${CMD_DIR}

# Chamar o script de instalação
install:
	sh $(SCRIPT) $(BINARY_NAME)

# Limpar o binário gerado
clean:
	rm -f $(BINARY_NAME)
