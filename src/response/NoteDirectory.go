package response

import "so-crm/src/entity"

type NoteDirectory struct {
	entity.NoteDirectories
	Notes uint `json:"notes"`
}
