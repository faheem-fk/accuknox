package repository

import (
	"accuknox/model"
	"accuknox/myerrors"
	"log"

	"gorm.io/gorm"
)

type noteRepository struct {
	db *gorm.DB
}

// NewNoteRepository creates a new NoteRepository with the given database connection.
func NewNoteRepository(db *gorm.DB) NoteRepository {
	return &noteRepository{db}
}

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
