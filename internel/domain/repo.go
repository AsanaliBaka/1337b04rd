package domain

type PostRepository interface {
	CreatePost(post *Post) error
	GetByIDPost(id string) (*Post, error)
	GetAllPost() ([]*Post, error)
	UpdatePost(post *Post) error
	DeletePost(id string) error
}

type CommentRepository interface {
	CreateComment(comment *Comment) error
	GetByPostIDComment(postID string) ([]*Comment, error)
}
