package main

import (
	"bufio"
	"encoding/json"
	"fmt"
	"io"
	"io/ioutil"
	"log"
	"math/rand"
	"net/http"
	"os"
	"time"
)

type ImageResponse struct {
	Payload struct {
		URL        string `json:"url"`
		PictureURL string `json:"picture_url"`
	} `json:"payload"`
}

func downloadImage(url string) error {
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

	log.Printf("Image saved as %s\n", fileName)
	return nil
}

func getRandomURLFromFile() (string, error) {
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

func uploadImage(imagePath string, accessToken string) (string, error) {
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
		fmt.Println("this might be the actual URL: ", imageResp.Payload.URL)
		return imageResp.Payload.PictureURL, nil
	} else {
		return "", fmt.Errorf("image upload failed. Response body: %s", string(responseBody))
	}
}

func main() {
	// Get a random URL from the url.txt file
	imageURL, err := getRandomURLFromFile()
	if err != nil {
		log.Fatal(err)
	}

	err = downloadImage(imageURL)
	if err != nil {
		log.Fatal(err)
	}

	// Set the access token
	accessToken := os.Getenv("GM_TOKEN")
	// Call the function to upload the image

	imageURL, err = uploadImage("image.jpg", accessToken)
	if err != nil {
		log.Fatal(err)
	}

	fmt.Println("GroupMe Image URL:", imageURL)

}
