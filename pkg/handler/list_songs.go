package handler

import (
	"net/http"
	"strconv"
	"strings"

	"github.com/gin-gonic/gin"
	"github.com/sirupsen/logrus"
	"github.com/skorpsrgvch/music-lib/models"
	"github.com/skorpsrgvch/music-lib/pkg/service"
)

// @Summary Get song info
// @Description Get song details
// @Tags info
// @Produce json
// @Param group query string true "Group name"
// @Param song query string true "Song name"
// @Success 200 {object} service.SongDetail "Ok"
// @Failure 400 "Bad request"
// @Failure 500 "Internal server error"
// @Router /info [get]
func (h *Handler) GetInfo(c *gin.Context) {
	group := c.Query("group")
	song := c.Query("song")
	logrus.Debugf("GetInfo request for group %s and song %s", group, song)

	songDetail := service.SongDetail{
		ReleaseDate: "16.07.2006",
		Text:        "Ooh baby, don't you know I suffer?\nOoh baby, can you hear me moan?\nYou caught me under false pretenses\nHow long before you let me go?\n\nOoh\nYou set my soul alight\nOoh\nYou set my soul alight",
		Link:        "https://www.youtube.com/watch?v=Xsp3_a-PMTw",
	}
	c.JSON(http.StatusOK, songDetail)
}

