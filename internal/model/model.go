package model

type ClaimedUser struct {
	ID    int32  `json:"id"`
	Name  string `json:"name"`
	Email string `json:"email"`
}

type PublicUser struct {
	ID     int32  `json:"id"`
	Name   string `json:"name"`
	Avatar string `json:"avatar"`
}

type Message struct {
	ID        int32      `json:"id"`
	Content   string     `json:"content"`
	CreatedAt string     `json:"createdAt"`
	User      PublicUser `json:"user"`
}

type Chat struct {
	ID      int32  `json:"id"`
	Name    string `json:"name"`
	Avatar  string `json:"avatar"`
	IsGroup bool   `json:"isGroup"`
}
