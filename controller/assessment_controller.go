package controller

import (
	"dhl/constant"
	"dhl/models"
	"dhl/services"
	"dhl/utils"
	"log"
	"net/http"

	"github.com/gin-gonic/gin"
)

type AssessmentController struct {
	assessmentService services.AssessmentService
	userService       services.UserService
	geminiService     services.GeminiService
}

func NewAssessmentController(assessmentService services.AssessmentService, userService services.UserService, geminiService services.GeminiService) *AssessmentController {
	return &AssessmentController{assessmentService: assessmentService, userService: userService, geminiService: geminiService}
}

func (ac *AssessmentController) GetAssessment(ctx *gin.Context) {
	assessmentSeq := ctx.Query("assessment_sequence")

	resp, err := ac.assessmentService.GetAssessment(assessmentSeq)
	if err != nil {
		models.ErrorResponse(ctx, "Failed to fetch assessment", http.StatusInternalServerError, err.Error(), nil, err)
		return
	}

	models.SuccessResponse(ctx, constant.Success, http.StatusOK, "assessment", resp, nil, nil)
	return
}

func (ac *AssessmentController) GetUserAssessment(ctx *gin.Context) {
	_, userId, _, err := utils.GetUserIDFromContext(ctx, ac.userService.FindUserIdBySub)
	if err != nil {
		models.ErrorResponse(ctx, constant.Failure, http.StatusUnauthorized, err.Error(), nil, err)
		return
	}
	var req models.GetAssessmentRequest
	if err := ctx.ShouldBindJSON(&req); err != nil {
		models.ErrorResponse(ctx, "Invalid input", http.StatusBadRequest, err.Error(), nil, err)
		return
	}
	req.UserId = userId

	resp, err := ac.assessmentService.GetUserAssessment(req)
	if err != nil {
		models.ErrorResponse(ctx, "Failed to fetch assessment", http.StatusInternalServerError, err.Error(), nil, err)
		return
	}

	models.SuccessResponse(ctx, constant.Success, http.StatusOK, "Loaded assessment", resp, nil, nil)
	return
}

func (ac *AssessmentController) GetUserAssessments(ctx *gin.Context) {
	_, userId, _, err := utils.GetUserIDFromContext(ctx, ac.userService.FindUserIdBySub)
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
	resp, totalRecords, err := ac.assessmentService.GetUserAssessments(userId, limit, offset, &filterData)
	if err != nil {
		models.ErrorResponse(ctx, "Failed to fetch assessment", http.StatusInternalServerError, err.Error(), nil, err)
		return
	}
	pagination := utils.GetPagination(limit, page, offset, totalRecords)

	models.SuccessResponse(ctx, constant.Success, http.StatusOK, "Loaded assessment", resp, pagination, nil)
	return
}

func (ac *AssessmentController) GetUserAssessmentCerficiate(ctx *gin.Context) {
	assessmentSession := ctx.Query("assessment_session")

	pdfBytes, err := ac.assessmentService.GenerateUserAssessmentCerficiate(assessmentSession)
	if err != nil {
		models.ErrorResponse(ctx, "Failed to generate assessment certificate", http.StatusInternalServerError, err.Error(), nil, err)
		return
	}
	ctx.Header("Content-Type", "application/pdf")
	ctx.Header("Content-Disposition", "attachment; filename=certificate.pdf")
	ctx.Data(http.StatusOK, "application/pdf", pdfBytes)
}

func (c *AssessmentController) DownloadAssessmentReport(ctx *gin.Context) {
	var filter models.AssessmentReportFilter

	if err := ctx.ShouldBindJSON(&filter); err != nil {
		models.ErrorResponse(ctx, "Invalid request", http.StatusBadRequest, err.Error(), nil, err)
		return
	}

	fileBytes, err := c.assessmentService.GenerateAssessmentExcel(ctx, filter)
	if err != nil {
		models.ErrorResponse(ctx, "Failed to generate report", http.StatusInternalServerError, err.Error(), nil, err)
		return
	}

	ctx.Header("Content-Disposition", "attachment; filename=assessment_report.xlsx")
	ctx.Data(http.StatusOK, "application/vnd.openxmlformats-officedocument.spreadsheetml.sheet", fileBytes)
}

func (ac *AssessmentController) SubmitUserAssessment(ctx *gin.Context) {
	_, userId, _, err := utils.GetUserIDFromContext(ctx, ac.userService.FindUserIdBySub)
	if err != nil {
		models.ErrorResponse(ctx, constant.Failure, http.StatusUnauthorized, err.Error(), nil, err)
		return
	}
	var req models.SubmitUserAssessmentRequest
	if err := ctx.ShouldBindJSON(&req); err != nil {
		models.ErrorResponse(ctx, "Invalid input", http.StatusBadRequest, err.Error(), nil, err)
		return
	}

	err = ac.assessmentService.SubmitAssessment(userId, req)
	if err != nil {
		models.ErrorResponse(ctx, "Failed to submit assessment", http.StatusInternalServerError, err.Error(), nil, err)
		return
	}

	models.SuccessResponse(ctx, constant.Success, http.StatusOK, "assessment submitted", nil, nil, nil)
	return
}

