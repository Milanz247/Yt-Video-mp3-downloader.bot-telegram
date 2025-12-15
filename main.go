package main

import (
	"context"
	"crypto/md5"
	"encoding/hex"
	"fmt"
	"log"
	"os"
	"os/exec"
	"path/filepath"
	"regexp"
	"strings"
	"sync"
	"time"

	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
	"github.com/joho/godotenv"
)

var (
	youtubeRegex         = regexp.MustCompile(`(?:https?://)?(?:www\.)?(?:youtube\.com/watch\?v=|youtu\.be/)([a-zA-Z0-9_-]{11})`)
	youtubePlaylistRegex = regexp.MustCompile(`(?:https?://)?(?:www\.)?youtube\.com/.*[?&]list=([a-zA-Z0-9_-]+)`)
)

type Bot struct {
	api          *tgbotapi.BotAPI
	downloadPath string
	urlCache     map[string]string
	cacheMutex   sync.RWMutex
}

func main() {
	// Load environment variables
	if err := godotenv.Load(); err != nil {
		log.Println("No .env file found, using system environment variables")
	}

	token := os.Getenv("TELEGRAM_BOT_TOKEN")
	if token == "" {
		log.Fatal("TELEGRAM_BOT_TOKEN is not set")
	}

	bot, err := tgbotapi.NewBotAPI(token)
	if err != nil {
		log.Panic(err)
	}

	bot.Debug = true
	log.Printf("Authorized on account %s", bot.Self.UserName)

	// Create downloads directory
	downloadPath := "downloads"
	if err := os.MkdirAll(downloadPath, 0755); err != nil {
		log.Fatal(err)
	}

	mediaBot := &Bot{
		api:          bot,
		downloadPath: downloadPath,
		urlCache:     make(map[string]string),
	}

	// Register bot commands (makes the bot interface modern in Telegram clients)
	commands := []tgbotapi.BotCommand{
		{Command: "start", Description: "Start the bot and show welcome"},
		{Command: "help", Description: "Show help and usage"},
		{Command: "latest", Description: "Show latest features"},
	}
	if _, err := bot.Request(tgbotapi.NewSetMyCommands(commands)); err != nil {
		log.Printf("Failed to set bot commands: %v", err)
	}

	// Check if yt-dlp is installed
	if !mediaBot.checkYtDlp() {
		log.Fatal("yt-dlp is not installed. Please install it: https://github.com/yt-dlp/yt-dlp")
	}

	u := tgbotapi.NewUpdate(0)
	u.Timeout = 60

	updates := bot.GetUpdatesChan(u)

	for update := range updates {
		if update.CallbackQuery != nil {
			mediaBot.handleCallbackQuery(update.CallbackQuery)
			continue
		}

		if update.Message == nil {
			continue
		}

		if update.Message.IsCommand() {
			mediaBot.handleCommand(update.Message)
		} else {
			mediaBot.handleMessage(update.Message)
		}
	}
}

func (b *Bot) checkYtDlp() bool {
	// Try multiple possible yt-dlp locations
	ytdlpPaths := []string{
		"yt-dlp",
		"/usr/local/bin/yt-dlp",
		"/usr/bin/yt-dlp",
		".venv/bin/yt-dlp",
	}

	for _, path := range ytdlpPaths {
		cmd := exec.Command(path, "--version")
		if cmd.Run() == nil {
			return true
		}
	}
	return false
}

func (b *Bot) getYtDlpPath() string {
	ytdlpPaths := []string{
		"yt-dlp",
		"/usr/local/bin/yt-dlp",
		"/usr/bin/yt-dlp",
		".venv/bin/yt-dlp",
	}

	for _, path := range ytdlpPaths {
		cmd := exec.Command(path, "--version")
		if cmd.Run() == nil {
			return path
		}
	}
	return "yt-dlp"
}

func (b *Bot) handleCommand(message *tgbotapi.Message) {
	switch message.Command() {
	case "start":
		b.sendWelcomeMessage(message.Chat.ID)
	case "help":
		b.sendHelpMessage(message.Chat.ID)
	default:
		msg := tgbotapi.NewMessage(message.Chat.ID, "Unknown command. Use /help for available commands.")
		b.api.Send(msg)
	}
}

