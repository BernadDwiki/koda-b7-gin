package repository

import (
	"context"

	"github.com/jackc/pgx/v5/pgxpool"
)

type WalletDashboard struct {
	Balance int64
	Income  int64
	Expense int64
}

type IWalletRepository interface {
	GetDashboard(ctx context.Context, userID int) (*WalletDashboard, error)
}

type WalletRepository struct {
	db *pgxpool.Pool
}

func NewWalletRepository(db *pgxpool.Pool) IWalletRepository {
	return &WalletRepository{db}
}

func (w *WalletRepository) GetDashboard(
	ctx context.Context,
	userID int,
) (*WalletDashboard, error) {
	query := `
	SELECT
		balance,
		income,
		expense
	FROM ewallets
	WHERE user_id = $1
	`

	var data WalletDashboard

	err := w.db.QueryRow(
		ctx,
		query,
		userID,
	).Scan(
		&data.Balance,
		&data.Income,
		&data.Expense,
	)

	if err != nil {
		return nil, err
	}

	return &data, nil
}
