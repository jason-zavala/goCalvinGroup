package main

import (
	"io"
	"log"
	"net/http"
	"os"
)

func main() {
	// Set the URL of the image
	imageURL := "https://assets.amuniversal.com/d6d11040deb701317193005056a9545d"

	// Send a GET request to the image URL
	resp, err := http.Get(imageURL)
	if err != nil {
		log.Fatalf("Failed to retrieve the image: %v", err)
	}
	defer resp.Body.Close()

	// Create a new file with the .jpg extension
	fileName := "image.jpg"
	file, err := os.Create(fileName)
	if err != nil {
		log.Fatalf("Failed to create the file: %v", err)
	}
	defer file.Close()

	// Copy the image data to the file
	_, err = io.Copy(file, resp.Body)
	if err != nil {
		log.Fatalf("Failed to save the image: %v", err)
	}

	log.Printf("Image saved as %s\n", fileName)
}
