package middleware

import (
	"net/http"
	"strings"

	"github.com/bernaddwiki/koda-b7-weekly10/internal/dto"
	"github.com/bernaddwiki/koda-b7-weekly10/internal/repository"
	"github.com/bernaddwiki/koda-b7-weekly10/internal/util"
	"github.com/gin-gonic/gin"
)

func AuthMiddleware(
	authRepo repository.IAuthRepository,
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

		isRevoked, err := authRepo.IsTokenRevoked(
			ctx.Request.Context(),
			token,
		)

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

		if isRevoked {
			ctx.AbortWithStatusJSON(
				http.StatusUnauthorized,
				dto.Response{
					Success: false,
					Message: "token revoked",
				},
			)
			return
		}

		claims, err := util.VerifyToken(token)
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
