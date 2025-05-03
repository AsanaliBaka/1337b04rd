package domain

import (
	"context"
	"fmt"
	"io"
	"log/slog"
	"time"
)

type postServer struct {
	postRepo    PostRepository
	commentRepo CommentRepository
	imageRepo   ImageStorage
	session     *SessionService
}

type PostServer interface {
	CreatePost(ctx context.Context, post *Post, imageData io.Reader) error
	CreateComment(ctx context.Context, postID string, comment *Comment) error
	GetAllPosts(ctx context.Context) ([]*Post, error)
	GetPost(ctx context.Context, postID string) (*Post, *[]*Comment, error)
}

func NewPostServer(
	postRepo PostRepository,
	commentRepo CommentRepository,
	imageRepo ImageStorage,
	sessionRepo *SessionService,
) PostServer {
	return &postServer{
		postRepo:    postRepo,
		commentRepo: commentRepo,
		imageRepo:   imageRepo,
		session:     sessionRepo,
	}
}

func (p *postServer) CreatePost(ctx context.Context, post *Post, imageData io.Reader) error {
	if imageData != nil {
		imageName := "post_" + post.ID
		imageURL, err := p.imageRepo.CreateImage(ctx, imageName, imageData, -1)
		if err != nil {
			return fmt.Errorf("failed to upload image for post %s: %w", post.ID, err)
		}
		post.ImageURL = imageURL
	}

	err := p.postRepo.CreatePost(ctx, post)
	if err != nil {
		return fmt.Errorf("failed to create post %s: %w", post.ID, err)
	}

	expirationTime := time.Now().Add(10 * time.Minute)
	p.session.AddSession(post.ID, expirationTime)
	slog.Info("post successfully created", "post_id", post.ID)

	return nil
}
func (p *postServer) GetAllPosts(ctx context.Context) ([]*Post, error) {
	posts, err := p.postRepo.GetAllPost(ctx)

	if err != nil {
		return nil, err
	}

	return posts, nil
}

func (p *postServer) CreateComment(ctx context.Context, postID string, comment *Comment, imageData io.Reader) error {
	if imageData != nil {
		imageName := "comment_" + comment.ID
		imageURL, err := p.imageRepo.CreateImage(ctx, imageName, imageData, -1)
		if err != nil {
			return fmt.Errorf("failed to upload image for comment %s: %w", comment.ID, err)
		}

		comment.ImageURL = imageURL

	}

	if err := p.commentRepo.CreateComment(ctx, comment); err != nil {
		return fmt.Errorf("failed to create comment %s: %w", comment.ID, err)
	}

	post, err := p.postRepo.GetByIDPost(ctx, postID)
	if err != nil {
		return fmt.Errorf("post not found for comment %s: %w", comment.ID, err)
	}

	post.UpdatedAt = time.Now().Add(15 * time.Minute)

	if err := p.postRepo.UpdatePost(ctx, postID, false, post.UpdatedAt); err != nil {
		return fmt.Errorf("failed to update post timestamp after comment %s: %w", comment.ID, err)
	}

	p.session.UpdateSession(postID, post.UpdatedAt)

	return nil
}


func (p *postServer) GetPost(ctx context.Context, postID string) (*Post, *[]*Comment, error) {
	post, err := p.postRepo.GetByIDPost(ctx, postID)
	if err != nil {
		return nil, nil, fmt.Errorf("post not found: %w", err)
	}
	comments, err := p.commentRepo.GetByPostIDComment(ctx, postID)
	if err != nil {
		return nil, nil, err
	}
	return post, &comments, nil
}

type SessionManager interface {
	GetSessionId(id string, ctx context.Context) (*UserRef, error)
	CreateSession(ctx context.Context, user *UserRef) error
}
type session struct {
	session SessionRepository
}

func NewSession(s SessionRepository) SessionManager {
	return &session{
		session: s,
	}
}

func (s *session) GetSessionId(id string, ctx context.Context) (*UserRef, error) {
	user, err := s.session.GetByIDSession(ctx, id)
	if err != nil {
		return nil, fmt.Errorf("failed to get session by ID %s: %w", id, err)
	}

	return user, nil
}

func (s *session) CreateSession(ctx context.Context, user *UserRef) error {
	err := s.session.CreateSession(ctx, user)
	if err != nil {
		return fmt.Errorf("failed to create session: %w", err)
	}
	return nil
}
