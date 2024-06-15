package requests

import "github.com/mjontop/synapse-api/models"

type UserRegisterRequestType struct {
	User models.User `json:"user"`
}

type loginUser struct {
	Username string `json:"username"`
	Email    string `json:"email"`
	Password string `json:"password"`
}

type UserLoginRequestType struct {
	User loginUser `json:"user"`
}
