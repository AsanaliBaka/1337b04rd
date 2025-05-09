package transport

import (
	"bytes"
	"context"
	"errors"
	"html/template"
	"log/slog"
	"net/http"
	"strings"
	"time"

	"1337b04rd/internal/domain"
	"1337b04rd/internal/ports/left"
	"1337b04rd/internal/ports/right"
	"1337b04rd/pkg"
	"1337b04rd/pkg/logger"
)

type Handler struct {
	service      left.APIPort
	templates    *template.Template
	imageStorage right.ImageStorage
	logger       *logger.CustomLogger
}

func NewPostHandler(postService left.APIPort, logger *logger.CustomLogger, imageStorage right.ImageStorage) *Handler {
	tmpl := template.Must(template.ParseGlob("web/templates/*.html"))
	return &Handler{
		service:      postService,
		templates:    tmpl,
		imageStorage: imageStorage,
		logger:       logger,
	}
}

func (h *Handler) HandleCatalog(w http.ResponseWriter, r *http.Request) {
	ctx, cancel := context.WithTimeout(r.Context(), 5*time.Second)
	defer cancel()

	data, err := h.service.GetCatalog(ctx)
	if err != nil {
		if errors.Is(ctx.Err(), context.DeadlineExceeded) {
			http.Error(w, "Request timed out", http.StatusGatewayTimeout)
			return
		}

		slog.Error("GetCatalog error", "error", err)
		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
		return
	}

	if err := h.templates.ExecuteTemplate(w, "catalog.html", data); err != nil {
		slog.Error("Failed to render template", "error", err)
		http.Error(w, "Render error", http.StatusInternalServerError)
		return
	}
}

func (h *Handler) HandleGetPost(w http.ResponseWriter, r *http.Request) {
	id := r.PathValue("id")

	ctx, cancel := context.WithTimeout(r.Context(), 5*time.Second)
	defer cancel()

	data, err := h.service.GetPostByID(ctx, id)
	if err != nil {
		if ctx.Err() == context.DeadlineExceeded {
			http.Error(w, "Request timed out", http.StatusGatewayTimeout)
			return
		}

		slog.Error(err.Error())
		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
		return
	}

	// Используем буфер для безопасного рендеринга шаблона
	var buf bytes.Buffer
	if err := h.templates.ExecuteTemplate(&buf, "post.html", data); err != nil {
		slog.Error("Failed to render template", "error", err)
		http.Error(w, "Render error", http.StatusInternalServerError)
		return
	}

	// Пишем рендеренный контент в ответ
	w.Header().Set("Content-Type", "text/html; charset=utf-8")
	_, err = w.Write(buf.Bytes())
	if err != nil {
		slog.Error("Failed to send response", "error", err)
	}
}

func (h *Handler) HandleCreatePostForm(w http.ResponseWriter, r *http.Request) {
	// Проверяем существование шаблона
	if h.templates.Lookup("create-post.html") == nil {
		slog.Error("Template not found", "template", "create-post.html")
		http.Error(w, "Template not found", http.StatusInternalServerError)
		return
	}

	// Добавляем данные, если нужно
	data := struct {
		Title string
	}{
		Title: "Create New Post",
	}

	// Рендерим шаблон
	err := h.templates.ExecuteTemplate(w, "create-post.html", data)
	if err != nil {
		slog.Error("Template execution error", "error", err)
		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
		return
	}
}

func (h *Handler) HandleSubmitPost(w http.ResponseWriter, r *http.Request) {
	if err := r.ParseMultipartForm(10 << 20); err != nil {
		http.Error(w, "Failed to parse form", http.StatusBadRequest)
		return
	}

	title := r.FormValue("title")
	content := r.FormValue("content")

	if title == "" || content == "" {
		http.Error(w, "Title and content are required", http.StatusBadRequest)
		return
	}

	file, header, err := r.FormFile("image")
	if err != nil {
		http.Error(w, "Failed to get image", http.StatusBadRequest)
		return
	}
	defer file.Close()

	session, ok := r.Context().Value(SessionKey).(*domain.Session)
	if !ok || session == nil {
		http.Error(w, "Unauthorized", http.StatusUnauthorized)
		return
	}

	ctx, cancel := context.WithTimeout(r.Context(), 5*time.Second)
	defer cancel()

	postID, err := pkg.GenerateUUID()
	if err != nil {
		slog.Error("UUID generation failed", "error", err)
		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
		return
	}

	post := &domain.Post{
		ID:        postID,
		Title:     title,
		Content:   content,
		Author:    session.UserID,
		CreatedAt: time.Now(),
	}

	// Загрузка изображения
	objectName, err := h.imageStorage.UploadImage(ctx, file, header)
	if err != nil {
		slog.Error("Image upload failed", "error", err)
		http.Error(w, "Failed to upload image", http.StatusInternalServerError)
		return
	}
	post.ImageURL = "/images/" + objectName

	if err := h.service.CreatePost(ctx, post); err != nil {
		slog.Error("Post creation failed", "error", err)
		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
		return
	}

	http.Redirect(w, r, "/catalog", http.StatusSeeOther)
}

