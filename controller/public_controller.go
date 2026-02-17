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
	service services.ContactService
}

func NewPublicController(service services.ContactService) *PublicController {
	return &PublicController{service}
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
