package main

import (
	"context"
	"os"
	"os/signal"
	"syscall"

	"github.com/joho/godotenv"
	"github.com/sirupsen/logrus"
	swaggerFiles "github.com/swaggo/files"
	ginSwagger "github.com/swaggo/gin-swagger"

	_ "github.com/lib/pq"
	ms "github.com/skorpsrgvch/music-lib"
	_ "github.com/skorpsrgvch/music-lib/docs" // Подключаем Swagger документацию
	"github.com/skorpsrgvch/music-lib/pkg/handler"
	"github.com/skorpsrgvch/music-lib/pkg/repository"
	"github.com/skorpsrgvch/music-lib/pkg/service"
	"github.com/spf13/viper"
)

// @title Music Info
// @version 0.0.1
// @description API for managing songs
// @host localhost:8000
// @BasePath /
func main() {
	// Настройка логирования
	logrus.SetFormatter(&logrus.TextFormatter{
		FullTimestamp: true,
	})
	logrus.SetLevel(logrus.DebugLevel) // Включаем debug-уровень логирования

	logrus.Info("Starting application...")

	// Инициализация конфигурации
	if err := initConfig(); err != nil {
		logrus.Fatalf("Error initializing configs: %s", err.Error())
	}
	logrus.Info("Config initialized successfully")

	// Загрузка переменных окружения
	if err := godotenv.Load(); err != nil {
		logrus.Fatalf("Error loading env variables: %s", err.Error())
	}
	logrus.Info("Environment variables loaded successfully")

	// Подключение к базе данных
	dbConfig := repository.Config{
		Host:     os.Getenv("DB_HOST"),
		Port:     os.Getenv("DB_PORT"),
		Username: os.Getenv("DB_USER"),
		DbName:   os.Getenv("DB_NAME"),
		Password: os.Getenv("DB_PASSWORD"),
		SSLMode:  os.Getenv("DB_SSLMODE"),
	}

	logrus.WithFields(logrus.Fields{
		"host": dbConfig.Host,
		"port": dbConfig.Port,
		"user": dbConfig.Username,
		"db":   dbConfig.DbName,
	}).Debug("Database configuration loaded")

	db, err := repository.NewPostgresDB(dbConfig)
	if err != nil {
		logrus.Fatalf("Error initializing DB: %s", err.Error())
	}
	logrus.Info("Database initialized successfully")

	// Инициализация репозитория, сервиса и обработчика
	repos := repository.NewRepository(db)
	logrus.Debug("Repository layer initialized")

	services := service.NewService(repos)
	logrus.Debug("Service layer initialized")

	handlers := handler.NewHandler(services)
	logrus.Debug("Handler layer initialized")

	// Запуск сервера
	server := new(ms.Server)
	port := viper.GetString("port")
	logrus.Infof("Starting HTTP server on port %s...", port)

	// Добавляем Swagger роут
	router := handlers.InitRoutes()
	router.GET("/swagger/*any", ginSwagger.WrapHandler(swaggerFiles.Handler))
	logrus.Infof("Swagger UI available at /swagger/index.html")

	go func() {
		if err := server.Run(port, router); err != nil {
			logrus.Fatalf("Error occurred while running HTTP server: %s", err.Error())
		}
	}()

	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGTERM, syscall.SIGINT)
	<-quit

	logrus.Infof("Application Shutting Down")

	if err := server.Shutdown(context.Background()); err != nil {
		logrus.Fatalf("error occured on server shutting down: %s", err.Error())
	}

	if err := db.Close(); err != nil {
		logrus.Fatalf("error occured on db connection close: %s", err.Error())
	}

}

func initConfig() error {
	viper.AddConfigPath("configs")
	viper.SetConfigName("config")
	err := viper.ReadInConfig()
	if err != nil {
		logrus.Error("Failed to read config file")
	} else {
		logrus.Info("Config file loaded successfully")
	}
	return err
}
