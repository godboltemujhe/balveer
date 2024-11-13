#!/bin/bash
set -e

echo "Updating package list..."

# macOS ke liye `brew` update command
if ! command -v brew &> /dev/null; then
    echo "Homebrew not found. Installing Homebrew..."
    /bin/bash -c "$(curl -fsSL https://raw.githubusercontent.com/Homebrew/install/HEAD/install.sh)"
else
    echo "Updating Homebrew..."
    brew update
fi

# Go language install karna agar nahi hai toh
if ! command -v go &> /dev/null; then
    echo "Installing Go..."
    brew install go
else
    echo "Go is already installed."
fi

# Go module initialize karna agar go.mod file nahi hai
if [ ! -f go.mod ]; then
    echo "Initializing Go module..."
    go mod init ranbal-telegram-bot
else
    echo "Go module already initialized."
fi

echo "Installing Telegram Bot API package..."
go get -u github.com/go-telegram-bot-api/telegram-bot-api

# macOS me sabhi files ko executable banane ke liye
chmod +x *

echo "Setup completed successfully!"