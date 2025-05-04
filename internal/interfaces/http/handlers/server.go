package handlers

import (
	"context"
	"net/http"

	"1337b04rd/pkg/logger"
)

type Server struct {
	addr   string
	router http.Handler
	logger *logger.CustomLogger // Добавлено логирование
}

func NewServer(port string, h *Handler) *Server {
	router := h.SetupRoutes()
	addr := ":" + port
	return &Server{
		addr:   addr,
		router: router,
		logger: h.logger, // Передаем логгер из Handler
	}
}

func (s *Server) Start() {
	s.logger.Info("Starting server", "address", s.addr)
	if err := http.ListenAndServe(s.addr, s.router); err != nil && err != http.ErrServerClosed {
		s.logger.Error("Server error", "error", err)
	}
}

func (s *Server) Shutdown(ctx context.Context) error {
	s.logger.Info("Shutting down server")
	srv := &http.Server{Addr: s.addr}
	return srv.Shutdown(ctx)
}
