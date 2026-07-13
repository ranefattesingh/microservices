CREATE TABLE users (
    id          SERIAL PRIMARY KEY,
    first_name  VARCHAR(20),
    last_name   VARCHAR(20),
    email       VARCHAR(250) UNIQUE,
    phone       VARCHAR(10) UNIQUE,
    access_type VARCHAR(6),
    created_at  DATETIME,
    updated_at  DATETIME
);
