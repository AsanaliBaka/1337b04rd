-- Создание таблицы Client (переименованная User)
CREATE TABLE Client (
    user_id UUID PRIMARY KEY,
    username TEXT NOT NULL,
    image_url TEXT,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
);

-- Создание таблицы Session
CREATE TABLE Session (
    session_id TEXT PRIMARY KEY,
    user_id UUID NOT NULL,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    expires_at TIMESTAMP NOT NULL,
    FOREIGN KEY (user_id) REFERENCES Client(user_id) ON DELETE CASCADE
);

-- Создание таблицы Post
CREATE TABLE Post (
    post_id UUID PRIMARY KEY,
    title TEXT NOT NULL,
    content TEXT,
    image_url TEXT,
    is_deleted BOOLEAN DEFAULT FALSE,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    user_id UUID NOT NULL,
    FOREIGN KEY (user_id) REFERENCES Client(user_id) ON DELETE CASCADE
);

-- Создание таблицы Comment с древовидной структурой
CREATE TABLE Comment (
    comment_id UUID PRIMARY KEY,
    content TEXT NOT NULL,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    avatar TEXT,
    post_id UUID NOT NULL,
    parent_comment_id UUID ,
    user_id UUID NOT NULL,
    FOREIGN KEY (post_id) REFERENCES Post(post_id) ON DELETE CASCADE,
    FOREIGN KEY (parent_comment_id) REFERENCES Comment(comment_id) ON DELETE CASCADE,
    FOREIGN KEY (user_id) REFERENCES Client(user_id) ON DELETE CASCADE
);

-- Добавление индексов для оптимизации
CREATE INDEX idx_session_user ON Session(user_id);
CREATE INDEX idx_post_user ON Post(user_id);
CREATE INDEX idx_comment_post ON Comment(post_id);
CREATE INDEX idx_comment_parent ON Comment(parent_comment_id);
CREATE INDEX idx_comment_user ON Comment(user_id);
CREATE INDEX idx_post_created ON Post(created_at);
CREATE INDEX idx_comment_created ON Comment(created_at);
