CREATE TABLE
  users (
    id SERIAL PRIMARY KEY,
    email VARCHAR(255) NOT NULL UNIQUE,
    name VARCHAR(255) NOT NULL UNIQUE,
    avatar VARCHAR(255) NOT NULL,
    password VARCHAR(255) NOT NULL
  );


CREATE TABLE
  chats (
    id SERIAL PRIMARY KEY,
    name VARCHAR(255),
    avatar VARCHAR(255),
    is_group BOOLEAN NOT NULL
  );


CREATE TABLE
  users_chats (
    user_id INTEGER NOT NULL REFERENCES users(id),
    chat_id INTEGER NOT NULL REFERENCES chats(id),
    PRIMARY KEY (user_id, chat_id)
  );


CREATE TABLE
  messages (
    id SERIAL PRIMARY KEY,
    content TEXT NOT NULL,
    created_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
    chat_id INTEGER NOT NULL REFERENCES chats(id) ON DELETE CASCADE,
    user_id INTEGER NOT NULL REFERENCES users(id)
  );