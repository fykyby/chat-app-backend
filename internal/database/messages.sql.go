// Code generated by sqlc. DO NOT EDIT.
// versions:
//   sqlc v1.26.0
// source: messages.sql

package database

import (
	"context"

	"github.com/jackc/pgx/v5/pgtype"
)

const createMessage = `-- name: CreateMessage :one
INSERT INTO
  messages (chat_id, user_id, content)
VALUES
  ($1, $2, $3)
RETURNING
  id,
  content,
  created_at
`

type CreateMessageParams struct {
	ChatID  int32
	UserID  int32
	Content string
}

type CreateMessageRow struct {
	ID        int32
	Content   string
	CreatedAt pgtype.Timestamp
}

func (q *Queries) CreateMessage(ctx context.Context, arg CreateMessageParams) (CreateMessageRow, error) {
	row := q.db.QueryRow(ctx, createMessage, arg.ChatID, arg.UserID, arg.Content)
	var i CreateMessageRow
	err := row.Scan(&i.ID, &i.Content, &i.CreatedAt)
	return i, err
}

const getMessages = `-- name: GetMessages :many
SELECT
  m.id,
  m.content,
  m.created_at,
  u.id AS user_id,
  u.name AS user_name,
  u.avatar AS user_avatar
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
  $3
`

type GetMessagesParams struct {
	ChatID int32
	Limit  int32
	Offset int32
}

type GetMessagesRow struct {
	ID         int32
	Content    string
	CreatedAt  pgtype.Timestamp
	UserID     int32
	UserName   string
	UserAvatar string
}

func (q *Queries) GetMessages(ctx context.Context, arg GetMessagesParams) ([]GetMessagesRow, error) {
	rows, err := q.db.Query(ctx, getMessages, arg.ChatID, arg.Limit, arg.Offset)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	var items []GetMessagesRow
	for rows.Next() {
		var i GetMessagesRow
		if err := rows.Scan(
			&i.ID,
			&i.Content,
			&i.CreatedAt,
			&i.UserID,
			&i.UserName,
			&i.UserAvatar,
		); err != nil {
			return nil, err
		}
		items = append(items, i)
	}
	if err := rows.Err(); err != nil {
		return nil, err
	}
	return items, nil
}
