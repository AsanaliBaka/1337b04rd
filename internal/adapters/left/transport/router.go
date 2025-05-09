package transport

import (
	"net/http"

	"1337b04rd/internal/ports/left"
	"1337b04rd/internal/ports/right"
	"1337b04rd/pkg/logger"
)

func SetupRoutes(service left.APIPort, logger *logger.CustomLogger, imageStorage right.ImageStorage, router *http.ServeMux) {
	h := NewPostHandler(service, logger, imageStorage)

	router.HandleFunc("GET /catalog", h.HandleCatalog)
	router.HandleFunc("GET /post/{id}", h.HandleGetPost)
	router.HandleFunc("GET /archive", h.HandleArchiveList)
	router.HandleFunc("GET /archive/post/{id}", h.HandleGetArchivedPost)
	router.HandleFunc("GET /create-post", h.HandleCreatePostForm) // форма создания
	router.HandleFunc("POST /submit-post", h.HandleSubmitPost)    // отправка формы
	router.HandleFunc("POST /post/submit-comment", h.HandleAddComment)
	router.HandleFunc("GET /images/", h.ServeImage)
}
