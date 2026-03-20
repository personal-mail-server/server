DO $$
BEGIN
    RAISE EXCEPTION '003_update_seed_login_user_password is not safely reversible; use a corrective forward migration instead';
END $$;
