CREATE TABLE users (
                       id SERIAL PRIMARY KEY,
                       username TEXT UNIQUE NOT NULL,
                       password TEXT NOT NULL

);

CREATE TABLE notes (
                       id SERIAL PRIMARY KEY,
                       user_id INTEGER NOT NULL,
                       title TEXT NOT NULL,
                       body TEXT NOT NULL,
                       FOREIGN KEY (user_id) REFERENCES users(id)
);