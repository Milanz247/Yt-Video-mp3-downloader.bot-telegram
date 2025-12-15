#!/bin/bash

# Start script for Telegram Media Downloader Bot

echo "🚀 Starting Telegram Media Downloader Bot..."

# Check if bot binary exists
if [ ! -f "./yt-bot" ]; then
    echo "⚠️  Bot binary not found. Building..."
    go build -o yt-bot
    
    if [ $? -ne 0 ]; then
        echo "❌ Build failed. Please check for errors."
        exit 1
    fi
    echo "✅ Build successful"
fi

# Check if .env exists
if [ ! -f ".env" ]; then
    echo "❌ .env file not found. Please create it from .env.example"
    exit 1
fi

# Create downloads directory
mkdir -p downloads

echo "✅ Starting bot..."
echo ""
./yt-bot
