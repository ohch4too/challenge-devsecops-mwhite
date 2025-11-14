package api

import "github.com/gin-gonic/gin"

func setupRouter() *gin.Engine {

	r := gin.Default()

	// route v1
	v1 := r.Group("/v1")
	{
		v1.GET("/users", ListUsers)
		v1.POST("/users", AddUser)
		v1.GET("/users/:id", GetUser)
		v1.DELETE("/users/:id", DelUser)
	}

	return r

}
