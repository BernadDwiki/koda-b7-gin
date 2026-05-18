package service

import (
	"context"

	"github.com/bernaddwiki/koda-b7-weekly10/internal/repository"
)

type IWalletService interface {
	GetDashboard(ctx context.Context, userID int) (*repository.WalletDashboard, error)
}

type WalletService struct {
	repo repository.IWalletRepository
}

func NewWalletService(repo repository.IWalletRepository) IWalletService {
	return &WalletService{repo}
}

func (w *WalletService) GetDashboard(
	ctx context.Context,
	userID int,
) (*repository.WalletDashboard, error) {
	return w.repo.GetDashboard(ctx, userID)
}
