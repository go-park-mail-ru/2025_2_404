-- Функция для автоматического обновления updated_at
CREATE OR REPLACE FUNCTION update_updated_at_column()
RETURNS TRIGGER AS $$
BEGIN
    NEW.updated_at = CURRENT_TIMESTAMP;
    RETURN NEW;
END;
$$ LANGUAGE plpgsql;

-- Функция для создания кошелька при добавлении нового клиента
CREATE OR REPLACE FUNCTION create_wallet_for_new_client()
RETURNS TRIGGER AS $$
BEGIN
    INSERT INTO client_wallet (client_id) VALUES (NEW.id);
    RETURN NEW;
END;
$$ LANGUAGE plpgsql;