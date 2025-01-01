package models

import (
	"time"
)

// Note represents a note in the system
type Note struct {
	Id        int       `json:"id"`
	Title     string    `json:"title"`
	Content   string    `json:"content"`
	UserId    int       `json:"user_id"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}

// NewNote creates a new Note instance
func NewNote(title string, content string, userId int) Note {
	return Note{
		Title:   title,
		Content: content,
		UserId:  userId,
	}
}
