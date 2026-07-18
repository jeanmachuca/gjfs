#!/bin/bash

# gjfs Installation Script
# Installs gjfs binary to /usr/local/bin or ~/.local/bin

set -e

VERSION="v0.1.0"
REPO="jeanmachuca/gjfs"
BINARY="gjfs"

# Colors
RED='\033[0;31m'
GREEN='\033[0;32m'
YELLOW='\033[1;33m'
NC='\033[0m' # No Color

info() {
    echo -e "${GREEN}[INFO]${NC} $1"
}

warn() {
    echo -e "${YELLOW}[WARN]${NC} $1"
}

error() {
    echo -e "${RED}[ERROR]${NC} $1"
    exit 1
}

# Detect OS and architecture
detect_platform() {
    OS=$(uname -s | tr '[:upper:]' '[:lower:]')
    ARCH=$(uname -m)

    case $ARCH in
        x86_64) ARCH="amd64" ;;
        aarch64|arm64) ARCH="arm64" ;;
        armv7l) ARCH="armv7" ;;
        *) error "Unsupported architecture: $ARCH" ;;
    esac

    case $OS in
        linux|darwin) ;;
        *) error "Unsupported OS: $OS" ;;
    esac

    info "Detected platform: $OS/$ARCH"
}

# Check if Go is installed
check_go() {
    if ! command -v go &> /dev/null; then
        error "Go is not installed. Please install Go 1.21+ first."
    fi
    GO_VERSION=$(go version | awk '{print $3}' | sed 's/go//')
    info "Go version: $GO_VERSION"
}

# Install from source
install_from_source() {
    info "Installing from source..."

    # Create temp directory
    TMP_DIR=$(mktemp -d)
    trap "rm -rf $TMP_DIR" EXIT

    cd "$TMP_DIR"

    # Clone repository
    info "Cloning repository..."
    git clone --depth 1 --branch "$VERSION" "https://github.com/$REPO.git" . 2>/dev/null || \
        git clone --depth 1 "https://github.com/$REPO.git" .

    # Build
    info "Building $BINARY..."
    go build -ldflags="-s -w -X main.version=$VERSION" -o "$BINARY" ./cmd/gjfs

    # Install
    install_binary "$TMP_DIR/$BINARY"
}

# Install binary
install_binary() {
    local binary_path=$1

    # Determine install directory
    if [ -w "/usr/local/bin" ]; then
        INSTALL_DIR="/usr/local/bin"
    elif [ -w "$HOME/.local/bin" ]; then
        INSTALL_DIR="$HOME/.local/bin"
    else
        INSTALL_DIR="/usr/local/bin"
        warn "Installing to $INSTALL_DIR requires sudo"
        sudo install -m 755 "$binary_path" "$INSTALL_DIR/$BINARY"
        info "Installed to $INSTALL_DIR/$BINARY (with sudo)"
        return
    fi

    install -m 755 "$binary_path" "$INSTALL_DIR/$BINARY"
    info "Installed to $INSTALL_DIR/$BINARY"

    # Check if install dir is in PATH
    if [[ ":$PATH:" != *":$INSTALL_DIR:"* ]]; then
        warn "$INSTALL_DIR is not in your PATH"
        warn "Add it to your PATH: export PATH=\$PATH:$INSTALL_DIR"
    fi
}

# Verify installation
verify_installation() {
    if command -v "$BINARY" &> /dev/null; then
        info "Installation verified!"
        "$BINARY" --version
    else
        error "Installation failed: $BINARY not found in PATH"
    fi
}

# Main
main() {
    echo "==================================="
    echo "  gjfs Installation Script"
    echo "==================================="
    echo ""

    detect_platform
    check_go
    install_from_source
    verify_installation

    echo ""
    info "Installation complete!"
    echo ""
    echo "Usage examples:"
    echo "  gjfs -schema schema.json"
    echo "  gjfs -schema-string '{\"type\": \"object\"}'"
    echo "  gjfs -schema schema.json -output example.json"
    echo "  gjfs -schema schema.json -validate data.json"
}

main "$@"
