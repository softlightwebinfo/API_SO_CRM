package response

import "so-crm/src/entity"

type NoteDirectory struct {
	entity.NoteDirectory
	Notes uint `json:"notes"`
}
