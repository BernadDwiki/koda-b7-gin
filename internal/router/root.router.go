package router

import (
	"github.com/bernaddwiki/koda-b7-weekly10/internal/controller"
	"github.com/bernaddwiki/koda-b7-weekly10/internal/middleware"
	"github.com/bernaddwiki/koda-b7-weekly10/internal/repository"
	"github.com/bernaddwiki/koda-b7-weekly10/internal/service"
	"github.com/gin-gonic/gin"
	"github.com/redis/go-redis/v9"
)

func RegisterRootRouter(r gin.IRouter, redisClient *redis.Client) {
	if authRepository == nil || authController == nil || sharedDB == nil {
		panic("RegisterAuthRouter must be called before RegisterRootRouter")
	}

	walletRepo := repository.NewWalletRepository(sharedDB)
	userRepo := repository.NewUserRepository(sharedDB)
	walletService := service.NewWalletService(walletRepo, userRepo, sharedDB, redisClient)
	walletController := controller.NewWalletController(walletService)

	protected := r.Group("/")
	protected.Use(middleware.AuthMiddleware(authRepository, redisClient))
	{
		protected.DELETE("auth/logout", authController.Logout)
		RegisterWalletRoutes(protected, walletController)
	}
}
