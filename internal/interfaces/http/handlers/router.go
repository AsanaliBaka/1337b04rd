package handlers

import (
	"1337b04rd/internal/domain"
	"net/http"
)

type Handler struct {
	handler *PostHandler
}

func NewHandler(h domain.PostServer) *Handler {
	return &Handler{
		handler: NewPostHandler(h),
	}
}

func newRouter(h *Handler) *http.ServeMux {
	router := http.NewServeMux()

	router.HandleFunc("/posts", h.handler.GetAllPosts)
	router.HandleFunc("/posts/{id}", h.handler.GetById)
	router.HandleFunc("/posts/create", h.handler.CreatePostHandler)
	router.HandleFunc("/comments/create", h.handler.CreateComment)

	return router
}
