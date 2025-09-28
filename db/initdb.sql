CREATE TABLE IF NOT EXISTS users (
	id INT GENERATED ALWAYS AS IDENTITY PRIMARY KEY,
	email TEXT UNIQUE NOT NULL CHECK (
		length(email) BETWEEN 5 AND 100 AND
		email LIKE '%@%.%'
	),
	password TEXT NOT NULL CHECK (
		length(password) BETWEEN 8 AND 50
	),
	user_name TEXT NOT NULL CHECK (
		length(user_name) BETWEEN 4 AND 50
	)
);

CREATE TABLE IF NOT EXISTS ads (
	id INT GENERATED ALWAYS AS IDENTITY PRIMARY KEY,
	creator_id INTEGER REFERENCES users(id),
	file_path TEXT NOT NULL CHECK (
		length(file_path) <= 40
	),
	title TEXT NOT NULL CHECK (
		length(title) BETWEEN 1 AND 40
	),
	text TEXT NOT NULL CHECK (
		length(text) BETWEEN 1 AND 200
	)
);

CREATE TABLE IF NOT EXISTS session (
	id INT GENERATED ALWAYS AS IDENTITY PRIMARY KEY,
	user_id INTEGER REFERENCES users(id),
	session_id TEXT UNIQUE NOT NULL CHECK (
		length(session_id) <= 64
	)
);