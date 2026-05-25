package service

import (
	"context"
	"errors"
	"fmt"
	"io"
	"mime/multipart"
	"os"
	"path/filepath"
	"strings"
	"time"

	"github.com/bernaddwiki/koda-b7-weekly10/internal/dto"
	"github.com/bernaddwiki/koda-b7-weekly10/internal/hash"
	"github.com/bernaddwiki/koda-b7-weekly10/internal/repository"
)

type IUserService interface {
	GetProfile(ctx context.Context, userID int) (*dto.ProfileResponse, error)
	SetPin(ctx context.Context, userID int, pin string) error
	CheckPin(ctx context.Context, userID int, pin string) (bool, error)
	EditPin(ctx context.Context, userID int, input dto.EditPinRequest) error
	FindReceivers(ctx context.Context, userID int, keyword string, page, limit int) (*dto.ReceiverListResponse, error)
	UpdateProfile(ctx context.Context, userID int, input dto.EditProfileRequest) (*dto.ProfileResponse, error)
	ChangePassword(ctx context.Context, userID int, input dto.ChangePasswordRequest) error
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

	phone := ""
	if user.PhoneNumber != nil {
		phone = *user.PhoneNumber
	}

	return &dto.ProfileResponse{
		ID:          user.ID,
		Name:        user.Name,
		Email:       user.Email,
		Picture:     user.Picture,
		PhoneNumber: phone,
	}, nil
}

func (u *UserService) SetPin(
	ctx context.Context,
	userID int,
	pin string,
) error {
	existingPin, err := u.repo.GetUserPin(ctx, userID)
	if err != nil {
		return err
	}

	if existingPin != "" {
		return errors.New("pin already set")
	}

	pinHash, err := hash.GenerateHash(pin)
	if err != nil {
		return err
	}

	return u.repo.UpdatePin(ctx, userID, pinHash)
}

func (u *UserService) CheckPin(
	ctx context.Context,
	userID int,
	pin string,
) (bool, error) {
	existingPin, err := u.repo.GetUserPin(ctx, userID)
	if err != nil {
		return false, err
	}

	if existingPin == "" {
		return false, nil
	}

	return hash.VerifyPassword(pin, existingPin), nil
}

func (u *UserService) EditPin(
	ctx context.Context,
	userID int,
	input dto.EditPinRequest,
) error {
	if input.CurrentPin == input.NewPin {
		return errors.New("new pin must be different from current pin")
	}

	currentPinHash, err := u.repo.GetUserPin(ctx, userID)
	if err != nil {
		return err
	}

	if currentPinHash == "" {
		return errors.New("pin is not set")
	}

	if !hash.VerifyPassword(input.CurrentPin, currentPinHash) {
		return errors.New("current pin is incorrect")
	}

	newPinHash, err := hash.GenerateHash(input.NewPin)
	if err != nil {
		return err
	}

	return u.repo.UpdatePin(ctx, userID, newPinHash)
}

func (u *UserService) FindReceivers(
	ctx context.Context,
	userID int,
	keyword string,
	page, limit int,
) (*dto.ReceiverListResponse, error) {
	users, total, err := u.repo.FindReceivers(ctx, userID, keyword, page, limit)
	if err != nil {
		return nil, err
	}

	items := make([]dto.ReceiverResponse, len(users))
	for i, user := range users {
		phone := ""
		if user.PhoneNumber != nil {
			phone = *user.PhoneNumber
		}
		items[i] = dto.ReceiverResponse{
			ID:          user.ID,
			Name:        user.Name,
			Email:       user.Email,
			PhoneNumber: phone,
		}
	}

	return &dto.ReceiverListResponse{
		Items: items,
		Page:  page,
		Limit: limit,
		Total: total,
	}, nil
}

func (u *UserService) UpdateProfile(
	ctx context.Context,
	userID int,
	input dto.EditProfileRequest,
) (*dto.ProfileResponse, error) {
	pictureURL := ""
	if input.ProfilePicture != nil {
		if input.ProfilePicture.Size > 2<<20 {
			return nil, errors.New("image max 2MB")
		}

		ext := strings.ToLower(filepath.Ext(input.ProfilePicture.Filename))
		allowed := map[string]bool{
			".jpg":  true,
			".jpeg": true,
			".png":  true,
		}

		if !allowed[ext] {
			return nil, errors.New("invalid image extension")
		}

		destinationDir := "public/profile"
		if err := os.MkdirAll(destinationDir, os.ModePerm); err != nil {
			return nil, err
		}

		fileName := fmt.Sprintf("%d%s", time.Now().UnixNano(), ext)
		destinationPath := filepath.Join(destinationDir, fileName)

		if err := saveUploadedFile(input.ProfilePicture, destinationPath); err != nil {
			return nil, err
		}

		pictureURL = "/img/" + fileName
	}

	if input.Name == "" && input.PhoneNumber == "" && pictureURL == "" {
		user, err := u.repo.GetProfile(ctx, userID)
		if err != nil {
			return nil, err
		}

		phone := ""
		if user.PhoneNumber != nil {
			phone = *user.PhoneNumber
		}

		return &dto.ProfileResponse{
			ID:          user.ID,
			Name:        user.Name,
			Email:       user.Email,
			Picture:     user.Picture,
			PhoneNumber: phone,
		}, nil
	}

	if err := u.repo.UpdateProfile(ctx, userID, input.Name, pictureURL, input.PhoneNumber); err != nil {
		return nil, err
	}

	user, err := u.repo.GetProfile(ctx, userID)
	if err != nil {
		return nil, err
	}

	phone := ""
	if user.PhoneNumber != nil {
		phone = *user.PhoneNumber
	}

	return &dto.ProfileResponse{
		ID:          user.ID,
		Name:        user.Name,
		Email:       user.Email,
		Picture:     user.Picture,
		PhoneNumber: phone,
	}, nil
}

func saveUploadedFile(fileHeader *multipart.FileHeader, destination string) error {
	src, err := fileHeader.Open()
	if err != nil {
		return err
	}
	defer src.Close()

	dst, err := os.Create(destination)
	if err != nil {
		return err
	}
	defer dst.Close()

	_, err = io.Copy(dst, src)
	return err
}

func (u *UserService) ChangePassword(
	ctx context.Context,
	userID int,
	input dto.ChangePasswordRequest,
) error {
	if input.CurrentPassword == input.NewPassword {
		return errors.New("new password must be different from current password")
	}

	currentPasswordHash, err := u.repo.GetPassword(ctx, userID)
	if err != nil {
		return err
	}

	if !hash.VerifyPassword(input.CurrentPassword, currentPasswordHash) {
		return errors.New("current password is incorrect")
	}

	newHash, err := hash.GenerateHash(input.NewPassword)
	if err != nil {
		return err
	}

	return u.repo.UpdatePassword(ctx, userID, newHash)
}
