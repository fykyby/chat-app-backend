-- name: CreateUser :one
INSERT INTO users (
  email, 
  name,
  password,
  avatar
) VALUES (
  $1, $2, $3, $4
)
RETURNING id, name, email, avatar;

-- name: GetUserByEmail :one
SELECT * FROM users WHERE email = $1;