func (ac *AssessmentController) SubmitUserAssessmentTypingResult(ctx *gin.Context) {
	_, userId, _, err := utils.GetUserIDFromContext(ctx, ac.userService.FindUserIdBySub)
	if err != nil {
		models.ErrorResponse(ctx, constant.Failure, http.StatusUnauthorized, err.Error(), nil, err)
		return
	}
	var req models.AssessmentTypingResult
	if err := ctx.ShouldBindJSON(&req); err != nil {
		models.ErrorResponse(ctx, "Invalid input", http.StatusBadRequest, err.Error(), nil, err)
		return
	}

	err = ac.assessmentService.SaveAssessmentTypingRespone(&req, userId)
	if err != nil {
		models.ErrorResponse(ctx, "Failed to submit typing results", http.StatusInternalServerError, err.Error(), nil, err)
		return
	}

	models.SuccessResponse(ctx, constant.Success, http.StatusOK, "results submitted", nil, nil, nil)
	return
}

func (ac *AssessmentController) DuplicateAssessment(ctx *gin.Context) {
	_, userId, _, err := utils.GetUserIDFromContext(ctx, ac.userService.FindUserIdBySub)
	if err != nil {
		models.ErrorResponse(ctx, constant.Failure, http.StatusUnauthorized, err.Error(), nil, err)
		return
	}
	var req models.DuplicateAssessmentRequest

	if err := ctx.ShouldBindJSON(&req); err != nil {
		models.ErrorResponse(ctx, "Invalid input", http.StatusBadRequest, err.Error(), nil, err)
		return
	}

	response, err := ac.assessmentService.CreateDuplicateAssessment(req.AssessmentSequence, userId)
	if err != nil {
		models.ErrorResponse(ctx, constant.Failure, http.StatusInternalServerError, "failed", nil, err)
		return
	}

	models.SuccessResponse(ctx, constant.Success, http.StatusOK, "Assessment", response, nil, nil)
	return
}

func (ac *AssessmentController) UploadAssessment(c *gin.Context) {
	_, userId, _, err := utils.GetUserIDFromContext(c, ac.userService.FindUserIdBySub)
	if err != nil {
		models.ErrorResponse(c, constant.Failure, http.StatusUnauthorized, err.Error(), nil, err)
		return
	}
	file, err := c.FormFile("file")
	if err != nil {
		models.ErrorResponse(c, constant.Failure, http.StatusBadRequest, "File missing", nil, err)
		return
	}

	f, err := file.Open()
	if err != nil {
		models.ErrorResponse(c, constant.Failure, http.StatusInternalServerError, "File open error", nil, err)
		return
	}
	defer f.Close()

	response, err := ac.assessmentService.CreateAssessmentViaFileUpload(f, file.Filename, userId, nil)
	if err != nil {
		models.ErrorResponse(c, constant.Failure, http.StatusInternalServerError, "Processing failed", nil, err)
		return
	}

	models.SuccessResponse(c, constant.Success, http.StatusOK, "Assessment", response, nil, nil)
}

func (ac *AssessmentController) CreateAssessment(ctx *gin.Context) {
	_, userId, _, err := utils.GetUserIDFromContext(ctx, ac.userService.FindUserIdBySub)
	if err != nil {
		models.ErrorResponse(ctx, constant.Failure, http.StatusUnauthorized, err.Error(), nil, err)
		return
	}
	var req models.ManualAssessmentRequest
	if err := ctx.ShouldBindJSON(&req); err != nil {
		models.ErrorResponse(ctx, "Invalid input", http.StatusBadRequest, err.Error(), nil, err)
		return
	}
	resp, err := ac.assessmentService.CreateAssessmentViaMaual(ctx, req, userId)
	if err != nil {
		models.ErrorResponse(ctx, constant.Failure, http.StatusInternalServerError, "failed", nil, err)
		return
	}
	models.SuccessResponse(ctx, constant.Success, http.StatusOK, "Assessment", resp, nil, nil)
}

