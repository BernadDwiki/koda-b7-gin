package service

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"time"

	"github.com/bernaddwiki/koda-b7-weekly10/internal/dto"
	"github.com/bernaddwiki/koda-b7-weekly10/internal/hash"
	"github.com/bernaddwiki/koda-b7-weekly10/internal/repository"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/redis/go-redis/v9"
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
	GetTransactionChart(ctx context.Context, userID int, start, end, flow string) ([]dto.TransactionDailyReportItem, error)
}

type WalletService struct {
	repo     repository.IWalletRepository
	userRepo repository.IUserRepository
	db       *pgxpool.Pool
	redis    *redis.Client
}

func NewWalletService(
	repo repository.IWalletRepository,
	userRepo repository.IUserRepository,
	db *pgxpool.Pool,
	redisClient *redis.Client,
) IWalletService {
	return &WalletService{repo: repo, userRepo: userRepo, db: db, redis: redisClient}
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

	// invalidate cache for sender and receiver history
	if w.redis != nil {
		patternSender := fmt.Sprintf("history:user:%d:*", senderID)
		senderKeys, _ := w.redis.Keys(ctx, patternSender).Result()
		if len(senderKeys) > 0 {
			w.redis.Del(ctx, senderKeys...)
		}

		patternReceiver := fmt.Sprintf("history:user:%d:*", input.ReceiverID)
		receiverKeys, _ := w.redis.Keys(ctx, patternReceiver).Result()
		if len(receiverKeys) > 0 {
			w.redis.Del(ctx, receiverKeys...)
		}
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

	// invalidate cache for receiver history
	if w.redis != nil {
		pattern := fmt.Sprintf("history:user:%d:*", receiverID)
		keys, _ := w.redis.Keys(ctx, pattern).Result()
		if len(keys) > 0 {
			w.redis.Del(ctx, keys...)
		}
	}

	return &dto.CreateTopUpResponse{
		TransactionID: transactionID,
		Status:        "success",
		Message:       "top up completed",
		Amount:        input.Amount,
		TaxPercent:    float64(paymentMethodConfig.TaxPercent) / 100,
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

func (w *WalletService) GetTransactionChart(
	ctx context.Context,
	userID int,
	start, end, flow string,
) ([]dto.TransactionDailyReportItem, error) {
	items, err := w.repo.GetTransactionChart(ctx, userID, start, end, flow)
	if err != nil {
		return nil, err
	}

	result := make([]dto.TransactionDailyReportItem, len(items))
	for i, item := range items {
		result[i] = dto.TransactionDailyReportItem{
			Date:             item.Date,
			Type:             item.Type,
			TotalTransaction: item.TotalTransaction,
		}
	}

	return result, nil
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

	cacheKey := fmt.Sprintf(
		"history:user:%d:search:%s:page:%d:limit:%d",
		userID,
		search,
		page,
		limit,
	)

	// CACHE HIT
	if w.redis != nil {
		cachedData, err := w.redis.Get(ctx, cacheKey).Result()
		if err == nil {
			var response struct {
				Data  []repository.TransactionReportItem `json:"data"`
				Total int                                `json:"total"`
			}

			if json.Unmarshal([]byte(cachedData), &response) == nil {
				fmt.Println("CACHE HIT")
				return response.Data, response.Total, nil
			}
		}
	}

	fmt.Println("CACHE MISS")

	offset := (page - 1) * limit

	items, total, err := w.repo.GetTransactionHistory(ctx, userID, search, limit, offset)
	if err != nil {
		return nil, 0, err
	}

	// SAVE CACHE
	if w.redis != nil {
		cachePayload := struct {
			Data  []repository.TransactionReportItem `json:"data"`
			Total int                                `json:"total"`
		}{
			Data:  items,
			Total: total,
		}

		jsonData, _ := json.Marshal(cachePayload)
		w.redis.Set(ctx, cacheKey, jsonData, 5*time.Minute)
	}

	return items, total, nil
}
