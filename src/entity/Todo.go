package entity

import (
	"so-crm/src/libs/ORM"
	"time"
)

type Todo struct {
	ORM.Model
	Title       string   `json:"title"`
	Description string   `json:"description"`
	Success     bool     `json:"success"`
	Priority    int      `json:"priority"`
	Tags        []string `json:"tags"`
	FkUserId    int64    `json:"fk_user_id"`
}

func (u *Todo) BeforeCreate() (err error) {
	u.CreatedAt = time.Now()
	u.UpdatedAt = time.Now()
	return
}

func (u *Todo) BeforeDelete() (err error) {
	newTime := time.Now()
	u.DeletedAt = &newTime
	u.UpdatedAt = time.Now()
	return
}
