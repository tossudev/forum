CREATE TABLE IF NOT EXISTS users (
	id INTEGER PRIMARY KEY,
	username TEXT COLLATE NOCASE NOT NULL UNIQUE,
	email TEXT COLLATE NOCASE NOT NULL UNIQUE,
	password_hash BLOB NOT NULL,
	password_salt BLOB NOT NULL,
	avatar TEXT,
	date_created INTEGER NOT NULL,
	role_id INTEGER NOT NULL,
	signature TEXT
	-- FOREIGN KEY(role_id) REFERENCES roles(id)
);

-- CREATE TABLE IF NOT EXISTS roles (
-- 	id INTEGER PRIMARY KEY,
-- 	name TEXT NOT NULL
-- );

CREATE TABLE IF NOT EXISTS threads (
	id INTEGER PRIMARY KEY,
	title TEXT NOT NULL,
	body TEXT NOT NULL,
	date_created TEXT NOT NULL,
	author_id INTEGER NOT NULL,
	category_id INTEGER NOT NULL,
	FOREIGN KEY(author_id) REFERENCES users(id) ON DELETE CASCADE,
	FOREIGN KEY(category_id) REFERENCES categories(id) ON DELETE CASCADE
);

CREATE TABLE IF NOT EXISTS thread_likes (
	thread_id INTEGER NOT NULL,
	user_id INTEGER NOT NULL,
	thread_like BOOLEAN NOT NULL,
	FOREIGN KEY(thread_id) REFERENCES threads(id) ON DELETE CASCADE,
	FOREIGN KEY(user_id) REFERENCES users(id) ON DELETE CASCADE,
	PRIMARY KEY (thread_id, user_id)
);

/*---------------------------threads fts table and triggers--------------------------*/
CREATE VIRTUAL TABLE IF NOT EXISTS threads_fts USING fts5(
	title, 
	body, 
	date_created, 
	content='threads', 
	content_rowid='id'
	);

CREATE TRIGGER IF NOT EXISTS threads_insert AFTER 
INSERT ON threads BEGIN
INSERT INTO threads_fts (rowid, title, body, date_created)
VALUES (new.id, new.title, new.body, new.date_created); END;

CREATE TRIGGER IF NOT EXISTS threads_update AFTER 
UPDATE ON threads BEGIN
INSERT INTO threads_fts (rowid, title, body, date_created)
VALUES (new.id, new.title, new.body, new.date_created); END;

CREATE TRIGGER IF NOT EXISTS threads_delete AFTER 
DELETE ON threads BEGIN
DELETE 
FROM threads_fts
WHERE rowid = old.id; END;
/*-----------------------------------------------------------------------------------*/

CREATE TABLE IF NOT EXISTS comments (
	id INTEGER PRIMARY KEY,
	body TEXT NOT NULL,
	date_created TEXT NOT NULL,
	thread_id INTEGER NOT NULL,
	author_id INTEGER NOT NULL,
	FOREIGN KEY(thread_id) REFERENCES threads(id) ON DELETE CASCADE,
	FOREIGN KEY(author_id) REFERENCES users(id) ON DELETE CASCADE
);

CREATE TABLE IF NOT EXISTS comment_likes (
	comment_id INTEGER NOT NULL,
	user_id INTEGER NOT NULL,
	comment_like BOOLEAN NOT NULL,
	FOREIGN KEY(comment_id) REFERENCES comments(id) ON DELETE CASCADE,
	FOREIGN KEY(user_id) REFERENCES users(id) ON DELETE CASCADE,
	PRIMARY KEY (comment_id, user_id)
);

/*---------------------------comments fts table and triggers--------------------------*/

CREATE VIRTUAL TABLE IF NOT EXISTS comments_fts USING fts5(
	body, 
	date_created, 
	content='comments', 
	content_rowid='id'
	);

CREATE TRIGGER IF NOT EXISTS comments_insert AFTER 
INSERT ON comments BEGIN
INSERT INTO comments_fts(rowid, body, date_created)
VALUES (new.id, new.body, new.date_created); END;

CREATE TRIGGER IF NOT EXISTS comments_update AFTER 
UPDATE ON comments BEGIN
INSERT INTO comments_fts(rowid, body, date_created)
VALUES (new.id, new.body, new.date_created); END;

CREATE TRIGGER IF NOT EXISTS comments_delete AFTER 
DELETE ON comments BEGIN
DELETE 
FROM comments_fts
WHERE rowid = old.id; END;
/*----------------------------------------------------------------------------------*/

CREATE TABLE IF NOT EXISTS categories (
	id INTEGER PRIMARY KEY,
	name TEXT NOT NULL
);

CREATE TABLE IF NOT EXISTS images (
	id INTEGER PRIMARY KEY,
	image_path TEXT NOT NULL,
	comment_id INTEGER,
	thread_id INTEGER,
	FOREIGN KEY(comment_id) REFERENCES comments(id) ON DELETE CASCADE,
	FOREIGN KEY(thread_id) REFERENCES threads(id) ON DELETE CASCADE
);

CREATE TABLE IF NOT EXISTS sessions (
	id TEXT PRIMARY KEY,
	user_id INTEGER NOT NULL,
	session_hash BLOB NOT NULL,
	csrf_token TEXT NOT NULL,
	date_created TEXT NOT NULL,
	expires_at TEXT NOT NULL,
	FOREIGN KEY(user_id) REFERENCES users(id) ON DELETE CASCADE
);
