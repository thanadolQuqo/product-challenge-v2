package models

type User struct {
	Username  string `json:"username"`
	Password  string `json:"password"`
	Token     string `json:"token"`
	ExpiresAt int64  `json:"expires"`
	CreatedAt string `json:"created_at"`
}

type UserAuthRequest struct {
	Username string `json:"username"`
	Password string `json:"password"`
}

type UserAuthResponse struct {
	Username string `json:"username"`
	Token    string `json:"token"`
}
