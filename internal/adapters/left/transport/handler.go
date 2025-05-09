package transport

import (
	"bytes"
	"context"
	"errors"
	"html/template"
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
			h.logger.Error("Request timed out", "error", err)
			h.renderError(w, http.StatusGatewayTimeout, "Request timed out")
			return
		}

		h.logger.Error("Failed to get catalog", "error", err)
		h.renderError(w, http.StatusInternalServerError, "Internal Server Error")
		return
	}

	if err := h.templates.ExecuteTemplate(w, "catalog.html", data); err != nil {
		h.logger.Error("Failed to render template", "error", err)
		h.renderError(w, http.StatusInternalServerError, "Render error")
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
			h.logger.Error("Request timed out", "error", err)
			h.renderError(w, http.StatusGatewayTimeout, "Request timed out")
			return
		}

		h.logger.Error("Failed to get post", "error", err)
		h.renderError(w, http.StatusInternalServerError, "Internal Server Error")
		return
	}

	var buf bytes.Buffer
	if err := h.templates.ExecuteTemplate(&buf, "post.html", data); err != nil {
		h.logger.Error("Render error", "error", err)
		h.renderError(w, http.StatusInternalServerError, "Render error")
		return
	}

	w.Header().Set("Content-Type", "text/html; charset=utf-8")
	_, err = w.Write(buf.Bytes())
	if err != nil {
		h.logger.Error("Write error", "error", err)
		h.renderError(w, http.StatusInternalServerError, "Write error")
		return
	}
}

func (h *Handler) HandleCreatePostForm(w http.ResponseWriter, r *http.Request) {
	if h.templates.Lookup("create-post.html") == nil {
		h.logger.Error("Template not found", "template", "create-post.html")
		h.renderError(w, http.StatusInternalServerError, "Template not found")
		return
	}

	data := struct {
		Title string
	}{
		Title: "Create New Post",
	}

	err := h.templates.ExecuteTemplate(w, "create-post.html", data)
	if err != nil {
		h.logger.Error("Failed to execute template", "error", err)
		h.renderError(w, http.StatusInternalServerError, "Internal Server Error")
		return
	}
}

func (h *Handler) HandleArchiveList(w http.ResponseWriter, r *http.Request) {
	ctx, cancel := context.WithTimeout(r.Context(), 5*time.Second)
	defer cancel()

	data, err := h.service.GetArchiveList(ctx)
	if err != nil {
		if ctx.Err() == context.DeadlineExceeded {
			h.logger.Error("Request timed out", "error", err)
			h.renderError(w, http.StatusGatewayTimeout, "Request timed out")
			return
		}

		h.logger.Error("Failed to get archive list", "error", err)
		h.renderError(w, http.StatusInternalServerError, "Internal Server Error")
		return
	}

	w.Header().Set("Content-Type", "text/html; charset=utf-8")
	if err := h.templates.ExecuteTemplate(w, "archive.html", data); err != nil {
		h.logger.Error("Failed to render template", "error", err)
		h.renderError(w, http.StatusInternalServerError, "Render error")
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
			h.logger.Error("Request timed out", "error", err)
			h.renderError(w, http.StatusGatewayTimeout, "Request timed out")
			return
		}

		h.logger.Error("Failed to get archived post", "error", err)
		h.renderError(w, http.StatusInternalServerError, "Internal Server Error")
		return
	}

	w.Header().Set("Content-Type", "application/json")

	if err := h.templates.ExecuteTemplate(w, "archive-post.html", data); err != nil {
		h.logger.Error("Render error", "error", err)
		h.renderError(w, http.StatusInternalServerError, "Render error")
		return
	}
}

