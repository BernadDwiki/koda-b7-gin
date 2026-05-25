package service

import (
	"context"
	"errors"
	"fmt"

	"github.com/bernaddwiki/koda-b7-weekly10/internal/dto"
	"github.com/bernaddwiki/koda-b7-weekly10/internal/hash"
	"github.com/bernaddwiki/koda-b7-weekly10/internal/repository"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
)

var (
	ErrInvalidTransferRequest = errors.New("invalid transfer request")
	ErrInvalidTopUpRequest    = errors.New("invalid top up request")
	ErrInvalidReceiver        = errors.New("invalid receiver")
	ErrReceiverNotFound       = errors.New("receiver not found")
	ErrInvalidPin             = errors.New("invalid pin")
	ErrPinNotSet              = errors.New("pin is not set")
	ErrPaymentMethodNotFound  = errors.New("payment method not found")
	ErrInsufficientBalance    = errors.New("insufficient balance")
	ErrEwalletNotFound        = errors.New("ewallet not found")
)

type IWalletService interface {
	GetDashboard(ctx context.Context, userID int) (*repository.WalletDashboard, error)
	CreateTransfer(ctx context.Context, senderID int, input dto.CreateTransferRequest) (*dto.CreateTransferResponse, error)
	CreateTopUp(ctx context.Context, receiverID int, input dto.CreateTopUpRequest) (*dto.CreateTopUpResponse, error)
	GetTransactionHistory(ctx context.Context, userID int, search string, page, limit int) ([]repository.TransactionReportItem, int, error)
	GetTransactionReport(ctx context.Context, userID int, start, end, flow string) ([]repository.TransactionReportItem, error)
}

type WalletService struct {
	repo     repository.IWalletRepository
	userRepo repository.IUserRepository
	db       *pgxpool.Pool
}

func NewWalletService(repo repository.IWalletRepository, userRepo repository.IUserRepository, db *pgxpool.Pool) IWalletService {
	return &WalletService{repo: repo, userRepo: userRepo, db: db}
}

func (w *WalletService) GetDashboard(
	ctx context.Context,
	userID int,
) (*repository.WalletDashboard, error) {
	return w.repo.GetDashboard(ctx, w.db, userID)
}

func (w *WalletService) CreateTransfer(
	ctx context.Context,
	senderID int,
	input dto.CreateTransferRequest,
) (*dto.CreateTransferResponse, error) {
	if senderID == input.ReceiverID {
		return nil, ErrInvalidReceiver
	}

	if input.Amount <= 0 {
		return nil, fmt.Errorf("%w: amount must be greater than zero", ErrInvalidTransferRequest)
	}

	if len(input.Note) > 255 {
		return nil, fmt.Errorf("%w: note cannot exceed 255 characters", ErrInvalidTransferRequest)
	}

	tx, err := w.db.BeginTx(ctx, pgx.TxOptions{})
	if err != nil {
		return nil, err
	}
	defer tx.Rollback(ctx)

	if input.Pin == "" {
		return nil, ErrInvalidPin
	}

	existingPinHash, err := w.userRepo.GetUserPin(ctx, senderID)
	if err != nil {
		return nil, err
	}
	if existingPinHash == "" {
		return nil, ErrPinNotSet
	}
	if !hash.VerifyPassword(input.Pin, existingPinHash) {
		return nil, ErrInvalidPin
	}

	senderExists, err := w.repo.EwalletExists(ctx, tx, senderID)
	if err != nil {
		return nil, err
	}
	if !senderExists {
		return nil, ErrEwalletNotFound
	}

	receiverExists, err := w.repo.EwalletExists(ctx, tx, input.ReceiverID)
	if err != nil {
		return nil, err
	}
	if !receiverExists {
		return nil, ErrReceiverNotFound
	}

	senderWallet, err := w.repo.GetDashboard(ctx, tx, senderID)
	if err != nil {
		return nil, err
	}
	if senderWallet.Balance < input.Amount {
		return nil, ErrInsufficientBalance
	}

	transactionID, err := w.repo.CreateTransfer(ctx, tx, senderID, input.ReceiverID, input.Amount, input.Note)
	if err != nil {
		return nil, err
	}

	if err = tx.Commit(ctx); err != nil {
		return nil, err
	}

	return &dto.CreateTransferResponse{
		TransactionID: transactionID,
		Status:        "success",
		Message:       "transfer completed",
	}, nil
}

func (w *WalletService) CreateTopUp(
	ctx context.Context,
	receiverID int,
	input dto.CreateTopUpRequest,
) (*dto.CreateTopUpResponse, error) {
	if input.Amount <= 0 {
		return nil, fmt.Errorf("%w: amount must be greater than zero", ErrInvalidTopUpRequest)
	}

	if len(input.Note) > 255 {
		return nil, fmt.Errorf("%w: note cannot exceed 255 characters", ErrInvalidTopUpRequest)
	}

	tx, err := w.db.BeginTx(ctx, pgx.TxOptions{})
	if err != nil {
		return nil, err
	}
	defer tx.Rollback(ctx)

	receiverExists, err := w.repo.EwalletExists(ctx, tx, receiverID)
	if err != nil {
		return nil, err
	}
	if !receiverExists {
		return nil, ErrEwalletNotFound
	}

	paymentMethodConfig, err := w.repo.GetPaymentMethodTopUpConfig(ctx, tx, input.PaymentMethodID)
	if err != nil {
		return nil, err
	}
	if paymentMethodConfig == nil {
		return nil, ErrPaymentMethodNotFound
	}

	taxAmount := input.Amount * paymentMethodConfig.TaxPercent / 100
	total := input.Amount + taxAmount + paymentMethodConfig.AdminFee

	transactionID, err := w.repo.CreateTopUp(
		ctx,
		tx,
		receiverID,
		input.PaymentMethodID,
		input.Amount,
		taxAmount,
		paymentMethodConfig.AdminFee,
		input.Note,
	)
	if err != nil {
		return nil, err
	}

	if err = tx.Commit(ctx); err != nil {
		return nil, err
	}

	return &dto.CreateTopUpResponse{
		TransactionID: transactionID,
		Status:        "success",
		Message:       "top up completed",
		Amount:        input.Amount,
		TaxPercent:    paymentMethodConfig.TaxPercent,
		TaxAmount:     taxAmount,
		AdminFee:      paymentMethodConfig.AdminFee,
		Total:         total,
	}, nil
}

func (w *WalletService) GetTransactionReport(
	ctx context.Context,
	userID int,
	start, end, flow string,
) ([]repository.TransactionReportItem, error) {
	return w.repo.GetTransactionReport(ctx, userID, start, end, flow)
}

func (w *WalletService) GetTransactionHistory(
	ctx context.Context,
	userID int,
	search string,
	page, limit int,
) ([]repository.TransactionReportItem, int, error) {
	if page <= 0 {
		page = 1
	}
	if limit <= 0 || limit > 100 {
		limit = 10
	}

	offset := (page - 1) * limit

	items, total, err := w.repo.GetTransactionHistory(ctx, userID, search, limit, offset)
	if err != nil {
		return nil, 0, err
	}

	return items, total, nil
}
