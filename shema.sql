DROP TABLE IF EXISTS "accounts";

CREATE TABLE IF NOT EXISTS "accounts" (
    id SERIAL PRIMARY KEY,
    username TEXT NOT NULL,
    password TEXT NOT NULL
);