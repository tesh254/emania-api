package main

import (
	"net/http"
	"log"

	"github.com/tesh254/emania-api/controllers"
	"github.com/gin-gonic/gin"
	"github.com/joho/godotenv"
)

func init() {
	if err := godotenv.Load(); err != nil {
		log.Print("No .env file found")
	}
}

func main() {
	router := gin.Default()

	v1 := router.Group("/api/v1") 
	{
		// User related endpoints
		users := new(controllers.UserController)
		v1.GET("/users", users.Find)
		v1.POST("/signup", users.Signup)
		v1.POST("/login", users.Login)
		v1.PUT("/email-verify", users.Verify)
		v1.POST("/password-reset", users.PasswordRequest)
		v1.PUT("/password-reset-submit", users.PasswordResetSubmit)
	}

	router.NoRoute(func(c *gin.Context) {
		c.String(http.StatusNotFound, "Not Found")
	})

	router.Run(":5000")
}