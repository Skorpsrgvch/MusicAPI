package models

type Song struct {
	ID          int    `json:"id"`
	GroupName   string `json:"group" binding:"required"`
	SongName    string `json:"song" binding:"required"`
	ReleaseDate string `json:"releaseDate" `
	Text        string `json:"text"`
	Lyrics      string `json:"lyrics"`
	Link        string `json:"link"`
}
