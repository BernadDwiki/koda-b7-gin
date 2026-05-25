package service

import (
	"context"
	"errors"
	"time"

	"github.com/bernaddwiki/koda-b7-weekly10/internal/dto"
	"github.com/bernaddwiki/koda-b7-weekly10/internal/hash"
	"github.com/bernaddwiki/koda-b7-weekly10/internal/jwt"
	"github.com/bernaddwiki/koda-b7-weekly10/internal/model"
	"github.com/bernaddwiki/koda-b7-weekly10/internal/repository"
)

type IAuthService interface {
	Register(ctx context.Context, req dto.RegisterRequest) (*model.User, error)
	Login(ctx context.Context, req dto.LoginRequest) (*dto.LoginResponse, error)
	Logout(
		ctx context.Context,
		userID int,
		token string,
		expiredAt time.Time,
	) error
}

type AuthService struct {
	repo repository.IAuthRepository
}

func NewAuthService(repo repository.IAuthRepository) IAuthService {
	return &AuthService{repo}
}

func (s *AuthService) Register(
	ctx context.Context,
	req dto.RegisterRequest,
) (*model.User, error) {
	emailTaken, err := s.repo.IsEmailTaken(ctx, req.Email)
	if err != nil {
		return nil, err
	}
	if emailTaken {
		return nil, errors.New("email already registered")
	}

	passwordHash, err := hash.GenerateHash(req.Password)
	if err != nil {
		return nil, err
	}

	user := model.User{
		Name:        "",
		Email:       req.Email,
		Password:    passwordHash,
		Pin:         "",
		PhoneNumber: nil,
	}

	return s.repo.CreateUser(ctx, user)
}

func (s *AuthService) Login(
	ctx context.Context,
	req dto.LoginRequest,
) (*dto.LoginResponse, error) {
	user, err := s.repo.FindUserByEmail(ctx, req.Email)
	if err != nil {
		return nil, errors.New("invalid credentials")
	}

	valid := hash.VerifyPassword(
		req.Password,
		user.Password,
	)

	if !valid {
		return nil, errors.New("invalid credentials")
	}

	token, err := jwt.GenerateToken(
		user.ID,
		user.Email,
	)

	if err != nil {
		return nil, err
	}

	hasPin := user.Pin != ""

	return &dto.LoginResponse{
		Token:  token,
		HasPin: hasPin,
	}, nil
}

func (s *AuthService) Logout(
	ctx context.Context,
	userID int,
	token string,
	expiredAt time.Time,
) error {
	return s.repo.StoreRevokedToken(
		ctx,
		userID,
		token,
		expiredAt,
	)
}
