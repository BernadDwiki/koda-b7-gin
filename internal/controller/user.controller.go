package controller

import (
	"net/http"
	"strconv"

	"github.com/bernaddwiki/koda-b7-weekly10/internal/dto"
	"github.com/bernaddwiki/koda-b7-weekly10/internal/errText"
	"github.com/bernaddwiki/koda-b7-weekly10/internal/jwt"
	"github.com/bernaddwiki/koda-b7-weekly10/internal/service"
	"github.com/gin-gonic/gin"
	"github.com/gin-gonic/gin/binding"
)

type UserController struct {
	service service.IUserService
}

func NewUserController(
	service service.IUserService,
) *UserController {
	return &UserController{service}
}

// GetProfile godoc
// @Summary Get Profile
// @Description Get user profile
// @Tags User
// @Accept json
// @Produce json
// @Security BearerAuth
// @Success 200 {object} dto.Response
// @Failure 401 {object} dto.Response
// @Failure 500 {object} dto.Response
// @Router /user/profile [get]
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

	claims := claimsRaw.(*jwt.JWTClaims)

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

// SetPin godoc
// @Summary Set PIN
// @Description Set a new user PIN
// @Tags User
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param request body dto.SetPinRequest true "Set PIN Request"
// @Success 200 {object} dto.Response
// @Failure 422 {object} dto.Response
// @Failure 401 {object} dto.Response
// @Router /user/set-pin [post]
func (u *UserController) SetPin(
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

	claims := claimsRaw.(*jwt.JWTClaims)

	var body dto.SetPinRequest

	if err := ctx.ShouldBindJSON(&body); err != nil {
		ctx.JSON(http.StatusUnprocessableEntity, dto.Response{
			Success: false,
			Message: errText.GetValidationErrorMessage(err),
		})
		return
	}

	err := u.service.SetPin(
		ctx.Request.Context(),
		claims.UserID,
		body.Pin,
	)

	if err != nil {
		ctx.JSON(http.StatusUnprocessableEntity, dto.Response{
			Success: false,
			Message: err.Error(),
		})
		return
	}

	ctx.JSON(http.StatusOK, dto.Response{
		Success: true,
		Message: "set pin success",
	})
}

// CheckPin godoc
// @Summary Check PIN
// @Description Verify current user PIN
// @Tags User
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param request body dto.CheckPinRequest true "Check PIN Request"
// @Success 200 {object} dto.Response
// @Failure 422 {object} dto.Response
// @Failure 401 {object} dto.Response
// @Router /user/check-pin [post]
func (u *UserController) CheckPin(
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

	claims := claimsRaw.(*jwt.JWTClaims)

	var body dto.CheckPinRequest
	if err := ctx.ShouldBindJSON(&body); err != nil {
		ctx.JSON(http.StatusUnprocessableEntity, dto.Response{
			Success: false,
			Message: errText.GetValidationErrorMessage(err),
		})
		return
	}

	valid, err := u.service.CheckPin(ctx.Request.Context(), claims.UserID, body.Pin)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, dto.Response{Success: false, Message: err.Error()})
		return
	}

	if !valid {
		ctx.JSON(http.StatusUnauthorized, dto.Response{Success: false, Message: "invalid pin"})
		return
	}

	ctx.JSON(http.StatusOK, dto.Response{Success: true, Message: "pin valid"})
}

// EditPin godoc
// @Summary Edit PIN
// @Description Update current user PIN
// @Tags User
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param request body dto.EditPinRequest true "Edit PIN Request"
// @Success 200 {object} dto.Response
// @Failure 422 {object} dto.Response
// @Failure 401 {object} dto.Response
// @Router /user/edit-pin [put]
func (u *UserController) EditPin(
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

	claims := claimsRaw.(*jwt.JWTClaims)

	var body dto.EditPinRequest
	if err := ctx.ShouldBindJSON(&body); err != nil {
		ctx.JSON(http.StatusUnprocessableEntity, dto.Response{
			Success: false,
			Message: errText.GetValidationErrorMessage(err),
		})
		return
	}

	if err := u.service.EditPin(ctx.Request.Context(), claims.UserID, body); err != nil {
		ctx.JSON(http.StatusUnprocessableEntity, dto.Response{Success: false, Message: err.Error()})
		return
	}

	ctx.JSON(http.StatusOK, dto.Response{Success: true, Message: "pin updated"})
}