// AddSong godoc
// @Summary Add a new song
// @Description Add a new song to the database
// @Tags songs
// @Accept json
// @Produce json
// @Param song body models.Song true "Song JSON"
// @Success 201 {object} map[string]string
// @Failure 400 {object} map[string]string
// @Failure 500 {object} map[string]string
// @Router /songs/ [post]
// Добавление песни
func (h *Handler) AddSong(c *gin.Context) {
	var song models.Song

	// Логируем запрос
	logrus.Info("Received request to add a new song")

	if err := c.ShouldBindJSON(&song); err != nil {
		logrus.Warnf("Invalid request body: %v", err)
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid request body"})
		return
	}

	logrus.WithFields(logrus.Fields{
		"group_name":   song.GroupName,
		"song":         song.SongName,
		"release_date": song.ReleaseDate,
	}).Info("Adding new song")

	if err := h.services.AddSong(song); err != nil {
		logrus.Errorf("Failed to add song: %v", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to add song"})
		return
	}

	logrus.Info("Song added successfully")
	c.JSON(http.StatusCreated, gin.H{"message": "Song added successfully"})
}

// GetSongs godoc
// @Summary Get all songs
// @Description Get a list of all songs with optional filtering
// @Tags songs
// @Accept json
// @Produce json
// @Param filter query string false "Filter by group_name, song or lyrics"
// @Param page query int false "Page number"
// @Param limit query int false "Number of results per page"
// @Success 200 {array} models.Song
// @Failure 500 {object} map[string]string
// @Router /songs/ [get]
// Получение списка песен с фильтрацией и пагинацией
func (h *Handler) GetSongs(c *gin.Context) {
	filter := c.Query("filter")
	page, _ := strconv.Atoi(c.DefaultQuery("page", "1"))    // По умолчанию page = 1
	limit, _ := strconv.Atoi(c.DefaultQuery("limit", "10")) // По умолчанию limit = 10

	// Логируем параметры запроса
	logrus.WithFields(logrus.Fields{
		"filter": filter,
		"page":   page,
		"limit":  limit,
	}).Info("Fetching songs with filters")

	songs, err := h.services.GetSongs(filter, page, limit)
	if err != nil {
		logrus.Errorf("Failed to get songs: %v", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to get songs"})
		return
	}

	logrus.Infof("Successfully retrieved %d songs", len(songs))
	c.JSON(http.StatusOK, songs)
}

// GetSongText godoc
// @Summary Get song text
// @Description Get the lyrics of a song by its ID with optional pagination
// @Tags songs
// @Accept json
// @Produce text/plain
// @Param id path int true "Song ID"
// @Param page query int false "Page number (default: 1)"
// @Param limit query int false "Number of verses per page (default: 5)"
// @Success 200 {string} string "Song text"
// @Failure 400 {object} map[string]string "Invalid song ID or page number"
// @Failure 500 {object} map[string]string "Failed to get song text"
// @Router /songs/{id}/text [get]
// Получение текста песни с пагинацией
func (h *Handler) GetSongText(c *gin.Context) {
	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		logrus.Warnf("Invalid song ID: %v", err)
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid song ID"})
		return
	}

	page, err := strconv.Atoi(c.DefaultQuery("page", "1"))
	if err != nil {
		logrus.Warnf("Invalid page number: %v", err)
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid page number"})
		return
	}
	pageSize, err := strconv.Atoi(c.DefaultQuery("limit", "5"))
	if err != nil {
		logrus.Warnf("Invalid page size: %v", err)
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid page size"})
		return
	}

	logrus.WithFields(logrus.Fields{
		"song_id":  id,
		"page":     page,
		"pageSize": pageSize,
	}).Info("Fetching song text")

	text, err := h.services.GetSongText(id)
	if err != nil {
		logrus.Errorf("Failed to get text from service: %v", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to get text"})
		return
	}

	text = strings.ReplaceAll(text, "\\n", "\n")
	verses := splitVerses(text)
	start := (page - 1) * pageSize
	end := start + pageSize

	if start >= len(verses) {
		logrus.Infof("Requested page %d exceeds available verses", page)
		c.JSON(http.StatusOK, gin.H{"text": ""})
		return
	}

	if end > len(verses) {
		end = len(verses)
	}

	logrus.Infof("Returning verses %d to %d", start, end)
	c.Header("Content-Type", "text/plain; charset=utf-8")
	c.String(http.StatusOK, strings.Join(verses[start:end], "\n\n")) // Разделение куплетов двойным переносом
}

// Разбивает текст на куплеты
func splitVerses(text string) []string {
	verses := strings.Split(text, "\\n\\n")
	for i := range verses {
		verses[i] = strings.TrimSuffix(verses[i], "\\n")
	}
	return verses
}

// UpdateSong godoc
// @Summary Update a song
// @Description Update details of an existing song by its ID
// @Tags songs
// @Accept json
// @Produce json
// @Param id path int true "Song ID"
// @Param song body models.Song true "Updated song data"
// @Success 200 {object} map[string]string "Song updated successfully"
// @Failure 400 {object} map[string]string "Invalid request body"
// @Failure 500 {object} map[string]string "Failed to update song"
// @Router /songs/{id} [put]
// Обновление информации о песне
func (h *Handler) UpdateSong(c *gin.Context) {
	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		logrus.Warnf("Invalid song ID for update: %v", err)
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid song ID"})
		return
	}

	var song models.Song
	if err := c.ShouldBindJSON(&song); err != nil {
		logrus.Warnf("Invalid request body: %v", err)
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid request body"})
		return
	}

	logrus.WithFields(logrus.Fields{
		"song_id":     id,
		"group_name":  song.GroupName,
		"song_name":   song.SongName,
		"releaseDate": song.ReleaseDate,
	}).Info("Updating song")

	if err := h.services.UpdateSong(id, song); err != nil {
		logrus.Errorf("Failed to update song: %v", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to update song"})
		return
	}

	logrus.Infof("Song with ID %d updated successfully", id)
	c.JSON(http.StatusOK, gin.H{"message": "Song updated successfully"})
}

// DeleteSong godoc
// @Summary Delete a song
// @Description Remove a song from the database by its ID
// @Tags songs
// @Accept json
// @Produce json
// @Param id path int true "Song ID"
// @Success 200 {object} map[string]string "Song deleted successfully"
// @Failure 500 {object} map[string]string "Failed to delete song"
// @Router /songs/{id} [delete]
// Удаление песни
func (h *Handler) DeleteSong(c *gin.Context) {
	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		logrus.Warnf("Invalid song ID for deletion: %v", err)
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid song ID"})
		return
	}

	logrus.Infof("Deleting song with ID %d", id)
	if err := h.services.DeleteSong(id); err != nil {
		logrus.Errorf("Failed to delete song: %v", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to delete song"})
		return
	}

	logrus.Infof("Song with ID %d deleted successfully", id)
	c.JSON(http.StatusOK, gin.H{"message": "Song deleted successfully"})
}
