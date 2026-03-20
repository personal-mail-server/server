DO $$
BEGIN
    IF NOT EXISTS (
        SELECT 1
        FROM users
        WHERE login_id = 'user-01'
    ) THEN
        RETURN;
    END IF;

    IF EXISTS (
        SELECT 1
        FROM users
        WHERE login_id = 'user-01'
          AND password_hash = '$2a$10$YaMUTHZOWQeF/1GTkaDIHuJ3s7F0pz5zsH.kx2.RFDSzTluLaxici'
          AND failed_attempts = 0
          AND locked_until IS NULL
    ) THEN
        DELETE FROM users
        WHERE login_id = 'user-01';
        RETURN;
    END IF;

    RAISE EXCEPTION 'cannot safely rollback 002_seed_login_user: user-01 state diverged from seed';
END $$;
