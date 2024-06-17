// Code generated by sqlc. DO NOT EDIT.
// versions:
//   sqlc v1.26.0
// source: chats.sql

package database

import (
	"context"
)

const createChat = `-- name: CreateChat :one
INSERT INTO chats (
  name,
  is_group
) VALUES (
  $1, 
  $2
)
RETURNING id, name, is_group
`

type CreateChatParams struct {
	Name    string
	IsGroup bool
}

func (q *Queries) CreateChat(ctx context.Context, arg CreateChatParams) (Chat, error) {
	row := q.db.QueryRow(ctx, createChat, arg.Name, arg.IsGroup)
	var i Chat
	err := row.Scan(&i.ID, &i.Name, &i.IsGroup)
	return i, err
}

const createUsersChat = `-- name: CreateUsersChat :one
INSERT INTO users_chats (
  user_id,
  chat_id
) VALUES (
  $1,
  $2
)
RETURNING user_id, chat_id
`

type CreateUsersChatParams struct {
	UserID int32
	ChatID int32
}

func (q *Queries) CreateUsersChat(ctx context.Context, arg CreateUsersChatParams) (UsersChat, error) {
	row := q.db.QueryRow(ctx, createUsersChat, arg.UserID, arg.ChatID)
	var i UsersChat
	err := row.Scan(&i.UserID, &i.ChatID)
	return i, err
}

const deleteChat = `-- name: DeleteChat :exec
DELETE FROM chats WHERE id = $1
`

func (q *Queries) DeleteChat(ctx context.Context, id int32) error {
	_, err := q.db.Exec(ctx, deleteChat, id)
	return err
}

const getUsersChat = `-- name: GetUsersChat :one
SELECT user_id, chat_id FROM users_chats WHERE user_id = $1 AND chat_id = $2
`

type GetUsersChatParams struct {
	UserID int32
	ChatID int32
}

func (q *Queries) GetUsersChat(ctx context.Context, arg GetUsersChatParams) (UsersChat, error) {
	row := q.db.QueryRow(ctx, getUsersChat, arg.UserID, arg.ChatID)
	var i UsersChat
	err := row.Scan(&i.UserID, &i.ChatID)
	return i, err
}
