// Code generated by sqlc. DO NOT EDIT.
// versions:
//   sqlc v1.26.0

package database

import (
	"github.com/jackc/pgx/v5/pgtype"
)

type Chat struct {
	ID      int32
	Name    string
	IsGroup bool
}

type Message struct {
	ID        int32
	ChatID    int32
	UserID    int32
	Content   string
	CreatedAt pgtype.Timestamp
}

type User struct {
	ID       int32
	Email    string
	Name     string
	Password string
	Avatar   pgtype.Text
}

type UsersChat struct {
	UserID int32
	ChatID int32
}