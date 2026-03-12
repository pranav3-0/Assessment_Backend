package controller

import (
	"dhl/constant"
	"dhl/models"
	"dhl/services"
	"fmt"
	"net/http"

	"github.com/gin-gonic/gin"
)

type PublicController struct {
	service     services.ContactService
	authService services.AuthService
}

func NewPublicController(
	service services.ContactService,
	authService services.AuthService,
) *PublicController {

	return &PublicController{
		service:     service,
		authService: authService,
	}
}

func (cc *PublicController) SubmitContactFormController(ctx *gin.Context) {
	var req models.ContactUsResponse

	if err := ctx.ShouldBindJSON(&req); err != nil {
		models.ErrorResponse(ctx, constant.Failure, http.StatusBadRequest, "invalid request", nil, err)
		return
	}

	if req.Name == "" || req.Email == "" || req.Subject == "" || req.Question == "" {
		models.ErrorResponse(ctx, constant.Failure, http.StatusBadRequest, "invalid request", nil, fmt.Errorf("Required fields missing"))
		return
	}

	err := cc.service.Submit(&req)
	if err != nil {
		models.ErrorResponse(ctx, constant.Failure, http.StatusInternalServerError, "error while submitting", nil, err)
		return
	}

	models.SuccessResponse(ctx, constant.Success, http.StatusOK, "Submitted successfully", nil, nil, nil)
	return
}

func (pc *PublicController) SendOTP(ctx *gin.Context) {

	var req models.SendOTPRequest

	if err := ctx.ShouldBindJSON(&req); err != nil {
		models.ErrorResponse(ctx, constant.Failure, http.StatusBadRequest, "Invalid request", nil, err)
		return
	}

	err := pc.authService.SendOTP(req.Email)
	if err != nil {
		models.ErrorResponse(ctx, constant.Failure, http.StatusInternalServerError, err.Error(), nil, err)
		return
	}

	models.SuccessResponse(ctx, constant.Success, http.StatusOK, "OTP sent successfully", nil, nil, nil)
}

func (pc *PublicController) VerifyOTP(ctx *gin.Context) {

	var req models.VerifyOTPRequest

	if err := ctx.ShouldBindJSON(&req); err != nil {
		models.ErrorResponse(ctx, constant.Failure, http.StatusBadRequest, "Invalid request", nil, err)
		return
	}

	token, userID, err := pc.authService.VerifyOTP(req.Email, req.OTP)
	if err != nil {
		models.ErrorResponse(ctx, constant.Failure, http.StatusUnauthorized, err.Error(), nil, err)
		return
	}

	res := models.VerifyOTPResponse{
		Token:  token,
		UserID: userID,
	}

	models.SuccessResponse(ctx, constant.Success, http.StatusOK, "OTP verified", res, nil, nil)
}

func (pc *PublicController) GetPublicToken(ctx *gin.Context) {
    token, err := services.GetPublicAssessmentToken()
    if err != nil {
        models.ErrorResponse(ctx, "failure", 500, "Failed to get public token", nil, err)
        return
    }
    models.SuccessResponse(ctx, "success", 200, "Token fetched", gin.H{"token": token}, nil, nil)
}