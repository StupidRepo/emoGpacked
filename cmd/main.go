package main

import (
	"bufio"
	"fmt"
	"github.com/StupidRepo/emoGPacked/pkg/models"
	"github.com/StupidRepo/emoGPacked/pkg/utils"
	"github.com/joho/godotenv"
	"github.com/sentinelb51/revoltgo"
	"os"
	"strconv"
	"strings"
	"time"
)

//var uploadGIFs *bool

func main() {
	//uploadGIFs = flag.Bool("g", false, "Whether to upload GIFs or not.")
	//flag.Parse()

	utils.Init()
	err := godotenv.Load()
	if err != nil {
		utils.Logger.Println("Failed to load the .env file: %v", err)
		// make the file
		_, err := os.Create(".env")
		if err != nil {
			utils.Logger.Println("Failed to create the .env file: %v", err)
		}
	}

	reader := bufio.NewReader(os.Stdin)

	token := os.Getenv("TOKEN")
	if token == "" {
		fmt.Print("Enter your user login token: ")
		token, _ = reader.ReadString('\n')
		token = strings.TrimSpace(token)

		if token != "" { // Save the token
			err := godotenv.Write(map[string]string{
				"TOKEN": token,
			}, ".env")
			if err != nil {
				utils.Logger.Printf("Failed to write the token to the .env file: %v", err)
			}
		}
	}

	utils.Session = revoltgo.New(token)
	err = utils.Session.Open()
	if err != nil {
		utils.Logger.Fatal(err)
	}

	time.Sleep(1 * time.Second)

	fmt.Print("Enter the URL of the pack to use: ")
	packURL, _ := reader.ReadString('\n')
	packURL = strings.TrimSpace(packURL)

	if token == "" {
		utils.Logger.Println("No token provided.")
	}

	if packURL == "" {
		utils.Logger.Println("No pack URL provided.")
	}

	err, servers := utils.GetServersList()
	if err != nil {
		utils.Logger.Fatalf("Failed to get the servers list from servers.json: %v", err)
	}

	fmt.Println("Servers:")
	for i, serverId := range servers {
		server, err := utils.Session.Server(serverId)
		if err != nil {
			utils.Logger.Fatalf("Failed to get the server %s: %v", serverId, err)
		}

		fmt.Printf("%d: %s\n", i, server.Name)
	}

	fmt.Print("Choose a server to upload the emojis to (default = 0): ")
	serverNumStr, _ := reader.ReadString('\n')
	serverNumStr = strings.TrimSpace(serverNumStr)

	serverNum := 0
	if serverNumStr != "" {
		serverNum, err = strconv.Atoi(serverNumStr)
		if err != nil {
			utils.Logger.Println("Failed to parse the server number: %v", err)
		}
	}

	//if !*uploadGIFs {
	//	utils.Logger.Println("GIFs will not be uploaded due to a weird Revolt issue where GIFs uploaded by this tool are not correctly displayed on the frontend." +
	//		"\nIf you want to upload GIFs, specify the -g flag.")
	//}

	StartEmojiProcess(packURL, servers[serverNum])
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
		//if !*uploadGIFs {
		//	mime, err := mimetype.DetectFile(fmt.Sprintf("%s/%s", *utils.TempDir, ending))
		//	if err != nil {
		//		utils.Logger.Printf("Failed to get the mimetype of the emoji %s: %v", emoji.Name, err)
		//		continue
		//	}
		//
		//	if mime.Is("image/gif") {
		//		utils.Logger.Printf("Skipping GIF emoji %s.", emoji.Name)
		//		continue
		//	}
		//}

		emojiSplit := strings.Split(emoji.Name, "-")
		emojiName := strings.Join(emojiSplit[1:], "")
		dlEmojis = append(dlEmojis, models.DownloadedEmoji{
			Name:   emojiName,
			Path:   fmt.Sprintf("%s/%s", *utils.TempDir, ending),
			Ending: ending,
		})
	}

	length := len(dlEmojis)
	if length == 0 {
		utils.Logger.Fatal("Failed to download any emojis.")
	}

	utils.Logger.Printf("Downloaded %d emojis.", length)

	for i, emoji := range dlEmojis {
		err, uplFile := utils.UploadEmojiToAutumn(emoji)
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
