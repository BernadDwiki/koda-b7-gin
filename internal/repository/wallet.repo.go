package repository

import (
	"context"
	"errors"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgconn"
	"github.com/jackc/pgx/v5/pgxpool"
)

type DBTX interface {
	Exec(ctx context.Context, sql string, args ...any) (pgconn.CommandTag, error)
	Query(ctx context.Context, sql string, args ...any) (pgx.Rows, error)
	QueryRow(ctx context.Context, sql string, args ...any) pgx.Row
}

type WalletDashboard struct {
	Balance int64
	Income  int64
	Expense int64
}

type PaymentMethodTopUpConfig struct {
	TaxPercent int64
	AdminFee   int64
}

type IWalletRepository interface {
	GetDashboard(ctx context.Context, db DBTX, userID int) (*WalletDashboard, error)
	EwalletExists(ctx context.Context, db DBTX, userID int) (bool, error)
	PaymentMethodExists(ctx context.Context, db DBTX, paymentMethodID int) (bool, error)
	GetPaymentMethodTopUpConfig(ctx context.Context, db DBTX, paymentMethodID int) (*PaymentMethodTopUpConfig, error)
	CreateTransfer(ctx context.Context, db DBTX, senderID, receiverID int, amount int64, note string) (int64, error)
	CreateTopUp(ctx context.Context, db DBTX, receiverID, paymentMethodID int, amount, tax, adminFee int64, note string) (int64, error)
	GetTransactionReport(ctx context.Context, userID int, start, end, flow string) ([]TransactionReportItem, error)
	GetTransactionChart(ctx context.Context, userID int, start, end, flow string) ([]TransactionChartItem, error)
	GetTransactionHistory(ctx context.Context, userID int, search string, limit, offset int) ([]TransactionReportItem, int, error)
}

type WalletRepository struct {
	db *pgxpool.Pool
}

func NewWalletRepository(db *pgxpool.Pool) IWalletRepository {
	return &WalletRepository{db}
}

func (w *WalletRepository) GetDashboard(
	ctx context.Context,
	db DBTX,
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

	err := db.QueryRow(
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

func (w *WalletRepository) CreateTransfer(
	ctx context.Context,
	db DBTX,
	senderID, receiverID int,
	amount int64,
	note string,
) (int64, error) {
	var senderBalance int64
	row := db.QueryRow(ctx, `
		SELECT balance
		FROM ewallets
		WHERE user_id = $1
		FOR UPDATE
	`, senderID)
	if err := row.Scan(&senderBalance); err != nil {
		return 0, err
	}

	var transactionID int64
	insertTx := `
		INSERT INTO transactions (amount, transaction_type, note, status)
		VALUES ($1, 'transfer', $2, 'success')
		RETURNING id
	`
	if err := db.QueryRow(ctx, insertTx, amount, note).Scan(&transactionID); err != nil {
		return 0, err
	}

	insertDetail := `
		INSERT INTO transfer_details (transaction_id, sender_id, receiver_id)
		VALUES ($1, $2, $3)
	`
	if _, err := db.Exec(ctx, insertDetail, transactionID, senderID, receiverID); err != nil {
		return 0, err
	}

	updateSender := `
		UPDATE ewallets
		SET balance = balance - $1,
		    expense = expense + $1,
		    updated_at = NOW()
		WHERE user_id = $2
	`
	if _, err := db.Exec(ctx, updateSender, amount, senderID); err != nil {
		return 0, err
	}

	updateReceiver := `
		UPDATE ewallets
		SET balance = balance + $1,
		    income = income + $1,
		    updated_at = NOW()
		WHERE user_id = $2
	`
	if _, err := db.Exec(ctx, updateReceiver, amount, receiverID); err != nil {
		return 0, err
	}

	return transactionID, nil
}

func (w *WalletRepository) CreateTopUp(
	ctx context.Context,
	db DBTX,
	receiverID, paymentMethodID int,
	amount, tax, adminFee int64,
	note string,
) (int64, error) {
	var receiverBalance int64
	row := db.QueryRow(ctx, `
		SELECT balance
		FROM ewallets
		WHERE user_id = $1
		FOR UPDATE
	`, receiverID)
	if err := row.Scan(&receiverBalance); err != nil {
		return 0, err
	}

	var transactionID int64
	insertTx := `
		INSERT INTO transactions (amount, transaction_type, note, status)
		VALUES ($1, 'top_up', $2, 'success')
		RETURNING id
	`
	if err := db.QueryRow(ctx, insertTx, amount, note).Scan(&transactionID); err != nil {
		return 0, err
	}

	insertDetail := `
		INSERT INTO top_up_details (transaction_id, receiver_id, payment_method_id, tax, admin_fee)
		VALUES ($1, $2, $3, $4, $5)
	`
	if _, err := db.Exec(ctx, insertDetail, transactionID, receiverID, paymentMethodID, tax, adminFee); err != nil {
		return 0, err
	}

	updateReceiver := `
		UPDATE ewallets
		SET balance = balance + $1,
		    updated_at = NOW()
		WHERE user_id = $2
	`
	if _, err := db.Exec(ctx, updateReceiver, amount, receiverID); err != nil {
		return 0, err
	}

	return transactionID, nil
}

func (w *WalletRepository) EwalletExists(
	ctx context.Context,
	db DBTX,
	userID int,
) (bool, error) {
	var exists bool
	query := `
	SELECT TRUE
	FROM ewallets
	WHERE user_id = $1
	`

	err := db.QueryRow(ctx, query, userID).Scan(&exists)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return false, nil
		}
		return false, err
	}

	return exists, nil
}

