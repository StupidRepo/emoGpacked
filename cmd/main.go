package main

import (
	"flag"
	"fmt"
	"github.com/StupidRepo/emoGPacked/pkg/models"
	"github.com/StupidRepo/emoGPacked/pkg/utils"
	"github.com/gabriel-vasile/mimetype"
	"github.com/sentinelb51/revoltgo"
	"os"
	"strings"
)

var uploadGIFs *bool

func main() {
	packURL := flag.String("p", "", "The URL of the pack to use.")
	serverID := flag.String("s", "", "The server ID to upload the emojis to.")
	token := flag.String("t", "", "Your user login token.")
	uploadGIFs = flag.Bool("g", false, "Whether to upload GIFs or not.")
	flag.Parse()

	utils.Init()

	if !*uploadGIFs {
		utils.Logger.Println("GIFs will not be uploaded due to a weird Revolt issue where GIFs uploaded by this tool are not correctly displayed on the frontend." +
			" If you want to upload GIFs, specify the -g flag.")
	}

	if *token == "" {
		utils.Logger.Fatal("No token provided. Please provide a token in the .env file.")
	}

	if *packURL == "" {
		utils.Logger.Fatal("No pack URL provided. Please provide a pack URL using the -pack flag.")
	}

	if *serverID == "" {
		utils.Logger.Fatal("No server ID provided. Please provide a server ID using the -server flag.")
	}

	utils.Session = revoltgo.New(*token)
	err := utils.Session.Open()
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

func StartEmojiProcess(url *string, serverID *string) {
	// Get the emoji pack
	var emojiPack models.EmojiPack
	err := utils.GETJson(*url, &emojiPack)
	if err != nil {
		utils.Logger.Fatalf("Failed to get the emoji pack: %v", err)
	}

	dlEmojiPack := models.DownloadEmojiPack{
		Name:   emojiPack.Name,
		Author: emojiPack.Author,
	}

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
		dlEmojiPack.Emojis = append(dlEmojiPack.Emojis, models.DownloadedEmoji{
			Name: emojiName,
			Path: fmt.Sprintf("%s/%s", *utils.TempDir, ending),
		})
	}

	length := len(dlEmojiPack.Emojis)
	if length == 0 {
		utils.Logger.Fatal("Failed to download any emojis.")
	}

	utils.Logger.Printf("Downloaded %d emojis.", length)

	for i, emoji := range dlEmojiPack.Emojis {
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
					ID:   *serverID,
				},
			},
		)

		if err != nil {
			utils.Logger.Printf("Failed to create emoji %s on server: %v", emoji.Name, err)
			continue
		}

		utils.Logger.Printf("Uploading emojis... (%.2f%%)", length, float64(i+1)/float64(length)*100)
	}

	utils.Logger.Println("Finished uploading emojis!")
}
