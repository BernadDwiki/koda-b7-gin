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
	"github.com/google/uuid"
	"github.com/redis/go-redis/v9"
)

var (
	ErrEmailAlreadyRegistered = errors.New("email already registered")
)

type IAuthService interface {
	Register(ctx context.Context, req dto.RegisterRequest) (*model.User, error)
	Login(ctx context.Context, req dto.LoginRequest) (*dto.LoginResponse, error)
	ForgotPassword(
		ctx context.Context,
		req dto.ForgotPasswordRequest,
	) (*dto.ForgotPasswordResponse, error)
	ResetPassword(
		ctx context.Context,
		req dto.ResetPasswordRequest,
	) error
	Logout(
		ctx context.Context,
		userID int,
		token string,
		expiredAt time.Time,
	) error
}

type AuthService struct {
	repo  repository.IAuthRepository
	redis *redis.Client
}

func NewAuthService(
	repo repository.IAuthRepository,
	redisClient *redis.Client,
) IAuthService {
	return &AuthService{repo: repo, redis: redisClient}
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
		return nil, ErrEmailAlreadyRegistered
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

func (s *AuthService) ForgotPassword(
	ctx context.Context,
	req dto.ForgotPasswordRequest,
) (*dto.ForgotPasswordResponse, error) {
	user, err := s.repo.FindUserByEmail(
		ctx,
		req.Email,
	)

	if err != nil {
		return nil, errors.New("email not found")
	}

	resetToken := uuid.NewString()

	err = s.redis.Set(
		ctx,
		"forgot-password:"+resetToken,
		user.ID,
		15*time.Minute,
	).Err()

	if err != nil {
		return nil, err
	}

	return &dto.ForgotPasswordResponse{
		ResetToken: resetToken,
	}, nil
}

func (s *AuthService) ResetPassword(
	ctx context.Context,
	req dto.ResetPasswordRequest,
) error {
	key := "forgot-password:" + req.Token

	userID, err := s.redis.Get(
		ctx,
		key,
	).Int()

	if err != nil {
		return errors.New("invalid or expired token")
	}

	passwordHash, err := hash.GenerateHash(
		req.NewPassword,
	)

	if err != nil {
		return err
	}

	err = s.repo.UpdatePassword(
		ctx,
		userID,
		passwordHash,
	)

	if err != nil {
		return err
	}

	_ = s.redis.Del(ctx, key)

	return nil
}

func (s *AuthService) Logout(
	ctx context.Context,
	userID int,
	token string,
	expiredAt time.Time,
) error {
	ttl := time.Until(expiredAt)
	if ttl <= 0 {
		ttl = time.Second
	}

	return s.redis.Set(
		ctx,
		"blacklist:"+token,
		"revoked",
		ttl,
	).Err()
}
