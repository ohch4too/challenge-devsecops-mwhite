package api

import (
	"challenge/internal/domain"
	"challenge/internal/service"

	"github.com/apex/log"
	"github.com/gin-gonic/gin"
)

// UserHandler handles HTTP requests for user operations
type UserHandler struct {
	svc service.UserService
}

// NewUserHandler creates a new user handler
func NewUserHandler(userService service.UserService) *UserHandler {
	return &UserHandler{svc: userService}
}

// ListUsers handles GET /users - retrieves all users
func (uh *UserHandler) ListUsers(c *gin.Context) {
	var users []domain.User

	users, err := uh.svc.ListUsers()
	if err != nil {
		errMsg := "Failed to fetch users"
		log.Errorf(errMsg)
		c.JSON(500, gin.H{"error": errMsg})
		return
	}

	log.Infof("List users returned %v values\n", len(users))
	c.JSON(200, users)
}

// AddUser handles POST /users - creates a new user
func (uh *UserHandler) AddUser(c *gin.Context) {
	var user domain.User

	if err := c.ShouldBindJSON(&user); err != nil {
		c.JSON(400, gin.H{"error": err.Error()})
		log.Errorf("Failed to bind JSON: %v", err)
		return
	}

	if err := user.Validate(); err != nil {
		c.JSON(400, gin.H{"error": err.Error()})
		log.Errorf("Validation failed: %v", err)
		return
	}

	err := uh.svc.AddUser(&user)
	if err != nil {
		errMsg := "Failed to add user"
		log.Errorf("%s: %v\n", errMsg, user)
		c.JSON(500, gin.H{"error": errMsg})
		return
	}

	c.JSON(201, user)
}

// GetUser handles GET /users/:id - retrieves a user by ID
func (uh *UserHandler) GetUser(c *gin.Context) {
	id := c.Params.ByName("id")
	user, err := uh.svc.GetUser(id)
	if err != nil {
		errMsg := "User not found"
		log.Errorf("%s: %v", errMsg, id)
		c.JSON(404, gin.H{"error": errMsg})
		return
	}

	c.JSON(200, user)
}

// DelUser handles DELETE /users/:id - deletes a user by ID
func (uh *UserHandler) DelUser(c *gin.Context) {
	id := c.Params.ByName("id")
	if err := uh.svc.DeleteUser(id); err != nil {
		errMsg := "User not found"
		log.Errorf("%s: %v", errMsg, id)
		c.JSON(404, gin.H{"error": errMsg})
		return
	}

	c.JSON(204, gin.H{"message": "User deleted"})
}