func (w *WalletRepository) PaymentMethodExists(
	ctx context.Context,
	db DBTX,
	paymentMethodID int,
) (bool, error) {
	var exists bool
	query := `
	SELECT TRUE
	FROM payment_methods
	WHERE id = $1
	`

	err := db.QueryRow(ctx, query, paymentMethodID).Scan(&exists)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return false, nil
		}
		return false, err
	}

	return exists, nil
}

func (w *WalletRepository) GetPaymentMethodTopUpConfig(
	ctx context.Context,
	db DBTX,
	paymentMethodID int,
) (*PaymentMethodTopUpConfig, error) {
	var config PaymentMethodTopUpConfig
	query := `
	SELECT tax_percent, admin_fee
	FROM payment_methods
	WHERE id = $1
	`

	err := db.QueryRow(ctx, query, paymentMethodID).Scan(&config.TaxPercent, &config.AdminFee)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, nil
		}
		return nil, err
	}

	return &config, nil
}

type TransactionReportItem struct {
	ID          int64  `json:"id"`
	Amount      int64  `json:"amount"`
	Type        string `json:"type"`
	Description string `json:"description"`
	Status      string `json:"status"`
	CreatedAt   string `json:"created_at"`
	Direction   string `json:"direction"`
}

type TransactionChartItem struct {
	Date             string `json:"date"`
	Type             string `json:"type"`
	TotalTransaction int64  `json:"total_transaction"`
}

func (w *WalletRepository) GetTransactionReport(
	ctx context.Context,
	userID int,
	start, end, flow string,
) ([]TransactionReportItem, error) {

	query := `
	SELECT *
	FROM (
		SELECT
			t.id,
			t.amount,
			t.transaction_type::text,
			t.note,
			t.status::text,
			to_char(t.created_at, 'YYYY-MM-DD"T"HH24:MI:SS"Z"') as created_at,
			CASE
				WHEN td.receiver_id = $1
				THEN 'income'
				ELSE 'expense'
			END as direction
		FROM transactions t
		JOIN transfer_details td
			ON td.transaction_id = t.id
		WHERE
			(td.sender_id = $1 OR td.receiver_id = $1)
			AND t.created_at BETWEEN $2 AND $3

		UNION ALL

		SELECT
			t.id,
			t.amount,
			t.transaction_type::text,
			t.note,
			t.status::text,
			to_char(t.created_at, 'YYYY-MM-DD"T"HH24:MI:SS"Z"') as created_at,
			'income' as direction
		FROM transactions t
		JOIN top_up_details tu
			ON tu.transaction_id = t.id
		WHERE
			tu.receiver_id = $1
			AND t.created_at BETWEEN $2 AND $3
	) trx
	WHERE
		$4 = 'both'
		OR trx.direction = $4
	ORDER BY created_at DESC
	`

	rows, err := w.db.Query(
		ctx,
		query,
		userID,
		start,
		end,
		flow,
	)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	result := []TransactionReportItem{}

	for rows.Next() {
		var item TransactionReportItem

		err := rows.Scan(
			&item.ID,
			&item.Amount,
			&item.Type,
			&item.Description,
			&item.Status,
			&item.CreatedAt,
			&item.Direction,
		)
		if err != nil {
			return nil, err
		}

		result = append(result, item)
	}

	return result, nil
}