// ChangePassword godoc
// @Summary Change Password
// @Description Change user password
// @Tags User
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param request body dto.ChangePasswordRequest true "Change Password Request"
// @Success 200 {object} dto.Response
// @Failure 422 {object} dto.Response
// @Failure 401 {object} dto.Response
// @Router /user/password [put]
func (u *UserController) ChangePassword(
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

	claims := claimsRaw.(*jwt.JWTClaims)

	var body dto.ChangePasswordRequest
	if err := ctx.ShouldBindJSON(&body); err != nil {
		ctx.JSON(http.StatusUnprocessableEntity, dto.Response{
			Success: false,
			Message: errText.GetValidationErrorMessage(err),
		})
		return
	}

	if err := u.service.ChangePassword(ctx.Request.Context(), claims.UserID, body); err != nil {
		ctx.JSON(http.StatusUnprocessableEntity, dto.Response{Success: false, Message: err.Error()})
		return
	}

	ctx.JSON(http.StatusOK, dto.Response{Success: true, Message: "password updated"})
}

// FindReceivers godoc
// @Summary Find Receivers
// @Description Search for receiver users by name or email
// @Tags User
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param search query string false "Search term"
// @Param page query int false "Page number"
// @Param limit query int false "Limit per page"
// @Success 200 {object} dto.Response
// @Failure 422 {object} dto.Response
// @Failure 401 {object} dto.Response
// @Router /user/receivers [get]
func (u *UserController) FindReceivers(
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

	claims := claimsRaw.(*jwt.JWTClaims)

	search := ctx.Query("search")
	pageQ := ctx.DefaultQuery("page", "1")
	limitQ := ctx.DefaultQuery("limit", "10")

	page, err := strconv.Atoi(pageQ)
	if err != nil || page < 1 {
		ctx.JSON(http.StatusUnprocessableEntity, dto.Response{Success: false, Message: "page must be a positive integer"})
		return
	}

	limit, err := strconv.Atoi(limitQ)
	if err != nil || limit < 1 || limit > 100 {
		ctx.JSON(http.StatusUnprocessableEntity, dto.Response{Success: false, Message: "limit must be a positive integer between 1 and 100"})
		return
	}

	result, err := u.service.FindReceivers(ctx.Request.Context(), claims.UserID, search, page, limit)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, dto.Response{Success: false, Message: err.Error()})
		return
	}

	ctx.JSON(http.StatusOK, dto.Response{Success: true, Message: "receivers fetched", Data: result})
}

// EditProfile godoc
// @Summary Update Profile
// @Description Update user profile including optional profile picture upload
// @Tags User
// @Accept multipart/form-data
// @Produce json
// @Security BearerAuth
// @Param name formData string false "Name"
// @Param phone_number formData string false "Phone Number"
// @Param profile_picture formData file false "Profile Picture"
// @Success 200 {object} dto.Response
// @Failure 422 {object} dto.Response
// @Failure 401 {object} dto.Response
// @Failure 500 {object} dto.Response
// @Router /user/profile [patch]
func (u *UserController) EditProfile(
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

	claims := claimsRaw.(*jwt.JWTClaims)

	var body dto.EditProfileRequest
	if err := ctx.ShouldBindWith(&body, binding.FormMultipart); err != nil {
		ctx.JSON(http.StatusUnprocessableEntity, dto.Response{Success: false, Message: errText.GetValidationErrorMessage(err)})
		return
	}

	profile, err := u.service.UpdateProfile(ctx.Request.Context(), claims.UserID, body)
	if err != nil {
		ctx.JSON(http.StatusUnprocessableEntity, dto.Response{Success: false, Message: err.Error()})
		return
	}

	ctx.JSON(http.StatusOK, dto.Response{Success: true, Message: "profile updated", Data: profile})
}
