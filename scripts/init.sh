#!/bin/bash
set -e

INTERNAL_DIR="/tmp/vimana/celestia"

# check if the binary is already installed
if [ -f "$INTERNAL_DIR/celestia" ]; then
    # Capture the version output
    VERSION_OUTPUT=$("$INTERNAL_DIR/celestia" version)
    
    # Check if the version matches "v0.11.0-rc15-dev"
    if [[ $VERSION_OUTPUT == *"v0.11.0-rc15"* ]]; then
        echo "🚀 Celestia is already installed with the correct version."
        exit 0
    else
        echo "🚀 Celestia is installed but with a different version."
    fi
fi

echo "🔍 Determining OS and architecture..."

OS=$(uname -s | tr '[:upper:]' '[:lower:]')
ARCH=$(uname -m)

if [[ "$ARCH" == "x86_64" ]]; then
    ARCH="amd64"
elif [[ "$ARCH" == "arm64" ]] || [[ "$ARCH" == "aarch64" ]]; then
    ARCH="arm64"
fi

echo "💻  OS: $OS"
echo "🏗️  ARCH: $ARCH"
TGZ_URL="https://github.com/Vistara-Labs/vimana/releases/download/celestia-v0.11.0-rc15/${OS}_${ARCH}.zip"

sudo mkdir -p "$INTERNAL_DIR"
sudo rm -rf "/tmp/vimcel"
sudo mkdir -p "/tmp/vimcel"
echo "💈 Downloading Celestia..."
sudo curl -o /tmp/vimcel/${OS}_${ARCH}.zip -L "$TGZ_URL" --progress-bar

# if OS is linux then install unzip
if [[ "$OS" == "linux" ]]; then
    # accept default yes
    sudo apt-get update > /dev/null
    sudo apt-get install unzip > /dev/null
fi
sudo unzip -q /tmp/vimcel/${OS}_${ARCH}.zip -d /tmp/vimcel/
sudo mv "/tmp/vimcel/${OS}_${ARCH}"/* "$INTERNAL_DIR"
sudo chmod +x "$INTERNAL_DIR"
sudo rm -rf "/tmp/vimcel"
mkdir -p ~/.vimana/celestia/light-node
chmod +x ~/.vimana/celestia/light-node
${INTERNAL_DIR}/celestia version
echo "💈 Celestia node version installed!"
