package models

type EmojiPack struct {
	Emojis []Emoji `json:"emotes"`
}

type Emoji struct {
	Name string `json:"name"`
	URL  string `json:"url"`
}

type DownloadedEmoji struct {
	Name string
	Path string
}

type UploadedFile struct {
	ID string `json:"id"`
}