func (b *Bot) sendWelcomeMessage(chatID int64) {
	// Prefer local welcome image if present, else use a default image URL
	welcomeLocal := "assets/welcome.jpg"
	welcomeURL := "https://images.unsplash.com/photo-1515879218367-8466d910aaa4?w=1200&q=80&auto=format&fit=crop"

	caption := `🎥 *Welcome to YouTube Downloader Bot!*

I can download YouTube videos and playlists. Send a link to get started.`

	// Build inline keyboard for a modern look
	keyboard := tgbotapi.NewInlineKeyboardMarkup(
		tgbotapi.NewInlineKeyboardRow(
			tgbotapi.NewInlineKeyboardButtonData("📥 Send a link", "help"),
			tgbotapi.NewInlineKeyboardButtonURL("📚 Examples", "https://github.com/Milanz247/Yt-Video-mp3-downloader.bot-telegram#example-links-to-test"),
		),
		tgbotapi.NewInlineKeyboardRow(
			tgbotapi.NewInlineKeyboardButtonData("❓ Help", "help"),
			tgbotapi.NewInlineKeyboardButtonData("⚙️ Settings", "settings"),
		),
	)

	if _, err := os.Stat(welcomeLocal); err == nil {
		photo := tgbotapi.NewPhoto(chatID, tgbotapi.FilePath(welcomeLocal))
		photo.Caption = caption
		photo.ParseMode = "Markdown"
		photo.ReplyMarkup = keyboard
		b.api.Send(photo)
		return
	}

	// Fallback to external image URL
	photo := tgbotapi.NewPhoto(chatID, tgbotapi.FileURL(welcomeURL))
	photo.Caption = caption
	photo.ParseMode = "Markdown"
	photo.ReplyMarkup = keyboard
	b.api.Send(photo)
}

func (b *Bot) sendHelpMessage(chatID int64) {
	text := `📖 *Help - How to use this bot*

*Supported Platform:*
• YouTube (videos, shorts, and playlists)

*Supported Formats:*
• Video: MP4 (various qualities: 360p, 480p, 720p, 1080p, best)
• Audio: MP3 (128kbps, 192kbps, 320kbps, best)
• Playlists: Download entire playlists or individual videos

*How to use:*
1. Copy a video or playlist link from YouTube
2. Send the link to me
3. Select your preferred quality from the options
4. Wait for the download to complete
5. Receive your media file!

*Playlist Options:*
• Download as single video (if playlist link)
• Download first 5 videos from playlist

*Commands:*
/start - Start the bot
/help - Show this help message

*Note:* Large files may take time to process. Please be patient! 🙏`

	msg := tgbotapi.NewMessage(chatID, text)
	msg.ParseMode = "Markdown"
	b.api.Send(msg)
}

func (b *Bot) handleMessage(message *tgbotapi.Message) {
	text := strings.TrimSpace(message.Text)

	// Check if message contains a URL
	if !strings.Contains(text, "http://") && !strings.Contains(text, "https://") && !strings.Contains(text, "www.") {
		msg := tgbotapi.NewMessage(message.Chat.ID, "Please send a valid YouTube video or playlist link.")
		b.api.Send(msg)
		return
	}

	// Detect platform
	platform := b.detectPlatform(text)
	if platform == "" {
		msg := tgbotapi.NewMessage(message.Chat.ID, "❌ Unsupported link. Please send a YouTube video or playlist link.")
		b.api.Send(msg)
		return
	}

	// Send quality selection keyboard
	b.sendQualityOptions(message.Chat.ID, text, platform)
}

func (b *Bot) detectPlatform(url string) string {
	if youtubeRegex.MatchString(url) || youtubePlaylistRegex.MatchString(url) || strings.Contains(url, "youtube.com") || strings.Contains(url, "youtu.be") {
		// Check if it's a playlist
		if youtubePlaylistRegex.MatchString(url) || strings.Contains(url, "list=") {
			return "youtube-playlist"
		}
		return "youtube"
	}
	return ""
}

