package repository

import (
	"database/sql"
	"fmt"
	"strings"

	"github.com/jmoiron/sqlx"
	"github.com/sirupsen/logrus"
	"github.com/skorpsrgvch/music-lib/models"
)

type SongPostgres struct {
	db *sqlx.DB
}

func NewSongPostgres(db *sqlx.DB) *SongPostgres {
	return &SongPostgres{db: db}
}

func (r *SongPostgres) AddSong(song models.Song) error {
	query := `INSERT INTO songs (group_name, song, release_date, text, lyrics, link) VALUES ($1, $2, $3, $4, $5, $6)`

	_, err := r.db.Exec(query, song.GroupName, song.SongName, song.ReleaseDate, song.Text, song.Lyrics, song.Link)
	if err != nil {
		logrus.WithFields(logrus.Fields{
			"group_name":   song.GroupName,
			"song":         song.SongName,
			"release_date": song.ReleaseDate,
		}).Errorf("Failed to add song: %v", err)
		return err
	}

	logrus.WithFields(logrus.Fields{
		"group_name":   song.GroupName,
		"song":         song.SongName,
		"release_date": song.ReleaseDate,
	}).Debug("Song added successfully")
	return nil
}

func (s *SongPostgres) GetSongs(filter string, page int, limit int) ([]models.Song, error) {
	offset := (page - 1) * limit

	query := `
        SELECT id, group_name, song, release_date, text, lyrics, link 
        FROM songs 
        WHERE ($1 = '' OR group_name ILIKE $1 OR song ILIKE $1 OR lyrics ILIKE $1) 
        LIMIT $2 OFFSET $3
    `

	logrus.WithFields(logrus.Fields{
		"filter": filter,
		"page":   page,
		"limit":  limit,
		"offset": offset,
	}).Debug("Executing query to fetch songs")

	rows, err := s.db.Query(query, "%"+filter+"%", limit, offset)
	if err != nil {
		logrus.Errorf("Failed to execute query: %v", err)
		return nil, err
	}
	defer rows.Close()

	songs := make([]models.Song, 0, limit)
	for rows.Next() {
		var song models.Song
		if err := rows.Scan(&song.ID, &song.GroupName, &song.SongName, &song.ReleaseDate, &song.Text, &song.Lyrics, &song.Link); err != nil {
			logrus.Errorf("Failed to scan song: %v", err)
			return nil, err
		}
		songs = append(songs, song)
	}

	if err := rows.Err(); err != nil {
		logrus.Errorf("Error after iterating rows: %v", err)
		return nil, err
	}

	logrus.WithFields(logrus.Fields{
		"retrieved_songs": len(songs),
		"filter":          filter,
		"page":            page,
		"limit":           limit,
	}).Debug("Successfully retrieved songs")

	return songs, nil
}

func (s *SongPostgres) GetSongText(id int) (string, error) {
	// Проверка наличия записи с указанным id
	var exists bool
	checkQuery := `SELECT EXISTS(SELECT 1 FROM songs WHERE id = $1)`
	err := s.db.QueryRow(checkQuery, id).Scan(&exists)
	if err != nil {
		logrus.WithFields(logrus.Fields{
			"song_id": id,
		}).Errorf("Failed to check if song exists: %v", err)
		return "", err
	}
	if !exists {
		logrus.WithFields(logrus.Fields{
			"song_id": id,
		}).Warn("Song does not exist")
		return "", fmt.Errorf("song with id %d does not exist", id)
	}

	// Запрос текста песни
	query := `SELECT text FROM songs WHERE id = $1`
	var text string
	err = s.db.QueryRow(query, id).Scan(&text)
	if err != nil {
		if err == sql.ErrNoRows {
			logrus.WithFields(logrus.Fields{
				"song_id": id,
			}).Warn("No text found for song")
			return "", fmt.Errorf("no text found for song id %d", id)
		}
		logrus.WithFields(logrus.Fields{
			"song_id": id,
		}).Errorf("Failed to get text: %v", err)
		return "", err
	}

	logrus.WithFields(logrus.Fields{
		"song_id": id,
		"text":    len(text),
	}).Debug("Successfully retrieved song text")
	return text, nil
}

func (s *SongPostgres) UpdateSong(id int, song models.Song) error {
	setClauses := make([]string, 0)
	values := make([]interface{}, 0)
	valueIndex := 1

	// Проверяем каждое поле song на наличие значения
	if song.GroupName != "" {
		setClauses = append(setClauses, fmt.Sprintf("group_name = $%d", valueIndex))
		values = append(values, song.GroupName)
		valueIndex++
	}
	if song.SongName != "" {
		setClauses = append(setClauses, fmt.Sprintf("song = $%d", valueIndex))
		values = append(values, song.SongName)
		valueIndex++
	}
	if song.ReleaseDate != "" {
		setClauses = append(setClauses, fmt.Sprintf("release_date = $%d", valueIndex))
		values = append(values, song.ReleaseDate)
		valueIndex++
	}
	if song.Text != "" {
		setClauses = append(setClauses, fmt.Sprintf("text = $%d", valueIndex))
		values = append(values, song.Text)
		valueIndex++
	}
	if song.Lyrics != "" {
		setClauses = append(setClauses, fmt.Sprintf("lyrics = $%d", valueIndex))
		values = append(values, song.Lyrics)
		valueIndex++
	}
	if song.Link != "" {
		setClauses = append(setClauses, fmt.Sprintf("link = $%d", valueIndex))
		values = append(values, song.Link)
		valueIndex++
	}

	if len(setClauses) == 0 {
		logrus.WithFields(logrus.Fields{
			"song_id": id,
		}).Warn("No fields provided for update")
		return nil
	}

	// Собираем SQL-запрос динамически
	query := fmt.Sprintf("UPDATE songs SET %s WHERE id = $%d", strings.Join(setClauses, ", "), valueIndex)
	values = append(values, id)

	logrus.WithFields(logrus.Fields{
		"song_id": id,
		"fields":  setClauses,
	}).Debug("Executing update query")

	_, err := s.db.Exec(query, values...)
	if err != nil {
		logrus.WithFields(logrus.Fields{
			"song_id": id,
		}).Errorf("Failed to update song: %v", err)
		return err
	}

	logrus.WithFields(logrus.Fields{
		"song_id":        id,
		"updated_fields": setClauses,
	}).Info("Song updated successfully")
	return nil
}

func (s *SongPostgres) DeleteSong(id int) error {
	query := `DELETE FROM songs WHERE id = $1`

	logrus.WithFields(logrus.Fields{
		"song_id": id,
	}).Debug("Attempting to delete song")

	res, err := s.db.Exec(query, id)
	if err != nil {
		logrus.WithFields(logrus.Fields{
			"song_id": id,
		}).Errorf("Failed to delete song: %v", err)
		return err
	}

	rowsAffected, err := res.RowsAffected()
	if err != nil {
		logrus.WithFields(logrus.Fields{
			"song_id": id,
		}).Errorf("Failed to retrieve affected rows: %v", err)
		return err
	}

	if rowsAffected == 0 {
		logrus.WithFields(logrus.Fields{
			"song_id": id,
		}).Warn("No song found with the given ID")
		return fmt.Errorf("no song found with id %d", id)
	}

	logrus.WithFields(logrus.Fields{
		"song_id": id,
	}).Info("Song deleted successfully")

	return nil
}
