package request

type NoteUpdate struct {
	Column    string
	Value     string
	Directory int64
	Note      int64
}
