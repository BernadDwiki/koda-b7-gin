package controller

import (
	"net/http"
	"strings"

	"github.com/bernaddwiki/koda-b7-weekly10/internal/dto"
	"github.com/bernaddwiki/koda-b7-weekly10/internal/service"
	"github.com/bernaddwiki/koda-b7-weekly10/internal/util"
	"github.com/gin-gonic/gin"
)

type AuthController struct {
	service service.IAuthService
}

func NewAuthController(
	service service.IAuthService,
) *AuthController {
	return &AuthController{service}
}

func (a *AuthController) Register(
	ctx *gin.Context,
) {
	var body dto.RegisterRequest

	if err := ctx.ShouldBindJSON(&body); err != nil {
		//dicabagkan lgi errornya, jangan pakai error mentahan
		ctx.JSON(http.StatusBadRequest, dto.Response{
			Success: false,
			Message: err.Error(),
		})
		return
	}

	user, err := a.service.Register(
		ctx.Request.Context(),
		body,
	)

	if err != nil {
		ctx.JSON(http.StatusBadRequest, dto.Response{
			Success: false,
			Message: err.Error(),
		})
		return
	}

	ctx.JSON(http.StatusCreated, dto.Response{
		Success: true,
		Message: "register success",
		Data:    user,
	})
}

func (a *AuthController) Login(
	ctx *gin.Context,
) {
	var body dto.LoginRequest

	if err := ctx.ShouldBindJSON(&body); err != nil {
		ctx.JSON(http.StatusBadRequest, dto.Response{
			Success: false,
			Message: err.Error(),
		})
		return
	}

	result, err := a.service.Login(
		ctx.Request.Context(),
		body,
	)

	if err != nil {
		ctx.JSON(http.StatusUnauthorized, dto.Response{
			Success: false,
			Message: err.Error(),
		})
		return
	}

	ctx.JSON(http.StatusOK, dto.Response{
		Success: true,
		Message: "login success",
		Data:    result,
	})
}

func (a *AuthController) Logout(
	ctx *gin.Context,
) {
	authHeader := ctx.GetHeader("Authorization")
	token := strings.TrimPrefix(authHeader, "Bearer ")

	claimsRaw, _ := ctx.Get("claims")
	claims := claimsRaw.(*util.JWTClaims)

	expiredAt := claims.ExpiresAt.Time

	err := a.service.Logout(
		ctx.Request.Context(),
		claims.UserID,
		token,
		expiredAt,
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
		Message: "logout success",
	})
}
