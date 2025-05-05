package handler

import (
	"net/http"

	"1337b04rd/internal/ports/left"
)

type PostHandlers struct {
	service left.PostService
}

func NewPostHandlers(service left.PostService) *PostHandlers {
	return &PostHandlers{service: service}
}

func (h *PostHandlers) CreatePost(w http.ResponseWriter, r *http.Request) {
	// Обработка HTTP запроса
}

func (h *PostHandlers) GetPost(w http.ResponseWriter, r *http.Request) {
	// Обработка HTTP запроса
}

func (h *PostHandlers) GetPosts(w http.ResponseWriter, r *http.Request) {
	// Обработка HTTP запроса
}

func (h *PostHandlers) AddComment(w http.ResponseWriter, r *http.Request) {
	// Обработка HTTP запроса
}
