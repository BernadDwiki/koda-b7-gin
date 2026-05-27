package router

import (
	"github.com/bernaddwiki/koda-b7-weekly10/internal/controller"
	"github.com/gin-gonic/gin"
)

func RegisterWalletRoutes(r gin.IRouter, walletController *controller.WalletController) {
	r.GET("/wallet/dashboard", walletController.Dashboard)
	r.GET("/wallet/transactions", walletController.TransactionHistory)
	r.GET("/wallet/transaction-report", walletController.TransactionReport)
	r.POST("/wallet/transfer", walletController.CreateTransfer)
	r.POST("/wallet/top-up", walletController.CreateTopUp)
}
