package forms

// SignupUserCommand defines user form structure
type SignupUserCommand struct {
	Name string `json:"name" binding:"required"`
	Email string `json:"email" binding:"required"`
	Password string `json:"password" binding:"required"`
	Role string `json:"role" binding:"required"`
}

// LoginUserCommand defines user form structure
type LoginUserCommand struct {
	Email string `json:"email" binding:"required"`
	Password string `json:"password" binding:"required"`
}

// PasswordRequestCommand defines user password request structure
type PasswordRequestCommand struct {
	Email string `json:"email" binding:"required"`
}

// PasswordSubmitCommand defines user password submit struct
type PasswordSubmitCommand struct {
	Password string `json:"password" binding:"required"`
	Confirm string `json:"confirm" binding:"required"`
}

// BlacklistCommand defines token blacklist submit struct
type BlacklistCommand struct {
	Token string `json:"token" binding:"required"`
}