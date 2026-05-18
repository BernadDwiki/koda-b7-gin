package controller

import (
	"net/http"

	"github.com/bernaddwiki/koda-b7-weekly10/internal/dto"
	"github.com/bernaddwiki/koda-b7-weekly10/internal/service"
	"github.com/bernaddwiki/koda-b7-weekly10/internal/util"
	"github.com/gin-gonic/gin"
)

type UserController struct {
	service service.IUserService
}

func NewUserController(
	service service.IUserService,
) *UserController {
	return &UserController{service}
}

func (u *UserController) GetProfile(
	ctx *gin.Context,
) {
	claimsRaw, exists := ctx.Get("claims")
	if !exists {
		ctx.JSON(http.StatusUnauthorized, dto.Response{
			Success: false,
			Message: "unauthorized",
		})
		return
	}

	claims := claimsRaw.(*util.JWTClaims)

	profile, err := u.service.GetProfile(
		ctx.Request.Context(),
		claims.UserID,
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
		Message: "get profile success",
		Data:    profile,
	})
}
