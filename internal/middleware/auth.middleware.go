package middleware

import (
	"net/http"
	"strings"

	"github.com/bernaddwiki/koda-b7-weekly10/internal/dto"
	"github.com/bernaddwiki/koda-b7-weekly10/internal/jwt"
	"github.com/bernaddwiki/koda-b7-weekly10/internal/repository"
	"github.com/gin-gonic/gin"
	"github.com/redis/go-redis/v9"
)

func AuthMiddleware(
	authRepo repository.IAuthRepository,
	redisClient *redis.Client,
) gin.HandlerFunc {
	return func(ctx *gin.Context) {

		authHeader := ctx.GetHeader("Authorization")

		if authHeader == "" {
			ctx.AbortWithStatusJSON(
				http.StatusUnauthorized,
				dto.Response{
					Success: false,
					Message: "authorization required",
				},
			)
			return
		}

		token := strings.TrimPrefix(
			authHeader,
			"Bearer ",
		)

		exists, err := redisClient.Exists(
			ctx.Request.Context(),
			"blacklist:"+token,
		).Result()

		if err != nil {
			ctx.AbortWithStatusJSON(
				http.StatusInternalServerError,
				dto.Response{
					Success: false,
					Message: err.Error(),
				},
			)
			return
		}

		if exists == 1 {
			ctx.AbortWithStatusJSON(
				http.StatusUnauthorized,
				dto.Response{
					Success: false,
					Message: "token revoked",
				},
			)
			return
		}

		claims, err := jwt.VerifyToken(token)
		if err != nil {
			ctx.AbortWithStatusJSON(
				http.StatusUnauthorized,
				dto.Response{
					Success: false,
					Message: "invalid token",
				},
			)
			return
		}

		ctx.Set("claims", claims)
		ctx.Next()
	}
}
