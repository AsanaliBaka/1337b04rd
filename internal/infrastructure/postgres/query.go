package postgres

const (
	qCreatePost = `INSERT INTO posts (id,title,content,image_url,author_id,created_at,updated_at,is_archived)
		VALUES ($1,$2,$3,$4,$5,$6,$7,$8)`
	qGetByIDPost = `SELECT id,title,content,image_url,author_id,created_at,updated_at,is_archived 
		FROM posts 
		WHERE id = $1`
	qGetAllPost = `SELECT * 
		FROM posts
		ORDER BY id`
	qUpdatePost = `UPDATE posts SET updated_at = $1, is_archived = $2 WHERE id = $3 `
	qDeletePost = `DELETE FROM posts where id=$1`
)

const (
	qCreateSession  = `INSERT INTO users (session_id,avatar_url,name) VALUES($1,$2,$3)`
	qGetByIDSession = `SELECT session_id, avatar_url, name FROM users WHERE session_id = $1`
	qDeleteSession  = `DELETE FROM users WHERE session_id = $1`
)

const (
	qCreateComment      = `INSERT INTO comments (id,texts, author_id,post_id,image_url,created_at) VALUES ($1,$2,$3,$4,$5,$6)`
	qGetByPostIDComment = `SELECT id, texts, author_id, post_id, image_url, created_at FROM comments WHERE post_id = $1`
)
