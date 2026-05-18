package model

import "time"

type User struct {
	ID          int
	Name        string
	Email       string
	Password    string
	Pin         string
	Picture     string
	PhoneNumber string
	CreatedAt   time.Time
	UpdatedAt   *time.Time
}