func (b *Bot) sendQualityOptions(chatID int64, url, platform string) {
	// Generate a short hash for the URL
	urlID := b.cacheURL(url)

	keyboard := tgbotapi.NewInlineKeyboardMarkup(
		tgbotapi.NewInlineKeyboardRow(
			tgbotapi.NewInlineKeyboardButtonData("🎬 Best Quality Video", fmt.Sprintf("v:best:%s", urlID)),
			tgbotapi.NewInlineKeyboardButtonData("🎬 1080p", fmt.Sprintf("v:1080:%s", urlID)),
		),
		tgbotapi.NewInlineKeyboardRow(
			tgbotapi.NewInlineKeyboardButtonData("🎬 720p", fmt.Sprintf("v:720:%s", urlID)),
			tgbotapi.NewInlineKeyboardButtonData("🎬 480p", fmt.Sprintf("v:480:%s", urlID)),
		),
		tgbotapi.NewInlineKeyboardRow(
			tgbotapi.NewInlineKeyboardButtonData("🎬 360p", fmt.Sprintf("v:360:%s", urlID)),
		),
		tgbotapi.NewInlineKeyboardRow(
			tgbotapi.NewInlineKeyboardButtonData("🎵 MP3 Best", fmt.Sprintf("a:best:%s", urlID)),
			tgbotapi.NewInlineKeyboardButtonData("🎵 MP3 320kbps", fmt.Sprintf("a:320:%s", urlID)),
		),
		tgbotapi.NewInlineKeyboardRow(
			tgbotapi.NewInlineKeyboardButtonData("🎵 MP3 192kbps", fmt.Sprintf("a:192:%s", urlID)),
			tgbotapi.NewInlineKeyboardButtonData("🎵 MP3 128kbps", fmt.Sprintf("a:128:%s", urlID)),
		),
	)

	// Choose message based on platform
	var messageText string
	if platform == "youtube-playlist" {
		messageText = "📋 *Playlist detected!*\n\nChoose what to download:"
		// Update keyboard for playlist
		keyboard = tgbotapi.NewInlineKeyboardMarkup(
			tgbotapi.NewInlineKeyboardRow(
				tgbotapi.NewInlineKeyboardButtonData("🎬 Single Video (Best)", fmt.Sprintf("v:best:%s", urlID)),
			),
			tgbotapi.NewInlineKeyboardRow(
				tgbotapi.NewInlineKeyboardButtonData("📋 First 5 Videos (Best)", fmt.Sprintf("p:5:best:%s", urlID)),
			),
			tgbotapi.NewInlineKeyboardRow(
				tgbotapi.NewInlineKeyboardButtonData("🎵 Single Audio (MP3)", fmt.Sprintf("a:best:%s", urlID)),
			),
			tgbotapi.NewInlineKeyboardRow(
				tgbotapi.NewInlineKeyboardButtonData("📋 First 5 Audios (MP3)", fmt.Sprintf("pa:5:best:%s", urlID)),
			),
			tgbotapi.NewInlineKeyboardRow(
				tgbotapi.NewInlineKeyboardButtonData("📋 View all items", fmt.Sprintf("list:%s", urlID)),
			),
		)
	} else {
		messageText = "📥 *Choose quality:*\n\nSelect the format and quality you prefer:"
	}

	msg := tgbotapi.NewMessage(chatID, messageText)
	msg.ParseMode = "Markdown"
	msg.ReplyMarkup = keyboard
	b.api.Send(msg)
}

func (b *Bot) cacheURL(url string) string {
	// Create a short hash from URL
	hash := md5.Sum([]byte(url))
	urlID := hex.EncodeToString(hash[:])[:12] // Use first 12 chars

	b.cacheMutex.Lock()
	b.urlCache[urlID] = url
	b.cacheMutex.Unlock()

	return urlID
}

func (b *Bot) getURLFromCache(urlID string) string {
	b.cacheMutex.RLock()
	defer b.cacheMutex.RUnlock()
	return b.urlCache[urlID]
}

