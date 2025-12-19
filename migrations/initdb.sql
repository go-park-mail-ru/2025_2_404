CREATE TABLE IF NOT EXISTS client (
	id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    name TEXT NOT NULL CHECK (
		length(name) >= 4 AND length(name) <= 50
	),
	email TEXT UNIQUE NOT NULL CHECK (
		length(email) >= 5 AND length(email) <= 100
	),
	password_hash TEXT NOT NULL CHECK (
		length(password_hash) <= 120
	),
	img_path VARCHAR(120),
	user_first_name VARCHAR(20),
	user_second_name VARCHAR(20),
	company VARCHAR(80),
	phone_number VARCHAR(30),
    created_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP
);

CREATE TABLE IF NOT EXISTS client_wallet (
	id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
	client_id UUID UNIQUE REFERENCES client(id) ON DELETE CASCADE,
	balance INT NOT NULL DEFAULT 0 CHECK (balance >= 0),
	created_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
	updated_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP
);

CREATE TABLE IF NOT EXISTS wallet_top_up (
	id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
	client_wallet_id UUID REFERENCES client_wallet(id) ON DELETE CASCADE,
	yoo_payment_id TEXT UNIQUE NOT NULL CHECK (
		length(yoo_payment_id) >= 1 AND length(yoo_payment_id) <= 100
	),
	amount INT NOT NULL CHECK (amount > 0),
    payment_method TEXT NOT NULL CHECK (
		length(payment_method) >= 1 AND length(payment_method) <= 40
	),
    status TEXT NOT NULL CHECK (
		length(status) >= 1 AND length(status) <= 40
	),
	created_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
	updated_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP
);

CREATE TABLE IF NOT EXISTS notification_user (
	id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
	client_id UUID REFERENCES client(id) ON DELETE CASCADE,
	notification_text TEXT NOT NULL CHECK (
		length(notification_text) >= 1 AND length(notification_text) <= 200
	),
    type TEXT NOT NULL CHECK (
		length(type) >= 1 AND length(type) <= 40
	),
	created_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP
);

CREATE TABLE IF NOT EXISTS ad (
	id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
	client_id UUID REFERENCES client(id) ON DELETE CASCADE,
	title TEXT NOT NULL CHECK (
		length(title) >= 1 AND length(title) <= 40
	),
	content TEXT NOT NULL CHECK (
		length(content) >= 1 AND length(content) <= 200
	),
    img_path TEXT CHECK(
		length(img_path) <= 100
	),
    target_url TEXT NOT NULL CHECK (
		length(target_url) >= 1 AND length(target_url) <= 200
	),
	created_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
	updated_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP
);

CREATE TABLE IF NOT EXISTS slots (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    user_id UUID NOT NULL REFERENCES client(id) ON DELETE CASCADE,
    slot_name VARCHAR(100) NOT NULL CHECK (length(slot_name) BETWEEN 1 AND 100),
    min_cost_adv INT NOT NULL CHECK (min_cost_adv >= 0),
    format_of_banner VARCHAR(20) NOT NULL CHECK (format_of_banner IN ('horizontal', 'vertical', 'square')),
    status VARCHAR(10) NOT NULL DEFAULT 'active' CHECK (status IN ('active', 'paused')),
    back_color VARCHAR(7) NOT NULL DEFAULT '#ffffff' CHECK (back_color ~ '^#[0-9A-Fa-f]{6}$'),
    text_color VARCHAR(7) NOT NULL DEFAULT '#000000' CHECK (text_color ~ '^#[0-9A-Fa-f]{6}$'),
    created_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP
);

CREATE TABLE IF NOT EXISTS ad_detail (
	id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
	ad_id UUID REFERENCES ad(id) ON DELETE CASCADE,
	budget INT NOT NULL CHECK (budget >= 0),
    status TEXT NOT NULL CHECK (
		status IN ('active', 'non-active')
	),
	created_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
	updated_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
	start_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
	end_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP
);

CREATE TABLE IF NOT EXISTS statistic (
	id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
	ad_detail_id UUID REFERENCES ad_detail(id) ON DELETE CASCADE,
	clicks INT NOT NULL DEFAULT 0 CHECK (clicks >= 0),
	impressions INT NOT NULL DEFAULT 0 CHECK (impressions >= 0),
	created_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
	updated_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP
);;

CREATE TABLE IF NOT EXISTS slot_event (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    slot_id UUID NOT NULL REFERENCES slots(id) ON DELETE CASCADE,
    ad_detail_id UUID NOT NULL REFERENCES ad_detail(id) ON DELETE CASCADE,
    event_type TEXT NOT NULL CHECK (event_type IN ('impression', 'click')),
    created_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP
);

CREATE INDEX IF NOT EXISTS idx_slot_event_slot_id_created_at 
ON slot_event (slot_id, created_at);

CREATE INDEX IF NOT EXISTS idx_slot_event_created_at 
ON slot_event (created_at);

CREATE INDEX IF NOT EXISTS idx_slot_event_slot_id_event_type_created_at 
ON slot_event (slot_id, event_type, created_at);

