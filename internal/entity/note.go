package entity

type Note struct {
	BaseEntity
	Title     string    `db:"title" json:"title"`
	Content   string    `db:"content" json:"content"`
	UserID    int64     `db:"user_id" json:"user_id"`
	NoteAttrs NoteAttrs `db:"note_attrs" json:"note_attrs"`
}

type NoteAttrs struct {
	Color string `json:"color"`
	Icon  string `json:"icon"`
}
