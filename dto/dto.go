package dto

// Define a SignUpRequest struct to represent the JSON request body
type SignUpRequest struct {
	Name     string `json:"name" binding:"required"`
	Email    string `json:"email" binding:"required,email"`
	Password string `json:"password" binding:"required"`
}

type LoginRequest struct {
	Email    string `json:"email" binding:"required,email"`
	Password string `json:"password" binding:"required"`
}

type LoginResponse struct {
	SID string `json:"sid"`
}

type AuthRequest struct {
	SID string `json:"sid"`
}

// NotesRequest defines the JSON request format for the /notes endpoint.
type NotesRequest struct {
	SID string `json:"sid"`
}

type DeleteNoteRequest struct {
	SID string `json:"sid"`
	ID  uint32 `json:"id"`
}

type CreateNoteRequest struct {
	SID  string `json:"sid"`
	Note string `json:"note"`
}
