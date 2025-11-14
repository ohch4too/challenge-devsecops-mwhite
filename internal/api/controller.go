package api

import (
	"challenge/internal/dummy"
	"fmt"

	"github.com/apex/log"
	"github.com/gin-gonic/gin"
)

func ListUsers(c *gin.Context) {
	users, _, err := dummy.UserList()
	if err != nil {
		errMsg := "Failed to fetch users"
		log.Errorf(errMsg)
		RespondJSON(c, 500, gin.H{"error": errMsg})
	} else {
		RespondJSON(c, 200, users)
	}

}

func AddUser(c *gin.Context) {

	var user dummy.User

	if err := c.ShouldBindJSON(&user); err != nil {
		c.JSON(400, gin.H{"error": err.Error()})
		fmt.Println("Failed to bind JSON")
	} else {
		fmt.Println(user)

		err := dummy.UserAdd(&user)
		if err != nil {
			errMsg := "Failed to add user"
			log.Errorf("%s: %v\n", errMsg, user)
			RespondJSON(c, 500, gin.H{"error": errMsg})
		} else {
			RespondJSON(c, 201, user)
		}
	}

}

func DelUser(c *gin.Context) {
	id := c.Params.ByName("id")
	if err := dummy.UserDel(id); err != nil {
		errMsg := "User not found"
		log.Errorf("%s: %v", errMsg, id)
		RespondJSON(c, 404, gin.H{"error": errMsg})
	} else {
		RespondJSON(c, 204, nil)
	}
}

func GetUser(c *gin.Context) {
	var user dummy.User
	id := c.Params.ByName("id")
	if err := dummy.UserGet(&user, id); err != nil {
		errMsg := "User not found"
		log.Errorf("%s: %v", errMsg, id)
		RespondJSON(c, 404, gin.H{"error": errMsg})
	} else {
		RespondJSON(c, 200, user)
	}
}
