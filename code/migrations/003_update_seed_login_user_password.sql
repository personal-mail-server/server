UPDATE users
SET password_hash = '$2a$10$YaMUTHZOWQeF/1GTkaDIHuJ3s7F0pz5zsH.kx2.RFDSzTluLaxici', updated_at = NOW()
WHERE login_id = 'user-01';
