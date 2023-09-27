package main

import (
	"accuknox/handler"
	"accuknox/model"
	"accuknox/repository"
	"accuknox/service"
	"context"
	"fmt"
	"log"
	"net/http"
	"os"
	"os/signal"
	"sync"
	"syscall"
	"time"

	"accuknox/dto"

	"github.com/gin-gonic/gin"
	"github.com/gin-gonic/gin/binding"
	"github.com/go-redis/redis"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

func main() {

	// Initialize Gin router
	router := gin.Default()
	// Database configuration
	dbHost := os.Getenv("POSTGRES_HOST")
	dbPort := "5432"
	dbName := os.Getenv("POSTGRES_DB")
	dbUser := os.Getenv("POSTGRES_USER")
	dbPassword := os.Getenv("POSTGRES_PASSWORD")
	redisHost := os.Getenv("REDIS_HOST")

	// Construct the database connection string
	dbConnectionString := fmt.Sprintf(
		"host=%s port=%s dbname=%s user=%s password=%s sslmode=disable",
		dbHost, dbPort, dbName, dbUser, dbPassword,
	)

	// Initialize database connection
	db, err := gorm.Open(postgres.Open(dbConnectionString), &gorm.Config{})
	if err != nil {
		panic("Failed to connect to the database")
	}

	db.AutoMigrate(&model.Note{}, &model.User{}, &model.UserSession{})

	rClient := redis.NewClient(&redis.Options{
		Addr: fmt.Sprintf("%v:6379", redisHost),
		DB:   0,
	})

	_, err = rClient.Ping().Result()
	if err != nil {
		log.Println(err)
		panic("Failed to connect to redis")
	}

	// Initialize repository implementations
	userRepo := repository.NewUserRepository(db, rClient)
	noteRepo := repository.NewNoteRepository(db)

	// Initialize service implementations with repositories
	userService := service.NewUserService(userRepo)
	noteService := service.NewNoteService(noteRepo)

	// Initialize handler implementations with services
	userHandler := handler.NewUserHandler(userService)
	noteHandler := handler.NewNoteHandler(noteService)

	sessionService := service.NewSessionService(userRepo)

	// Register routes using the handler implementations
	v1 := router.Group("/v1")
	{
		// Public endpoints (signup and login)
		v1.POST("/signup", userHandler.SignUpHandler)
		v1.POST("/login", userHandler.LoginHandler)

		// Notes-related endpoints that require authorization
		v1.POST("/notes", authorizeMiddleware(sessionService), noteHandler.CreateNoteHandler)
		v1.GET("/notes", authorizeMiddleware(sessionService), noteHandler.GetAllUserNotesHandler)
		v1.DELETE("/notes", authorizeMiddleware(sessionService), noteHandler.DeleteNoteHandler)
		// Add more routes as needed
	}
	// Create a context with cancellation support
	ctx, cancel := context.WithCancel(context.Background())

	// WaitGroup to wait for all goroutines to finish
	var wg sync.WaitGroup

	// Start the HTTP server in a separate goroutine
	srv := &http.Server{
		Addr:    ":8080",
		Handler: router,
	}
	wg.Add(1)
	go func() {
		defer wg.Done()
		if err := srv.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			fmt.Printf("HTTP server error: %v\n", err)
		}
	}()

	// Listen for OS signals to initiate graceful shutdown
	c := make(chan os.Signal, 1)
	signal.Notify(c, syscall.SIGINT, syscall.SIGTERM)
	<-c

	// Cancel the context to initiate shutdown
	cancel()

	// Give some time for ongoing requests to finish (adjust as needed)
	timeout := 15 * time.Second
	ctx, cancel = context.WithTimeout(context.Background(), timeout)
	defer cancel()

	// Shutdown the HTTP server
	if err := srv.Shutdown(ctx); err != nil {
		fmt.Printf("HTTP server shutdown error: %v\n", err)
	}

	// Wait for all goroutines to finish
	wg.Wait()

	fmt.Println("Server gracefully shut down")
}

// authorizeMiddleware is a custom middleware to check the session ID (SID) for authorization.
func authorizeMiddleware(sessionService service.SessionService) gin.HandlerFunc {
	return func(c *gin.Context) {
		// Inside your handler where you need to extract the SID
		var request dto.AuthRequest

		if err := c.ShouldBindBodyWith(&request, binding.JSON); err != nil {
			log.Println("[authorizeMiddleware] ", err)
			c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid request format"})
			return
		}

		sid := request.SID

		// Check if the session is valid using the SessionService
		userId, valid := sessionService.IsValidSession(sid)
		if valid {
			c.Set("userId", userId) // Store the session ID in the context for later use
			c.Next()                // Continue to the next middleware or handler
		} else {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "Unauthorized"})
			c.Abort() // Abort further processing
		}
	}
}
