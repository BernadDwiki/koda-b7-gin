package router

import (
	"github.com/bernaddwiki/koda-b7-weekly10/internal/controller"
	"github.com/gin-gonic/gin"
)

func RegisterWalletRoutes(r gin.IRouter, walletController *controller.WalletController) {
	wallet := r.Group("/wallet")
	{
		wallet.GET("/dashboard", walletController.Dashboard)
		wallet.GET("/transactions", walletController.TransactionHistory)
		wallet.GET("/transaction-report", walletController.TransactionReport)
		wallet.POST("/transfer", walletController.CreateTransfer)
		wallet.POST("/top-up", walletController.CreateTopUp)
	}
}
