# Calvin and Hobbes Project

The Calvin project is a Go project that provides functionality for generating random URLs, downloading images from those URLs, retrieving random quotes, uploading images, and sending messages. It is designed to work with GroupMe, a group messaging platform.

## Installation

To use the Calvin project, follow these steps:

1. Download the project files.
2. Install the Go programming language on your system.
3. Build the project by running the following command in the project directory:

   ```bash
   go build
   ```
4. Exectute the binary with: 

```bash
./calvin
```

## Functionality

The Calvin project provides the following functionality:

### URL Generator
- Implements the URLGenerator interface
- Generates a random URL by reading from a file (url.txt)

### ImageDownloader
- Implements the ImageDownloader interface.
- Downloads an image from a given URL using HTTP.
### QuoteProvider
- Implements the QuoteProvider interface.
- Retrieves a random quote from a text file (calvin/quotes.txt).
### ImageUploader
- Implements the ImageUploader interface.
- Uploads an image file to the GroupMe Image Service.
### MessageSender
- Implements the MessageSender interface.
- Sends a message with an attached image to a GroupMe group using the GroupMe Bots API.
### Calvin Struct
- Combines all the functionality into a single struct.
- Stores configuration data such as access tokens and command mappings.
- Provides methods to handle callbacks and send messages with images


## Usages
To use the Calvin project, you need to integrate it with a GroupMe bot and set up the necessary access tokens. See this example [here](https://dev.groupme.com/tutorials/bots) for generating the correct tokens. 

Calvin Struct setup example:

```go 

// Create a new instance of the Calvin struct
calvin := &Calvin{
    URLGenerator:    &DefaultURLGenerator{},
    ImageDownloader: &DefaultImageDownloader{},
    QuoteProvider:   &DefaultQuoteProvider{},
    ImageUploader:   &DefaultImageUploader{},
    MessageSender:   &DefaultMessageSender{},
    Command: map[string]bool{
        "next":  true,
        "!next": true,
    },
    BotID:       "<YOUR_BOT_ID>",
    AccessToken: "<YOUR_ACCESS_TOKEN>",
    Logger:      log.New(os.Stdout, "[Calvin]", log.LstdFlags),
}

// Set up an HTTP route to handle callbacks from GroupMe
http.HandleFunc("/callback", calvin.HandleCallback)

// Start the HTTP server
http.ListenAndServe(":8080", nil)

```


