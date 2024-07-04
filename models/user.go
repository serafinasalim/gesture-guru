package models

type User struct {
	Id       int    `json:"id"`
	Username string `json:"username"`
	Email    string `json:"email"`
	Password string `json:"password"`
	Profile  string `json:"profile"`
}

type UserRegister struct {
	Username        string `json:"username" validate:"required,max=20"`
	Email           string `json:"email" validate:"required,max=50"`
	Password        string `json:"password" validate:"required,max=50"`
	ConfirmPassword string `json:"confirmPassword" validate:"required,max=50"`
}

type UserLogin struct {
	Username string `json:"username" validate:"required,max=20"`
	Password string `json:"password" validate:"required,max=50"`
}
