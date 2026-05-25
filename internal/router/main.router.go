package router

import (
	"net/http"

	"github.com/bernaddwiki/koda-b7-weekly10/internal/dto"
	"github.com/bernaddwiki/koda-b7-weekly10/internal/middleware"
	"github.com/gin-gonic/gin"
	"github.com/jackc/pgx/v5/pgxpool"

	swaggerFiles "github.com/swaggo/files"
	ginSwagger "github.com/swaggo/gin-swagger"
)

func SetupRouter(r *gin.Engine, db *pgxpool.Pool) {
	r.Use(middleware.CORSMiddleware)
	r.Static("/img", "public/profile")

	r.GET("/swagger/*any",
		ginSwagger.WrapHandler(swaggerFiles.Handler))

	RegisterAuthRouter(r, db)
	RegisterUserRouter(r, db)
	RegisterRootRouter(r)

	r.NoRoute(func(ctx *gin.Context) {
		ctx.JSON(http.StatusNotFound, dto.Response{
			Message: "invalid route not found",
			Success: false,
		})
	})
}
