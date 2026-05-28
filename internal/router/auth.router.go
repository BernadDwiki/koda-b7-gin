package router

import (
	"github.com/bernaddwiki/koda-b7-weekly10/internal/controller"
	"github.com/bernaddwiki/koda-b7-weekly10/internal/repository"
	"github.com/bernaddwiki/koda-b7-weekly10/internal/service"
	"github.com/gin-gonic/gin"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/redis/go-redis/v9"
)

var (
	authRepository repository.IAuthRepository
	authController *controller.AuthController
	sharedDB       *pgxpool.Pool
)

func RegisterAuthRouter(r gin.IRouter, db *pgxpool.Pool, redisClient *redis.Client) {
	sharedDB = db
	authRepository = repository.NewAuthRepository(db)
	authService := service.NewAuthService(authRepository, redisClient)
	authController = controller.NewAuthController(authService)

	auth := r.Group("/auth")
	{
		auth.POST("/register", authController.Register)
		auth.POST("/login", authController.Login)

		auth.POST("/forgot-password", authController.ForgotPassword)
		auth.POST("/reset-password", authController.ResetPassword)
	}
}
