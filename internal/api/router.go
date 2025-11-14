package api

import (
	"challenge/internal/repository"
	"challenge/internal/service"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

func SetupRouter(db *gorm.DB) *gin.Engine {

	r := gin.Default()

	// Initialize repository and service
	userRepo := repository.NewUserRepository(db)
	userService := service.NewUserService(userRepo)
	userHandler := NewUserHandler(userService)

	// route v1
	v1 := r.Group("/v1")
	{
		v1.GET("/users", userHandler.ListUsers)
		v1.POST("/users", userHandler.AddUser)
		v1.GET("/users/:id", userHandler.GetUser)
		v1.DELETE("/users/:id", userHandler.DelUser)
	}

	return r

}
