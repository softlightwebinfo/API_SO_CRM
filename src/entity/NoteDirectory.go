package entity

import (
	"so-crm/src/libs/ORM"
	"time"
)

type NoteDirectories struct {
	ORM.Model
	Title    string `json:"title"`
	FkUserId int64  `json:"fk_user_id"`
}

func (u *NoteDirectories) BeforeCreate() (err error) {
	u.CreatedAt = time.Now()
	u.UpdatedAt = time.Now()

	return
}
