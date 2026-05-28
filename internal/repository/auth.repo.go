package repository

import (
	"context"

	"github.com/bernaddwiki/koda-b7-weekly10/internal/model"
	"github.com/jackc/pgx/v5/pgxpool"
)

type IAuthRepository interface {
	CreateUser(ctx context.Context, user model.User) (*model.User, error)
	FindUserByEmail(ctx context.Context, email string) (*model.User, error)
	IsEmailTaken(ctx context.Context, email string) (bool, error)
	IsPhoneTaken(ctx context.Context, phoneNumber string) (bool, error)
	UpdatePassword(ctx context.Context, userID int, password string) error
}

type AuthRepository struct {
	db *pgxpool.Pool
}

func NewAuthRepository(db *pgxpool.Pool) IAuthRepository {
	return &AuthRepository{db}
}

func (r *AuthRepository) CreateUser(ctx context.Context, user model.User) (*model.User, error) {
	query := `
	INSERT INTO users (
		name,
		email,
		password,
		pin,
		phone_number
	)
	VALUES ($1,$2,$3,$4,$5)
	RETURNING id, name, email, phone_number, created_at
	`

	var created model.User

	err := r.db.QueryRow(
		ctx,
		query,
		user.Name,
		user.Email,
		user.Password,
		user.Pin,
		user.PhoneNumber,
	).Scan(
		&created.ID,
		&created.Name,
		&created.Email,
		&created.PhoneNumber,
		&created.CreatedAt,
	)

	if err != nil {
		return nil, err
	}

	_, err = r.db.Exec(
		ctx,
		`INSERT INTO ewallets (user_id) VALUES ($1)`,
		created.ID,
	)

	if err != nil {
		return nil, err
	}

	return &created, nil
}

func (r *AuthRepository) FindUserByEmail(ctx context.Context, email string) (*model.User, error) {
	query := `
	SELECT
		id,
		name,
		email,
		password,
		pin,
		phone_number
	FROM users
	WHERE email = $1
	`

	var user model.User

	err := r.db.QueryRow(ctx, query, email).Scan(
		&user.ID,
		&user.Name,
		&user.Email,
		&user.Password,
		&user.Pin,
		&user.PhoneNumber,
	)

	if err != nil {
		return nil, err
	}

	return &user, nil
}

func (r *AuthRepository) IsEmailTaken(ctx context.Context, email string) (bool, error) {
	query := `
	SELECT EXISTS(
		SELECT 1
		FROM users
		WHERE email = $1
	)
	`

	var exists bool
	err := r.db.QueryRow(ctx, query, email).Scan(&exists)
	return exists, err
}

func (r *AuthRepository) IsPhoneTaken(ctx context.Context, phoneNumber string) (bool, error) {
	query := `
	SELECT EXISTS(
		SELECT 1
		FROM users
		WHERE phone_number = $1
	)
	`

	var exists bool
	err := r.db.QueryRow(ctx, query, phoneNumber).Scan(&exists)
	return exists, err
}

func (r *AuthRepository) UpdatePassword(
	ctx context.Context,
	userID int,
	password string,
) error {
	query := `
	UPDATE users
	SET password = $1
	WHERE id = $2
	`

	_, err := r.db.Exec(
		ctx,
		query,
		password,
		userID,
	)

	return err
}
