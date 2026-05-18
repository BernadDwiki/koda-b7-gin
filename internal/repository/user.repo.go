package repository

import (
	"context"

	"github.com/bernaddwiki/koda-b7-weekly10/internal/model"
	"github.com/jackc/pgx/v5/pgxpool"
)

type IUserRepository interface {
	GetProfile(ctx context.Context, userID int) (*model.User, error)
}

type UserRepository struct {
	db *pgxpool.Pool
}

func NewUserRepository(
	db *pgxpool.Pool,
) IUserRepository {
	return &UserRepository{db}
}

func (u *UserRepository) GetProfile(
	ctx context.Context,
	userID int,
) (*model.User, error) {
	query := `
	SELECT
		id,
		name,
		email,
		COALESCE(picture, ''),
		COALESCE(phone_number, ''),
		created_at,
		updated_at
	FROM users
	WHERE id = $1
	`

	var user model.User

	err := u.db.QueryRow(
		ctx,
		query,
		userID,
	).Scan(
		&user.ID,
		&user.Name,
		&user.Email,
		&user.Picture,
		&user.PhoneNumber,
		&user.CreatedAt,
		&user.UpdatedAt,
	)

	if err != nil {
		return nil, err
	}

	return &user, nil
}
