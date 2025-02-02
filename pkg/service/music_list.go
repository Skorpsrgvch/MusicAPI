package service

import (
	"github.com/skorpsrgvch/music-lib/models"
	"github.com/skorpsrgvch/music-lib/pkg/repository"
)

type SongService struct {
	repo repository.Song
}

func NewSongService(repo repository.Song) *SongService {
	return &SongService{repo: repo}
}

func (s *SongService) AddSong(list models.Song) error {
	return s.repo.AddSong(list)
}

func (s *SongService) GetSongs(filter string, page int, limit int) ([]models.Song, error) {
	return s.repo.GetSongs(filter, page, limit)
}
func (s *SongService) GetSongText(id int) (string, error) {
	return s.repo.GetSongText(id)
}
func (s *SongService) UpdateSong(id int, song models.Song) error {
	return s.repo.UpdateSong(id, song)
}
func (s *SongService) DeleteSong(id int) error {
	return s.repo.DeleteSong(id)
}
