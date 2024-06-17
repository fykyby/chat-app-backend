-- name: CreateChat :one
INSERT INTO chats (
  name,
  avatar,
  is_group
) VALUES (
  $1, 
  $2,
  $3
)
RETURNING 
  id, 
  name, 
  avatar,
  is_group;

-- name: CreateUsersChat :one
INSERT INTO users_chats (
  user_id,
  chat_id
) VALUES (
  $1,
  $2
)
RETURNING 
  user_id, 
  chat_id;

-- name: GetUsersChat :one
SELECT 
  user_id, 
  chat_id 
FROM 
  users_chats 
WHERE 
  user_id = $1 
AND 
  chat_id = $2;

-- name: GetUserChatList :many
SELECT 
  c.id, 
  c.name, 
  c.avatar
FROM 
  chats c
JOIN 
  users_chats uc ON c.id = uc.chat_id
WHERE 
  uc.user_id = $1;

-- name: DeleteChat :exec
DELETE FROM 
  chats 
WHERE 
  id = $1;
