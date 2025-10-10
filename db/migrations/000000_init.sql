
-- Создание таблицы client
CREATE TABLE IF NOT EXISTS client (
	id INT GENERATED ALWAYS AS IDENTITY PRIMARY KEY,
    name TEXT UNIQUE NOT NULL CHECK (
		length(name) >= 4 AND length(name) <= 50
	),
	email TEXT UNIQUE NOT NULL CHECK (
		length(email) >= 5 AND length(email) <= 100
	),
	password_hash TEXT NOT NULL CHECK (
		length(password_hash) >= 8 AND length(password_hash) <= 120
	),
    created_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP
);

CREATE TABLE IF NOT EXISTS client_wallet (
	id INT GENERATED ALWAYS AS IDENTITY PRIMARY KEY,
	client_id INT UNIQUE REFERENCES client(id) ON DELETE CASCADE,
	balance NUMERIC(12, 2) NOT NULL DEFAULT 0 CHECK (balance >= 0),
	created_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
	updated_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP
);

CREATE TABLE IF NOT EXISTS wallet_top_up (
	id INT GENERATED ALWAYS AS IDENTITY PRIMARY KEY,
	client_wallet_id INT REFERENCES client_wallet(id) ON DELETE CASCADE,
	amount NUMERIC(12, 2) NOT NULL CHECK (amount > 0),
    payment_method TEXT NOT NULL CHECK (
		length(payment_method) >= 1 AND length(payment_method) <= 40
	),
    status TEXT NOT NULL CHECK (
		length(status) >= 1 AND length(status) <= 40
	),
	created_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP
);

CREATE TABLE IF NOT EXISTS notification_user (
	id INT GENERATED ALWAYS AS IDENTITY PRIMARY KEY,
	client_id INTEGER REFERENCES client(id) ON DELETE CASCADE,
	notification_text TEXT NOT NULL CHECK (
		length(notification_text) >= 1 AND length(notification_text) <= 200
	),
    type TEXT NOT NULL CHECK (
		length(type) >= 1 AND length(type) <= 40
	),
	created_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP
);

CREATE TABLE IF NOT EXISTS platform (
	id INT GENERATED ALWAYS AS IDENTITY PRIMARY KEY,
	platform_name TEXT UNIQUE NOT NULL CHECK (
		length(platform_name) >= 1 AND length(platform_name) <= 100
	),
	created_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP
);

CREATE TABLE IF NOT EXISTS ad (
	id INT GENERATED ALWAYS AS IDENTITY PRIMARY KEY,
	client_id INT REFERENCES client(id) ON DELETE CASCADE,
	title TEXT NOT NULL CHECK (
		length(title) >= 1 AND length(title) <= 40
	),
	content TEXT NOT NULL CHECK (
		length(content) >= 1 AND length(content) <= 200
	),
    img_bin BYTEA,
    target_url TEXT NOT NULL CHECK (
		length(target_url) >= 1 AND length(target_url) <= 200
	),
	created_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
	updated_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP
);

CREATE TABLE IF NOT EXISTS ad_detail (
	id INT GENERATED ALWAYS AS IDENTITY PRIMARY KEY,
	ad_id INT REFERENCES ad(id) ON DELETE CASCADE,
	platform_id INT REFERENCES platform(id) ON DELETE CASCADE,
	amount_for_ad NUMERIC(12, 2) NOT NULL CHECK (amount_for_ad > 0),
    status TEXT NOT NULL CHECK (
		length(status) >= 1 AND length(status) <= 40
	),
	created_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
	updated_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP
);

CREATE TABLE IF NOT EXISTS statistic (
	id INT GENERATED ALWAYS AS IDENTITY PRIMARY KEY,
	ad_detail_id INT REFERENCES ad_detail(id) ON DELETE CASCADE,
	clicks INT NOT NULL DEFAULT 0 CHECK (clicks >= 0),
	impressions INT NOT NULL DEFAULT 0 CHECK (impressions >= 0),
	created_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
	updated_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP
);

CREATE TABLE IF NOT EXISTS session (
	id INT GENERATED ALWAYS AS IDENTITY PRIMARY KEY,
	user_id INTEGER REFERENCES client(id) ON DELETE CASCADE,
	session_id TEXT UNIQUE NOT NULL CHECK (
		length(session_id) >= 1 AND length(session_id) <= 64
	),
	created_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
	updated_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP
);