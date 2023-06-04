package main

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"os"
	"strings"
)

type BotCommand struct {
	Name       string `json:"name"`
	Message    string `json:"text"`
	AvatarURL  string `json:"avatar_url"`
	CreatedAt  int64  `json:"created_at"`
	GroupID    string `json:"group_id"`
	ID         string `json:"id"`
	SenderID   string `json:"sender_id"`
	SenderType string `json:"sender_type"`
	SourceGUID string `json:"source_guid"`
	System     bool   `json:"system"`
	UserID     string `json:"user_id"`
}

func (b *Calvin) handleBotImage() {

	msg, err := b.QuoteProvider.GetRandomQuote()
	if err != nil {
		log.Fatal(err)
	}

	imageURL, err := b.URLGenerator.GetRandomURL()
	if err != nil {
		log.Fatal(err)
	}

	err = b.ImageDownloader.DownloadImage(imageURL)
	if err != nil {
		log.Fatal(err)
	}

	imageURL, err = b.ImageUploader.UploadImage("image.jpg", b.AccessToken)
	if err != nil {
		log.Fatal(err)
	}

	err = b.MessageSender.SendMessage(imageURL, msg, b.BotID)

	if err != nil {
		fmt.Printf("Error: %s\n", err)
	} else {
		fmt.Println("Message sent successfully!")
	}
}

func (b *Calvin) handleCallback(w http.ResponseWriter, r *http.Request) {

	if r.Method != http.MethodPost {
		http.Error(w, "Method Not Allowed", http.StatusMethodNotAllowed)
		return
	}

	var command BotCommand
	err := json.NewDecoder(r.Body).Decode(&command)

	if err != nil {
		return
	}

	words := strings.Fields(command.Message)

	if len(words) == 0 || words[0] != b.Command {
		return
	}

	b.handleBotImage()
}

func main() {
	botID := os.Getenv("BOT_ID")
	accessToken := os.Getenv("GM_TOKEN")

	bot := &Calvin{
		URLGenerator:    &DefaultURLGenerator{},
		ImageDownloader: &DefaultImageDownloader{},
		QuoteProvider:   &DefaultQuoteProvider{},
		ImageUploader:   &DefaultImageUploader{},
		MessageSender:   &DefaultMessageSender{},
		Command:         "!comic",
		BotID:           botID,
		AccessToken:     accessToken,
	}

	// Set up the HTTP server
	http.HandleFunc("/callback", bot.handleCallback)

	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		http.ServeFile(w, r, "index.html")
	})

	// Start the server
	port := os.Getenv("PORT")
	if port == "" {
		port = "9001"
	}

	fmt.Printf("Server listening on port %s...\n", port)
	err := http.ListenAndServe(":"+port, nil)
	if err != nil {
		log.Fatal(err)
	}
}
