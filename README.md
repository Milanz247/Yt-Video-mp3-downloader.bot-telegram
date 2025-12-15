# Telegram Media Downloader Bot

A powerful Telegram bot built with Go that downloads videos and audio from YouTube, Facebook, and Instagram with multiple quality options.

## Features

- ✅ **Multi-platform support**: YouTube, Facebook, Instagram
- 🎬 **Video downloads**: Multiple quality options (360p, 480p, 720p, 1080p, Best)
- 🎵 **Audio downloads**: MP3 with multiple bitrates (128kbps, 192kbps, 320kbps, Best)
- 🚀 **Fast and efficient**: Built with Go for optimal performance
- 💬 **User-friendly**: Interactive buttons for quality selection
- 🔒 **Reliable**: Uses yt-dlp for robust media extraction

## Prerequisites

Before running the bot, ensure you have the following installed:

1. **Go** (version 1.21 or higher)
   ```bash
   # Check Go version
   go version
   ```

2. **yt-dlp** - A powerful video downloader
   ```bash
   # Install on Linux/macOS
   sudo wget https://github.com/yt-dlp/yt-dlp/releases/latest/download/yt-dlp -O /usr/local/bin/yt-dlp
   sudo chmod a+rx /usr/local/bin/yt-dlp
   
   # Or using pip
   pip install yt-dlp
   
   # Or using package manager
   # Ubuntu/Debian
   sudo apt install yt-dlp
   
   # macOS
   brew install yt-dlp
   
   # Verify installation
   yt-dlp --version
   ```

3. **FFmpeg** - Required for audio conversion and merging
   ```bash
   # Ubuntu/Debian
   sudo apt install ffmpeg
   
   # macOS
   brew install ffmpeg
   
   # Verify installation
   ffmpeg -version
   ```

## Installation

1. **Clone or navigate to the project directory**
   ```bash
   cd /home/milanmadusanka/Projects/yt-bot
   ```

2. **Create environment file**
   ```bash
   cp .env.example .env
   ```

3. **Edit the `.env` file** with your bot token (already configured):
   ```
   TELEGRAM_BOT_TOKEN=8266088482:AAE1zZjqDQ4puqZKa5DgW-6tbJaeVE5YN6Q
   ```

4. **Install Go dependencies**
   ```bash
   go mod download
   ```

5. **Build the bot**
   ```bash
   go build -o yt-bot
   ```

## Usage

### Running the Bot

```bash
# Run directly
go run main.go

# Or run the compiled binary
./yt-bot
```

### Using the Bot on Telegram

1. **Start the bot**: Send `/start` to receive a welcome message
2. **Get help**: Send `/help` to see usage instructions
3. **Download media**:
   - Send a YouTube, Facebook, or Instagram video link
   - Select your preferred quality and format from the interactive buttons
   - Wait for the bot to process and send your file

### Supported Quality Options

**Video Formats:**
- 🎬 Best Quality (highest available)
- 🎬 1080p (Full HD)
- 🎬 720p (HD)
- 🎬 480p (SD)
- 🎬 360p (Low)

**Audio Formats (MP3):**
- 🎵 Best Quality
- 🎵 320kbps (High)
- 🎵 192kbps (Medium)
- 🎵 128kbps (Low)

## Example Links to Test

- **YouTube**: `https://www.youtube.com/watch?v=VIDEO_ID`
- **Facebook**: `https://www.facebook.com/watch?v=VIDEO_ID`
- **Instagram**: `https://www.instagram.com/p/POST_ID/` or `https://www.instagram.com/reel/REEL_ID/`

## Project Structure

```
yt-bot/
├── main.go           # Main bot application
├── go.mod            # Go module dependencies
├── .env              # Environment variables (not in git)
├── .env.example      # Example environment file
├── .gitignore        # Git ignore rules
├── README.md         # This file
└── downloads/        # Temporary download directory (created automatically)
```

## Configuration

The bot can be configured through environment variables in the `.env` file:

- `TELEGRAM_BOT_TOKEN`: Your Telegram bot token (required)

## Limitations

- Maximum file size: 50MB (Telegram's limit for bot uploads)
- Private/restricted videos may not be downloadable
- Some platforms may have rate limiting

## Troubleshooting

### Bot doesn't start
- Verify your bot token is correct in `.env`
- Check that all dependencies are installed
- Ensure yt-dlp and ffmpeg are in your PATH

### Downloads fail
- Update yt-dlp: `yt-dlp -U` or `pip install -U yt-dlp`
- Check if the video link is accessible
- Try a different quality option

### File too large
- Select a lower quality option
- yt-dlp will download the best available format within the size limit

## Development

### Running in Debug Mode

The bot runs in debug mode by default (set in code: `bot.Debug = true`). Check console output for detailed logs.

### Code Structure

- **main()**: Initializes bot and starts message polling
- **handleCommand()**: Processes bot commands (/start, /help)
- **handleMessage()**: Detects and processes video links
- **handleCallbackQuery()**: Handles quality selection buttons
- **downloadMedia()**: Downloads video/audio using yt-dlp
- **sendFile()**: Sends downloaded file to user

## Technologies Used

- **Language**: Go (Golang)
- **Telegram API**: github.com/go-telegram-bot-api/telegram-bot-api/v5
- **Media Downloader**: yt-dlp
- **Media Processing**: FFmpeg
- **Environment Config**: github.com/joho/godotenv

## License

MIT License - Feel free to use and modify as needed.

## Contributing

Contributions are welcome! Feel free to submit issues or pull requests.

## Support

If you encounter any issues or have questions:
1. Check the Troubleshooting section
2. Ensure all prerequisites are installed
3. Verify your bot token is correct
4. Check yt-dlp is up to date

---

**Made with ❤️ using Go and yt-dlp**
