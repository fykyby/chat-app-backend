-- name: CreateMessage :one
INSERT INTO
  messages (chat_id, user_id, content)
VALUES
  ($1, $2, $3)
RETURNING
  id,
  content,
  created_at;


-- name: GetMessages :many
SELECT
  m.id,
  m.content,
  m.created_at,
  u.name,
  u.avatar
FROM
  messages m
  JOIN users u ON u.id = user_id
WHERE
  chat_id = $1
ORDER BY
  created_at DESC
LIMIT
  $2
OFFSET
  $3;