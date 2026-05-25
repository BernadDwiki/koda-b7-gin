package dto

import "mime/multipart"

type ProfileResponse struct {
	ID          int    `json:"id"`
	Name        string `json:"name"`
	Email       string `json:"email"`
	Picture     string `json:"picture"`
	PhoneNumber string `json:"phone_number"`
}

type ReceiverResponse struct {
	ID          int    `json:"id"`
	Name        string `json:"name"`
	Email       string `json:"email"`
	PhoneNumber string `json:"phone_number"`
}

type ReceiverListResponse struct {
	Items []ReceiverResponse `json:"items"`
	Page  int                `json:"page"`
	Limit int                `json:"limit"`
	Total int                `json:"total"`
}

type SetPinRequest struct {
	Pin string `json:"pin" binding:"required,len=6,numeric"`
}

type CheckPinRequest struct {
	Pin string `json:"pin" binding:"required,len=6,numeric"`
}

type EditProfileRequest struct {
	Name           string                `form:"name" binding:"omitempty,min=1,max=255"`
	ProfilePicture *multipart.FileHeader `form:"profile_picture" binding:"-"`
	PhoneNumber    string                `form:"phone_number" binding:"omitempty"`
}

type ChangePasswordRequest struct {
	CurrentPassword string `json:"current_password" binding:"required,min=8"`
	NewPassword     string `json:"new_password" binding:"required,min=8"`
}

type EditPinRequest struct {
	CurrentPin string `json:"current_pin" binding:"required,len=6,numeric"`
	NewPin     string `json:"new_pin" binding:"required,len=6,numeric"`
}
