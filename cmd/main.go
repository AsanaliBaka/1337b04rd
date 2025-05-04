package main

import (
	"context"
	"log"
	"os"
	"os/signal"
	"syscall"
	"time"

	"1337b04rd/internal/domain"
	"1337b04rd/internal/infrastructure/api"
	"1337b04rd/internal/infrastructure/config"
	"1337b04rd/internal/infrastructure/minio"
	"1337b04rd/internal/infrastructure/postgres"
	"1337b04rd/internal/interfaces/http/handlers"
	"1337b04rd/pkg/logger"

	"github.com/jackc/pgx/v5/pgxpool"
)

func main() {
	// Инициализация конфигурации
	cfg, err := config.LoadConfig()
	if err != nil {
		log.Fatalf("Failed to load config: %v", err)
	}

	// Инициализация логгера
	appLogger, err := logger.NewCustomLogger()
	if err != nil {
		log.Fatalf("Failed to initialize logger: %v", err)
	}
	appLogger.Info("Application starting...", "version", "1.0.0")

	// Инициализация подключения к PostgreSQL
	dbConnString := buildDBConnectionString(cfg)
	dbPool, err := pgxpool.New(context.Background(), dbConnString)
	if err != nil {
		appLogger.Error("Failed to connect to database", "error", err)
		os.Exit(1)
	}
	defer dbPool.Close()

	if err := dbPool.Ping(context.Background()); err != nil {
		appLogger.Error("Failed to ping database", "error", err)
		os.Exit(1)
	}
	appLogger.Info("Successfully connected to database")

	// Инициализация хранилища изображений (MinIO/S3)
	imageStorage, err := minio.NewImageStrorage(cfg, context.Background())
	if err != nil {
		appLogger.Error("Failed to initialize image storage", "error", err)
		os.Exit(1)
	}
	appLogger.Info("Image storage initialized")

	// Инициализация Rick and Morty API
	rickAndMortyAPI := api.NewRickAndMortyAPI(appLogger)
	appLogger.Info("Rick and Morty API client initialized")

	// Инициализация репозиториев
	postRepo := postgres.NewPostRepo(dbPool)
	commentRepo := postgres.NewCommentRepo(dbPool)
	sessionRepo := postgres.NewSessionRepo(dbPool)

	// Инициализация сервисов
	sessionService := domain.NewSessionService(postRepo)
	postService := domain.NewPostServer(postRepo, commentRepo, imageStorage, sessionService)
	sessionManager := domain.NewSession(sessionRepo)

	// Инициализация HTTP обработчиков
	handler := handlers.NewHandler(
		postService,
		sessionManager,
		rickAndMortyAPI,
		appLogger,
	)

	// Создание HTTP сервера с использованием вашей структуры Server
	server := handlers.NewServer(cfg.Port, handler)

	// Канал для graceful shutdown
	shutdownChan := make(chan os.Signal, 1)
	signal.Notify(shutdownChan, syscall.SIGINT, syscall.SIGTERM)

	// Запуск сервера в отдельной горутине
	go func() {
		appLogger.Info("Starting HTTP server", "port", cfg.Port)
		server.Start()
	}()

	// Ожидание сигнала завершения
	sig := <-shutdownChan
	appLogger.Info("Received shutdown signal", "signal", sig)

	// Настройка graceful shutdown
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	if err := server.Shutdown(ctx); err != nil {
		appLogger.Error("Server shutdown failed", "error", err)
	} else {
		appLogger.Info("Server stopped gracefully")
	}
}

func buildDBConnectionString(cfg *config.Config) string {
	return "postgres://" + cfg.DBUser + ":" + cfg.DBPassword + "@" + cfg.DBHost + ":" + cfg.DBPort + "/" + cfg.DBName
}
