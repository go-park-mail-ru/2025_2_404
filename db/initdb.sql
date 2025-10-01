CREATE TABLE IF NOT EXISTS app_user (
	id INT GENERATED ALWAYS AS IDENTITY PRIMARY KEY,
	email TEXT UNIQUE NOT NULL CHECK (
		length(email) >= 5 AND length(email) <= 100
	),
	password TEXT NOT NULL CHECK (
		length(password) >= 8 AND length(password) <= 50
	),
	user_name TEXT NOT NULL CHECK (
		length(user_name) >= 4 AND length(user_name) <= 50
	)
);

CREATE TABLE IF NOT EXISTS ad (
	id INT GENERATED ALWAYS AS IDENTITY PRIMARY KEY,
	creator_id INTEGER REFERENCES app_user(id) ON DELETE CASCADE,
	file_path TEXT CHECK (
		length(file_path) <= 40
	),
	title TEXT NOT NULL CHECK (
		length(title) >= 1 AND length(title) <= 40
	),
	text_ad TEXT NOT NULL CHECK (
		length(text_ad) >= 1 AND length(text_ad) <= 200
	)
);

CREATE TABLE IF NOT EXISTS session (
	id INT GENERATED ALWAYS AS IDENTITY PRIMARY KEY,
	user_id INTEGER REFERENCES app_user(id) ON DELETE CASCADE,
	session_id TEXT UNIQUE NOT NULL CHECK (
		length(session_id) >= 1 AND length(session_id) <= 64
	)
);