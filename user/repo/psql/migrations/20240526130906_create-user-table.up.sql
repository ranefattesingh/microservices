CREATE TABLE users (
    id          UUID        PRIMARY KEY,
    name        VARCHAR(50) NOT NULL,
    email       VARCHAR(50) NOT NULL,
    password    VARCHAR(20) NOT NULL,
    is_admin    BOOLEAN     NOT NULL,
    create_date DATE        DEFAULT NOW(),
    update_date DATE        NOT NULL
);