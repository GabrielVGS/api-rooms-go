package dtos

type AuthLoginRequest struct {
	Email    string `json:"email"`
	Password string `json:"password"`
}

type AuthLoginResponse struct {
	Token string       `json:"token"`
	User  UserResponse `json:"user"`
}

type AuthRegisterRequest struct {
	Username string `json:"username"`
	Email    string `json:"email"`
	Password string `json:"password"`
}

type AuthProfileResponse struct {
	UserID    uint   `json:"user_id"`
	Username  string `json:"username"`
	Email     string `json:"email"`
	CreatedAt string `json:"created_at"`
	UpdatedAt string `json:"updated_at"`
}
