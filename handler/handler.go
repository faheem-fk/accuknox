package handler

import (
	"accuknox/dto"
	"accuknox/model"
	"accuknox/myerrors"
	"accuknox/service"
	"log"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/gin-gonic/gin/binding"
)

// UserServiceHandler defines methods for user-related handlers.
type UserServiceHandler interface {
	SignUpHandler(c *gin.Context)
	LoginHandler(c *gin.Context)
	// Add more user-related handlers here
}

// NoteServiceHandler defines methods for note-related handlers.
type NoteServiceHandler interface {
	CreateNoteHandler(c *gin.Context)
	GetAllUserNotesHandler(c *gin.Context)
	DeleteNoteHandler(c *gin.Context)
	// Add more note-related handlers here
}

// userHandler implements UserServiceHandler.
type userHandler struct {
	userService service.UserService
}

// NewUserHandler creates a new userHandler with the provided UserService.
func NewUserHandler(userService service.UserService) UserServiceHandler {
	return &userHandler{userService}
}

func (h *userHandler) SignUpHandler(c *gin.Context) {
	var req dto.SignUpRequest

	// Bind the request body to the SignUpRequest struct
	if err := c.ShouldBindJSON(&req); err != nil {
		log.Println("[LoginHandler] ", err)
		c.JSON(http.StatusBadRequest, gin.H{"error": "Bad Request"})
		return
	}

	// Call the UserService to create the user
	sessionID, err := h.userService.CreateUser(req.Name, req.Email, req.Password)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to create user"})
		return
	}

	// Respond with the created user
	c.JSON(http.StatusOK, gin.H{"sid": sessionID})
}

func (h *userHandler) LoginHandler(c *gin.Context) {
	var req dto.LoginRequest

	// Bind the request body to the LoginRequest struct
	if err := c.ShouldBindJSON(&req); err != nil {
		log.Println("[LoginHandler] ", err)
		c.JSON(http.StatusBadRequest, gin.H{"error": "Bad Request"})
		return
	}

	// Call the UserService's Login method to authenticate the user and obtain the session ID (SID)
	sessionID, err := h.userService.Login(req.Email, req.Password)
	if err != nil {
		if err == myerrors.ErrAuthentication {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "Authentication failed"})
		} else {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to log in"})
		}
		return
	}

	// Respond with the session ID (SID)
	c.JSON(http.StatusOK, gin.H{"sid": sessionID})
}

// Add more user-related handlers here

// noteHandler implements NoteServiceHandler.
type noteHandler struct {
	noteService service.NoteService
}

// NewNoteHandler creates a new noteHandler with the provided NoteService.
func NewNoteHandler(noteService service.NoteService) NoteServiceHandler {
	return &noteHandler{noteService}
}

func (h *noteHandler) CreateNoteHandler(c *gin.Context) {
	// Extract the user ID from the context
	userID, _ := c.Get("userId")

	var req dto.CreateNoteRequest

	// Bind the request body to the CreateNoteRequest struct
	if err := c.ShouldBindBodyWith(&req, binding.JSON); err != nil {
		log.Println("[CreateNoteHandler] ", err)
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// Create a new Note model based on the request data
	newNote := &model.Note{
		UserID:  userID.(uint),
		Content: req.Note,
	}

	// Call the NoteService to create the note
	createdNote, err := h.noteService.CreateNote(newNote)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to create note"})
		return
	}

	// Respond with the created note
	c.JSON(http.StatusOK, gin.H{"note": createdNote})
}

func (h *noteHandler) GetAllUserNotesHandler(c *gin.Context) {
	// Extract the user ID from the context
	userID, _ := c.Get("userId")

	// Call the NoteService to get all notes for the user
	notes, err := h.noteService.GetAllNotesOfUser(userID.(uint))
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to get user notes"})
		return
	}

	// Respond with the user's notes
	c.JSON(http.StatusOK, gin.H{"notes": notes})
}

func (h *noteHandler) DeleteNoteHandler(c *gin.Context) {
	userId, _ := c.Get("userId")

	var requestBody dto.DeleteNoteRequest

	// Bind the request body to the CreateNoteRequest struct
	if err := c.ShouldBindBodyWith(&requestBody, binding.JSON); err != nil {
		log.Println("[DeleteNoteHandler] ", err)
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// Delete the note associated with the provided ID if it belongs to the authenticated user
	err := h.noteService.DeleteNote(userId.(uint), uint(requestBody.ID))
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to delete note"})
		return
	}

	// Respond with a success message if the note was deleted successfully
	c.JSON(http.StatusOK, gin.H{"message": "Note deleted successfully"})
}
