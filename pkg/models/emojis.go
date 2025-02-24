package models

type EmojiPack struct {
	Name   string `json:"name"`
	Author string `json:"author"`

	Emojis []Emoji `json:"emotes"`
}

type DownloadEmojiPack struct {
	Name   string
	Author string

	Emojis []DownloadedEmoji
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
