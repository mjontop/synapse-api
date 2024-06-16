package responses

type user struct {
	Email    string `json:"email"`
	Token    string `json:"token"`
	Username string `json:"username"`
	Bio      string `json:"bio"`
	Image    string `json:"image"`
}

type LoggedInUserResponse struct {
	User user `json:"user"`
}

func NewLoggedInUserResponse(email, token, username, bio, image string) LoggedInUserResponse {
	user := user{
		Email:    email,
		Token:    token,
		Username: username,
		Bio:      bio,
		Image:    image,
	}
	return LoggedInUserResponse{User: user}
}
