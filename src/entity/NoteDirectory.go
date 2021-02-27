package entity

import (
	"so-crm/src/libs/ORM"
	"time"
)

type NoteDirectory struct {
	ORM.Model
	Title    string `json:"title"`
	FkUserId uint   `json:"fk_user_id"`
}

func (u *NoteDirectory) BeforeCreate() (err error) {
	u.CreatedAt = time.Now()
	u.UpdatedAt = time.Now()
	return
}
