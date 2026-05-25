package router

import (
	"github.com/bernaddwiki/koda-b7-weekly10/internal/controller"
	"github.com/bernaddwiki/koda-b7-weekly10/internal/middleware"
	"github.com/bernaddwiki/koda-b7-weekly10/internal/repository"
	"github.com/bernaddwiki/koda-b7-weekly10/internal/service"
	"github.com/gin-gonic/gin"
	"github.com/jackc/pgx/v5/pgxpool"
)

func RegisterUserRouter(r gin.IRouter, db *pgxpool.Pool) {
	if authRepository == nil {
		panic("RegisterAuthRouter must be called before RegisterUserRouter")
	}

	userRepo := repository.NewUserRepository(db)
	userService := service.NewUserService(userRepo)
	userController := controller.NewUserController(userService)

	protected := r.Group("/")
	protected.Use(middleware.AuthMiddleware(authRepository))
	{
		user := protected.Group("/user")
		{
			user.GET("/profile", userController.GetProfile)
			user.PUT("/profile", userController.EditProfile)
			user.PATCH("/profile", userController.EditProfile)
			user.PUT("/password", userController.ChangePassword)
			user.POST("/set-pin", userController.SetPin)
			user.POST("/check-pin", userController.CheckPin)
			user.PUT("/edit-pin", userController.EditPin)
			user.GET("/receivers", userController.FindReceivers)
		}
	}
}
