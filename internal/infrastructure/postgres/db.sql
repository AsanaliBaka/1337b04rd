CREATE TABLE users(
    session_id VARCHAR (50) PRIMARY KEY,
    avatar_url TEXT,
    name VARCHAR (100) NOT NULL
)

CREATE TABLE posts(
    id VARCHAR (50) PRIMARY KEY , 
    title TEXT NOT NULL, 
    content TEXT NOT NULL,
    image_url TEXT NOT NULL,
    author_id VARCHAR (36) NOT NULL REFERENCES users(session_id),
    created_at TIMESTAMP WITH TIME ZONE NOT NULL,
    updated_at TIMESTAMP WITH TIME ZONE ,
    is_archived BOOLEAN DEFAULT FALSE

)

CREATE TABLE comments (
    id VARCHAR(50) PRIMARY KEY, 
    texts TEXT, 
    author_id VARCHAR(100) NOT NULL REFERENCES users(session_id),
    post_id VARCHAR(100) NOT NULL REFERENCES posts (id),
    image_url TEXT, 
    created_at TIMESTAMP WITH TIME ZONE NOT NULL
)