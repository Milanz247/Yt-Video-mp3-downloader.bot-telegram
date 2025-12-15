# Quick Start Guide

## 🚀 Fast Setup (3 steps)

### 1. Install Prerequisites

**yt-dlp** (video downloader):
```bash
# Linux
sudo wget https://github.com/yt-dlp/yt-dlp/releases/latest/download/yt-dlp -O /usr/local/bin/yt-dlp
sudo chmod a+rx /usr/local/bin/yt-dlp

# Or with pip
pip install yt-dlp
```

**FFmpeg** (media processor):
```bash
# Ubuntu/Debian
sudo apt install ffmpeg

# Fedora
sudo dnf install ffmpeg

# macOS
brew install ffmpeg
```

### 2. Configure Bot Token

Your bot token is already set in `.env`:
```
TELEGRAM_BOT_TOKEN=8266088482:AAE1zZjqDQ4puqZKa5DgW-6tbJaeVE5YN6Q
```

### 3. Run the Bot

```bash
# Option 1: Run directly
go run main.go

# Option 2: Build and run
go build -o yt-bot
./yt-bot
```

## ✅ Verify Installation

Check if prerequisites are installed:
```bash
yt-dlp --version   # Should show version number
ffmpeg -version    # Should show version info
go version         # Should show Go 1.21+
```

## 📱 Using the Bot

1. Open Telegram and search for your bot
2. Send `/start` to begin
3. Send a video link (YouTube, Facebook, or Instagram)
4. Choose quality from the buttons
5. Receive your downloaded file!

## 🔥 Supported Platforms

- ✅ **YouTube**: Regular videos, shorts, live streams
- ✅ **Facebook**: Public videos and posts
- ✅ **Instagram**: Posts, reels, and stories (public)

## 🎯 Quality Options

**Video:**
- Best Quality (highest available)
- 1080p (Full HD)
- 720p (HD)
- 480p (SD)
- 360p (Low bandwidth)

**Audio (MP3):**
- Best Quality
- 320kbps (Highest)
- 192kbps (Good)
- 128kbps (Standard)

## 🆘 Quick Troubleshooting

**Bot doesn't start?**
- Check `.env` has correct token
- Run `go mod download`

**Download fails?**
- Update yt-dlp: `yt-dlp -U`
- Check if link is accessible
- Try different quality

**File too large?**
- Select lower quality
- Telegram limit: 50MB

## 📊 Example Commands

```bash
# Run with logs
go run main.go

# Build optimized binary
go build -ldflags="-s -w" -o yt-bot

# Run in background
nohup ./yt-bot > bot.log 2>&1 &

# Stop background process
pkill yt-bot
```

## 🔧 Advanced Configuration

Run as systemd service (Linux):

1. Create `/etc/systemd/system/yt-bot.service`:
```ini
[Unit]
Description=Telegram Media Downloader Bot
After=network.target

[Service]
Type=simple
User=milanmadusanka
WorkingDirectory=/home/milanmadusanka/Projects/yt-bot
ExecStart=/home/milanmadusanka/Projects/yt-bot/yt-bot
Restart=always
RestartSec=10

[Install]
WantedBy=multi-user.target
```

2. Enable and start:
```bash
sudo systemctl enable yt-bot
sudo systemctl start yt-bot
sudo systemctl status yt-bot
```

---

**Need help?** Check `README.md` for detailed documentation!
