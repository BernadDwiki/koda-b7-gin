package router

import (
	"net/http"

	"github.com/bernaddwiki/koda-b7-weekly10/internal/dto"
	"github.com/bernaddwiki/koda-b7-weekly10/internal/middleware"
	"github.com/gin-gonic/gin"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/redis/go-redis/v9"

	swaggerFiles "github.com/swaggo/files"
	ginSwagger "github.com/swaggo/gin-swagger"
)

func SetupRouter(r *gin.Engine, db *pgxpool.Pool, redisClient *redis.Client) {
	r.Use(middleware.CORSMiddleware)
	r.Static("/img", "public/profile")

	r.GET("/swagger/*any",
		ginSwagger.WrapHandler(swaggerFiles.Handler))

	RegisterAuthRouter(r, db, redisClient)
	RegisterUserRouter(r, db, redisClient)
	RegisterRootRouter(r, redisClient)

	r.NoRoute(func(ctx *gin.Context) {
		ctx.JSON(http.StatusNotFound, dto.Response{
			Message: "invalid route not found",
			Success: false,
		})
	})
}
