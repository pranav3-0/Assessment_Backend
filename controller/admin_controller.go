package controller

import (
	"dhl/constant"
	"dhl/models"
	"dhl/services"
	"dhl/utils"
	"log"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

type AdminController struct {
	userService         services.UserService
	authService         services.AuthService
	assessmentService   services.AssessmentService
	notificationService services.NotificationService
	contactService      services.ContactService
	jobService          services.JobDescriptionService
	questionService     services.QuestionService
}

func NewAdminController(userService services.UserService, authService services.AuthService, assessmentService services.AssessmentService, notificationService services.NotificationService, contactService services.ContactService, jobService services.JobDescriptionService, questionService services.QuestionService) *AdminController {
	return &AdminController{userService: userService, authService: authService, assessmentService: assessmentService, notificationService: notificationService, contactService: contactService, jobService: jobService, questionService: questionService}
}

func (uc *AdminController) GetAssessments(ctx *gin.Context) {
	page, limit, offset := utils.GetPaginationParams(ctx)
	assessmentSeq := ctx.Query("assessment_sequence")
	assessmentSession := ctx.Query("assessment_session")
	filterData := models.AssessmentFilter{}
	if assessmentSeq != "" {
		filterData.AssessmentSequence = &assessmentSeq
	}
	if assessmentSession != "" {
		filterData.AssessmentSessionId = &assessmentSession
	}

	assessments, totalRecords, err := uc.assessmentService.GetAssessments(limit, offset, &filterData)
	if err != nil {
		models.ErrorResponse(ctx, constant.Failure, http.StatusInternalServerError, "Failed to retrieve assessments", nil, err)
		return
	}
	pagination := utils.GetPagination(limit, page, offset, totalRecords)

	models.SuccessResponse(ctx, constant.Success, http.StatusOK, "assessments", assessments, pagination, nil)
	return
}

func (uc *AdminController) GetManagerAssessments(ctx *gin.Context) {
	_, userId, _, err := utils.GetUserIDFromContext(ctx, uc.userService.FindUserIdBySub)
	if err != nil {
		models.ErrorResponse(ctx, constant.Failure, http.StatusUnauthorized, err.Error(), nil, err)
		return
	}
	page, limit, offset := utils.GetPaginationParams(ctx)
	assessmentSeq := ctx.Query("assessment_sequence")
	assessmentSession := ctx.Query("assessment_session")
	filterData := models.AssessmentFilter{}
	if assessmentSeq != "" {
		filterData.AssessmentSequence = &assessmentSeq
	}
	if assessmentSession != "" {
		filterData.AssessmentSessionId = &assessmentSession
	}
	assessments, totalRecords, err := uc.assessmentService.GetManagerAssessments(userId, limit, offset, &filterData)
	if err != nil {
		models.ErrorResponse(ctx, constant.Failure, http.StatusInternalServerError, "Failed to retrieve assessments", nil, err)
		return
	}
	pagination := utils.GetPagination(limit, page, offset, totalRecords)

	models.SuccessResponse(ctx, constant.Success, http.StatusOK, "assessments", assessments, pagination, nil)
	return
}

func (uc *AdminController) UpdateAssessment(ctx *gin.Context) {
	var request models.UpdateAssessmentRequest
	if err := ctx.ShouldBindJSON(&request); err != nil {
		models.ErrorResponse(ctx, constant.Failure, http.StatusBadRequest, "Invalid request body", nil, err)
		return
	}

	if err := uc.assessmentService.UpdateAssessmentService(request); err != nil {
		models.ErrorResponse(ctx, constant.Failure, http.StatusInternalServerError, "failed to update assessment", nil, err)
		return
	}
	models.SuccessResponse(ctx, constant.Success, http.StatusOK, "assessment updated", nil, nil, nil)
	return
}

func (uc *AdminController) UpdateAssessmentStatusController(ctx *gin.Context) {
	var req models.UpdateAssessmentStatusRequest
	if err := ctx.ShouldBindJSON(&req); err != nil {
		models.ErrorResponse(ctx, constant.Failure, http.StatusBadRequest, "Invalid request body", nil, err)
		return
	}
	err := uc.assessmentService.UpdateAssessmentStatusService(req)
	if err != nil {
		models.ErrorResponse(ctx, constant.Failure, http.StatusInternalServerError, "Failed to update status", nil, err)
		return
	}
	models.SuccessResponse(ctx, constant.Success, http.StatusOK, "Status updated", nil, nil, nil)
	return
}

func (uc *AdminController) DistributeAssessmentToUserController(ctx *gin.Context) {
	var req models.DistributeAssessmentRequest
	if err := ctx.ShouldBindJSON(&req); err != nil {
		models.ErrorResponse(ctx, constant.Failure, http.StatusBadRequest, "Invalid request body", nil, err)
		return
	}
	err := uc.assessmentService.DistributeAssessmentUser(req.AssessmentSequence, req.UserIds)
	if err != nil {
		models.ErrorResponse(ctx, constant.Failure, http.StatusInternalServerError, "Failed to distribute", nil, err)
		return
	}

	go func(userIds []string, seq string) {
		if err := uc.notificationService.SendDistributeAssessmentMail(userIds, seq, false); err != nil {
			log.Println("Error while sending mails:", err)
		}
	}(req.UserIds, req.AssessmentSequence)

	models.SuccessResponse(ctx, constant.Success, http.StatusOK, "Assessments distributed", nil, nil, nil)
	return
}

func (uc *AdminController) DistributeAssessmentToManagerController(ctx *gin.Context) {
	var req models.DistributeAssessmentRequest
	if err := ctx.ShouldBindJSON(&req); err != nil {
		models.ErrorResponse(ctx, constant.Failure, http.StatusBadRequest, "Invalid request body", nil, err)
		return
	}
	err := uc.assessmentService.DistributeAssessmentManager(req.AssessmentSequence, req.UserIds)
	if err != nil {
		models.ErrorResponse(ctx, constant.Failure, http.StatusInternalServerError, "Failed to distribute", nil, err)
		return
	}
	go func(userIds []string, seq string) {
		if err := uc.notificationService.SendDistributeAssessmentMail(userIds, seq, true); err != nil {
			log.Println("Error while sending mails:", err)
		}
	}(req.UserIds, req.AssessmentSequence)
	models.SuccessResponse(ctx, constant.Success, http.StatusOK, "Assessments distributed", nil, nil, nil)
	return
}

func (uc *AdminController) MapUsersToManagerController(ctx *gin.Context) {
	_, userId, _, err := utils.GetUserIDFromContext(ctx, uc.userService.FindUserIdBySub)
	if err != nil {
		models.ErrorResponse(ctx, constant.Failure, http.StatusUnauthorized, err.Error(), nil, err)
		return
	}
	var req models.MapUsersToManagerRequest
	if err := ctx.ShouldBindJSON(&req); err != nil {
		models.ErrorResponse(ctx, constant.Failure, http.StatusBadRequest, "Invalid request body", nil, err)
		return
	}
	err = uc.userService.MapUsersToManager(req, userId)
	if err != nil {
		models.ErrorResponse(ctx, constant.Failure, http.StatusInternalServerError, "Failed to add mapping", nil, err)
		return
	}
	models.SuccessResponse(ctx, constant.Success, http.StatusOK, "Users assigned to manager", nil, nil, nil)
	return
}

func (uc *AdminController) GetQuestionsController(ctx *gin.Context) {
	page, limit, offset := utils.GetPaginationParams(ctx)
	questions, totalRecords, err := uc.assessmentService.GetQuestions(limit, offset)
	if err != nil {
		models.ErrorResponse(ctx, constant.Failure, http.StatusInternalServerError, "Failed to retrieve questions", nil, err)
		return
	}
	pagination := utils.GetPagination(limit, page, offset, totalRecords)

	models.SuccessResponse(ctx, constant.Success, http.StatusOK, "questions", questions, pagination, nil)
	return
}

func (uc *AdminController) GetUsers(ctx *gin.Context) {
	highestRole, userId, _, err := utils.GetUserIDFromContext(ctx, uc.userService.FindUserIdBySub)
	if err != nil {
		models.ErrorResponse(ctx, constant.Failure, http.StatusUnauthorized, err.Error(), nil, err)
		return
	}
	ctx.Get("roles")
	page, limit, offset := utils.GetPaginationParams(ctx)
	userRole := ctx.Query("role")
	response, totalRecords, err := uc.userService.GetAllUsers(limit, offset, &userRole, highestRole, userId)
	if err != nil {
		models.ErrorResponse(ctx, constant.Failure, http.StatusInternalServerError, "Failed to load users", nil, err)
		return
	}
	pagination := utils.GetPagination(limit, page, offset, totalRecords)
	models.SuccessResponse(ctx, constant.Success, http.StatusOK, "users loaded", response, pagination, nil)
	return
}

func (uc *AdminController) UpdateUserProfile(ctx *gin.Context) {
	userIDStr := ctx.Query("id")
	userID, err := uuid.Parse(userIDStr)
	if err != nil {
		models.ErrorResponse(ctx, constant.Failure, http.StatusBadRequest, "Invalid user ID", nil, err)
		return
	}

	var req models.UserProfileUpdate
	if err := ctx.ShouldBindJSON(&req); err != nil {
		models.ErrorResponse(ctx, constant.Failure, http.StatusBadRequest, "Invalid request body", nil, err)
		return
	}

	if err := uc.userService.UpdateUserProfile(userID, req); err != nil {
		models.ErrorResponse(ctx, constant.Failure, http.StatusInternalServerError, "Failed to update profile", nil, err)
		return
	}

	models.SuccessResponse(ctx, constant.Success, http.StatusOK, "Profile updated successfully", nil, nil, nil)
	return
}

func (uc *AdminController) MigrateUsersToKeycloak(ctx *gin.Context) {
	var input struct {
		UserIds []*string `json:"user_ids"`
	}
	if err := ctx.ShouldBindJSON(&input); err != nil {
		models.ErrorResponse(ctx, constant.Failure, http.StatusBadRequest, "Missing required fields.", nil, err)
		return
	}
	err := uc.userService.MigrateUsersToKeycloak(input.UserIds)
	if err != nil {
		models.ErrorResponse(ctx, constant.Failure, http.StatusInternalServerError, "Failed to migrate user", nil, err)
		return
	}

	models.SuccessResponse(ctx, constant.Success, http.StatusOK, "User registered successfully", nil, nil, nil)
	return
}

func (uc *AdminController) MigrateRoleToKeycloak(ctx *gin.Context) {
	var input struct {
		Roles []string `json:"roles"`
	}
	if err := ctx.ShouldBindJSON(&input); err != nil {
		models.ErrorResponse(ctx, constant.Failure, http.StatusBadRequest, "Missing required fields.", nil, err)
		return
	}
	uc.userService.MigrateRoleToKeycloak(input.Roles)

	models.SuccessResponse(ctx, constant.Success, http.StatusOK, "Roles migrated successfully", nil, nil, nil)
	return
}

func (ac *AdminController) CreateUsersOnNotify(ctx *gin.Context) {
	var input struct {
		UserIds []*string `json:"user_ids"`
	}
	if err := ctx.ShouldBindJSON(&input); err != nil {
		models.ErrorResponse(ctx, constant.Failure, http.StatusBadRequest, "Missing required fields.", nil, err)
		return
	}

	err := ac.notificationService.AddUsersToNotify(input.UserIds)
	if err != nil {
		models.ErrorResponse(ctx, constant.Failure, http.StatusInternalServerError, "failed to register user(s)", nil, err)
		return
	}
	models.SuccessResponse(ctx, constant.Success, http.StatusOK, "success", nil, nil, nil)
	return
}

func (cc *AdminController) ListContactRespController(ctx *gin.Context) {
	responses, err := cc.contactService.List()
	if err != nil {
		models.ErrorResponse(ctx, constant.Failure, http.StatusInternalServerError, "error getting responses", nil, err)
		return
	}
	models.SuccessResponse(ctx, constant.Success, http.StatusOK, "messages", responses, nil, nil)
	return
}

func (ac *AdminController) CreateJobDescription(ctx *gin.Context) {

	var req models.JobDescription

	if err := ctx.ShouldBindJSON(&req); err != nil {
		models.ErrorResponse(ctx, constant.Failure, http.StatusBadRequest, "Invalid request body", nil, err)
		return
	}

	_, userId, _, err := utils.GetUserIDFromContext(ctx, ac.userService.FindUserIdBySub)
	if err != nil {
		models.ErrorResponse(ctx, constant.Failure, http.StatusUnauthorized, err.Error(), nil, err)
		return
	}

	if err := ac.jobService.CreateJob(req, userId); err != nil {
		models.ErrorResponse(ctx, constant.Failure, http.StatusInternalServerError, "Failed to create job description", nil, err)
		return
	}

	models.SuccessResponse(ctx, constant.Success, http.StatusOK, "Job Description created", nil, nil, nil)
}

func (ac *AdminController) GetJobDescriptions(ctx *gin.Context) {

	page, limit, offset := utils.GetPaginationParams(ctx)

	jobs, total, err := ac.jobService.GetJobs(limit, offset)
	if err != nil {
		models.ErrorResponse(ctx, constant.Failure, http.StatusInternalServerError, "Failed to retrieve job descriptions", nil, err)
		return
	}

	pagination := utils.GetPagination(limit, page, offset, total)

	models.SuccessResponse(ctx, constant.Success, http.StatusOK, "job_descriptions", jobs, pagination, nil)
}

func (ac *AdminController) CreateQuestion(ctx *gin.Context) {

	var req models.CreateQuestionRequest

	if err := ctx.ShouldBindJSON(&req); err != nil {
		models.ErrorResponse(ctx, constant.Failure, http.StatusBadRequest, "Invalid request body", nil, err)
		return
	}

	_, userId, _, err := utils.GetUserIDFromContext(ctx, ac.userService.FindUserIdBySub)
	if err != nil {
		models.ErrorResponse(ctx, constant.Failure, http.StatusUnauthorized, err.Error(), nil, err)
		return
	}

	questionID, err := ac.questionService.CreateQuestion(req, userId)
	if err != nil {
		models.ErrorResponse(ctx, constant.Failure, http.StatusInternalServerError, "Failed to create question", nil, err)
		return
	}

	models.SuccessResponse(ctx, constant.Success, http.StatusOK, "Question created successfully",
		map[string]interface{}{"question_id": questionID}, nil, nil)
}
