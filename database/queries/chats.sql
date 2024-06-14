-- name: CreateChat :one
INSERT INTO chats (
  name
) VALUES (
  $1
)
RETURNING id, name;

-- name: CreateUsersChat :one
INSERT INTO users_chats (
  user_id,
  chat_id
) VALUES (
  $1,
  $2
)
RETURNING user_id, chat_id;

-- name: GetUsersChat :one
SELECT user_id, chat_id FROM users_chats WHERE user_id = $1 AND chat_id = $2;

-- name: DeleteChat :exec
DELETE FROM chats WHERE id = $1;
