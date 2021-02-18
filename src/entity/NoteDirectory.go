package entity

type NoteDirectory struct {
	ID        uint   `json:"id" gorm:"primary_key"`
	Title     string `json:"title"`
	CreatedAt string `json:"created_at"`
	UpdatedAt string `json:"updated_at"`
}
