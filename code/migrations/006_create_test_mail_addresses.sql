CREATE TABLE IF NOT EXISTS test_mail_addresses (
    id BIGSERIAL PRIMARY KEY,
    owner_user_id BIGINT NOT NULL REFERENCES users(id) ON DELETE CASCADE,
    email TEXT NOT NULL UNIQUE,
    created_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    updated_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    deleted_at TIMESTAMPTZ NULL
);

CREATE INDEX IF NOT EXISTS idx_test_mail_addresses_owner_user_id ON test_mail_addresses (owner_user_id);
CREATE INDEX IF NOT EXISTS idx_test_mail_addresses_owner_user_id_deleted_at ON test_mail_addresses (owner_user_id, deleted_at);
