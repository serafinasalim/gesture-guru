package models

type User struct {
	Id       int     `json:"id"`
	Username string  `json:"username"`
	Email    string  `json:"email"`
	Password string  `json:"password"`
	Profile  *string `json:"profile"`
	Bio      *string `json:"bio"`
	Verified bool    `json:"verified"`
}

type UserRegister struct {
	Username        string `json:"username" validate:"required,max=20"`
	Email           string `json:"email" validate:"required,email,max=50"`
	Password        string `json:"password" validate:"required,max=50"`
	ConfirmPassword string `json:"confirmPassword" validate:"required,max=50"`
}

type UserVerify struct {
	Otp string `json:"otp" validate:"required,min=6,max=6"`
}

type UserLogin struct {
	Username string `json:"username" validate:"required,max=20"`
	Password string `json:"password" validate:"required,max=50"`
}

type UserUpdate struct {
	Profile  string `json:"profile"`
	Username string `json:"username" validate:"max=20"`
	Bio      string `json:"bio"`
}

type UserUpdateEmail struct {
	Email string `json:"email" validate:"required,email,max=50"`
}

type UserDetail struct {
	UserId int `json:"userId" validate:"required"`
}
