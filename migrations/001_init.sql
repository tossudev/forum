CREATE TABLE IF NOT EXISTS users (
	id INTEGER PRIMARY KEY,
	username TEXT COLLATE NOCASE NOT NULL UNIQUE,
	password TEXT NOT NULL,
	email TEXT NOT NULL UNIQUE,
	avatar TEXT,
	date_created TEXT NOT NULL, 
	role_id INTEGER NOT NULL,
	signature TEXT,
	FOREIGN KEY(role_id) REFERENCES roles(id)
);

CREATE TABLE IF NOT EXISTS roles (
	id INTEGER PRIMARY KEY,
	name TEXT NOT NULL
);

CREATE TABLE IF NOT EXISTS threads (
	id INTEGER PRIMARY KEY,
	title TEXT NOT NULL,
	body TEXT NOT NULL,
	date_created TEXT NOT NULL,
	author_id INTEGER NOT NULL,	
	category_id INTEGER NOT NULL,
	FOREIGN KEY(author_id) REFERENCES users(id),
	FOREIGN KEY(category_id) REFERENCES categories(id)
);

CREATE TABLE IF NOT EXISTS comments (
	id INTEGER PRIMARY KEY,
	body TEXT NOT NULL,
	date_created TEXT NOT NULL,
	thread_id INTEGER NOT NULL,
	author_id INTEGER NOT NULL,
	FOREIGN KEY(thread_id) REFERENCES threads(id),
	FOREIGN KEY(author_id) REFERENCES users(id)
);

CREATE TABLE IF NOT EXISTS thread_likes (
	thread_id INTEGER NOT NULL,
	user_id INTEGER NOT NULL,
	thread_like BOOLEAN NOT NULL, 
	FOREIGN KEY(thread_id) REFERENCES threads(id),
	FOREIGN KEY(user_id) REFERENCES users(id),
	PRIMARY KEY (thread_id, user_id)
);

CREATE TABLE IF NOT EXISTS comment_likes (
	comment_id INTEGER NOT NULL,
	user_id INTEGER NOT NULL,
	thread_like BOOLEAN NOT NULL,
	FOREIGN KEY(comment_id) REFERENCES comments(id),
	FOREIGN KEY(user_id) REFERENCES users(id),
	PRIMARY KEY (comment_id, user_id)
);

CREATE TABLE IF NOT EXISTS categories (
	id INTEGER PRIMARY KEY,
	name TEXT NOT NULL
);
 
CREATE TABLE IF NOT EXISTS images (
	id INTEGER PRIMARY KEY,
	image_path TEXT NOT NULL,
	comment_id INTEGER,
	thread_id INTEGER, 
	FOREIGN KEY(comment_id) REFERENCES comments(id),
	FOREIGN KEY(thread_id) REFERENCES threads(id)
);

CREATE TABLE IF NOT EXISTS sessions (
	id INTEGER PRIMARY KEY,
	session_id TEXT NOT NULL,
	user_id INTEGER NOT NULL,
	FOREIGN KEY(user_id) REFERENCES users(id)
);









