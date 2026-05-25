package repository

import (
	"context"
	"fmt"
	"strings"

	"github.com/bernaddwiki/koda-b7-weekly10/internal/model"
	"github.com/jackc/pgx/v5/pgxpool"
)

type IUserRepository interface {
	GetProfile(ctx context.Context, userID int) (*model.User, error)
	GetUserPin(ctx context.Context, userID int) (string, error)
	GetPassword(ctx context.Context, userID int) (string, error)
	UpdatePin(ctx context.Context, userID int, pin string) error
	UpdatePassword(ctx context.Context, userID int, password string) error
	UpdateProfile(ctx context.Context, userID int, name, picture, phoneNumber string) error
	FindReceivers(ctx context.Context, userID int, keyword string, page, limit int) ([]model.User, int, error)
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

func (u *UserRepository) GetUserPin(
	ctx context.Context,
	userID int,
) (string, error) {
	query := `
	SELECT
		COALESCE(pin, '')
	FROM users
	WHERE id = $1
	`

	var pin string

	err := u.db.QueryRow(ctx, query, userID).Scan(&pin)

	if err != nil {
		return "", err
	}

	return pin, nil
}

func (u *UserRepository) GetPassword(
	ctx context.Context,
	userID int,
) (string, error) {
	query := `
	SELECT
		password
	FROM users
	WHERE id = $1
	`

	var password string
	err := u.db.QueryRow(ctx, query, userID).Scan(&password)
	if err != nil {
		return "", err
	}

	return password, nil
}

func (u *UserRepository) UpdatePin(
	ctx context.Context,
	userID int,
	pin string,
) error {
	query := `
	UPDATE users
	SET pin = $1, updated_at = NOW()
	WHERE id = $2
	`

	_, err := u.db.Exec(ctx, query, pin, userID)

	return err
}

func (u *UserRepository) UpdatePassword(
	ctx context.Context,
	userID int,
	password string,
) error {
	query := `
	UPDATE users
	SET password = $1, updated_at = NOW()
	WHERE id = $2
	`

	_, err := u.db.Exec(ctx, query, password, userID)
	return err
}

func (u *UserRepository) FindReceivers(
	ctx context.Context,
	userID int,
	keyword string,
	page, limit int,
) ([]model.User, int, error) {
	offset := (page - 1) * limit
	search := "%"
	if keyword != "" {
		search = "%" + keyword + "%"
	}

	countQuery := `
	SELECT COUNT(*)
	FROM users
	WHERE id != $1
	AND (
		name ILIKE $2
		OR email ILIKE $2
		OR COALESCE(phone_number, '') ILIKE $2
	)
	`

	var total int
	if err := u.db.QueryRow(ctx, countQuery, userID, search).Scan(&total); err != nil {
		return nil, 0, err
	}

	query := `
	SELECT id, name, email, COALESCE(picture, ''), COALESCE(phone_number, '')
	FROM users
	WHERE id != $1
	AND (
		name ILIKE $2
		OR email ILIKE $2
		OR COALESCE(phone_number, '') ILIKE $2
	)
	ORDER BY name ASC
	LIMIT $3 OFFSET $4
	`

	rows, err := u.db.Query(ctx, query, userID, search, limit, offset)
	if err != nil {
		return nil, 0, err
	}
	defer rows.Close()

	var users []model.User
	for rows.Next() {
		var udata model.User
		if err := rows.Scan(&udata.ID, &udata.Name, &udata.Email, &udata.Picture, &udata.PhoneNumber); err != nil {
			return nil, 0, err
		}
		users = append(users, udata)
	}

	return users, total, nil
}

func (u *UserRepository) UpdateProfile(
	ctx context.Context,
	userID int,
	name, picture, phoneNumber string,
) error {
	setClauses := []string{}
	args := []any{}

	if name != "" {
		setClauses = append(setClauses, fmt.Sprintf("name = $%d", len(args)+1))
		args = append(args, name)
	}

	if picture != "" {
		setClauses = append(setClauses, fmt.Sprintf("picture = $%d", len(args)+1))
		args = append(args, picture)
	}

	if phoneNumber != "" {
		setClauses = append(setClauses, fmt.Sprintf("phone_number = $%d", len(args)+1))
		args = append(args, phoneNumber)
	}

	if len(setClauses) == 0 {
		return nil
	}

	setClauses = append(setClauses, "updated_at = NOW()")
	query := fmt.Sprintf(
		"UPDATE users SET %s WHERE id = $%d",
		strings.Join(setClauses, ", "),
		len(args)+1,
	)
	args = append(args, userID)

	_, err := u.db.Exec(ctx, query, args...)
	return err
}
