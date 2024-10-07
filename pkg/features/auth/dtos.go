package auth

type LoginRequest struct {
	Email    string
	Password string
}

type LoginResponse struct {
	AccessToken  string
	RefreshToken string
}

type RegisterRequest struct {
	Email    string
	Password string
}
