INSERT INTO users (login_id, password_hash)
VALUES ('user-01', '$2a$10$YaMUTHZOWQeF/1GTkaDIHuJ3s7F0pz5zsH.kx2.RFDSzTluLaxici')
ON CONFLICT (login_id) DO NOTHING;
