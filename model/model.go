package model

type UserBean struct {
	UserName string
	Password string
}

type UserResponse struct {
	Id       string
	UserName string
}

type LoginResponse struct {
	Url string
}