func (b *Bot) handleCallbackQuery(query *tgbotapi.CallbackQuery) {
	parts := strings.Split(query.Data, ":")

	// Handle short two-part callbacks (list/open) early
	if len(parts) == 2 {
		if parts[0] == "list" {
			urlID := parts[1]
			b.presentPlaylistItems(query.Message.Chat.ID, urlID)
			callback := tgbotapi.NewCallback(query.ID, "Opening playlist items...")
			b.api.Request(callback)
			return
		}
		if parts[0] == "open" {
			videoID := parts[1]
			videoURL := b.getURLFromCache(videoID)
			if videoURL == "" {
				callback := tgbotapi.NewCallback(query.ID, "❌ Link expired. Please view playlist again.")
				b.api.Request(callback)
				return
			}
			callback := tgbotapi.NewCallback(query.ID, "Opening video options...")
			b.api.Request(callback)
			// Send quality options for this specific video
			b.sendQualityOptions(query.Message.Chat.ID, videoURL, "youtube")
			return
		}
		if parts[0] == "help" {
			// Show help message
			callback := tgbotapi.NewCallback(query.ID, "Opening help...")
			b.api.Request(callback)
			b.sendHelpMessage(query.Message.Chat.ID)
			return
		}
		if parts[0] == "settings" {
			callback := tgbotapi.NewCallback(query.ID, "Opening settings...")
			b.api.Request(callback)
			msg := tgbotapi.NewMessage(query.Message.Chat.ID, "⚙️ Settings are minimal for now — send /help for instructions.")
			b.api.Send(msg)
			return
		}
	}

	if len(parts) < 3 {
		return
	}
	formatType := parts[0] // "v" (video), "a" (audio), "p" (playlist video), "pa" (playlist audio)

	var quality, urlID string
	var playlistCount int

	// Handle playlist format: "p:5:best:urlID" or "pa:5:best:urlID"
	if formatType == "p" || formatType == "pa" {
		if len(parts) != 4 {
			return
		}
		playlistCount, _ = fmt.Sscanf(parts[1], "%d", &playlistCount)
		quality = parts[2]
		urlID = parts[3]
	} else {
		// Regular format: "v:best:urlID" or "a:best:urlID"
		quality = parts[1]
		urlID = parts[2]
	}

	// Get URL from cache
	url := b.getURLFromCache(urlID)
	if url == "" {
		callback := tgbotapi.NewCallback(query.ID, "❌ Link expired. Please send the link again.")
		b.api.Request(callback)
		return
	}

	// Convert format type to full name
	format := "video"
	isPlaylist := false
	if formatType == "a" || formatType == "pa" {
		format = "audio"
	}
	if formatType == "p" || formatType == "pa" {
		isPlaylist = true
	}

	// Answer callback query
	callback := tgbotapi.NewCallback(query.ID, "Processing your request...")
	b.api.Request(callback)

	// Send processing message
	var processingText string
	if isPlaylist {
		processingText = fmt.Sprintf("⏳ Downloading %d items from playlist... This may take a few minutes.", playlistCount)
	} else {
		processingText = "⏳ Downloading... This may take a few moments."
	}
	processingMsg := tgbotapi.NewMessage(query.Message.Chat.ID, processingText)
	sentMsg, _ := b.api.Send(processingMsg)

	// Download the media
	if isPlaylist {
		log.Printf("Starting playlist download: format=%s, quality=%s, count=%d, url=%s", format, quality, playlistCount, url)
		b.downloadPlaylist(query.Message.Chat.ID, url, format, quality, playlistCount, sentMsg.MessageID)
		return
	}

	log.Printf("Starting download: format=%s, quality=%s, url=%s", format, quality, url)
	filePath, title, err := b.downloadMedia(url, format, quality)
	if err != nil {
		log.Printf("Download error: %v", err)
		errorMsg := tgbotapi.NewMessage(query.Message.Chat.ID, fmt.Sprintf("❌ Error: %v", err))
		b.api.Send(errorMsg)
		b.api.Request(tgbotapi.NewDeleteMessage(query.Message.Chat.ID, sentMsg.MessageID))
		return
	}

	log.Printf("Download successful: %s", filePath)

	// Delete processing message
	b.api.Request(tgbotapi.NewDeleteMessage(query.Message.Chat.ID, sentMsg.MessageID))

	// Send the file
	log.Printf("Sending file to user... (title=%s)", title)
	if err := b.sendFile(query.Message.Chat.ID, filePath, format, title); err != nil {
		log.Printf("Failed to send file: %v", err)
		// Keep file so user can retry later or for debugging
	} else {
		// Clean up only after successful send
		log.Printf("Cleaning up: %s", filePath)
		os.Remove(filePath)
	}
}

