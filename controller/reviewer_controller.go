package controller

import (
	"dhl/constant"
	"dhl/models"
	"dhl/services"
	"dhl/utils"
	"net/http"

	"github.com/gin-gonic/gin"
)

type ReviewerController struct {
	assessmentService services.AssessmentService
	
}

func NewReviewerController(assessmentService services.AssessmentService) *ReviewerController {
	return &ReviewerController{
		assessmentService: assessmentService,
	}
}

func (rc *ReviewerController) GetReviewerAssessments(ctx *gin.Context) {

	page, limit, offset := utils.GetPaginationParams(ctx)

	reviewerID := ctx.Query("user_id")

	if reviewerID == "" {
		models.ErrorResponse(ctx, constant.Failure, http.StatusBadRequest, "user_id is required", nil, nil)
		return
	}

	assessments, total, err := rc.assessmentService.GetReviewerAssessments(
		reviewerID,
		limit,
		offset,
	)

	if err != nil {
		models.ErrorResponse(ctx, constant.Failure, http.StatusInternalServerError, "Failed to fetch reviewer assessments", nil, err)
		return
	}

	pagination := utils.GetPagination(limit, page, offset, total)

	models.SuccessResponse(ctx, constant.Success, http.StatusOK, "reviewer assessments", assessments, pagination, nil)
}


func (rc *ReviewerController) ReviewAssessment(ctx *gin.Context) {

	var req models.ReviewAssessmentRequest

	if err := ctx.ShouldBindJSON(&req); err != nil {
		models.ErrorResponse(ctx, constant.Failure, http.StatusBadRequest, "Invalid request", nil, err)
		return
	}

	sub, exists := ctx.Get("sub")
	if !exists {
		models.ErrorResponse(ctx, constant.Failure, http.StatusUnauthorized, "user not found", nil, nil)
		return
	}

	reviewerID := sub.(string)

	err := rc.assessmentService.ReviewAssessment(
		reviewerID,
		req.AssessmentSequence,
		req.Status,
		req.Note,
	)

	if err != nil {
		models.ErrorResponse(ctx, constant.Failure, http.StatusInternalServerError, "Failed to update review status", nil, err)
		return
	}

	models.SuccessResponse(ctx, constant.Success, http.StatusOK, "assessment review updated", nil, nil, nil)
}