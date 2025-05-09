package db

import (
	"context"
	"database/sql"
	"errors"

	"1337b04rd/internal/domain"
)

type Repo struct {
	Conn *sql.DB
}

func NewRepo(db *sql.DB) *Repo {
	return &Repo{Conn: db}
}

//  PostRepository --------------------

func (r *Repo) ListCatalog(ctx context.Context) ([]*domain.PostSummary, error) {
	rows, err := r.Conn.QueryContext(ctx, `
		SELECT 
			p.post_id, 
			p.title, 
			p.image_url, 
			p.created_at, 
			c.username
		FROM 
			Post p
		JOIN 
			Client c ON p.user_id = c.user_id
		WHERE 
			p.is_deleted = FALSE
		ORDER BY 
			p.created_at DESC;
	`)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var posts []*domain.PostSummary
	for rows.Next() {
		var post domain.PostSummary
		if err := rows.Scan(&post.ID, &post.Title, &post.ImageURL, &post.CreatedAt, &post.Author); err != nil {
			return nil, err
		}
		posts = append(posts, &post)
	}
	return posts, rows.Err()
}

func (r *Repo) GetPostByID(ctx context.Context, id string) (*domain.Post, error) {
	row := r.Conn.QueryRowContext(ctx, `
		SELECT p.post_id, p.title, p.content, p.image_url, p.created_at, u.username, u.user_id
		FROM Post p
		JOIN Client u ON p.user_id = u.user_id
		WHERE p.post_id = $1
	`, id)

	var post domain.Post
	if err := row.Scan(&post.ID, &post.Title, &post.Content, &post.ImageURL, &post.CreatedAt, &post.Author, &post.AuthorID); err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, err
		}
		return nil, err
	}

	comments, err := r.getCommentsByPostID(ctx, post.ID)
	if err != nil {
		return nil, err
	}
	post.Comments = comments

	return &post, nil
}

func (r *Repo) CreatePost(ctx context.Context, post *domain.Post) error {
	_, err := r.Conn.ExecContext(ctx, `
		INSERT INTO Post (post_id, title, content, image_url, user_id)
		VALUES ($1, $2, $3, $4, $5)
	`, post.ID, post.Title, post.Content, post.ImageURL, post.Author) // Author = user_id
	return err
}

func (r *Repo) GetPosts(ctx context.Context) ([]domain.Post, error) {
	rows, err := r.Conn.QueryContext(ctx, `
		SELECT 
    p.post_id, 
    p.title, 
    p.content, 
    p.image_url, 
    p.created_at, 
    u.username 
FROM 
    Post p
JOIN 
    Client u ON p.user_id = u.user_id
WHERE 
    p.is_deleted = FALSE
ORDER BY 
    p.created_at DESC;
	`)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var posts []domain.Post
	for rows.Next() {
		var post domain.Post
		if err := rows.Scan(&post.ID, &post.Title, &post.Content, &post.ImageURL, &post.CreatedAt, &post.Author); err != nil {
			return nil, err
		}
		comments, err := r.getCommentsByPostID(ctx, post.ID)
		if err != nil {
			return nil, err
		}
		post.Comments = comments
		posts = append(posts, post)
	}
	return posts, rows.Err()
}

//  ArchiveRepository --------------------

func (r *Repo) ListArchiveCatalog(ctx context.Context) ([]*domain.PostSummary, error) {
	rows, err := r.Conn.QueryContext(ctx, `
		SELECT 
			p.post_id, 
			p.title, 
			p.image_url, 
			p.created_at, 
			c.username
		FROM 
			Post p
		JOIN 
			Client c ON p.user_id = c.user_id
		WHERE 
			p.is_deleted = TRUE
		ORDER BY 
			p.created_at DESC;
	`)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var posts []*domain.PostSummary
	for rows.Next() {
		var post domain.PostSummary
		if err := rows.Scan(&post.ID, &post.Title, &post.ImageURL, &post.CreatedAt, &post.Author); err != nil {
			return nil, err
		}
		posts = append(posts, &post)
	}
	return posts, rows.Err()
}