func (b *Bot) downloadPlaylist(chatID int64, url, format, quality string, count, processingMsgID int) {
	ytdlp := b.getYtDlpPath()
	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Minute)
	defer cancel()

	// Download playlist info to get video URLs
	timestamp := time.Now().UnixNano()
	playlistInfoFile := filepath.Join(b.downloadPath, fmt.Sprintf("playlist_%d.txt", timestamp))

	// Get first N video IDs from playlist
	cmd := exec.CommandContext(ctx, ytdlp,
		"--flat-playlist",
		"--print", "url",
		"--playlist-end", fmt.Sprintf("%d", count),
		url)

	output, err := cmd.CombinedOutput()
	if err != nil {
		log.Printf("Playlist fetch error: %v - %s", err, string(output))
		errorMsg := tgbotapi.NewMessage(chatID, "❌ Failed to fetch playlist. Please try again.")
		b.api.Send(errorMsg)
		b.api.Request(tgbotapi.NewDeleteMessage(chatID, processingMsgID))
		return
	}

	videoURLs := strings.Split(strings.TrimSpace(string(output)), "\n")
	log.Printf("Found %d videos in playlist", len(videoURLs))

	// Download each video
	successCount := 0
	for i, videoURL := range videoURLs {
		if strings.TrimSpace(videoURL) == "" {
			continue
		}

		// Update status
		statusMsg := tgbotapi.NewEditMessageText(chatID, processingMsgID,
			fmt.Sprintf("⏳ Downloading item %d/%d from playlist...", i+1, len(videoURLs)))
		b.api.Send(statusMsg)

		// Download single video
		filePath, title, err := b.downloadMedia(videoURL, format, quality)
		if err != nil {
			log.Printf("Failed to download video %d: %v", i+1, err)
			continue
		}

		// Send the file
		if err := b.sendFile(chatID, filePath, format, title); err != nil {
			log.Printf("Failed to send playlist item %d: %v", i+1, err)
			// don't remove file; continue to next
		} else {
			os.Remove(filePath)
			successCount++
		}
		successCount++
	}

	// Delete processing message and send completion message
	b.api.Request(tgbotapi.NewDeleteMessage(chatID, processingMsgID))
	completionMsg := tgbotapi.NewMessage(chatID,
		fmt.Sprintf("✅ Downloaded %d/%d items from playlist!", successCount, len(videoURLs)))
	b.api.Send(completionMsg)

	os.Remove(playlistInfoFile)
}

// fetchPlaylistEntries fetches up to max entries (title and URL) from a playlist
func (b *Bot) fetchPlaylistEntries(url string, max int) ([][2]string, error) {
	ytdlp := b.getYtDlpPath()
	args := []string{"--flat-playlist", "--no-warnings", "--print", "%(title)s||%(url)s", url}
	cmd := exec.Command(ytdlp, args...)
	output, err := cmd.Output()
	if err != nil {
		return nil, fmt.Errorf("playlist fetch failed: %v - %s", err, string(output))
	}
	lines := strings.Split(strings.TrimSpace(string(output)), "\n")
	entries := make([][2]string, 0, len(lines))
	for i, line := range lines {
		if max > 0 && i >= max {
			break
		}
		parts := strings.SplitN(line, "||", 2)
		if len(parts) != 2 {
			continue
		}
		title := strings.TrimSpace(parts[0])
		link := strings.TrimSpace(parts[1])
		entries = append(entries, [2]string{title, link})
	}
	return entries, nil
}

// presentPlaylistItems sends a message with playlist items and buttons to open each item
func (b *Bot) presentPlaylistItems(chatID int64, playlistURLID string) {
	url := b.getURLFromCache(playlistURLID)
	if url == "" {
		msg := tgbotapi.NewMessage(chatID, "❌ Playlist link expired. Please send the playlist again.")
		b.api.Send(msg)
		return
	}

	// Fetch up to 25 items
	entries, err := b.fetchPlaylistEntries(url, 25)
	if err != nil || len(entries) == 0 {
		msg := tgbotapi.NewMessage(chatID, "❌ Failed to fetch playlist items or playlist is empty.")
		b.api.Send(msg)
		return
	}

	// Build keyboard
	rows := [][]tgbotapi.InlineKeyboardButton{}
	for i, e := range entries {
		title := e[0]
		link := e[1]
		id := b.cacheURL(link)
		display := fmt.Sprintf("%02d. %s", i+1, truncateString(title, 50))
		rows = append(rows, tgbotapi.NewInlineKeyboardRow(tgbotapi.NewInlineKeyboardButtonData(display, fmt.Sprintf("open:%s", id))))
	}
	// Add a back button
	rows = append(rows, tgbotapi.NewInlineKeyboardRow(tgbotapi.NewInlineKeyboardButtonData("◀ Back", fmt.Sprintf("v:best:%s", playlistURLID))))

	keyboard := tgbotapi.NewInlineKeyboardMarkup(rows...)
	msg := tgbotapi.NewMessage(chatID, "📋 *Playlist items*\n\nSelect an item to open quality options:")
	msg.ParseMode = "Markdown"
	msg.ReplyMarkup = keyboard
	b.api.Send(msg)
}

