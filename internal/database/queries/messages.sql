-- name: CreateMessage :one
INSERT INTO messages (
  chat_id, 
  user_id,
  content
) VALUES (
  $1, 
  $2, 
  $3
)
RETURNING 
  id, 
  content, 
  created_at;

-- name: GetMessagesPage :many
SELECT 
  chat_id, 
  user_id, 
  content, 
  created_at
FROM 
  messages
WHERE 
  chat_id = $1
ORDER BY 
  created_at DESC
LIMIT 
  $2 
OFFSET 
  $3;
