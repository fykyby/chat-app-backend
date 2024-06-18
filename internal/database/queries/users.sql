-- name: CreateUser :one
INSERT INTO
  users (email, name, password, avatar)
VALUES
  ($1, $2, $3, $4)
RETURNING
  id,
  name,
  email,
  avatar;


-- name: GetUserByEmail :one
SELECT
  *
FROM
  users
WHERE
  email = $1;


-- name: GetPublicUser :one
SELECT
  id,
  name,
  avatar
FROM
  users
WHERE
  id = $1;


-- name: SearchPublicUsers :many
SELECT
  id,
  name,
  avatar
FROM
  users
WHERE
  name ILIKE $1
  AND id != $2
LIMIT
  $3
OFFSET
  $4;