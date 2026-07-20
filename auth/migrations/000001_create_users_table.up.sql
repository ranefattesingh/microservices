CREATE TABLE auth (
    id              BIGSERIAL PRIMARY KEY,
    email           VARCHAR(250) UNIQUE NOT NULL,
    password_hash   VARCHAR(250) NOT NULL,
    created_at      TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    updated_at      TIMESTAMPTZ NOT NULL DEFAULT NOW()
);
