DO $$
DECLARE
    test_email TEXT := 'test@example.com';
    test_client_id UUID;
    wallet_id UUID;
    ad_ids UUID[];
    slot_ids UUID[];
    ad_detail_id UUID;
    i INT;
    j INT;
    start_date TIMESTAMP := '2025-12-01 00:00:00';
    end_date TIMESTAMP := '2026-01-01 00:00:00';
    curr_date DATE;
BEGIN
    -- Проверяем, существует ли уже тестовый пользователь
    SELECT id INTO test_client_id FROM client WHERE email = test_email;
    IF test_client_id IS NOT NULL THEN
        RAISE NOTICE 'Тестовый пользователь уже существует. Пропускаем создание.';
        RETURN;
    END IF;

    -- Создаём клиента
    INSERT INTO client (name, email, password_hash)
    VALUES ('Test User', test_email, '$2a$10$tR9Au8E0W3yFcDsuGTpVjeesvF5yRf7Fwvg8IN7Kl3k/ESwgW6s8e')
    RETURNING id INTO test_client_id;

    -- Кошелёк
    INSERT INTO client_wallet (client_id, balance)
    VALUES (test_client_id, 10000)
    RETURNING id INTO wallet_id;

    -- 10 объявлений
    FOR i IN 1..10 LOOP
        INSERT INTO ad (client_id, title, content, target_url)
        VALUES (
            test_client_id,
            'Ad Title ' || i,
            'Ad Content ' || i,
            'https://example.com/ad' || i
        )
        RETURNING id INTO ad_ids[i];
    END LOOP;

    -- 2 слота
    FOR i IN 1..2 LOOP
        INSERT INTO slots (user_id, slot_name, min_cost_adv, format_of_banner)
        VALUES (
            test_client_id,
            'Slot ' || i,
            100,
            CASE WHEN i % 2 = 1 THEN 'horizontal' ELSE 'vertical' END
        )
        RETURNING id INTO slot_ids[i];
    END LOOP;

    -- ad_detail + slot_event (упрощённо: без детальных событий, только суммарная статистика)
    FOR i IN 1..10 LOOP
        INSERT INTO ad_detail (ad_id, budget, status, start_at, end_at)
        VALUES (ad_ids[i], 1000, 'active', start_date, end_date)
        RETURNING id INTO ad_detail_id;

        -- Простая статистика: 3100 показов, 310 кликов за месяц
        INSERT INTO statistic (ad_detail_id, clicks, impressions)
        VALUES (ad_detail_id, 310, 3100);
    END LOOP;

    RAISE NOTICE '✅ Тестовые данные успешно созданы для клиента %', test_client_id;
END $$;