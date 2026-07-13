CREATE TABLE users (
    id           BIGSERIAL PRIMARY KEY,
    first_name   VARCHAR(20) NOT NULL,
    last_name    VARCHAR(20) NOT NULL,
    email        VARCHAR(250) UNIQUE NOT NULL,
    phone        VARCHAR(10) UNIQUE,
    access_type  VARCHAR(20) NOT NULL,
    created_at   TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    updated_at   TIMESTAMPTZ NOT NULL DEFAULT NOW()
);
