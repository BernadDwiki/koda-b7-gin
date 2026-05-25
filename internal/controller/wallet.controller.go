package controller

import (
	"errors"
	"net/http"
	"strconv"
	"time"

	"github.com/bernaddwiki/koda-b7-weekly10/internal/dto"
	"github.com/bernaddwiki/koda-b7-weekly10/internal/jwt"
	"github.com/bernaddwiki/koda-b7-weekly10/internal/service"
	"github.com/gin-gonic/gin"
)

type WalletController struct {
	service service.IWalletService
}

func NewWalletController(service service.IWalletService) *WalletController {
	return &WalletController{service}
}

// Dashboard godoc
// @Summary Wallet Dashboard
// @Description Get authenticated user wallet summary data
// @Tags Wallet
// @Accept json
// @Produce json
// @Security BearerAuth
// @Success 200 {object} dto.Response
// @Failure 401 {object} dto.Response
// @Failure 500 {object} dto.Response
// @Router /wallet/dashboard [get]
func (w *WalletController) Dashboard(ctx *gin.Context) {
	claimsRaw, _ := ctx.Get("claims")
	claims := claimsRaw.(*jwt.JWTClaims)

	data, err := w.service.GetDashboard(
		ctx.Request.Context(),
		claims.UserID,
	)

	if err != nil {
		ctx.JSON(http.StatusInternalServerError, dto.Response{
			Success: false,
			Message: err.Error(),
		})
		return
	}

	ctx.JSON(http.StatusOK, dto.Response{
		Success: true,
		Message: "get dashboard success",
		Data:    data,
	})
}

// CreateTransfer godoc
// @Summary Create Transfer
// @Description Create a wallet transfer to another user
// @Tags Wallet
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param request body dto.CreateTransferRequest true "Transfer Request"
// @Success 200 {object} dto.Response
// @Failure 400 {object} dto.Response
// @Failure 401 {object} dto.Response
// @Failure 404 {object} dto.Response
// @Failure 500 {object} dto.Response
// @Router /wallet/transfer [post]
func (w *WalletController) CreateTransfer(ctx *gin.Context) {
	claimsRaw, _ := ctx.Get("claims")
	claims := claimsRaw.(*jwt.JWTClaims)

	var request dto.CreateTransferRequest
	if err := ctx.ShouldBindJSON(&request); err != nil {
		ctx.JSON(http.StatusBadRequest, dto.Response{Success: false, Message: err.Error()})
		return
	}

	response, err := w.service.CreateTransfer(ctx.Request.Context(), claims.UserID, request)
	if err != nil {
		if errors.Is(err, service.ErrInvalidTransferRequest) {
			ctx.JSON(http.StatusBadRequest, dto.Response{Success: false, Message: err.Error()})
			return
		}
		if errors.Is(err, service.ErrInvalidReceiver) {
			ctx.JSON(http.StatusBadRequest, dto.Response{Success: false, Message: err.Error()})
			return
		}
		if errors.Is(err, service.ErrInvalidPin) || errors.Is(err, service.ErrPinNotSet) {
			ctx.JSON(http.StatusBadRequest, dto.Response{Success: false, Message: err.Error()})
			return
		}
		if errors.Is(err, service.ErrReceiverNotFound) {
			ctx.JSON(http.StatusNotFound, dto.Response{Success: false, Message: err.Error()})
			return
		}
		if errors.Is(err, service.ErrInsufficientBalance) {
			ctx.JSON(http.StatusBadRequest, dto.Response{Success: false, Message: err.Error()})
			return
		}
		if errors.Is(err, service.ErrEwalletNotFound) {
			ctx.JSON(http.StatusNotFound, dto.Response{Success: false, Message: err.Error()})
			return
		}

		ctx.JSON(http.StatusInternalServerError, dto.Response{Success: false, Message: err.Error()})
		return
	}

	ctx.JSON(http.StatusOK, dto.Response{Success: true, Message: "transfer created", Data: response})
}

