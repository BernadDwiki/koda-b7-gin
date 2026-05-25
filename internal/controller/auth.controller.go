package controller

import (
	"net/http"
	"strings"

	"github.com/bernaddwiki/koda-b7-weekly10/internal/dto"
	"github.com/bernaddwiki/koda-b7-weekly10/internal/errText"
	"github.com/bernaddwiki/koda-b7-weekly10/internal/jwt"
	"github.com/bernaddwiki/koda-b7-weekly10/internal/service"
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

// Register godoc
// @Summary Register User
// @Description Register a new user using name, email, and password
// @Tags Auth
// @Accept json
// @Produce json
// @Param request body dto.RegisterRequest true "Register Request"
// @Success 201 {object} dto.Response
// @Failure 400 {object} dto.Response
// @Router /auth/register [post]
func (a *AuthController) Register(
	ctx *gin.Context,
) {
	var body dto.RegisterRequest

	if err := ctx.ShouldBindJSON(&body); err != nil {
		ctx.JSON(http.StatusBadRequest, dto.Response{
			Success: false,
			Message: errText.GetValidationErrorMessage(err),
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
		Data: dto.RegisterResponse{
			ID:    user.ID,
			Email: user.Email,
		},
	})
}

// Login godoc
// @Summary Login User
// @Description Login menggunakan email dan password
// @Tags Auth
// @Accept json
// @Produce json
// @Param request body dto.LoginRequest true "Login Request"
// @Success 200 {object} dto.Response
// @Failure 400 {object} dto.Response
// @Failure 401 {object} dto.Response
// @Router /auth/login [post]
func (a *AuthController) Login(
	ctx *gin.Context,
) {
	var body dto.LoginRequest

	if err := ctx.ShouldBindJSON(&body); err != nil {
		ctx.JSON(http.StatusBadRequest, dto.Response{
			Success: false,
			Message: errText.GetValidationErrorMessage(err),
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

// Logout godoc
// @Summary Logout User
// @Description Invalidate current JWT token
// @Tags Auth
// @Accept json
// @Produce json
// @Security BearerAuth
// @Success 200 {object} dto.Response
// @Failure 401 {object} dto.Response
// @Failure 500 {object} dto.Response
// @Router /auth/logout [post]
func (a *AuthController) Logout(
	ctx *gin.Context,
) {
	authHeader := ctx.GetHeader("Authorization")
	token := strings.TrimPrefix(authHeader, "Bearer ")

	claimsRaw, _ := ctx.Get("claims")
	claims := claimsRaw.(*jwt.JWTClaims)

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
