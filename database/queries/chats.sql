-- name: CreateChat :one
INSERT INTO chats (
  name
) VALUES (
  $1
)
RETURNING id, name;
