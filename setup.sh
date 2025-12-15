#!/bin/bash

# Installation script for Telegram Media Downloader Bot

echo "🚀 Setting up Telegram Media Downloader Bot..."
echo ""

# Check if Go is installed
if ! command -v go &> /dev/null; then
    echo "❌ Go is not installed. Please install Go 1.21 or higher."
    echo "Visit: https://golang.org/dl/"
    exit 1
fi

echo "✅ Go is installed: $(go version)"

# Check if yt-dlp is installed
if ! command -v yt-dlp &> /dev/null; then
    echo "⚠️  yt-dlp is not installed. Installing..."
    
    if [[ "$OSTYPE" == "linux-gnu"* ]]; then
        sudo wget https://github.com/yt-dlp/yt-dlp/releases/latest/download/yt-dlp -O /usr/local/bin/yt-dlp
        sudo chmod a+rx /usr/local/bin/yt-dlp
        echo "✅ yt-dlp installed successfully"
    elif [[ "$OSTYPE" == "darwin"* ]]; then
        if command -v brew &> /dev/null; then
            brew install yt-dlp
            echo "✅ yt-dlp installed successfully"
        else
            echo "❌ Homebrew not found. Please install yt-dlp manually:"
            echo "   pip install yt-dlp"
            exit 1
        fi
    else
        echo "❌ Please install yt-dlp manually:"
        echo "   pip install yt-dlp"
        exit 1
    fi
else
    echo "✅ yt-dlp is installed: $(yt-dlp --version)"
fi

# Check if ffmpeg is installed
if ! command -v ffmpeg &> /dev/null; then
    echo "⚠️  ffmpeg is not installed. Installing..."
    
    if [[ "$OSTYPE" == "linux-gnu"* ]]; then
        sudo apt-get update
        sudo apt-get install -y ffmpeg
        echo "✅ ffmpeg installed successfully"
    elif [[ "$OSTYPE" == "darwin"* ]]; then
        if command -v brew &> /dev/null; then
            brew install ffmpeg
            echo "✅ ffmpeg installed successfully"
        else
            echo "❌ Homebrew not found. Please install ffmpeg manually"
            exit 1
        fi
    else
        echo "❌ Please install ffmpeg manually"
        exit 1
    fi
else
    echo "✅ ffmpeg is installed"
fi

# Create .env file if it doesn't exist
if [ ! -f .env ]; then
    echo ""
    echo "📝 Creating .env file..."
    cp .env.example .env
    echo "✅ .env file created"
    echo "⚠️  Please update the .env file with your Telegram bot token if needed"
else
    echo "✅ .env file already exists"
fi

# Download Go dependencies
echo ""
echo "📦 Downloading Go dependencies..."
go mod download
echo "✅ Dependencies downloaded"

# Build the bot
echo ""
echo "🔨 Building the bot..."
go build -o yt-bot
echo "✅ Bot built successfully"

# Create downloads directory
mkdir -p downloads

echo ""
echo "🎉 Setup complete!"
echo ""
echo "To run the bot:"
echo "  ./yt-bot"
echo ""
echo "Or:"
echo "  go run main.go"
echo ""