func (ac *AssessmentController) CreateSessionImage(ctx *gin.Context) {
	log.Println("=== CreateSessionImage API called ===")

	// Get user from context
	_, userId, _, err := utils.GetUserIDFromContext(ctx, ac.userService.FindUserIdBySub)
	if err != nil {
		log.Printf("ERROR: Failed to get user ID from context: %v", err)
		models.ErrorResponse(ctx, constant.Failure, http.StatusUnauthorized, err.Error(), nil, err)
		return
	}
	log.Printf("Authenticated user ID: %s", userId)

	// Get session_id from form data
	sessionID := ctx.PostForm("session_id")
	log.Printf("Received session_id: %s", sessionID)
	if sessionID == "" {
		log.Println("ERROR: session_id is empty or missing")
		models.ErrorResponse(ctx, "Invalid input", http.StatusBadRequest, "session_id is required", nil, nil)
		return
	}

	// Get image file from form data
	file, err := ctx.FormFile("image")
	if err != nil {
		log.Printf("ERROR: Failed to get image file from form: %v", err)
		models.ErrorResponse(ctx, constant.Failure, http.StatusBadRequest, "Image file is required", nil, err)
		return
	}
	log.Printf("Received image file: %s, size: %d bytes", file.Filename, file.Size)

	// Open and read the file
	fileContent, err := file.Open()
	if err != nil {
		log.Printf("ERROR: Failed to open image file: %v", err)
		models.ErrorResponse(ctx, constant.Failure, http.StatusInternalServerError, "Failed to read image file", nil, err)
		return
	}
	defer fileContent.Close()

	// Read file bytes
	imageData := make([]byte, file.Size)
	bytesRead, err := fileContent.Read(imageData)
	if err != nil {
		log.Printf("ERROR: Failed to read image data: %v", err)
		models.ErrorResponse(ctx, constant.Failure, http.StatusInternalServerError, "Failed to read image data", nil, err)
		return
	}
	log.Printf("Successfully read %d bytes from image file", bytesRead)

	// Call service to save the image
	log.Printf("Calling CreateSessionImage service for user: %s, session: %s", userId, sessionID)
	resp, err := ac.assessmentService.CreateSessionImage(userId, sessionID, imageData)
	if err != nil {
		log.Printf("ERROR: Service failed to save session image: %v", err)
		models.ErrorResponse(ctx, "Failed to save session image", http.StatusInternalServerError, err.Error(), nil, err)
		return
	}

	log.Printf("SUCCESS: Session image saved with ID: %d", resp.ID)
	models.SuccessResponse(ctx, constant.Success, http.StatusOK, "Session image saved successfully", resp, nil, nil)
}

func (ac *AssessmentController) GenerateAssessmentWithAI(ctx *gin.Context) {
	var req models.GenerateAssessmentRequest
	if err := ctx.ShouldBindJSON(&req); err != nil {
		models.ErrorResponse(ctx, "Invalid input", http.StatusBadRequest, err.Error(), nil, err)
		return
	}

	assessment, err := ac.geminiService.GenerateAssessment(req)
	if err != nil {
		models.ErrorResponse(ctx, "Failed to generate assessment", http.StatusInternalServerError, err.Error(), nil, err)
		return
	}

	response := models.GenerateAssessmentResponse{
		Assessment: *assessment,
	}

	models.SuccessResponse(ctx, constant.Success, http.StatusOK, "Assessment generated successfully", response, nil, nil)
}

func (ac *AssessmentController) SaveGeneratedAssessment(ctx *gin.Context) {
	_, userId, _, err := utils.GetUserIDFromContext(ctx, ac.userService.FindUserIdBySub)
	if err != nil {
		models.ErrorResponse(ctx, constant.Failure, http.StatusUnauthorized, err.Error(), nil, err)
		return
	}

	var req models.SaveGeneratedAssessmentRequest
	if err := ctx.ShouldBindJSON(&req); err != nil {
		models.ErrorResponse(ctx, "Invalid input", http.StatusBadRequest, err.Error(), nil, err)
		return
	}

	// Validate the assessment
	if req.Assessment.AssessmentName == "" {
		models.ErrorResponse(ctx, "Invalid input", http.StatusBadRequest, "assessment_name is required", nil, nil)
		return
	}

	if len(req.Assessment.Questions) == 0 {
		models.ErrorResponse(ctx, "Invalid input", http.StatusBadRequest, "at least one question is required", nil, nil)
		return
	}

	// Validate each question
	for i, q := range req.Assessment.Questions {
		if q.Title == "" {
			models.ErrorResponse(ctx, "Invalid input", http.StatusBadRequest,
				"question title is required for question "+string(rune(i+1)), nil, nil)
			return
		}
		if len(q.Options) < 2 {
			models.ErrorResponse(ctx, "Invalid input", http.StatusBadRequest,
				"at least 2 options required for question "+string(rune(i+1)), nil, nil)
			return
		}

		hasCorrectAnswer := false
		for _, opt := range q.Options {
			if opt.IsCorrect {
				hasCorrectAnswer = true
				break
			}
		}
		if !hasCorrectAnswer {
			models.ErrorResponse(ctx, "Invalid input", http.StatusBadRequest,
				"question "+string(rune(i+1))+" must have at least one correct answer", nil, nil)
			return
		}
	}

	response, err := ac.assessmentService.CreateAssessmentViaFileUpload(nil, "", userId, &req.Assessment)
	if err != nil {
		models.ErrorResponse(ctx, constant.Failure, http.StatusInternalServerError, "Failed to save assessment", nil, err)
		return
	}

	models.SuccessResponse(ctx, constant.Success, http.StatusOK, "Assessment saved successfully", response, nil, nil)
}