func truncateString(s string, n int) string {
	if len(s) <= n {
		return s
	}
	return s[:n-1] + "…"
}

func (b *Bot) downloadMedia(url, format, quality string) (string, string, error) {
	// Use timestamp with nanoseconds for uniqueness fallback
	timestamp := time.Now().UnixNano()
	var outputFile string
	ytdlp := b.getYtDlpPath()

	// Try to get title and id for nicer filenames (avoid warnings in output)
	title := ""
	id := fmt.Sprintf("%d", timestamp)
	if out, err := exec.Command(ytdlp, "--no-warnings", "--get-title", url).Output(); err == nil {
		// take the last non-empty line to avoid warnings appearing before/after
		lines := strings.Split(strings.TrimSpace(string(out)), "\n")
		for i := len(lines) - 1; i >= 0; i-- {
			l := strings.TrimSpace(lines[i])
			if l != "" {
				title = l
				break
			}
		}
	}
	if out, err := exec.Command(ytdlp, "--no-warnings", "--get-id", url).Output(); err == nil {
		id = strings.TrimSpace(string(out))
	}
	if title == "" {
		title = fmt.Sprintf("media_%d", timestamp)
	}
	// Sanitize title for filesystem
	safeTitle := sanitizeFilename(title)

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Minute)
	defer cancel()

	// Common args for better compatibility
	commonArgs := []string{
		"--no-playlist",
		"--no-warnings",
		"--user-agent", "Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/120.0.0.0 Safari/537.36",
	}

	// Add cookies if file exists (for Facebook/Instagram)
	cookiesFile := "cookies.txt"
	if _, err := os.Stat(cookiesFile); err == nil {
		commonArgs = append(commonArgs, "--cookies", cookiesFile)
	}

	var cmd *exec.Cmd
	var ext string
	if format == "video" {
		ext = "mp4"
		outputFile = filepath.Join(b.downloadPath, fmt.Sprintf("%s - %s.%s", safeTitle, id, ext))
		formatStr := b.getVideoFormat(quality)
		args := append([]string{"-f", formatStr, "--merge-output-format", ext, "-o", outputFile}, commonArgs...)
		args = append(args, url)
		cmd = exec.CommandContext(ctx, ytdlp, args...)
	} else {
		ext = "mp3"
		outputFile = filepath.Join(b.downloadPath, fmt.Sprintf("%s - %s.%s", safeTitle, id, ext))
		bitrateStr := b.getAudioBitrate(quality)
		args := append([]string{"-x", "--audio-format", "mp3", "--audio-quality", bitrateStr, "-o", outputFile}, commonArgs...)
		args = append(args, url)
		cmd = exec.CommandContext(ctx, ytdlp, args...)
	}

	log.Printf("Running: %s (output: %s)", cmd.String(), outputFile)

	output, err := cmd.CombinedOutput()
	log.Printf("yt-dlp finished with error: %v", err)
	log.Printf("yt-dlp output: %s", string(output))

	if err != nil {
		// Extract meaningful error from output
		outputStr := string(output)
		errorMsg := "Download failed"

		if strings.Contains(outputStr, "ERROR:") {
			// Find the error line
			lines := strings.Split(outputStr, "\n")
			for _, line := range lines {
				if strings.Contains(line, "ERROR:") {
					// Clean up the error message
					errorMsg = strings.TrimSpace(strings.TrimPrefix(line, "ERROR:"))
					// Simplify common errors
					if strings.Contains(errorMsg, "SSL") || strings.Contains(errorMsg, "handshake") || strings.Contains(errorMsg, "timed out") {
						errorMsg = "Connection timeout. Facebook/Instagram may be blocking downloads. Try a YouTube link instead."
					} else if strings.Contains(errorMsg, "Unable to download webpage") {
						errorMsg = "Cannot access this video. It may be private or region-locked."
					} else if strings.Contains(errorMsg, "Video unavailable") {
						errorMsg = "Video is unavailable or has been removed."
					}
					break
				}
			}
		}

		return "", "", fmt.Errorf("%s", errorMsg)
	}

	// Check if file exists
	if _, err := os.Stat(outputFile); os.IsNotExist(err) {
		return "", "", fmt.Errorf("output file not found: %s", outputFile)
	}

	log.Printf("Successfully downloaded to: %s", outputFile)
	return outputFile, title, nil
}

