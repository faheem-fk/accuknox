// model/model.go
package model

import "github.com/jinzhu/gorm"

// Note represents a note in the application.
type Note struct {
	ID      uint   `json:"id, omitempty"`
	UserID  uint   `json:"-"`
	Content string `json:"note"`
}

// User represents a user in the application.
type User struct {
	ID           uint   `json:"-"`
	Name         string `json:"name"`
	Email        string `json:"email" gorm:"uniqueIndex"`
	PasswordHash string `json:"-"`
}

// UserSession represents a user session with a unique session ID (sid).
type UserSession struct {
	gorm.Model
	UserID uint   `json:"user_id"`
	SID    string `json:"sid" redis:"column:sid"`
}
