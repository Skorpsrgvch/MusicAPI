package handler

import (
	"github.com/gin-gonic/gin"
	"github.com/sirupsen/logrus"
	"github.com/skorpsrgvch/music-lib/pkg/service"
)

type Handler struct {
	services *service.Service
}

func NewHandler(services *service.Service) *Handler {
	logrus.Info("Initializing handler layer...")
	return &Handler{services: services}
}

func (h *Handler) InitRoutes() *gin.Engine {
	logrus.Info("Initializing routes...")

	router := gin.New()
	router.Use(gin.Logger())
	router.Use(gin.Recovery())

	songs := router.Group("/songs")
	{
		// @Summary Add a new song
		// @Description Add a new song to the library.
		// @Tags songs
		// @Accept json
		// @Produce json
		// @Param input body service.Song  true "Song info"
		// @Success 201 {object} string
		// @Failure 400 {string} string
		songs.POST("/", h.AddSong)
		// @Summary Get all songs
		// @Description Get a list of all songs
		// @Tags songs
		// @Produce json
		// @Success 200 {array} service.Song
		// @Failure 500 {string} string
		songs.GET("/", h.GetSongs)
		// @Summary Get song text by ID
		// @Description Get the lyrics of a song by ID
		// @Tags songs
		// @Produce text/plain
		// @Param id path int true "Song ID"
		// @Param page query int false "Page number"
		// @Param limit query int false "Page size"
		// @Success 200 {string} string
		// @Failure 400 {string} string
		// @Failure 500 {string} string
		songs.GET("/:id/text", h.GetSongText)
		// @Summary Update song by ID
		// @Description Update existing song.
		// @Tags songs
		// @Accept json
		// @Produce json
		// @Param id path int true "Song ID"
		// @Param input body service.Song  true "Song info"
		// @Success 200 {string} string
		// @Failure 400 {string} string
		// @Failure 500 {string} string
		songs.PUT("/:id", h.UpdateSong)
		// @Summary Delete song by ID
		// @Description Delete existing song.
		// @Tags songs
		// @Param id path int true "Song ID"
		// @Success 200 {string} string
		// @Failure 400 {string} string
		songs.DELETE("/:id", h.DeleteSong)
	}

	logrus.Info("Routes initialized successfully")
	return router
}

// logRequest — middleware, логирующее каждый вызов обработчика
func (h *Handler) logRequest(route string, handlerFunc gin.HandlerFunc) gin.HandlerFunc {
	return func(c *gin.Context) {
		logrus.WithFields(logrus.Fields{
			"method": c.Request.Method,
			"path":   c.Request.URL.Path,
			"route":  route,
		}).Info("Incoming request")
		handlerFunc(c)
	}
}
