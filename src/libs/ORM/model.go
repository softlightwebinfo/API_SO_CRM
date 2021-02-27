package ORM

import "time"

type Model struct {
	ID        uint       `json:"id"`
	CreatedAt time.Time  `json:"created_at"`
	UpdatedAt time.Time  `json:"updated_at"`
	DeletedAt *time.Time `sql:"index" json:"deleted_at"`
}

func (m *Model) BeforeCreate() (err error) {
	m.CreatedAt = time.Now()
	m.UpdatedAt = time.Now()
	return nil
}

type IModel interface {
	BeforeCreate() (err error)
}
