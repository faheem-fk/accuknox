// service/service.go
package service

import (
	"accuknox/model"
	"accuknox/myerrors"
	"accuknox/repository"

	"golang.org/x/crypto/bcrypt"
)

// NoteService provides methods for managing notes.
type NoteService interface {
	CreateNote(note *model.Note) (*model.Note, error)
	GetNoteByID(id uint) (*model.Note, error)
	GetAllNotesOfUser(userID uint) ([]*model.Note, error)
	DeleteNote(userId, noteId uint) error
	// Add more note-related methods here
}

// UserService provides methods for user management.
type UserService interface {
	CreateUser(name, email, password string) (string, error)
	Login(email, password string) (string, error) // Add the Login method
	// Add more user-related methods here
}

// SessionService defines methods for session management.
type SessionService interface {
	IsValidSession(sid string) (uint, bool)
}

type noteService struct {
	noteRepo repository.NoteRepository
}

// userService struct
type userService struct {
	userRepo repository.UserRepository
}

// SessionServiceImpl implements SessionService.
type SessionServiceImpl struct {
	userRepo repository.UserRepository
}

// NewNoteService creates a new NoteService with the provided NoteRepository.
func NewNoteService(noteRepo repository.NoteRepository) NoteService {
	return &noteService{noteRepo}
}

// NewUserService creates a new UserService with the provided UserRepository.
func NewUserService(userRepo repository.UserRepository) UserService {
	return &userService{userRepo}
}

// NewSessionService creates a new SessionService with the provided UserRepository.
func NewSessionService(userRepo repository.UserRepository) SessionService {
	return &SessionServiceImpl{userRepo}
}

// CreateNote creates a new note.
func (s *noteService) CreateNote(note *model.Note) (*model.Note, error) {
	// Implement the CreateNote method using the noteRepo
	return s.noteRepo.CreateNote(note)
}

// GetNoteByID retrieves a note by its ID.
func (s *noteService) GetNoteByID(id uint) (*model.Note, error) {
	// Implement the GetNoteByID method using the noteRepo
	return s.noteRepo.GetNoteByID(id)
}

// GetAllNotesOfUser retrieves all notes of a user by their UserID.
func (s *noteService) GetAllNotesOfUser(userID uint) ([]*model.Note, error) {
	// Implement the GetAllNotesOfUser method using the noteRepo
	return s.noteRepo.GetAllNotesOfUser(userID)
}

// DeleteNote deletes a note by its ID.
func (s *noteService) DeleteNote(userId, noteId uint) error {
	// Implement the DeleteNote method using the noteRepo
	return s.noteRepo.DeleteNote(userId, noteId)
}

// ...

// CreateUser creates a new user and returns the created user object.
func (s *userService) CreateUser(name, email, password string) (string, error) {
	// Generate a password hash for the provided password
	hashedPassword, err := generatePasswordHash(password)
	if err != nil {
		return "", err
	}
	// Create a User struct with the request data
	user := model.User{
		Name:         name,
		Email:        email,
		PasswordHash: hashedPassword,
	}

	// Call the repository method to create the user
	s.userRepo.CreateUser(&user)
	//On success, create new uinque session
	newSession := &model.UserSession{UserID: user.ID}
	newSession, err = s.userRepo.CreateSession(newSession)
	if err != nil {
		return "", err
	}
	return newSession.SID, nil
}

// Login authenticates a user with their email and password.
func (s *userService) Login(email, password string) (string, error) {
	// Implement the Login method using the userRepo
	user, err := s.userRepo.GetUserByEmail(email)
	if err != nil {
		return "", err
	}

	// Check if the provided password matches the user's stored password
	if !checkPasswordHash(password, user.PasswordHash) {
		return "", myerrors.ErrAuthentication // Custom authentication error
	}

	//On success, create new uinque session
	newSession := &model.UserSession{UserID: user.ID}
	newSession, err = s.userRepo.CreateSession(newSession)
	if err != nil {
		return "", err
	}

	return newSession.SID, nil
}

// IsValidSession checks if the session ID (SID) is valid.
func (s *SessionServiceImpl) IsValidSession(sid string) (uint, bool) {
	// Delegate the session validation to the UserRepository or your session store
	// You should implement your session validation logic here
	return s.userRepo.IsValidSession(sid)
}

// Password hash checking function
func checkPasswordHash(password, hash string) bool {
	err := bcrypt.CompareHashAndPassword([]byte(hash), []byte(password))
	return err == nil
}

// generatePasswordHash generates a password hash for the given password.
func generatePasswordHash(password string) (string, error) {
	bytes, err := bcrypt.GenerateFromPassword([]byte(password), 14)
	return string(bytes), err
}