func (h *Handler) HandleAddComment(w http.ResponseWriter, r *http.Request) {
	if err := r.ParseForm(); err != nil {
		http.Error(w, "Invalid form data", http.StatusBadRequest)
		return
	}

	session, ok := r.Context().Value(SessionKey).(*domain.Session)
	if !ok || session == nil {
		http.Error(w, "Unauthorized", http.StatusUnauthorized)
		return
	}

	postID := r.URL.Query().Get("id")
	parentID := r.FormValue("parent_comment_id")
	content := r.FormValue("content")

	if content == "" {
		http.Error(w, "Content are required", http.StatusBadRequest)
		return
	}

	uuid, err := pkg.GenerateUUID()
	if err != nil {
		slog.Error(err.Error())
	}

	comment := &domain.Comment{
		ID:        uuid,
		Author:    session.UserID,
		Content:   content,
		CreatedAt: time.Now(),
	}

	ctx, cancel := context.WithTimeout(r.Context(), 5*time.Second)
	defer cancel()

	if parentID != "" {
		if err := h.service.ReplyToComment(ctx, parentID, comment); err != nil {
			http.Error(w, "Failed to add reply: "+err.Error(), http.StatusInternalServerError)
			return
		}
	} else {
		if err := h.service.AddComment(ctx, postID, comment); err != nil {
			http.Error(w, "Failed to add comment: "+err.Error(), http.StatusInternalServerError)
			return
		}
	}

	http.Redirect(w, r, "/post/"+postID, http.StatusSeeOther)
}

func (h *Handler) HandleArchiveList(w http.ResponseWriter, r *http.Request) {
	ctx, cancel := context.WithTimeout(r.Context(), 5*time.Second)
	defer cancel()

	data, err := h.service.GetArchiveList(ctx)
	if err != nil {
		if ctx.Err() == context.DeadlineExceeded {
			http.Error(w, "Request timed out", http.StatusGatewayTimeout)
			return
		}

		slog.Error(err.Error())
		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "text/html; charset=utf-8")
	if err := h.templates.ExecuteTemplate(w, "archive.html", data); err != nil {
		slog.Error("Failed to render template", "error", err)
		http.Error(w, "Render error", http.StatusInternalServerError)
		return
	}
}

func (h *Handler) HandleGetArchivedPost(w http.ResponseWriter, r *http.Request) {
	id := r.PathValue("id")

	ctx, cancel := context.WithTimeout(r.Context(), 5*time.Second)
	defer cancel()

	data, err := h.service.GetArchivedPostByID(ctx, id)
	if err != nil {
		if ctx.Err() == context.DeadlineExceeded {
			http.Error(w, "Request timed out", http.StatusGatewayTimeout)
			return
		}

		slog.Error(err.Error())
		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")

	if err := h.templates.ExecuteTemplate(w, "archive-post.html", data); err != nil {
		slog.Error("Failed to render template", "error", err)
		http.Error(w, "Render error", http.StatusInternalServerError)
		return
	}
}

func (h *Handler) ServeImage(w http.ResponseWriter, r *http.Request) {
	parts := strings.Split(r.URL.Path, "/")
	if len(parts) < 3 {
		http.Error(w, "Invalid image path", http.StatusBadRequest)
		return
	}
	imageName := parts[2]

	data, contentType, err := h.imageStorage.GetImage(r.Context(), imageName)
	if err != nil {
		http.Error(w, "Image not found", http.StatusNotFound)
		return
	}

	w.Header().Set("Content-Type", contentType)
	w.Write(data)
}
