package controllers

import (
	"os"

	"github.com/gin-gonic/gin"
	"github.com/tesh254/emania-api/forms"
	"github.com/tesh254/emania-api/helpers"
	"github.com/tesh254/emania-api/models"
	"github.com/tesh254/emania-api/services"
)

var userModel = new(models.UserModel)

var blacklistModel = new(models.BlacklistModel)

// SignupRequestBody defines the req.body of signup endpoint
type SignupRequestBody struct {
	name     string `form:"name" bson:"name"`
	email    string `form:"email" bson:"emali"`
	password string `form:"password" bson:"password"`
	role     string `form:"role" bson:"role"`
}

// UserController defines the controller's structure
type UserController struct{}

// Find controller handle fetching of users
func (user *UserController) Find(c *gin.Context) {
	list, err := userModel.Find()

	if err != nil {
		c.JSON(404, gin.H{"message": "Fetching users error", "error": err.Error()})
		c.Abort()
	} else {
		c.JSON(200, gin.H{"data": list})
	}
}

// Signup controller handles registering a user
func (user *UserController) Signup(c *gin.Context) {
	var data forms.SignupUserCommand

	if c.BindJSON(&data) != nil {
		c.JSON(406, gin.H{"message": "Details not filled", "form": data})
		c.Abort()
		return
	}

	result, _ := userModel.GetUserByEmailCount(data.Email)

	if result.Email != "" {
		c.JSON(403, gin.H{"message": "Email is already in use"})
		c.Abort()
		return
	}

	err := userModel.Signup(data)

	if err != nil {
		c.JSON(400, gin.H{"message": "Problem Creating user account"})
		c.Abort()
		return
	}

	// Handle email verification sendout
	verificationToken, err := services.GenerateToken(data.Email)

	link := os.Getenv("BASE_URL") + "/email-verify?token=" + verificationToken

	subject := "Account Verification"
	body := "Thank you for creating a Packit account, click this link to verify <a href='" + link + "'>here</a>."

	go services.SendMail(data.Email, data.Name, subject, body)

	if err != nil {
		c.JSON(406, gin.H{"message": "Account could not be created", "error": err.Error()})
		c.Abort()
		return
	}

	c.JSON(201, gin.H{"message": "Account has been created"})
}

// Login controller handles logging in a user
func (user *UserController) Login(c *gin.Context) {
	var data forms.LoginUserCommand

	if c.BindJSON(&data) != nil {
		c.JSON(406, gin.H{"message": "Details not filled", "data": data})
		c.Abort()
		return
	}

	result, err := userModel.Login(data.Email)

	if result.Email == "" {
		c.JSON(404, gin.H{"message": "Invalid login details"})
		c.Abort()
		return
	}

	if result.IsVerified == false {
		c.JSON(403, gin.H{"message": "Account has not been verified"})
		c.Abort()
		return
	}

	if err != nil {
		c.JSON(400, gin.H{"message": "Problem logging in to the account"})
		c.Abort()
		return
	}

	hashedPassword := []byte(result.Password)
	password := []byte(data.Password)

	err = helpers.PasswordCompare(password, hashedPassword)

	if err != nil {
		c.JSON(403, gin.H{"message": "Invalid login credentials"})
		c.Abort()
		return
	}

	jwtToken, err2 := services.GenerateToken(data.Email)

	if err2 != nil {
		c.JSON(403, gin.H{"message": "There was a problem creating a login session"})
		c.Abort()
		return
	}

	c.JSON(200, gin.H{"message": "Log in success", "token": jwtToken})
}

