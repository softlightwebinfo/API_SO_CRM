package entity

import (
	"so-crm/src/libs/ORM"
	"time"
)

type User struct {
	ORM.Model
	Email    string `json:"email"`
	Password string `json:"password"`
}

func (u *User) BeforeCreate() (err error) {
	u.CreatedAt = time.Now()
	u.UpdatedAt = time.Now()
	return
}
