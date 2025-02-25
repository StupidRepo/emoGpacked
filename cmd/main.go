package main

import (
	"bufio"
	"encoding/json"
	"flag"
	"fmt"
	"github.com/StupidRepo/emoGPacked/pkg/models"
	"github.com/StupidRepo/emoGPacked/pkg/utils"
	"github.com/gabriel-vasile/mimetype"
	"github.com/sentinelb51/revoltgo"
	"os"
	"strings"
	"time"
)

var uploadGIFs *bool

func main() {
	uploadGIFs = flag.Bool("g", false, "Whether to upload GIFs or not.")
	flag.Parse()

	utils.Init()

	reader := bufio.NewReader(os.Stdin)

	var token string
	var serverID string

	err, defs := TryGetDefaults()
	if err != nil {
		utils.Logger.Printf("Failed to get defaults: %v", err)
	} else {
		if time.Since(defs.LastStoredToken) < time.Minute*3 {
			token = defs.Token
		} else {
			err := os.Remove("defaults.json")
			if err != nil {
				utils.Logger.Printf("Failed to remove the defaults file: %v", err)
			}
		}

		serverID = defs.ServerId
	}

	if token == "" {
		fmt.Print("Enter your user login token: ")
		token, _ = reader.ReadString('\n')
		token = strings.TrimSpace(token)
	}

	if serverID == "" {
		fmt.Print("Enter the server ID to upload the emojis to: ")
		serverID, _ = reader.ReadString('\n')
		serverID = strings.TrimSpace(serverID)
	}

	fmt.Print("Enter the URL of the pack to use: ")
	packURL, _ := reader.ReadString('\n')
	packURL = strings.TrimSpace(packURL)

	SaveDefaults(models.DefaultOptions{
		Token:    token,
		ServerId: serverID,

		LastStoredToken: time.Now(),
	})

	if !*uploadGIFs {
		utils.Logger.Println("GIFs will not be uploaded due to a weird Revolt issue where GIFs uploaded by this tool are not correctly displayed on the frontend." +
			"\nIf you want to upload GIFs, specify the -g flag.")
	}

	if token == "" {
		utils.Logger.Println("No token provided.")
	}

	if serverID == "" {
		utils.Logger.Println("No server ID provided.")
	}

	if packURL == "" {
		utils.Logger.Println("No pack URL provided.")
	}

	utils.Session = revoltgo.New(token)
	err = utils.Session.Open()
	if err != nil {
		utils.Logger.Fatal(err)
	}

	StartEmojiProcess(packURL, serverID)
	defer func() {
		utils.Logger.Println("Gracefully shutting down! :D")

		err = os.RemoveAll(*utils.TempDir)
		if err != nil {
			utils.Logger.Println("Failed to remove the temporary directory: %v", err)
		}

		err = utils.Session.Close()
		if err != nil {
			utils.Logger.Println("Lmao, failed to shutdown: %v", err)
		}
	}()
}

func TryGetDefaults() (error, *models.DefaultOptions) {
	var defaults models.DefaultOptions

	file, err := os.Open("defaults.json")
	if err != nil {
		return err, nil
	}
	defer func(file *os.File) {
		err := file.Close()
		if err != nil {
			utils.Logger.Fatalf("Failed to close the defaults file: %v", err)
		}
	}(file)

	err = json.NewDecoder(file).Decode(&defaults)
	if err != nil {
		utils.Logger.Fatalf("Failed to decode the defaults file: %v", err)
	}

	return nil, &defaults
}

func SaveDefaults(options models.DefaultOptions) {
	file, err := os.Create("defaults.json")
	if err != nil {
		utils.Logger.Fatalf("Failed to create the defaults file: %v", err)
	}
	defer func(file *os.File) {
		err := file.Close()
		if err != nil {
			utils.Logger.Fatalf("Failed to close the defaults file: %v", err)
		}
	}(file)

	err = json.NewEncoder(file).Encode(options)
	if err != nil {
		utils.Logger.Fatalf("Failed to encode the defaults file: %v", err)
	}
}

func StartEmojiProcess(url string, serverID string) {
	// Get the emoji pack
	var emojiPack models.EmojiPack
	err := utils.GETJson(url, &emojiPack)
	if err != nil {
		utils.Logger.Fatalf("Failed to get the emoji pack: %v", err)
	}

	var dlEmojis []models.DownloadedEmoji

	// Download the emojis
	for _, emoji := range emojiPack.Emojis {
		err, data := utils.GET(emoji.URL)
		if err != nil {
			utils.Logger.Printf("Failed to get emoji %s: %v", emoji.Name, err)
			continue
		}

		paths := strings.Split(emoji.URL, "/")
		ending := paths[len(paths)-1]

		err = utils.Download(data, ending)
		if err != nil {
			utils.Logger.Printf("Failed to download emoji %s: %v", emoji.Name, err)
			continue
		}

		// Get mimetype of data and check if it's a GIF
		if !*uploadGIFs {
			mime, err := mimetype.DetectFile(fmt.Sprintf("%s/%s", *utils.TempDir, ending))
			if err != nil {
				utils.Logger.Printf("Failed to get the mimetype of the emoji %s: %v", emoji.Name, err)
				continue
			}

			if mime.Is("image/gif") {
				utils.Logger.Printf("Skipping GIF emoji %s.", emoji.Name)
				continue
			}
		}

		emojiSplit := strings.Split(emoji.Name, "-")
		emojiName := strings.Join(emojiSplit[1:], "")
		dlEmojis = append(dlEmojis, models.DownloadedEmoji{
			Name: emojiName,
			Path: fmt.Sprintf("%s/%s", *utils.TempDir, ending),
		})
	}

	length := len(dlEmojis)
	if length == 0 {
		utils.Logger.Fatal("Failed to download any emojis.")
	}

	utils.Logger.Printf("Downloaded %d emojis.", length)

	for i, emoji := range dlEmojis {
		err, uplFile := utils.UploadEmojiToAutumn(emoji.Path)
		if err != nil {
			utils.Logger.Printf("Failed to upload emoji %s: %v", emoji.Name, err)
			continue
		}

		_, err = utils.Session.EmojiCreate(
			uplFile.ID,
			revoltgo.EmojiCreateData{
				Name: emoji.Name,
				Parent: &revoltgo.EmojiParent{
					Type: "Server",
					ID:   serverID,
				},
			},
		)

		if err != nil {
			utils.Logger.Printf("Failed to create emoji %s on server: %v", emoji.Name, err)
			continue
		}

		utils.Logger.Printf("Uploading emojis... (%.2f%%)", float64(i+1)/float64(length)*100)
	}

	utils.Logger.Println("Finished uploading emojis!")
}