// Verify controller handles account verification
func (user *UserController) Verify(c *gin.Context) {
	token := c.Query("token")

	doc, _ := blacklistModel.FindToken(token)

	if doc.Token != "" {
		c.JSON(403, gin.H{"message": "Token expired, renew one"})
		c.Abort()
		return
	}

	email, msg, _ := services.DecodeToken(token)

	if msg == "Token is invalid" {
		c.JSON(403, gin.H{"message": msg})
		c.Abort()
		return
	}

	if msg == "Error decoding token" {
		c.JSON(400, gin.H{"message": msg})
		c.Abort()
		return
	}

	if msg == "Token is Invalid" {
		c.JSON(403, gin.H{"message": msg})
		c.Abort()
		return
	}

	findUser, _ := userModel.GetUserByEmail(email)

	if findUser.Email == "" {
		c.JSON(404, gin.H{"message": "Account not found"})
		c.Abort()
		return
	}

	if findUser.IsVerified == true {
		c.JSON(400, gin.H{"message": "Account has already been verified"})
		c.Abort()
		return
	}

	_, err := userModel.VerifyAccount(email)

	if err != nil {
		c.JSON(403, gin.H{"message": "Problem verifying yout account "})
		c.Abort()
		return
	}

	err = blacklistModel.Add(token)

	if err != nil {
		c.JSON(403, gin.H{"message": "There was a problem blacklisting yor session"})
		c.Abort()
		return
	}

	c.JSON(201, gin.H{"message": "Account has been verified"})
}

// PasswordRequest controller handles password reset request
func (user *UserController) PasswordRequest(c *gin.Context) {
	var data forms.PasswordRequestCommand

	if c.BindJSON(&data) != nil {
		c.JSON(406, gin.H{"message": "Provide email"})
		c.Abort()
		return
	}

	// Handle email password reset sendout
	verificationToken, err := services.GenerateToken(data.Email)

	link := os.Getenv("BASE_URL") + "/password-reset-submit?token=" + verificationToken

	subject := "Password Reset"
	body := "You are receiving this email because you requested for a password reset. Click this link to verify <a href='" + link + "'>here</a>."

	go services.SendMail(data.Email, "", subject, body)

	if err != nil {
		c.JSON(406, gin.H{"message": "Error generating password request sessoin"})
	}

	c.JSON(201, gin.H{"message": "Password request sent to your email"})
}

// PasswordResetSubmit controller handles updating user password
func (user *UserController) PasswordResetSubmit(c *gin.Context) {
	var data forms.PasswordSubmitCommand

	token := c.Query("token")

	email, msg, _ := services.DecodeToken(token)

	if msg == "Token is invalid" {
		c.JSON(403, gin.H{"message": msg})
		c.Abort()
		return
	}

	if msg == "Error decoding token" {
		c.JSON(400, gin.H{"message": "Token expired, resend email"})
		c.Abort()
		return
	}

	if msg == "Token is Invalid" {
		c.JSON(403, gin.H{"message": msg})
		c.Abort()
		return
	}

	doc, _ := blacklistModel.FindToken(token)

	if doc.Token != "" {
		c.JSON(400, gin.H{"message": "Token has expired, request a new one"})
		c.Abort()
		return
	}

	findUser, _ := userModel.GetUserByEmail(email)

	if findUser.Email == "" {
		c.JSON(404, gin.H{"message": "Account not found"})
		c.Abort()
		return
	}

	// Bind request body to define data struct
	if c.BindJSON(&data) != nil {
		c.JSON(406, gin.H{"message": "Provide passwords"})
		c.Abort()
		return
	}

	// Blacklist the token
	err := blacklistModel.Add(token)

	if err != nil {
		c.JSON(403, gin.H{"message": "Error blacklisting session"})
		c.Abort()
		return
	}

	// Check if passwords match
	if data.Password != data.Confirm {
		c.JSON(400, gin.H{"message": "Passwords do not match"})
		c.Abort()
		return
	}

	// Generate hash for password
	hashedPassword := helpers.GeneratePasswordHash([]byte(data.Password))

	// Update user password
	_, err2 := userModel.PasswordUpdate(email, hashedPassword)

	// Check if we encounter an error updating an account
	if err2 != nil {
		c.JSON(404, gin.H{"message": "Account does not exist"})
		c.Abort()
		return
	}

	c.JSON(201, gin.H{"message": "Password update successfully"})
}
