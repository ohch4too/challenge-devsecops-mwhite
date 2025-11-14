package api

import (
	"fmt"
	"os"

	"github.com/gin-gonic/gin"
)

type ResponseData struct {
	Status int
	Data   interface{}
}

func RespondJSON(w *gin.Context, status int, payload interface{}) {
	fmt.Println("status ", status)
	var res ResponseData

	res.Status = status
	res.Data = payload

	w.JSON(status, res)
}

func Start() {
	router := setupRouter()

	certFile := os.Getenv("TLS_CERT_FILE")
	keyFile := os.Getenv("TLS_KEY_FILE")

	if certFile != "" && keyFile != "" {
		router.RunTLS(":10000", certFile, keyFile)
	} else {
		router.Run(":10000")
	}
}
