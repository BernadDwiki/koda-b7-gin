package dto

type RegisterRequest struct {
	Email    string `json:"email" binding:"required,email"`
	Password string `json:"password" binding:"required,min=8"`
}

type LoginRequest struct {
	Email    string `json:"email" binding:"required,email"`
	Password string `json:"password" binding:"required"`
}

type LoginResponse struct {
	Token  string `json:"token"`
	HasPin bool   `json:"has_pin"`
}

type RegisterResponse struct {
	ID    int    `json:"id"`
	Email string `json:"email"`
}
