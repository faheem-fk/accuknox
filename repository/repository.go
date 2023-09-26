// repository/repository.go
package repository

import (
	"accuknox/model"
	"accuknox/myerrors"
	"fmt"
	"log"

	"github.com/google/uuid"
	"gorm.io/gorm"
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

type noteRepository struct {
	db *gorm.DB
}

type userRepository struct {
	db *gorm.DB
}

// NewNoteRepository creates a new NoteRepository with the given database connection.
func NewNoteRepository(db *gorm.DB) NoteRepository {
	return &noteRepository{db}
}

// NewUserRepository creates a new UserRepository with the given database connection.
func NewUserRepository(db *gorm.DB) UserRepository {
	return &userRepository{db}
}

// CreateNote creates a new note.
func (r *noteRepository) CreateNote(note *model.Note) (*model.Note, error) {
	// Use GORM's Create method to insert the note into the database
	result := r.db.Create(note)

	// Check for errors during the creation process
	if result.Error != nil {
		log.Println("[Repo:CreateNote] ", result.Error)
		return nil, myerrors.ErrInternalServer
	}

	// Return the created note
	return note, nil
}

// GetNoteByID retrieves a note by its ID.
func (r *noteRepository) GetNoteByID(id uint) (*model.Note, error) {
	// Implement the GetNoteByID method
	return nil, nil
}

// GetAllNotesOfUser retrieves all notes of a user by their UserID.
func (r *noteRepository) GetAllNotesOfUser(userID uint) ([]*model.Note, error) {
	// Create a slice to hold the retrieved notes
	var notes []*model.Note

	// Use GORM's Find method to retrieve all notes for the specified userID
	result := r.db.Where("user_id = ?", userID).Find(&notes)

	// Check for errors during the retrieval process
	if result.Error != nil {
		log.Println("[Repo:GetAllNotesOfUser] ", result.Error)
		return nil, result.Error
	}

	// Return the retrieved notes
	return notes, nil
}

// DeleteNote deletes a note by its ID for a specific user.
func (r *noteRepository) DeleteNote(userID, noteID uint) error {
	// Use GORM's Delete method to delete the note
	result := r.db.Where("id = ? AND user_id = ?", noteID, userID).Delete(&model.Note{})

	// Check for errors during the deletion process
	if result.Error != nil {
		log.Println("[Repo:DeleteNote] ", result.Error)
		return result.Error
	}

	// Check if any record was deleted (0 records deleted means note not found)
	if result.RowsAffected == 0 {
		log.Println("[Repo:DeleteNote] ", result.Error)
		return myerrors.ErrRecordNotFound
	}

	return nil // Note deleted successfully
}

// CreateUser creates a new user.
func (r *userRepository) CreateUser(user *model.User) (*model.User, error) {
	// Use GORM's Create method to insert the user into the database
	result := r.db.Create(user)

	// Check for errors during the creation process
	if result.Error != nil {
		log.Println("[Repo:CreateUser] ", result.Error)
		return nil, result.Error
	}

	// Return the created user
	return user, nil
}

// CreateSession creates a new user session.
func (r *userRepository) CreateSession(session *model.UserSession) (*model.UserSession, error) {
	// Use GORM's Create method to insert the session into the database
	sid, _ := generateSessionID()
	session.SID = sid
	result := r.db.Create(session)

	// Check for errors during the creation process
	if result.Error != nil {
		log.Println("[Repo:CreateSession] ", result.Error)
		return nil, result.Error
	}

	// Return the created session
	return session, nil
}

// GetSessionBySID retrieves a user session by its SID.
func (r *userRepository) GetSessionBySID(sid string) (*model.UserSession, error) {
	// Implement the GetSessionBySID method
	return nil, nil
}

// GetUserByEmail retrieves a user by their email.
func (r *userRepository) GetUserByEmail(email string) (*model.User, error) {
	var user model.User
	if err := r.db.Where("email = ?", email).First(&user).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			log.Println("[Repo:GetUserByEmail] ", err)
			return nil, myerrors.ErrRecordNotFound // User not found
		}

		return nil, err // Database error
	}
	return &user, nil
}

// IsValidSession checks if the session ID (SID) is valid and returns the userID if valid.
func (r *userRepository) IsValidSession(sid string) (uint, bool) {
	var session model.UserSession
	if err := r.db.Where("sid = ?", sid).First(&session).Error; err != nil {
		log.Println("[Repo:IsValidSession] ", err)
		return 0, false // Session not found or expired
	}

	// // Check if the session is expired
	// if session.CreatedAt.Before(time.Now().Add(-24 * time.Hour)) {
	// 	return 0, false // Session is expired
	// }

	return session.UserID, true // Session is valid, return the userID
}

// Rest of the UserRepository methods...

func generateSessionID() (string, error) {
	id := uuid.New()
	return fmt.Sprintf("%v", id.String()), nil
}
