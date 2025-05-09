package main

import (
	"log"

	"1337b04rd/internal/adapters/left/transport"
	"1337b04rd/internal/adapters/right/api"
	"1337b04rd/internal/adapters/right/db"
	"1337b04rd/internal/adapters/right/minio"
	"1337b04rd/internal/application"
	"1337b04rd/pkg/logger"
)

func main() {
	// Логгер
	logger, err := logger.NewCustomLogger()
	if err != nil {
		log.Fatalf("Failed to initialize logger: %v", err)
	}
	logger.Info("Logger initialized successfully")

	postgres := db.NewPostgres()
	defer postgres.Close()
	logger.Info("Database connection established successfully")

	// Инициализация MinIO
	minioClient, err := minio.NewImageStorage(
		"minio:9000",
		"minioadmin",
		"minioadmin",
		"images",
		false,
	)
	if err != nil {
		logger.Error("Failed to initialize MinIO:", err)
		log.Fatalf("MinIO initialization error: %v", err)
	}
	logger.Info("MinIO client initialized successfully")

	// Инициализация Rick and Morty API
	rickAndMortyAPI, err := api.NewRickAndMortyAPI()
	if err != nil {
		logger.Error("Failed to initialize Rick and Morty API:", err)
		log.Fatalf("Rick and Morty API initialization error: %v", err)
	}
	logger.Info("Rick and Morty API initialized successfully")
	user_service := application.NewUser()

	service := application.NewApp(postgres, rickAndMortyAPI, minioClient, *user_service)

	logger.Info("Service initialized successfully")
	// Запуск сервера
	server := transport.NewHTTPServer(service, logger, minioClient)
	if err := server.Serve(); err != nil {
		logger.Error("Failed to start server:", err)
		log.Fatalf("Server error: %v", err)
	}
}
