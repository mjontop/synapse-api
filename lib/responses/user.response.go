package responses

type UserDto struct {
	Email    string `json:"email"`
	Token    string `json:"token,omitempty"`
	Username string `json:"username"`
	Bio      string `json:"bio"`
	Image    string `json:"image"`
}

type LoggedInUserResponse struct {
	User UserDto `json:"user"`
}

func NewLoggedInUserResponse(userDto UserDto) LoggedInUserResponse {
	user := UserDto{
		Email:    userDto.Email,
		Token:    userDto.Token,
		Username: userDto.Username,
		Bio:      userDto.Bio,
		Image:    userDto.Image,
	}
	return LoggedInUserResponse{User: user}
}

type UsersProfileResponse = LoggedInUserResponse // both have same feilds except token

func NewUsersProfileResponse(userDto UserDto) UsersProfileResponse {

	user := UserDto{
		Email:    userDto.Email,
		Token:    "", // token sent empty intentionlly
		Username: userDto.Username,
		Bio:      userDto.Bio,
		Image:    userDto.Image,
	}

	usersProfile := NewLoggedInUserResponse(user)
	return usersProfile
}