func (h *Handler) HandleSubmitPost(w http.ResponseWriter, r *http.Request) {
	if err := r.ParseMultipartForm(10 << 20); err != nil {
		h.logger.Error("Failed to parse form", "error", err)
		h.renderError(w, http.StatusBadRequest, "Failed to parse form")
		return
	}

	title := r.FormValue("title")
	content := r.FormValue("content")

	if title == "" || content == "" {
		h.logger.Error("Title and content are required")
		h.renderError(w, http.StatusBadRequest, "Title and content are required")
		return
	}

	file, header, err := r.FormFile("image")
	if err != nil {
		h.logger.Error("Failed to get file from form", "error", err)
		h.renderError(w, http.StatusBadRequest, "Failed to get file from form")
		return
	}
	defer file.Close()

	session, ok := r.Context().Value(SessionKey).(*domain.Session)
	if !ok || session == nil {
		h.logger.Error("Session not found or expired")
		h.renderError(w, http.StatusUnauthorized, "Unauthorized")
		return
	}

	ctx, cancel := context.WithTimeout(r.Context(), 5*time.Second)
	defer cancel()

	postID, err := pkg.GenerateUUID()
	if err != nil {
		h.logger.Error("Failed to generate UUID", "error", err)
		h.renderError(w, http.StatusInternalServerError, "Failed to generate UUID")
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
		h.logger.Error("Failed to upload image", "error", err)
		h.renderError(w, http.StatusInternalServerError, "Failed to upload image")
		return
	}
	post.ImageURL = "/images/" + objectName

	if err := h.service.CreatePost(ctx, post); err != nil {
		h.logger.Error("Failed to create post", "error", err)
		h.renderError(w, http.StatusInternalServerError, "Internal Server Error")
		return
	}

	http.Redirect(w, r, "/catalog", http.StatusSeeOther)
}

func (h *Handler) HandleAddComment(w http.ResponseWriter, r *http.Request) {
	if err := r.ParseForm(); err != nil {
		h.logger.Error("Failed to parse form", "error", err)
		h.renderError(w, http.StatusBadRequest, "Invalid form data")
		return
	}

	session, ok := r.Context().Value(SessionKey).(*domain.Session)
	if !ok || session == nil {
		h.logger.Error("Session not found or expired")
		h.renderError(w, http.StatusUnauthorized, "Unauthorized")
		return
	}

	postID := r.URL.Query().Get("id")
	parentID := r.FormValue("parent_comment_id")
	content := r.FormValue("content")

	if content == "" {
		h.logger.Error("Content is required")
		h.renderError(w, http.StatusBadRequest, "Content is required")
		return
	}

	uuid, err := pkg.GenerateUUID()
	if err != nil {
		h.logger.Error("Failed to generate UUID", "error", err)
		h.renderError(w, http.StatusInternalServerError, "Failed to generate UUID")
		return
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
			h.logger.Error("Failed to add reply", "error", err)
			h.renderError(w, http.StatusInternalServerError, "Failed to add reply: "+err.Error())
			return
		}
	} else {
		if err := h.service.AddComment(ctx, postID, comment); err != nil {
			h.logger.Error("Failed to add comment", "error", err)
			h.renderError(w, http.StatusInternalServerError, "Failed to add comment: "+err.Error())
			return
		}
	}

	http.Redirect(w, r, "/post/"+postID, http.StatusSeeOther)
}

func (h *Handler) ServeImage(w http.ResponseWriter, r *http.Request) {
	parts := strings.Split(r.URL.Path, "/")
	if len(parts) < 3 {
		h.logger.Error("Invalid image path", "path", r.URL.Path)
		h.renderError(w, http.StatusBadRequest, "Invalid image path")
		return
	}
	imageName := parts[2]

	data, contentType, err := h.imageStorage.GetImage(r.Context(), imageName)
	if err != nil {
		h.logger.Error("Failed to get image", "error", err)
		h.renderError(w, http.StatusNotFound, "Image not found")
		return
	}

	w.Header().Set("Content-Type", contentType)
	w.Write(data)
}

func (h *Handler) renderError(w http.ResponseWriter, status int, message string) {
	w.Header().Set("Content-Type", "text/html; charset=utf-8")

	var buf bytes.Buffer
	err := h.templates.ExecuteTemplate(&buf, "error.html", map[string]interface{}{
		"Code":    status,
		"Message": message,
	})
	if err != nil {
		h.logger.Error("Failed to render error template", "error", err)
		http.Error(w, "An unexpected error occurred", http.StatusInternalServerError)
		return
	}

	_, writeErr := w.Write(buf.Bytes())
	if writeErr != nil {
		h.logger.Error("Failed to write response", "error", writeErr)
		return
	}
}
