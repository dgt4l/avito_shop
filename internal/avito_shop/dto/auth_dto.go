package dto

type AuthRequest struct {
	Username string `json:"username" db:"username"`
	Password string `json:"password" db:"password_salt"`
}

type AuthResponse struct {
	Token string `json:"token"`
}
