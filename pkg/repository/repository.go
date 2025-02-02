package repository

import (
	"github.com/jmoiron/sqlx"
	"github.com/skorpsrgvch/music-lib/models"
)

type Song interface {
	AddSong(list models.Song) error
	GetSongs(filter string, page int, limit int) ([]models.Song, error)
	GetSongText(id int) (string, error)
	UpdateSong(id int, song models.Song) error
	DeleteSong(id int) error
}

type Repository struct {
	Song
}

func NewRepository(db *sqlx.DB) *Repository {
	return &Repository{
		Song: NewSongPostgres(db),
	}
}
