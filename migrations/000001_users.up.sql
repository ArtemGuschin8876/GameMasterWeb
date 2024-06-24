CREATE TABLE IF NOT EXISTS users (
    id bigserial PRIMARY KEY,
    name TEXT NOT NULL, 
    nickname TEXT NOT NULL UNIQUE,
    email TEXT NOT NULL UNIQUE, 
    city TEXT NOT NULL, 
    about TEXT NOT NULL,
    image TEXT NOT NULL
);

ALTER TABLE users
ADD CONSTRAINT unique_email UNIQUE (email),
ADD CONSTRAINT unique_nickname UNIQUE (nickname);
