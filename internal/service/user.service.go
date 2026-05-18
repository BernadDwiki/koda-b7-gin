package service

import (
	"context"

	"github.com/bernaddwiki/koda-b7-weekly10/internal/dto"
	"github.com/bernaddwiki/koda-b7-weekly10/internal/repository"
)

type IUserService interface {
	GetProfile(ctx context.Context, userID int) (*dto.ProfileResponse, error)
}

type UserService struct {
	repo repository.IUserRepository
}

func NewUserService(
	repo repository.IUserRepository,
) IUserService {
	return &UserService{repo}
}

func (u *UserService) GetProfile(
	ctx context.Context,
	userID int,
) (*dto.ProfileResponse, error) {
	user, err := u.repo.GetProfile(ctx, userID)
	if err != nil {
		return nil, err
	}

	return &dto.ProfileResponse{
		ID:          user.ID,
		Name:        user.Name,
		Email:       user.Email,
		Picture:     user.Picture,
		PhoneNumber: user.PhoneNumber,
	}, nil
}
