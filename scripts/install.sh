#!/bin/sh

CMD_PATH="$(realpath "$(dirname "$0")")"
BUILD_PATH="$(dirname "$CMD_PATH")"
BINARY=${1:-"$BUILD_PATH/logz"}
LOCAL_BIN="$HOME/.local/bin"
GLOBAL_BIN="/usr/local/bin"

GREEN="\033[0;32m"
YELLOW="\033[0;33m"
RED="\033[0;31m"
NC="\033[0m"

detect_shell_rc() {
    shell_rc_file=""
    user_shell=$(basename "$SHELL")
    case "$user_shell" in
        bash) shell_rc_file="$HOME/.bashrc" ;;
        zsh) shell_rc_file="$HOME/.zshrc" ;;
        sh) shell_rc_file="$HOME/.profile" ;;
        fish) shell_rc_file="$HOME/.config/fish/config.fish" ;;
        *)
            echo "${YELLOW}Warning: Unsupported shell ($user_shell). Modify PATH manually.${NC}"
            return 1
            ;;
    esac
    echo "$shell_rc_file"
}

add_to_path() {
    target_path="$1"
    shell_rc_file=$(detect_shell_rc)
    if [ -z "$shell_rc_file" ]; then
        echo "${RED}Error: Could not determine shell configuration file.${NC}"
        return 1
    fi

    if grep -q "export PATH=.*$target_path" "$shell_rc_file" 2>/dev/null; then
        echo "${GREEN}✅ $target_path is already in $shell_rc_file.${NC}"
        return 0
    fi

    echo "export PATH=$target_path:\$PATH" >> "$shell_rc_file"
    echo "${GREEN}Added $target_path to PATH in $shell_rc_file.${NC}"
    echo "Run 'source $shell_rc_file' to apply changes."
}

install_binary() {
    if [ "$(id -u)" -ne 0 ]; then
        echo "You are not root. Installing in $LOCAL_BIN..."
        mkdir -p "$LOCAL_BIN"
        cp "$BINARY" "$LOCAL_BIN"
        add_to_path "$LOCAL_BIN"
    else
        echo "Root detected. Installing in $GLOBAL_BIN..."
        cp "$BINARY" "$GLOBAL_BIN"
        add_to_path "$GLOBAL_BIN"
    fi
}

install_upx() {
    if ! command -v upx > /dev/null; then
        echo "${YELLOW}Installing UPX...${NC}"
        if [ "$(uname)" = "Darwin" ]; then
            brew install upx
        elif command -v apt-get > /dev/null; then
            sudo apt-get install -y upx
        else
            echo "${RED}Install UPX manually from https://upx.github.io/${NC}"
            exit 1
        fi
    else
        echo "${GREEN}✅ UPX is already installed.${NC}"
    fi
}

check_dependencies() {
    for dep in "$@"; do
        if ! command -v "$dep" > /dev/null; then
            echo "${RED}Error: $dep is not installed.${NC}"
            exit 1
        else
            echo "${GREEN}✅ $dep is installed.${NC}"
        fi
    done
}

build_binary() {
    echo "Building the binary..."
    go build -ldflags "-s -w -X main.version=$(git describe --tags) -X main.commit=$(git rev-parse HEAD) -X main.date=$(date +%Y-%m-%d)" -trimpath -o "$BINARY" "$CMD_PATH"
    install_upx
    upx "$BINARY" --force-overwrite --lzma --no-progress --no-color
}

clean() {
    echo "Cleaning up build artifacts..."
    rm -f "$BINARY"
}

validate_versions() {
    REQUIRED_GO_VERSION="1.18"
    GO_VERSION=$(go version | awk '{print $3}' | sed 's/go//')
    if [ "$(printf '%s\n' "$REQUIRED_GO_VERSION" "$GO_VERSION" | sort -V | head -n1)" != "$REQUIRED_GO_VERSION" ]; then
        echo "${RED}Error: Go version must be >= $REQUIRED_GO_VERSION. Detected: $GO_VERSION${NC}"
        exit 1
    fi
    echo "${GREEN}✅ Go version is valid: $GO_VERSION${NC}"
}

summary() {
    install_dir=$(dirname "$BINARY")
    echo "${GREEN}Build and installation complete!${NC}"
    echo "Binary: $BINARY"
    echo "Installed in: $install_dir"
    check_path "$install_dir"
}

build_and_validate() {
    validate_versions
    build_binary
}

case "$1" in
    build)
        build_and_validate
        ;;
    install)
        build_and_validate
        install_binary
        summary
        ;;
    clean)
        clean
        ;;
    *)
        echo "Usage: $0 {build|install|clean}"
        exit 1
        ;;
esac
