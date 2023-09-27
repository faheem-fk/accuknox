package repository

import (
	"accuknox/model"
	"accuknox/myerrors"
	"fmt"
	"log"
	"strconv"
	"time"

	"github.com/go-redis/redis"
	"github.com/google/uuid"
	"gorm.io/gorm"
)

type userRepository struct {
	db      *gorm.DB
	rClient *redis.Client
}

// NewUserRepository creates a new UserRepository with the given database connection.
func NewUserRepository(db *gorm.DB, rClient *redis.Client) UserRepository {
	return &userRepository{db, rClient}
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
	sid, err := generateSessionID()
	if err != nil {
		log.Println("[Repo:GetUserByEmail] ", err)
		return nil, myerrors.ErrInternalServer
	}
	session.SID = sid

	exp := time.Duration(24 * time.Hour) // 24 Hours
	r.rClient.Set(sid, session.UserID, exp).Err()

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
	val, err := r.rClient.Get(sid).Result()
	if err != nil {
		return 0, false
	}

	userid, err := strconv.Atoi(val)
	if err != nil {
		return 0, false
	}

	return uint(userid), true
}

// Rest of the UserRepository methods...

func generateSessionID() (string, error) {
	id := uuid.New()
	return fmt.Sprintf("%v", id.String()), nil
}