func (r *Repo) GetArchivedPostByID(ctx context.Context, id string) (*domain.Post, error) {
	row := r.Conn.QueryRowContext(ctx, `
		SELECT p.post_id, p.title, p.content, p.image_url, p.created_at, u.username, u.user_id
		FROM Post p
		JOIN Client u ON p.user_id = u.user_id
		WHERE p.post_id = $1 AND is_deleted = TRUE
	`, id)

	var post domain.Post
	if err := row.Scan(&post.ID, &post.Title, &post.Content, &post.ImageURL, &post.CreatedAt, &post.Author, &post.AuthorID); err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, domain.ErrNotFound
		}
		return nil, err
	}

	comments, err := r.getCommentsByPostID(ctx, post.ID)
	if err != nil {
		return nil, err
	}
	post.Comments = comments

	return &post, nil
}

func (r *Repo) ArchivePostByID(ctx context.Context, id string) (*domain.Post, error) {
	_, err := r.Conn.ExecContext(ctx, `
		UPDATE Post SET is_deleted = TRUE WHERE post_id = $1
	`, id)
	if err != nil {
		return nil, err
	}
	return r.GetArchivedPostByID(ctx, id)
}

// CommentRepository --------------------

func (r *Repo) AddComment(ctx context.Context, postID string, comment *domain.Comment) error {
	_, err := r.Conn.ExecContext(ctx, `
		INSERT INTO Comment (comment_id, content,avatar, post_id, user_id)
		 VALUES ($1, $2, $3, $4, $5)
	`, comment.ID, comment.Content, comment.AvaterLink, postID, comment.Author)
	return err
}

func (r *Repo) ReplyToComment(ctx context.Context, postID string, parentID string, comment *domain.Comment) error {
	_, err := r.Conn.ExecContext(ctx, `
		INSERT INTO Comment (content, avatar, post_id, parent_comment_id, user_id)
		VALUES ($1, $2, $3, $4, $5)
	`, comment.Content, comment.AvaterLink, postID, parentID, comment.Author)
	return err
}

// UserRepository --------------------

func (r *Repo) CreateUser(ctx context.Context, user *domain.User) error {
	_, err := r.Conn.ExecContext(ctx, `
		INSERT INTO Client (user_id, username, image_url, created_at)
		VALUES ($1, $2, $3, $4)
	`, user.ID, user.Username, user.ImageURL, user.CreatedAt)
	return err
}

func (r *Repo) GetUserByID(ctx context.Context, userID string) (*domain.User, error) {
	row := r.Conn.QueryRowContext(ctx, `
		SELECT user_id, username, image_url
		FROM Client
		WHERE user_id = $1
	`, userID)

	var user domain.User
	if err := row.Scan(&user.ID, &user.Username, &user.ImageURL); err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, domain.ErrNotFound
		}
		return nil, err
	}
	return &user, nil
}

func (r *Repo) GetMaxCharacterID(ctx context.Context) (int, error) {
	row := r.Conn.QueryRowContext(ctx, `SELECT COUNT(*) FROM Client`)
	var count int
	if err := row.Scan(&count); err != nil {
		return 0, err
	}
	return count, nil
}

// SessionRepository --------------------

func (r *Repo) GetSession(ctx context.Context, sessionID string) (*domain.Session, error) {
	row := r.Conn.QueryRowContext(ctx, `
		SELECT session_id, user_id, expires_at
		FROM Session
		WHERE session_id = $1
	`, sessionID)

	var session domain.Session
	if err := row.Scan(&session.ID, &session.UserID, &session.ExpiresAt); err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, domain.ErrNotFound
		}
		return nil, err
	}
	return &session, nil
}

func (r *Repo) SaveSession(ctx context.Context, session *domain.Session) error {
	_, err := r.Conn.ExecContext(ctx, `
		INSERT INTO Session (session_id, user_id, expires_at)
		VALUES ($1, $2, $3)
	`, session.ID, session.UserID, session.ExpiresAt)
	return err
}

// Вспомогательная --------------------

func (r *Repo) getCommentsByPostID(ctx context.Context, postID string) ([]domain.Comment, error) {
	rows, err := r.Conn.QueryContext(ctx, `
		SELECT c.comment_id, c.content, c.created_at, u.username, c.avatar
		FROM Comment c
		JOIN Client u ON c.user_id = u.user_id
		WHERE c.post_id = $1
		ORDER BY c.created_at ASC
	`, postID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var comments []domain.Comment
	for rows.Next() {
		var c domain.Comment
		if err := rows.Scan(&c.ID, &c.Content, &c.CreatedAt, &c.Author, &c.AvaterLink); err != nil {
			return nil, err
		}
		comments = append(comments, c)
	}
	return comments, rows.Err()
}
