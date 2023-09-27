// repository/repository.go
package repository

import (
	"accuknox/model"
)

// NoteRepository defines methods for managing notes.
type NoteRepository interface {
	CreateNote(note *model.Note) (*model.Note, error)
	GetNoteByID(id uint) (*model.Note, error)
	GetAllNotesOfUser(userID uint) ([]*model.Note, error)
	DeleteNote(userId, noteId uint) error
	// Add more note-related methods here
}

// UserRepository defines methods for user management.
type UserRepository interface {
	CreateUser(user *model.User) (*model.User, error)
	CreateSession(session *model.UserSession) (*model.UserSession, error)
	GetSessionBySID(sid string) (*model.UserSession, error)
	GetUserByEmail(email string) (*model.User, error)
	IsValidSession(sid string) (uint, bool)
	// Add more user-related methods here
}
