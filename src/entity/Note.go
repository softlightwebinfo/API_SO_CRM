package entity

import (
	"so-crm/src/libs/ORM"
	"time"
)

type Notes struct {
	ORM.Model
	Title       string `json:"title"`
	Description string `json:"description"`
	FkDirectory int64  `json:"fk_directory"`
}

func (u *Notes) BeforeCreate() (err error) {
	u.CreatedAt = time.Now()
	u.UpdatedAt = time.Now()
	return
}

func (u *Notes) BeforeDelete() (err error) {
	newTime := time.Now()
	u.DeletedAt = &newTime
	u.UpdatedAt = time.Now()
	return
}
