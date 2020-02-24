package models

import (
	"time"
	
	"github.com/tesh254/emania-api/forms"
	"github.com/tesh254/emania-api/helpers"
	"gopkg.in/mgo.v2/bson"
)

// User defines user structure
type User struct {
	ID         bson.ObjectId       `json:"_id,omitempty" bson:"_id,omitempty"`
	Name       string              `json:"name" bson:"name"`
	Email      string              `json:"email" bson:"email"`
	Password   string              `json:"password" bson:"password"`
	Role       string              `json:"role" bson:"role"`
	IsVerified bool                `json:"is_verified" bson:"is_verified"`
	CreatedAt  bson.MongoTimestamp `json:"created_at" bson:"created_at"`
	UpdatedAt  bson.MongoTimestamp `json:"updated_at" bson:"updated_at"`
}

// UserModel defines the model structure
type UserModel struct{}

// Find handles fetching all users
func (u *UserModel) Find() (list []User, err error) {
	collection := dbConnect.Use(databaseName, "account")
	err = collection.Find(bson.M{}).All(&list)
	return list, err
}

// Signup handle registering a user
func (u *UserModel) Signup(data forms.SignupUserCommand) error {
	collection := dbConnect.Use(databaseName, "account")
	err := collection.Insert(bson.M{
		"name":        data.Name,
		"email":       data.Email,
		"password":    helpers.GeneratePasswordHash([]byte(data.Password)),
		"role":        data.Role,
		"user_id":     helpers.String(8),
		"is_verified": false,
		"created_at":  time.Now().Format("2006-01-02"),
		"update_at":   time.Now().Format("2006-01-02"),
	})

	return err
}

// Login handles finding user by email and account verified
func (u *UserModel) Login(email string) (user User, err error) {
	collection := dbConnect.Use(databaseName, "account")
	err = collection.Find(bson.M{"email": email}).One(&user)

	return user, err
}

// GetUserByEmail handle getting a single user
func (u *UserModel) GetUserByEmail(email string) (user User, err error) {
	collection := dbConnect.Use(databaseName, "account")
	err = collection.Find(bson.M{"email": email}).One(&user)

	return user, err
}

// VerifyAccount handle verifying a user account
func (u *UserModel) VerifyAccount(email string) (user User, err error) {
	collection := dbConnect.Use(databaseName, "account")
	err = collection.Update(bson.M{"email": email}, bson.M{"$set": bson.M{"is_verified": true}})

	return user, err
}

// PasswordUpdate handle updating user password
func (u *UserModel) PasswordUpdate(email string, hashedPassword string) (user User, err error) {
	collection := dbConnect.Use(databaseName, "account")
	err = collection.Update(bson.M{"email": email}, bson.M{"$set": bson.M{"password": hashedPassword }})

	return user, err
}
 
// GetUserByEmailCount handle getting a single user
func (u *UserModel) GetUserByEmailCount(email string) (user User, err error) {
	collection := dbConnect.Use(databaseName, "account")
	err = collection.Find(bson.M{"email": email}).One(&user)

	return user, err
}

// GetUserByUserID handle getting a user by userid
func (u *UserModel) GetUserByUserID(userID string) (user User, err error) {
	collection := dbConnect.Use(databaseName, "account")
	err = collection.Find(bson.M{"user_id": userID}).One(&user)
	return user, err
}
