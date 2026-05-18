package dto

type ProfileResponse struct {
	ID          int    `json:"id"`
	Name        string `json:"name"`
	Email       string `json:"email"`
	Picture     string `json:"picture"`
	PhoneNumber string `json:"phone_number"`
}
