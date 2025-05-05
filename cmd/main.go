package main

import (
	"context"
	"fmt"
	"log/slog"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"1337b04rd/internal/adapters/left/transport/handler"
	postgres "1337b04rd/internal/adapters/right/db"
	"1337b04rd/internal/adapters/right/minio"
	"1337b04rd/internal/application"
)

func main() {
	// Инициализация логгера
	logger := slog.New(slog.NewJSONHandler(os.Stdout, nil))
	slog.SetDefault(logger)

	// Инициализация приложения
	app, cleanup, err := setupApplication()
	if err != nil {
		slog.Error("Application setup failed", "error", err)
		os.Exit(1)
	}
	defer cleanup()

	// Настройка HTTP сервера
	server := setupServer(app)

	// Запуск сервера и воркера
	startServer(server)
}

func setupApplication() (*application.Application, func(), error) {
	// Инициализация БД
	dbRepo, err := postgres.NewPostgresRepository(os.Getenv("DB_DSN"))
	if err != nil {
		return nil, nil, fmt.Errorf("database initialization failed: %w", err)
	}

	// Инициализация MinIO
	minioClient, err := minio.New(minio.DefaultConfig())
	if err != nil {
		dbRepo.Close()
		return nil, nil, fmt.Errorf("minio initialization failed: %w", err)
	}

	// Создание зависимостей приложения
	app := application.New(dbRepo, minioClient)

	// Очистка ресурсов
	cleanup := func() {
		dbRepo.Close()
	}

	return app, cleanup, nil
}

func setupServer() *http.Server {
	router := handler.SetupRoutes()
	server := &http.Server{
		Addr:         ":8080",
		Handler:      router,
		ReadTimeout:  10 * time.Second,
		WriteTimeout: 30 * time.Second,
		IdleTimeout:  60 * time.Second,
	}
	return server
}

func startServer(server *http.Server) {
	done := make(chan os.Signal, 1)
	signal.Notify(done, os.Interrupt, syscall.SIGINT, syscall.SIGTERM)

	go func() {
		slog.Info("Starting server", "address", server.Addr)
		if err := server.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			slog.Error("Server failed", "error", err)
		}
	}()

	<-done
	slog.Info("Shutting down server...")

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	if err := server.Shutdown(ctx); err != nil {
		slog.Error("Server shutdown failed", "error", err)
	}
}