// CreateTopUp godoc
// @Summary Create Top Up
// @Description Create a wallet top-up request using payment method
// @Tags Wallet
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param request body dto.CreateTopUpRequest true "Top Up Request"
// @Success 200 {object} dto.Response
// @Failure 400 {object} dto.Response
// @Failure 404 {object} dto.Response
// @Failure 500 {object} dto.Response
// @Router /wallet/top-up [post]
func (w *WalletController) CreateTopUp(ctx *gin.Context) {
	claimsRaw, _ := ctx.Get("claims")
	claims := claimsRaw.(*jwt.JWTClaims)

	var request dto.CreateTopUpRequest
	if err := ctx.ShouldBindJSON(&request); err != nil {
		ctx.JSON(http.StatusBadRequest, dto.Response{Success: false, Message: err.Error()})
		return
	}

	response, err := w.service.CreateTopUp(ctx.Request.Context(), claims.UserID, request)
	if err != nil {
		if errors.Is(err, service.ErrInvalidTopUpRequest) {
			ctx.JSON(http.StatusBadRequest, dto.Response{Success: false, Message: err.Error()})
			return
		}
		if errors.Is(err, service.ErrPaymentMethodNotFound) {
			ctx.JSON(http.StatusNotFound, dto.Response{Success: false, Message: err.Error()})
			return
		}
		if errors.Is(err, service.ErrEwalletNotFound) {
			ctx.JSON(http.StatusNotFound, dto.Response{Success: false, Message: err.Error()})
			return
		}

		ctx.JSON(http.StatusInternalServerError, dto.Response{Success: false, Message: err.Error()})
		return
	}

	ctx.JSON(http.StatusOK, dto.Response{Success: true, Message: "top up created", Data: response})
}

// TransactionsReport godoc
// @Summary Transaction Report
// @Description Get transaction report or paginated history for authenticated user
// @Tags Wallet
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param days query int false "Number of days to include in report" default(7)
// @Param flow query string false "Filter by flow: income, expense, or both" default(both)
// @Param search query string false "Search term for transaction history"
// @Param page query int false "Page number for paginated history"
// @Param limit query int false "Limit per page for paginated history"
// @Success 200 {object} dto.Response
// @Failure 400 {object} dto.Response
// @Failure 401 {object} dto.Response
// @Failure 500 {object} dto.Response
// @Router /wallet/transactions [get]
// TransactionsReport returns transaction list between start and end query params (RFC3339).
// If not provided, defaults to last 7 days.
func (w *WalletController) TransactionsReport(ctx *gin.Context) {
	claimsRaw, _ := ctx.Get("claims")
	claims := claimsRaw.(*jwt.JWTClaims)

	// Use `days` query param to specify how many days back to include (default 7)
	daysQ := ctx.DefaultQuery("days", "7")
	days, err := strconv.Atoi(daysQ)
	if err != nil || days <= 0 {
		ctx.JSON(http.StatusBadRequest, dto.Response{Success: false, Message: "invalid days value, must be positive integer"})
		return
	}

	end := time.Now().UTC()
	start := end.AddDate(0, 0, -days)

	flow := ctx.DefaultQuery("flow", "both")
	if flow != "income" && flow != "expense" && flow != "both" {
		ctx.JSON(http.StatusBadRequest, dto.Response{Success: false, Message: "invalid flow value, must be 'income', 'expense' or 'both'"})
		return
	}

	// If pagination params are provided, treat this as a paginated/search request
	pageQ := ctx.Query("page")
	limitQ := ctx.Query("limit")
	search := ctx.DefaultQuery("search", "")

	if pageQ != "" || limitQ != "" || search != "" {
		// parse page and limit
		page := 1
		limit := 10
		if pageQ != "" {
			if p, err := strconv.Atoi(pageQ); err == nil && p > 0 {
				page = p
			}
		}
		if limitQ != "" {
			if l, err := strconv.Atoi(limitQ); err == nil && l > 0 {
				limit = l
			}
		}

		items, total, err := w.service.GetTransactionHistory(ctx.Request.Context(), claims.UserID, search, page, limit)
		if err != nil {
			ctx.JSON(http.StatusInternalServerError, dto.Response{Success: false, Message: err.Error()})
			return
		}

		resp := map[string]interface{}{
			"items": items,
			"page":  page,
			"limit": limit,
			"total": total,
		}

		ctx.JSON(http.StatusOK, dto.Response{Success: true, Message: "transaction history fetched", Data: resp})
		return
	}

	items, err := w.service.GetTransactionReport(ctx.Request.Context(), claims.UserID, start.Format(time.RFC3339), end.Format(time.RFC3339), flow)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, dto.Response{Success: false, Message: err.Error()})
		return
	}

	ctx.JSON(http.StatusOK, dto.Response{Success: true, Message: "transaction report fetched", Data: items})
}