// sanitizeFilename removes or replaces characters not safe for filenames
func sanitizeFilename(name string) string {
	// Replace newlines and slashes
	name = strings.ReplaceAll(name, "\n", " ")
	name = strings.ReplaceAll(name, "/", "-")
	name = strings.ReplaceAll(name, "\\", "-")
	// Trim excessive whitespace and control chars
	name = strings.TrimSpace(name)
	// Remove characters that are problematic
	forbidden := []string{":", "*", "?", "\"", "<", ">", "|"}
	for _, c := range forbidden {
		name = strings.ReplaceAll(name, c, "-")
	}
	// Limit length
	max := 120
	if len(name) > max {
		name = name[:max]
	}
	return name
}

func (b *Bot) getVideoFormat(quality string) string {
	switch quality {
	case "best":
		return "bestvideo[ext=mp4]+bestaudio[ext=m4a]/best[ext=mp4]/best"
	case "1080":
		return "bestvideo[height<=1080][ext=mp4]+bestaudio[ext=m4a]/best[height<=1080][ext=mp4]/best"
	case "720":
		return "bestvideo[height<=720][ext=mp4]+bestaudio[ext=m4a]/best[height<=720][ext=mp4]/best"
	case "480":
		return "bestvideo[height<=480][ext=mp4]+bestaudio[ext=m4a]/best[height<=480][ext=mp4]/best"
	case "360":
		return "bestvideo[height<=360][ext=mp4]+bestaudio[ext=m4a]/best[height<=360][ext=mp4]/best"
	default:
		return "best[ext=mp4]/best"
	}
}

func (b *Bot) getAudioBitrate(quality string) string {
	switch quality {
	case "best":
		return "0"
	case "320":
		return "320K"
	case "192":
		return "192K"
	case "128":
		return "128K"
	default:
		return "0"
	}
}

func (b *Bot) sendFile(chatID int64, filePath, format, title string) error {
	fileInfo, err := os.Stat(filePath)
	if err != nil {
		msg := tgbotapi.NewMessage(chatID, "❌ Error reading file")
		b.api.Send(msg)
		return err
	}

	// Telegram file size limit is 50MB
	const maxSize = 50 * 1024 * 1024
	if fileInfo.Size() > maxSize {
		msg := tgbotapi.NewMessage(chatID, "❌ File is too large (>50MB). Try a lower quality.")
		b.api.Send(msg)
		return fmt.Errorf("file too large")
	}

	// Try sending with retries for transient network issues
	var lastErr error
	maxAttempts := 3
	for attempt := 1; attempt <= maxAttempts; attempt++ {
		if format == "video" {
			video := tgbotapi.NewVideo(chatID, tgbotapi.FilePath(filePath))
			if title != "" {
				video.Caption = fmt.Sprintf("✅ %s", title)
			} else {
				video.Caption = "✅ Here's your video!"
			}
			_, lastErr = b.api.Send(video)
		} else {
			audio := tgbotapi.NewAudio(chatID, tgbotapi.FilePath(filePath))
			if title != "" {
				audio.Caption = fmt.Sprintf("✅ %s", title)
			} else {
				audio.Caption = "✅ Here's your audio!"
			}
			_, lastErr = b.api.Send(audio)
		}

		if lastErr == nil {
			return nil
		}

		// Detect likely transient network errors by inspecting error text
		errStr := strings.ToLower(lastErr.Error())
		if strings.Contains(errStr, "connection reset") || strings.Contains(errStr, "client connection force closed") || strings.Contains(errStr, "timeout") || strings.Contains(errStr, "temporary") {
			log.Printf("Transient send error (attempt %d/%d): %v - retrying...", attempt, maxAttempts, lastErr)
			time.Sleep(time.Duration(attempt) * 2 * time.Second)
			continue
		}

		// Non-transient error, break early
		break
	}

	// If we're here, send failed
	msg := tgbotapi.NewMessage(chatID, fmt.Sprintf("❌ Error sending file after %d attempts: %v", maxAttempts, lastErr))
	b.api.Send(msg)
	return lastErr
}
