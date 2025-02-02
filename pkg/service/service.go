package service

import (
	"github.com/skorpsrgvch/music-lib/models"
	"github.com/skorpsrgvch/music-lib/pkg/repository"
)

type SongDetail struct {
	ReleaseDate string `json:"releaseDate" example:"16.07.2006"`
	Text        string `json:"text" example:"Ooh baby, don't you know I suffer?\nOoh baby, can you hear me moan?\nYou caught me under false pretenses\nHow long before you let me go?\n\nOoh\nYou set my soul alight\nOoh\nYou set my soul alight"`
	Link        string `json:"link" example:"https://www.youtube.com/watch?v=Xsp3_a-PMTw"`
}

type Song interface {
	AddSong(list models.Song) error
	GetSongs(filter string, page int, limit int) ([]models.Song, error)
	GetSongText(id int) (string, error)
	UpdateSong(id int, song models.Song) error
	DeleteSong(id int) error
}

type Service struct {
	Song
}

func NewService(repos *repository.Repository) *Service {
	return &Service{
		Song: NewSongService(repos.Song),
	}
}
