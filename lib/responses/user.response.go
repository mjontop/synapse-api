package responses

type user struct {
	Email    string `json:"email"`
	Token    string `json:"token"`
	Username string `json:"username"`
	Bio      string `json:"bio"`
	Image    string `json:"image"`
}

type LoginUserResponse struct {
	User user `json:"user"`
}

func NewLoginUserResponse(email, token, username, bio, image string) LoginUserResponse {
	user := user{
		Email:    email,
		Token:    token,
		Username: username,
		Bio:      bio,
		Image:    image,
	}
	return LoginUserResponse{User: user}
}
