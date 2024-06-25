CREATE TABLE IF NOT EXISTS users (
    id bigserial PRIMARY KEY,
    name TEXT NOT NULL, 
    nickname TEXT NOT NULL UNIQUE,
    email TEXT NOT NULL UNIQUE, 
    city TEXT NOT NULL, 
    about TEXT NOT NULL,
    image TEXT NOT NULL
);

CREATE UNIQUE INDEX ON users (email);
CREATE UNIQUE INDEX ON users (nickname);
