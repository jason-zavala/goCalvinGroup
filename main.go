package main

import (
	"bufio"
	"encoding/json"
	"fmt"
	"log"
	"math/rand"
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

	var imageURL string

	if len(b.URLMap) == 0 {
		imageURL, err = b.URLGenerator.GetRandomURL()
		if err != nil {
			log.Fatal(err)
		}
	} else {
		randomIndex := rand.Intn(len(b.URLMap))
		imageURL = b.URLMap[randomIndex]
		b.PreviousURLIndex = randomIndex
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

func (b *Calvin) handleNextBotImage() {

	msg, err := b.QuoteProvider.GetRandomQuote()
	if err != nil {
		log.Fatal(err)
	}

	updatedIndex := b.PreviousURLIndex
	updatedIndex++
	imageURL := b.URLMap[updatedIndex]
	b.PreviousURLIndex = updatedIndex

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

	if len(words) == 0 || !b.Command[words[0]] {
		return
	}

	if b.Command["!next"] || b.Command["next"] {
		b.handleNextBotImage()
		return
	}
	b.handleBotImage()
}

func buildURLMap() (map[int]string, error) {
	file, err := os.Open("url.txt")

	if err != nil {
		return nil, err
	}

	defer file.Close()

	urlMap := make(map[int]string)
	scanner := bufio.NewScanner(file)
	count := 0

	for scanner.Scan() {
		urlMap[count] = scanner.Text()
		count++
	}

	if len(urlMap) == 0 {
		return nil, nil
	}
	fmt.Println("finished building map")
	return urlMap, nil
}

func main() {
	botID := os.Getenv("BOT_ID")
	accessToken := os.Getenv("GM_TOKEN")
	urlMap, err := buildURLMap()

	if err != nil {
		log.Fatal(err)
	}
	commands := map[string]bool{
		"!comic": true,
		"!next":  true,
		"comic":  true,
		"next":   true,
	}

	bot := &Calvin{
		URLGenerator:    &DefaultURLGenerator{},
		ImageDownloader: &DefaultImageDownloader{},
		QuoteProvider:   &DefaultQuoteProvider{},
		ImageUploader:   &DefaultImageUploader{},
		MessageSender:   &DefaultMessageSender{},
		Command:         commands,
		BotID:           botID,
		AccessToken:     accessToken,
		URLMap:          urlMap,
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
	err = http.ListenAndServe(":"+port, nil)
	if err != nil {
		log.Fatal(err)
	}
}
