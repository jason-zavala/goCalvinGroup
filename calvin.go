package main

import (
	"bufio"
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"io/ioutil"
	"math/rand"
	"net/http"
	"os"
	"strings"
	"time"
)

type URLGenerator interface {
	GetRandomURL() (string, error)
}

type ImageDownloader interface {
	DownloadImage(url string) error
}

type QuoteProvider interface {
	GetRandomQuote() (string, error)
}

type ImageUploader interface {
	UploadImage(imagePath, accessToken string) (string, error)
}

type MessageSender interface {
	SendMessage(url, msg, botID string) error
}

type Calvin struct {
	URLGenerator     URLGenerator
	ImageDownloader  ImageDownloader
	QuoteProvider    QuoteProvider
	ImageUploader    ImageUploader
	MessageSender    MessageSender
	Command          map[string]bool
	BotID            string
	AccessToken      string
	PreviousURLIndex int
	URLMap           map[int]string
}

type ImageResponse struct {
	Payload struct {
		URL        string `json:"url"`
		PictureURL string `json:"picture_url"`
	} `json:"payload"`
}

type DefaultURLGenerator struct{}

func (d *DefaultURLGenerator) GetRandomURL() (string, error) {
	// Open the file
	file, err := os.Open("url.txt")
	if err != nil {
		return "", err
	}
	defer file.Close()

	// Create a new random seed based on the current time
	rand.Seed(time.Now().UnixNano())

	// Read all the lines from the file
	scanner := bufio.NewScanner(file)
	var lines []string
	for scanner.Scan() {
		lines = append(lines, scanner.Text())
	}

	// Check if any lines were read
	if len(lines) == 0 {
		return "", nil
	}

	// Select a random line from the file
	randomIndex := rand.Intn(len(lines))
	randomLine := lines[randomIndex]

	return randomLine, nil
}

type DefaultImageDownloader struct{}

func (d *DefaultImageDownloader) DownloadImage(url string) error {
	// Send a GET request to the image URL
	resp, err := http.Get(url)
	if err != nil {
		return fmt.Errorf("failed to retrieve the image: %v", err)
	}
	defer resp.Body.Close()

	// Create a new file with the .jpg extension
	fileName := "image.jpg"
	file, err := os.Create(fileName)
	if err != nil {
		return fmt.Errorf("failed to create the file: %v", err)
	}
	defer file.Close()

	// Copy the image data to the file
	_, err = io.Copy(file, resp.Body)
	if err != nil {
		return fmt.Errorf("failed to save the image: %v", err)
	}

	//log.Printf("Image saved as %s\n", fileName)
	return nil
}

type DefaultQuoteProvider struct{}

func (q *DefaultQuoteProvider) GetRandomQuote() (string, error) {
	// Read the contents of the text file
	content, err := ioutil.ReadFile("quotes.txt")
	if err != nil {
		return "", err
	}

	// Split the content into lines
	lines := strings.Split(string(content), "\n")

	// Seed the random number generator
	rand.Seed(time.Now().UnixNano())

	// Generate a random index
	randomIndex := rand.Intn(len(lines))

	// Get the random line from the lines slice
	randomLine := lines[randomIndex]

	return randomLine, nil
}

type DefaultImageUploader struct{}

func (u *DefaultImageUploader) UploadImage(imagePath, accessToken string) (string, error) {
	// Set the GroupMe Image Service endpoint URL
	url := "https://image.groupme.com/pictures"

	// Read the image file
	file, err := os.Open(imagePath)
	if err != nil {
		return "", fmt.Errorf("failed to open image file: %s", err.Error())
	}
	defer file.Close()

	// Create the HTTP client
	client := &http.Client{}

	// Create the request
	req, err := http.NewRequest("POST", url, file)
	if err != nil {
		return "", fmt.Errorf("failed to create HTTP request: %v", err)
	}

	// Set the request headers
	req.Header.Set("Content-Type", "image/jpeg")
	req.Header.Set("X-Access-Token", accessToken)

	// Send the request
	resp, err := client.Do(req)
	if err != nil {
		return "", fmt.Errorf("failed to send HTTP request: %v", err)
	}
	defer resp.Body.Close()

	// Read the response body
	responseBody, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return "", fmt.Errorf("failed to read response body: %v", err)
	}

	// Process the response
	if resp.StatusCode == http.StatusOK {
		var imageResp ImageResponse
		if err := json.Unmarshal(responseBody, &imageResp); err != nil {
			return "", fmt.Errorf("failed to parse response: %v", err)
		}
		return imageResp.Payload.PictureURL, nil
	} else {
		return "", fmt.Errorf("image upload failed. Response body: %s", string(responseBody))
	}
}

type DefaultMessageSender struct{}

func (s *DefaultMessageSender) SendMessage(url, msg, botID string) error {
	type Attachment struct {
		Type string `json:"type"`
		URL  string `json:"url"`
	}

	data := struct {
		BotID       string       `json:"bot_id"`
		Text        string       `json:"text"`
		Attachments []Attachment `json:"attachments"`
	}{
		BotID: botID,
		Text:  msg,
		Attachments: []Attachment{
			{
				Type: "image",
				URL:  url,
			},
		},
	}

	jsonData, err := json.Marshal(data)
	if err != nil {
		return fmt.Errorf("failed to marshal JSON: %w", err)
	}

	resp, err := http.Post("https://api.groupme.com/v3/bots/post", "application/json", bytes.NewBuffer(jsonData))
	if err != nil {
		return fmt.Errorf("failed to send HTTP request: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK && resp.StatusCode != http.StatusAccepted {
		return fmt.Errorf("received non-200/202 status code: %s", resp.Status)
	}

	return nil
}
