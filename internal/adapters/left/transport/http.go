package transport

import (
	"net/http"

	"1337b04rd/internal/ports/left"
	"1337b04rd/internal/ports/right"
	"1337b04rd/pkg/logger"
)

type Server struct {
	addr    string
	router  *http.ServeMux
	service left.APIPort
}

func NewHTTPServer(service left.APIPort, logger *logger.CustomLogger, imageUploader right.ImageStorage) *Server {
	router := newRouter(service, logger, imageUploader)

	addr := ":8080"
	return &Server{
		addr:    addr,
		router:  router,
		service: service,
	}
}

func newRouter(service left.APIPort, logger *logger.CustomLogger, imageUploader right.ImageStorage) *http.ServeMux {
	router := http.NewServeMux()

	SetupRoutes(service, logger, imageUploader, router)
	return router
}

func (s *Server) Serve() error {
	wrappedRouter := Chain(s.router, WithSession(s.service))

	server := &http.Server{
		Addr:    s.addr,
		Handler: wrappedRouter,
	}

	return server.ListenAndServe()
}
