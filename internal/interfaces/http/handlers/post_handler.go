package http

import (
	"encoding/json"
	"net/http"

	"1337b04rd/internal/application"
	"1337b04rd/internal/domain"
)

type PostHandler struct {
	postService *application.PostService
}

func NewPostHandler(postService *application.PostService) *PostHandler {
	return &PostHandler{postService: postService}
}

func (h *PostHandler) CreatePost(w http.ResponseWriter, r *http.Request) {
	// 1. Парсинг формы (title, content, image)
	if err := r.ParseMultipartForm(10 << 20); err != nil { // 10 MB
		http.Error(w, "Failed to parse form", http.StatusBadRequest)
		return
	}

	title := r.FormValue("title")
	content := r.FormValue("content")
	file, header, err := r.FormFile("image")
	if err != nil {
		http.Error(w, "Image is required", http.StatusBadRequest)
		return
	}
	defer file.Close()

	// 2. Получаем пользователя из сессии (пока заглушка)
	author := &domain.User{SessionID: "session123"} // Замените на реальную логику

	// 3. Создаём пост через сервис
	post, err := h.postService.CreatePost(r.Context(), title, content, file, header.Size, author)
	if err != nil {
		http.Error(w, "Failed to create post", http.StatusInternalServerError)
		return
	}

	// 4. Возвращаем ответ (например, редирект или JSON)
	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(post)
}
