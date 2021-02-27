package entity

import (
	"so-crm/src/libs/ORM"
	"time"
)

type Note struct {
	ORM.Model
	Title       string `json:"title"`
	Description string `json:"description"`
	FkDirectory uint   `json:"fk_directory"`
}

func (u *Note) BeforeCreate() (err error) {
	u.CreatedAt = time.Now()
	u.UpdatedAt = time.Now()
	return
}
