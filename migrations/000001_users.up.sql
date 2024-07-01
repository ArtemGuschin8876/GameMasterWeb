CREATE TABLE IF NOT EXISTS users (
    id bigserial PRIMARY KEY,
    name TEXT NOT NULL, 
    nickname TEXT NOT NULL,
    email TEXT NOT NULL, 
    city TEXT NOT NULL, 
    about TEXT NOT NULL,
    image TEXT NOT NULL
);