func (w *WalletRepository) GetTransactionChart(
	ctx context.Context,
	userID int,
	start, end, flow string,
) ([]TransactionChartItem, error) {
	query := `
SELECT
	TO_CHAR(DATE(created_at), 'YYYY-MM-DD') AS date,
	direction AS type,
	SUM(amount) AS total_transaction
FROM (
	SELECT t.created_at, t.amount,
		CASE
			WHEN td.receiver_id = $1
			THEN 'income'
			ELSE 'expense'
		END as direction
	FROM transactions t
	JOIN transfer_details td
		ON td.transaction_id = t.id
	WHERE
		(td.sender_id = $1 OR td.receiver_id = $1)
		AND t.created_at BETWEEN $2 AND $3

	UNION ALL

	SELECT t.created_at, t.amount,
		'income' as direction
	FROM transactions t
	JOIN top_up_details tu
		ON tu.transaction_id = t.id
	WHERE
		tu.receiver_id = $1
		AND t.created_at BETWEEN $2 AND $3
) trx
WHERE
	$4 = 'both'
	OR trx.direction = $4
GROUP BY DATE(created_at), direction
ORDER BY DATE(created_at), direction
`

	rows, err := w.db.Query(ctx, query, userID, start, end, flow)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	result := []TransactionChartItem{}
	for rows.Next() {
		var item TransactionChartItem
		if err := rows.Scan(&item.Date, &item.Type, &item.TotalTransaction); err != nil {
			return nil, err
		}
		result = append(result, item)
	}

	return result, nil
}

func (w *WalletRepository) GetTransactionHistory(
	ctx context.Context,
	userID int,
	search string,
	limit, offset int,
) ([]TransactionReportItem, int, error) {

	// Count total matching rows
	countQuery := `
SELECT COUNT(*) FROM (
	SELECT
		t.id,
		t.amount,
		t.transaction_type::text,
		t.note,
		t.status::text,
		to_char(t.created_at, 'YYYY-MM-DD"T"HH24:MI:SS"Z"') as created_at,
		CASE
			WHEN td.receiver_id = $1
			THEN 'income'
			ELSE 'expense'
		END as direction
	FROM transactions t
	JOIN transfer_details td
		ON td.transaction_id = t.id
	WHERE
		(td.sender_id = $1 OR td.receiver_id = $1)

	UNION ALL

	SELECT
		t.id,
		t.amount,
		t.transaction_type::text,
		t.note,
		t.status::text,
		to_char(t.created_at, 'YYYY-MM-DD"T"HH24:MI:SS"Z"') as created_at,
		'income' as direction
	FROM transactions t
	JOIN top_up_details tu
		ON tu.transaction_id = t.id
	WHERE
		tu.receiver_id = $1
) trx
WHERE ($2 = '' OR lower(trx.note) LIKE '%' || lower($2) || '%')
`

	var total int
	if err := w.db.QueryRow(ctx, countQuery, userID, search).Scan(&total); err != nil {
		return nil, 0, err
	}

	query := `
SELECT * FROM (
	SELECT
		t.id,
		t.amount,
		t.transaction_type::text,
		t.note,
		t.status::text,
		to_char(t.created_at, 'YYYY-MM-DD"T"HH24:MI:SS"Z"') as created_at,
		CASE
			WHEN td.receiver_id = $1
			THEN 'income'
			ELSE 'expense'
		END as direction
	FROM transactions t
	JOIN transfer_details td
		ON td.transaction_id = t.id
	WHERE
		(td.sender_id = $1 OR td.receiver_id = $1)

	UNION ALL

	SELECT
		t.id,
		t.amount,
		t.transaction_type::text,
		t.note,
		t.status::text,
		to_char(t.created_at, 'YYYY-MM-DD"T"HH24:MI:SS"Z"') as created_at,
		'income' as direction
	FROM transactions t
	JOIN top_up_details tu
		ON tu.transaction_id = t.id
	WHERE
		tu.receiver_id = $1
) trx
WHERE ($2 = '' OR lower(trx.note) LIKE '%' || lower($2) || '%')
ORDER BY created_at DESC
LIMIT $3 OFFSET $4
`

	rows, err := w.db.Query(ctx, query, userID, search, limit, offset)
	if err != nil {
		return nil, 0, err
	}
	defer rows.Close()

	result := []TransactionReportItem{}

	for rows.Next() {
		var item TransactionReportItem
		if err := rows.Scan(
			&item.ID,
			&item.Amount,
			&item.Type,
			&item.Description,
			&item.Status,
			&item.CreatedAt,
			&item.Direction,
		); err != nil {
			return nil, 0, err
		}
		result = append(result, item)
	}

	return result, total, nil
}
