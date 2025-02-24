package utils

import (
	"bytes"
	"encoding/json"
	"errors"
	"github.com/StupidRepo/emoGPacked/pkg/models"
	"github.com/sentinelb51/revoltgo"
	"io"
	"log"
	"mime/multipart"
	"net/http"
	"os"
	"time"
)

var Logger *log.Logger
var Session *revoltgo.Session
var Client *http.Client

var TempDir *string

func Init() {
	Logger = log.New(log.Writer(), "emoGPacked ", log.LstdFlags)
	Client = &http.Client{
		Timeout: time.Second * 20,
	}

	tempDir, err := os.MkdirTemp("", "emoGPacked")
	if err != nil {
		Logger.Fatalf("Failed to create a temporary directory: %v", err)
	}
	TempDir = &tempDir
}

func GET(url string) (error, io.ReadCloser) {
	// Get the response
	resp, err := Client.Get(url)
	if err != nil {
		return err, nil
	}

	return nil, resp.Body
}

func GETJson(url string, target interface{}) error {
	// Get the response
	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		return err
	}

	req.Header.Set("x-session-token", Session.Token)

	resp, err := Client.Do(req)
	if err != nil {
		return err
	}
	defer func(Body io.ReadCloser) {
		err := Body.Close()
		if err != nil {
			Logger.Printf("Failed to close the response body: %v", err)
		}
	}(resp.Body)

	// Decode the response
	return json.NewDecoder(resp.Body).Decode(target)
}

func GetServersList() (error, []string) {
	// Get from servers.json
	file, err := os.Open("servers.json")
	if err != nil {
		return err, nil
	}

	// Decode the JSON
	var servers *[]string
	err = json.NewDecoder(file).Decode(&servers)

	return err, *servers
}

// UploadEmojiToAutumn uploads an emoji to Revolt Autumn.
func UploadEmojiToAutumn(emoji models.DownloadedEmoji) (error, *models.UploadedFile) {
	// Open the file
	file, err := os.Open(emoji.Path)
	if err != nil {
		return err, nil
	}
	defer func(file *os.File) {
		err := file.Close()
		if err != nil {
			Logger.Fatalf("Failed to close the file: %v", err)
		}
	}(file)

	// Detect the MIME type of the file
	//mime, err := mimetype.DetectFile(emoji.Path)
	//if err != nil {
	//	return err, nil
	//}

	var requestBody bytes.Buffer
	writer := multipart.NewWriter(&requestBody)

	part, err := writer.CreateFormFile("file", emoji.Ending)
	if err != nil {
		return err, nil
	}

	// Copy the file to the writer
	_, err = io.Copy(part, file)
	if err != nil {
		return err, nil
	}

	// Close the writer before creating the request
	err = writer.Close()
	if err != nil {
		return err, nil
	}

	// Create the request
	req, err := http.NewRequest("POST", "https://autumn.revolt.chat/emojis", &requestBody)
	if err != nil {
		return err, nil
	}

	// Set the headers
	req.Header.Set("x-session-token", Session.Token)
	req.Header.Set("Content-Type", writer.FormDataContentType())

	resp, err := Client.Do(req)
	if err != nil {
		return err, nil
	}
	defer func(Body io.ReadCloser) {
		err := Body.Close()
		if err != nil {
			Logger.Fatalf("Failed to close the response body: %v", err)
		}
	}(resp.Body)

	bodyBytes, err := io.ReadAll(resp.Body)
	if err != nil {
		return err, nil
	}

	if resp.StatusCode != 200 {
		err, erruh := models.Error{}.MakeErrorFromResponse(bodyBytes)
		if err != nil {
			return err, nil
		}
		return errors.New(erruh.Type), nil
	}

	var uploadedFile models.UploadedFile
	err = json.Unmarshal(bodyBytes, &uploadedFile)
	if err != nil {
		return err, nil
	}

	return nil, &uploadedFile
}

func Download(reader io.Reader, fileName string) error {
	file, err := os.Create(*TempDir + "/" + fileName)
	if err != nil {
		return err
	}
	defer func(file *os.File) {
		err := file.Close()
		if err != nil {
			Logger.Fatalf("Failed to close the file: %v", err)
		}
	}(file)

	_, err = io.Copy(file, reader)
	if err != nil {
		return err
	}

	return nil
}
