package handlers

import (
	"1337b04rd/internal/domain"
	"encoding/json"
	"fmt"
	"net/http"
)

type PostHandler struct {
	post domain.PostServer
}

func NewPostHandler(p domain.PostServer) *PostHandler {
	return &PostHandler{
		post: p,
	}
}

func (p *PostHandler) CreatePostHandler(w http.ResponseWriter, r *http.Request) {
	user, ok := r.Context().Value("user").(*domain.UserRef)

	if !ok || user == nil {
		http.Error(w, "User not authenticated", http.StatusUnauthorized)
		return
	}

	err := r.ParseMultipartForm(10 << 20)

	if err != nil {
		http.Error(w, "failed to parse form", http.StatusBadRequest)
		return
	}

	title := r.FormValue("title")
	content := r.FormValue("content")
	file, _, err := r.FormFile("image")

	if err != nil && err != http.ErrMissingFile {
		http.Error(w, "failed to read image", http.StatusBadRequest)
		return
	}

	defer func() {
		if file != nil {
			file.Close()
		}
	}()

	post := domain.NewPost(title, content, *user)

	err = p.post.CreatePost(r.Context(), post, file)

	if err != nil {
		http.Error(w, "failed to create post", http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(post)
}

func (p *PostHandler) GetAllPosts(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	posts, err := p.post.GetAllPosts(ctx)

	if err != nil {
		http.Error(w, "failed to get posts", http.StatusInternalServerError)
		return
	}

	if err := json.NewEncoder(w).Encode(posts); err != nil {
		http.Error(w, "failed to encode posts", http.StatusInternalServerError)
		return
	}
}

func (p *PostHandler) GetById(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	id := r.PathValue("id")
	if id == "" {
		http.Error(w, "missing post id", http.StatusBadRequest)
		return
	}

	post, comments, err := p.post.GetPost(ctx, id)

	if err != nil {
		http.Error(w, fmt.Sprintf("failed to retrieve post: %v", err), http.StatusInternalServerError)
		return
	}

	response := struct {
		Post     *domain.Post       `json:"post"`
		Comments *[]*domain.Comment `json:"comments"`
	}{
		Post:     post,
		Comments: comments,
	}

	w.WriteHeader(http.StatusOK)

	if err := json.NewEncoder(w).Encode(response); err != nil {
		http.Error(w, fmt.Sprintf("failed to encode response: %v", err), http.StatusInternalServerError)
		return
	}

}

func (p *PostHandler) CreateComment(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	postID := r.PathValue("id")
	if postID == "" {
		http.Error(w, "missing post id", http.StatusBadRequest)
		return
	}

	userRef, ok := ctx.Value("user").(*domain.UserRef)
	if !ok || userRef == nil {
		http.Error(w, "user not found in context", http.StatusUnauthorized)
		return
	}

	var req struct {
		Text string `json:"text"`
	}
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "invalid request body", http.StatusBadRequest)
		return
	}

	comment := domain.NewComments(req.Text, *userRef)

	if err := p.post.CreateComment(ctx, postID, comment); err != nil {
		http.Error(w, "failed to create comment", http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(comment)

}
