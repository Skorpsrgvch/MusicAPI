package musiclib

import (
	"context"
	"net/http"
	"time"

	"github.com/sirupsen/logrus"
)

type Server struct {
	httpServer *http.Server
}

func (s *Server) Run(port string, handler http.Handler) error {
	logrus.Infof("Starting HTTP server on port %s...", port)

	s.httpServer = &http.Server{
		Addr:           ":" + port,
		Handler:        handler,
		MaxHeaderBytes: 1 << 20, // 1MB
		ReadTimeout:    10 * time.Second,
		WriteTimeout:   10 * time.Second,
	}

	err := s.httpServer.ListenAndServe()
	if err != nil {
		logrus.Errorf("HTTP server stopped with error: %v", err)
		return err
	}

	logrus.Info("HTTP server stopped successfully")
	return nil
}

func (s *Server) Shutdown(ctx context.Context) error {
	logrus.Info("Shutting down HTTP server...")

	err := s.httpServer.Shutdown(ctx)
	if err != nil {
		logrus.Errorf("Error while shutting down server: %v", err)
		return err
	}

	logrus.Info("HTTP server shut down successfully")
	return nil
}
