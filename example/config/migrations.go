package config

const UserTable = `
CREATE TABLE IF NOT EXISTS users (
	id INTEGER PRIMARY KEY AUTOINCREMENT,
	username TEXT NOT NULL UNIQUE,
	age INTEGER DEFAULT 0,
	gender TEXT,
	firstname TEXT,
	lastname TEXT,
	email TEXT NOT NULL UNIQUE,
	password TEXT NOT NULL,
	provider TEXT NOT NULL
);
`

const PostTable = `
CREATE TABLE IF NOT EXISTS posts (
	id INTEGER PRIMARY KEY AUTOINCREMENT,
	title TEXT NOT NULL,
	content TEXT NOT NULL,
	created DATETIME DEFAULT CURRENT_TIMESTAMP,
	user_id INTEGER,
	foreign key (user_id) REFERENCES users (id)
);		
`

const CategoryTable = `
CREATE TABLE IF NOT EXISTS categories (
	id INTEGER PRIMARY KEY AUTOINCREMENT,
	name TEXT NOT NULL,
	name_slug TEXT NOT NULL
);
`

const PostCategoryTable = `
CREATE TABLE IF NOT EXISTS posts_category (
	id INTEGER PRIMARY KEY AUTOINCREMENT,
	post_id INTEGER,
	category_id INTEGER,
	foreign key (post_id) REFERENCES posts (id),
	foreign key (category_id) REFERENCES categories (id)
);
`

const PostRatingTable = `
CREATE TABLE IF NOT EXISTS posts_rating (
	id INTEGER PRIMARY KEY AUTOINCREMENT,
	post_id INTEGER,
	user_id INTEGER,
	rating INTEGER,
	foreign key (post_id) REFERENCES posts (id),
	foreign key (user_id) REFERENCES users (id)
);
`

const PostRepliesTable = `
CREATE TABLE IF NOT EXISTS posts_replies (
	id INTEGER PRIMARY KEY AUTOINCREMENT,
	post_id INTEGER,
	user_id INTEGER,
	content TEXT,
	created DATETIME DEFAULT CURRENT_TIMESTAMP,
	foreign key (post_id) REFERENCES posts (id),
	foreign key (user_id) REFERENCES users (id)
);
`
const PostRepliesRatingTable = `
CREATE TABLE IF NOT EXISTS posts_replies_rating (
	id INTEGER PRIMARY KEY AUTOINCREMENT,
	post_reply_id INTEGER,
	user_id INTEGER,
	rating INTEGER,
	foreign key (post_reply_id) REFERENCES posts_replies (id),
	foreign key (user_id) REFERENCES users (id)
);
`

const SessionTable = `
CREATE TABLE IF NOT EXISTS sessions (
	id INTEGER PRIMARY KEY AUTOINCREMENT,
	user_id INTEGER,
	name TEXT,
	value TEXT,
	expiration DATETIME,
	foreign key (user_id) REFERENCES users (id)
);
`
