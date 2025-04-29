package postgres

const (
	qCreatePost = `INSERT INTO posts (id,title,content,image_url,author_id,created_at,upadted_at,is_archived)
		VALUES ($1,$2,$3,$4,$5,$6,$7,$8)`
	qGetByIDPost = `SELECT id,title,content,image_url,author_id,created_at,upadted_at,is_archived 
		FROM posts 
		WHERE id = $1`
	qGetAllPost = `SELECT * 
		FROM posts
		ORDER BY id`
	qUpdatePost
)
