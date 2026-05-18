package controller

import (
	"net/http"

	"github.com/bernaddwiki/koda-b7-weekly10/internal/dto"
	"github.com/bernaddwiki/koda-b7-weekly10/internal/service"
	"github.com/bernaddwiki/koda-b7-weekly10/internal/util"
	"github.com/gin-gonic/gin"
)

type WalletController struct {
	service service.IWalletService
}

func NewWalletController(service service.IWalletService) *WalletController {
	return &WalletController{service}
}

func (w *WalletController) Dashboard(ctx *gin.Context) {
	claimsRaw, _ := ctx.Get("claims")
	claims := claimsRaw.(*util.JWTClaims)

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
