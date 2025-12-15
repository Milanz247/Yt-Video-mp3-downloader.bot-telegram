# Telegram YouTube Downloader Bot

A focused Telegram bot built with Go that downloads videos and audio from YouTube (single videos and playlists) with multiple quality options.

## Features

- ✅ **Platform support**: YouTube (videos, shorts, playlists)
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

2. **yt-dlp** - A powerful video downloader (required)
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

1. **Clone the repository**
   ```bash
   git clone https://github.com/Milanz247/Yt-Video-mp3-downloader.bot-telegram.git
   cd Yt-Video-mp3-downloader.bot-telegram
   ```

2. **Create environment file** and set your Telegram bot token
   ```bash
   cp .env.example .env
   # Edit .env and set TELEGRAM_BOT_TOKEN to your bot token
   ```

4. **Install Go dependencies**
   ```bash
   go mod download
   ```

5. **Build the bot**
   ```bash
   go build -o yt-bot
   ```

Or run the setup script which checks/install prerequisites and builds the bot:

```bash
chmod +x setup.sh
./setup.sh
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
   - Send a YouTube video or playlist link
   - If a playlist is detected you'll get playlist options (single item, first N items, or "View all items")
   - Use the interactive buttons to choose format/quality and start download
   - The bot will send files back to you when ready

### Modern UI

- When you send `/start` the bot will send a welcome image (if `assets/welcome.jpg` exists it will be used, otherwise a hosted image is used) and a modern inline keyboard with quick actions:
   - "Send a link" – quick entry to send a link
   - "Examples" – link to README examples
   - "Help" and "Settings" quick actions

You can provide your own `assets/welcome.jpg` (recommended size ~1200x600) to brand the bot.

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

- **YouTube (video)**: `https://www.youtube.com/watch?v=VIDEO_ID`
- **YouTube (playlist)**: `https://www.youtube.com/playlist?list=PLAYLIST_ID` or `https://www.youtube.com/watch?v=VIDEO_ID&list=PLAYLIST_ID`

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

The bot is configured through the `.env` file. At minimum set:

- `TELEGRAM_BOT_TOKEN`: Your Telegram bot token (required)

Optionally you can add `cookies.txt` (exported from your browser) in the project root if you want to try downloading geo-restricted or protected content (may not be necessary for YouTube).

## Limitations

- Maximum file size: 50MB (Telegram's standard limit for bot uploads). The bot will tell you if a file is too large and suggest lower quality.
- Private/restricted videos or region-restricted content may not be downloadable without cookies or special handling.
- Playlist downloads of large playlists are limited (the UI fetches and lists the first 25 items for selection); you can download the first N items using the playlist buttons.

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

## Creator & Support

Created by: **Milan Madusanka** (GitHub: @Milanz247)

If you need help:
1. Check Troubleshooting above
2. Ensure prerequisites are installed and `yt-dlp` is up to date (`yt-dlp -U`)
3. Open an issue in the GitHub repo

---

**Made with ❤️ by Milan Madusanka using Go and yt-dlp**
