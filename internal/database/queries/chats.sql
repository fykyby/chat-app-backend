-- name: CreateChat :one
INSERT INTO
  chats (name, avatar, is_group)
VALUES
  ($1, $2, $3)
RETURNING
  id,
  name,
  avatar,
  is_group;


-- name: GetUserChats :many
SELECT
  c.id,
  c.name,
  c.avatar,
  c.is_group
FROM
  chats c
  JOIN users_chats uc ON c.id = uc.chat_id
WHERE
  uc.user_id = $1;


-- name: DeleteChat :exec
DELETE FROM
  chats
WHERE
  id = $1;


-- name: GetChatOfTwoUsers :one
SELECT
  c.id,
  c.name,
  c.avatar,
  c.is_group
FROM
  users_chats uc1
  JOIN users_chats uc2 ON uc1.chat_id = uc2.chat_id
  JOIN chats c ON uc1.chat_id = c.id
WHERE
  uc1.user_id = $1
  AND uc2.user_id = $2;