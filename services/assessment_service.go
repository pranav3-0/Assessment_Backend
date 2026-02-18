package services

import (
	"context"
	"dhl/constant"
	"dhl/models"
	"dhl/repository"
	"dhl/utils"
	"errors"
	"log"
	"mime/multipart"
	"time"

	"github.com/google/uuid"
	"github.com/xuri/excelize/v2"
	"gorm.io/gorm"
)

type AssessmentService interface {
	// GET
	GetAssessment(assessmentSeq string) (*models.AssessmentResponse, error)
	GetUserAssessment(req models.GetAssessmentRequest) (*models.AssessmentResponse, error)
	GetUserAssessments(userId string, limit, offset int, filters *models.AssessmentFilter) (interface{}, int64, error)
	GetManagerAssessments(managerID string, limit, offset int, filters *models.AssessmentFilter) (interface{}, int64, error)
	GetAssessments(limit, offset int, filters *models.AssessmentFilter) (interface{}, int64, error)
	GetQuestions(limit, offset int) (interface{}, int64, error)
	GenerateUserAssessmentCerficiate(assessmentSession string) ([]byte, error)
	GenerateAssessmentExcel(ctx context.Context, filter models.AssessmentReportFilter) ([]byte, error)

	// CREATE
	CreateDuplicateAssessment(assessmentSequence, userId string) (interface{}, error)
	CreateAssessmentViaFileUpload(file multipart.File, filename, userId string, assessment *models.SheetAssessment) (interface{}, error)
	CreateAssessmentViaMaual(ctx context.Context, request models.ManualAssessmentRequest, userId string) (interface{}, error)
	SubmitAssessment(userID string, req models.SubmitUserAssessmentRequest) error
	DistributeAssessmentUser(assessmentSeq string, userIDs []string) error
	DistributeAssessmentManager(assessmentSeq string, userIDs []string) error
	SaveAssessmentTypingRespone(input *models.AssessmentTypingResult, userId string) error
	CreateSessionImage(userID, sessionID string, imageData []byte) (*models.SessionImageResponse, error)

	// UPDATE
	UpdateAssessmentService(models.UpdateAssessmentRequest) error
	UpdateAssessmentStatusService(models.UpdateAssessmentStatusRequest) error
}

type AssessmentServiceImpl struct {
	assessmentRepo repository.AssessmentRepository
	db             *gorm.DB
}

func NewAssessmentService(assessmentRepo repository.AssessmentRepository, db *gorm.DB) AssessmentService {
	return &AssessmentServiceImpl{assessmentRepo: assessmentRepo, db: db}
}

func (s *AssessmentServiceImpl) GetAssessment(assessmentSeq string) (*models.AssessmentResponse, error) {
	if assessmentSeq == "" {
		return nil, errors.New("invalid input: assessmentSeq")
	}

	assessment, err := s.assessmentRepo.GetAssessmentMstByAssmtSeq(assessmentSeq)
	if err != nil {
		return nil, err
	}
	if assessment == nil {
		return nil, errors.New("assessment not found")
	}

	assessmentExt, err := s.assessmentRepo.GetDhlSurveyExtByAssmtSeq(assessmentSeq)
	if err != nil {
		return nil, err
	}

	questions, err := s.assessmentRepo.GetAssessmentQuestions(assessmentSeq)
	if err != nil {
		return nil, err
	}

	// Collect question IDs for bulk tag fetching
	questionIDs := make([]int64, len(questions))
	for i, q := range questions {
		questionIDs[i] = q.QuestionID
	}
	questionTagsMap, _ := s.assessmentRepo.GetTagRequestsByQuestionIDs(questionIDs)

	var questionResponses []models.AssessmentQuestion
	for _, q := range questions {
		questionContent, err := s.assessmentRepo.GetContentByQuestionID(q.QuestionID)
		if err != nil {
			return nil, err
		}

		options, err := s.assessmentRepo.GetOptionsByQuestionID(q.QuestionID)
		if err != nil {
			return nil, err
		}

		var answers []models.Answer
		for _, opt := range options {
			optContents, _ := s.assessmentRepo.GetContentByID(opt.ContentID)
			answers = append(answers, models.Answer{
				AnswerID:      int(opt.OptionID),
				Sequence:      &opt.SequenceID,
				OptionLabel:   optContents.Value,
				CorrectAnswer: opt.IsAnswer,
			})
		}

		questionTags := questionTagsMap[q.QuestionID]
		questionResponses = append(questionResponses, models.AssessmentQuestion{
			QuestionID:      int(q.QuestionID),
			Sequence:        int(q.SequenceID),
			Title:           questionContent.Value,
			Answers:         answers,
			AttemptedAnswer: nil,
			QuestionTime:    int(q.DurationInSeconds),
			QuestionTypeId:  questionContent.QuestionTypeId,
			QuestionType:    questionContent.QuestionType,
			Tags:            questionTags,
		})
	}

	tags, _ := s.assessmentRepo.GetTagRequestsByAssessmentSequence(assessment.AssessmentSequence)

	resp := &models.AssessmentResponse{
		AssessmentID:        assessment.AssessmentID,
		AssessmentSequence:  assessment.AssessmentSequence,
		AssessmentName:      assessment.AssessmentDesc,
		AssessmentStatus:    assessmentExt.State,
		NoOfAttempts:        assessmentExt.AttemptsLimit,
		QuestionsCount:      len(questionResponses),
		AssessmentType:      assessment.AssessmentType,
		AssessmentDuration:  &assessment.Duration,
		Questions:           questionResponses,
		Instruction:         assessment.Instructions,
		TimeLimit:           assessmentExt.TimeLimit,
		Marks:               assessment.Marks,
		Certificate:         assessmentExt.Certificate,
		Deadline:            &assessmentExt.Deadline,
		CenterName:          assessmentExt.CenterName,
		ServiceLineName:     assessmentExt.ServiceLineName,
		BusinessPartnerName: assessmentExt.BusinessPartnerName,
		SubBusinessPartner:  assessmentExt.SubBusinessPartnerName,
		ServiceGroupName:    assessmentExt.ServiceGroupName,
		ServiceName:         assessmentExt.ServiceName,
		Tags:                tags,
	}

	return resp, nil
}

func (s *AssessmentServiceImpl) GetUserAssessment(req models.GetAssessmentRequest) (*models.AssessmentResponse, error) {
	if req.AssessmentId == "" || req.UserId == "" {
		return nil, errors.New("invalid input: assessmentId and userId required")
	}
	assessment, err := s.assessmentRepo.GetAssessmentMstByAssmtSeq(req.AssessmentId)
	if err != nil {
		return nil, err
	}
	assessmentExt, err := s.assessmentRepo.GetDhlSurveyExtByAssmtSeq(req.AssessmentId)
	if err != nil {
		return nil, err
	}
	if assessment == nil {
		return nil, errors.New("assessment not found")
	}

	tx := s.db.Begin()
	session, err := s.assessmentRepo.GetOrCreateUserSession(tx, req.UserId, req.AssessmentId, assessment.AssessmentType, assessment.PartnerID)
	if err != nil {
		tx.Rollback()
		return nil, err
	}

	status, err := s.assessmentRepo.UpdateAssessmentStatus(tx, req.UserId, req.AssessmentId, "STARTED")
	if err != nil {
		tx.Rollback()
		return nil, err
	}
	attemptedTypingTest := false
	typingRes, _ := s.assessmentRepo.GetAssessmentTypingResult(req.UserId, req.AssessmentId)
	if typingRes != nil {
		attemptedTypingTest = true
	}

	questions, err := s.assessmentRepo.GetAssessmentQuestions(req.AssessmentId)
	if err != nil {
		tx.Rollback()
		return nil, err
	}

	// Collect question IDs for bulk tag fetching
	questionIDs := make([]int64, len(questions))
	for i, q := range questions {
		questionIDs[i] = q.QuestionID
	}
	questionTagsMap, _ := s.assessmentRepo.GetTagRequestsByQuestionIDs(questionIDs)

	var questionResponses []models.AssessmentQuestion
	for _, q := range questions {
		questionContent, err := s.assessmentRepo.GetContentByQuestionID(q.QuestionID)
		if err != nil {
			tx.Rollback()
			return nil, err
		}

		options, err := s.assessmentRepo.GetOptionsByQuestionID(q.QuestionID)
		if err != nil {
			tx.Rollback()
			return nil, err
		}

		var answers []models.Answer
		for _, opt := range options {
			optContents, _ := s.assessmentRepo.GetContentByID(opt.ContentID)
			answers = append(answers, models.Answer{
				AnswerID:    int(opt.OptionID),
				Sequence:    &opt.SequenceID,
				OptionLabel: optContents.Value,
			})
		}

		questionTags := questionTagsMap[q.QuestionID]
		questionResponses = append(questionResponses, models.AssessmentQuestion{
			QuestionID:      int(q.QuestionID),
			Sequence:        int(q.SequenceID),
			Title:           questionContent.Value,
			Answers:         answers,
			AttemptedAnswer: nil,
			QuestionTime:    int(q.DurationInSeconds),
			QuestionTypeId:  questionContent.QuestionTypeId,
			QuestionType:    questionContent.QuestionType,
			SkippingAllowed: q.SkippingAllowed,
			Tags:            questionTags,
		})
	}

	if err := tx.Commit().Error; err != nil {
		return nil, err
	}

	// Fetch session image if exists
	log.Printf("Fetching session image for session: %s", session.SessionID.String())
	sessionImage, err := s.assessmentRepo.GetSessionImageBySessionID(session.SessionID.String())
	if err != nil {
		log.Printf("WARNING: Failed to fetch session image: %v", err)
		// Continue without image - don't fail the entire request
	}
	if sessionImage != nil {
		log.Printf("Session image found, size: %d bytes", len(sessionImage))
	} else {
		log.Println("No session image found, returning empty image data")
	}

	tags, _ := s.assessmentRepo.GetTagRequestsByAssessmentSequence(assessment.AssessmentSequence)

	resp := &models.AssessmentResponse{
		AssessmentID:          assessment.AssessmentID,
		AssessmentSequence:    assessment.AssessmentSequence,
		AssessmentName:        assessment.AssessmentDesc,
		AssessmentUsersStatus: status.AssessmentStatus,
		AssessmentStatus:      assessmentExt.State,
		QuestionsCount:        len(questionResponses),
		AssessmentType:        assessment.AssessmentType,
		AssessmentDuration:    &assessment.Duration,
		Certificate:           assessmentExt.Certificate,
		NoOfAttempts:          assessmentExt.AttemptsLimit,
		Questions:             questionResponses,
		SessionID:             session.SessionID.String(),
		AssessmentReport:      false,
		Instruction:           assessment.Instructions,
		TimeLimit:             assessmentExt.TimeLimit,
		Marks:                 assessment.Marks,
		AttemptedTypingTest:   attemptedTypingTest,
		IsTypingTest:          assessmentExt.IsTypingTest,
		SessionImage:          sessionImage,
		Tags:                  tags,
	}

	return resp, nil
}

func (s *AssessmentServiceImpl) GetUserAssessments(
	userId string,
	limit, offset int,
	filters *models.AssessmentFilter,
) (interface{}, int64, error) {

	// Case 1: Specific session + sequence
	if filters.AssessmentSessionId != nil && filters.AssessmentSequence != nil {
		asmt, err := s.assessmentRepo.GetUserAssessmentResponse(
			&userId,
			*filters.AssessmentSequence,
			*filters.AssessmentSessionId,
			string(constant.User),
		)
		if err != nil {
			return nil, 0, err
		}
		return asmt, 0, nil
	}

	// Case 2: Get sessions of specific assessment
	if filters.AssessmentSequence != nil {
		return s.assessmentRepo.GetUserAssessmentSessions(
			userId,
			*filters.AssessmentSequence,
		)
	}

	// Case 3: Paginated list (MAIN LIST API)
	assessments, total, err := s.assessmentRepo.GetAssessmentsForUserWithPagination(
		userId,
		limit,
		offset,
	)
	if err != nil {
		return nil, 0, err
	}

	// 🔥 ADD TAGS HERE
	for i, a := range assessments {
		tags, err := s.assessmentRepo.GetTagsByAssessmentSequence(a.AssessmentSequence)
		if err == nil {
			assessments[i].Tags = tags
		}
	}

	return assessments, total, nil
}

func (s *AssessmentServiceImpl) GetManagerAssessments(
	managerID string,
	limit, offset int,
	filters *models.AssessmentFilter,
) (interface{}, int64, error) {

	if filters.AssessmentSessionId != nil && filters.AssessmentSequence != nil {
		asmt, err := s.assessmentRepo.GetUserAssessmentResponse(
			nil,
			*filters.AssessmentSequence,
			*filters.AssessmentSessionId,
			string(constant.Manager),
		)
		if err != nil {
			return nil, 0, err
		}
		return asmt, 0, nil
	}

	if filters.AssessmentSequence != nil {
		return s.assessmentRepo.GetAssessmentsAttendeeInfo(
			*filters.AssessmentSequence,
			limit,
			offset,
		)
	}

	return s.assessmentRepo.GetAssessmentsForManagerWithPagination(
		managerID,
		limit,
		offset,
	)
}

func (s *AssessmentServiceImpl) GetAssessments(limit, offset int, filters *models.AssessmentFilter) (interface{}, int64, error) {
	if filters.AssessmentSessionId != nil && filters.AssessmentSequence != nil {
		asmt, err := s.assessmentRepo.GetUserAssessmentResponse(nil, *filters.AssessmentSequence, *filters.AssessmentSessionId, string(constant.Admin))
		if err != nil {
			return nil, 0, err
		}
		return asmt, 0, nil
	}
	if filters.AssessmentSequence != nil {
		return s.assessmentRepo.GetAssessmentsAttendeeInfo(*filters.AssessmentSequence, limit, offset)
	}
	return s.assessmentRepo.GetAssessmentsWithPagination(limit, offset)
}


func (s *AssessmentServiceImpl) GetQuestions(limit, offset int) (interface{}, int64, error) {
	return s.assessmentRepo.GetPaginatedQuestionsWithOptions(limit, offset)
}

func (s *AssessmentServiceImpl) GenerateUserAssessmentCerficiate(assessmentSession string) ([]byte, error) {
	details, err := s.assessmentRepo.GetCertificateDetailsBySessionId(assessmentSession)
	if err != nil {
		return nil, err
	}
	userScore := (float64(details.MarksObtained) / float64(details.TotalMarks)) * 100
	passed := userScore >= details.PassingScore
	pdfBytes, err := utils.GenerateCertificatePDF(details.UserName, details.AssessmentTitle, details.TotalMarks, details.MarksObtained, details.PassingScore, userScore, details.CompletedAt, passed)
	if err != nil {
		return nil, err
	}
	return pdfBytes, nil
}

func (s *AssessmentServiceImpl) CreateDuplicateAssessment(assessmentSequence, userId string) (interface{}, error) {
	asmtMst, err := s.assessmentRepo.GetAssessmentMstByAssmtSeq(assessmentSequence)
	if err != nil {
		return nil, err
	}
	asmtQtns, err := s.assessmentRepo.GetAssessmentQuestions(asmtMst.AssessmentSequence)
	if err != nil {
		return nil, err
	}

	// Get all question IDs for bulk tag fetching
	questionIDs := make([]int64, len(asmtQtns))
	for i, q := range asmtQtns {
		questionIDs[i] = q.QuestionID
	}
	questionTagsMap, _ := s.assessmentRepo.GetTagRequestsByQuestionIDs(questionIDs)

	var questionsToAdd []models.SheetQuestion
	for _, assmtQtn := range asmtQtns {
		var optionsToAdd []models.SheetOption
		question, err := s.assessmentRepo.GetContentByQuestionID(assmtQtn.QuestionID)
		if err != nil {
			return nil, err
		}
		options, err := s.assessmentRepo.GetOptionsByQuestionID(assmtQtn.QuestionID)
		if err != nil {
			return nil, err
		}
		for _, op := range options {
			opContent, err := s.assessmentRepo.GetContentByID(op.ContentID)
			if err != nil {
				return nil, err
			}
			optionsToAdd = append(optionsToAdd, models.SheetOption{
				IsCorrect: op.IsAnswer,
				Score:     op.AnswerScore,
				Label:     opContent.Value,
			})
		}
		questionTags := questionTagsMap[assmtQtn.QuestionID]
		questionsToAdd = append(questionsToAdd, models.SheetQuestion{
			Title:   question.Value,
			Options: optionsToAdd,
			Tags:    questionTags,
		})
	}

	// Get tags from original assessment
	tags, _ := s.assessmentRepo.GetTagRequestsByAssessmentSequence(asmtMst.AssessmentSequence)

	newAssmt := models.SheetAssessment{
		AssessmentName: asmtMst.AssessmentDesc,
		Questions:      questionsToAdd,
		Tags:           tags,
	}
	tx := s.db.Begin()
	newAssmtSeq, err := s.assessmentRepo.SaveAssessmentWithQuestions(context.Background(), tx, newAssmt, userId)
	if err != nil {
		tx.Rollback()
		return nil, err
	}
	if err := tx.Commit().Error; err != nil {
		return nil, err
	}

	response := map[string]string{
		"assessment_sequence": newAssmtSeq,
	}
	return response, nil
}
func (s *AssessmentServiceImpl) CreateAssessmentViaFileUpload(file multipart.File, filename, userId string, assessment *models.SheetAssessment) (interface{}, error) {
	var jsonResponse *models.SheetAssessment
	var err error

	// If assessment is provided directly (from AI generation), use it
	// Otherwise parse from file upload
	if assessment != nil {
		jsonResponse = assessment
	} else {
		jsonResponse, err = utils.ParseQuestionnaireExcelToJSON(file, filename)
		if err != nil {
			return nil, err
		}
	}

	tx := s.db.Begin()
	assessmentSequence, err := s.assessmentRepo.SaveAssessmentWithQuestions(context.Background(), tx, *jsonResponse, userId)
	if err != nil {
		tx.Rollback()
		return nil, err
	}
	tx.Commit()
	jsonResponse.AssessmentSequence = assessmentSequence
	return jsonResponse, nil
}

func (s *AssessmentServiceImpl) CreateAssessmentViaMaual(ctx context.Context, req models.ManualAssessmentRequest, userId string) (interface{}, error) {
	tx := s.db.Begin()
	if tx.Error != nil {
		return nil, tx.Error
	}

	assessment := &models.AssessmentMst{
		AssessmentDesc: req.AssessmentName,
		Duration:       req.Duration,
		Marks:          req.Marks,
		StartTime:      req.StartTime,
		ValidFrom:      req.ValidFrom,
		ValidTo:        req.ValidTo,
		Instructions:   req.Instructions,
		AssessmentType: req.AssessmentType,
		CreatedBy:      userId,
		CreatedOn:      time.Now(),
		JobID:          req.JobID,
		IsActive:       true,
		IsDeleted:      false,
	}
	asmt, err := s.assessmentRepo.CreateAssessment(ctx, tx, assessment)
	if err != nil {
		tx.Rollback()
		return nil, err
	}
	ext := &models.DhlSurveySurveyExt{
		SurveyID:               asmt.AssessmentID,
		AssessmentSequence:     asmt.AssessmentSequence,
		State:                  string(constant.Draft),
		Color:                  1,
		Certificate:            req.Certificate,
		CenterID:               req.CenterID,
		ServiceLineID:          req.ServiceLineID,
		BusinessPartnerID:      req.BusinessPartnerID,
		SubBusinessPartnerID:   req.SubBusinessPartnerID,
		ServiceGroupID:         req.ServiceGroupID,
		ServiceID:              req.ServiceID,
		UsersLoginRequired:     true,
		AttemptsLimit:          1,
		TimeLimit:              30,
		ShowResult:             true,
		IsTypingTest:           false,
		CertificationGiveBadge: true,
	}

	createdExt, err := s.assessmentRepo.CreateAssessmentExt(ctx, tx, ext)
	if err != nil {
		tx.Rollback()
		return nil, err
	}

	err = s.assessmentRepo.CreateAssessmentQuestions(ctx, tx, asmt.AssessmentID, asmt.AssessmentSequence, req.Questions, userId)
	if err != nil {
		tx.Rollback()
		return nil, err
	}

	// Handle tags if provided (with parent-child tag support)
	if len(req.Tags) > 0 {
		for _, tagReq := range req.Tags {
			tagIDs, err := s.assessmentRepo.ProcessTagRequest(tx, tagReq, userId)
			if err != nil {
				tx.Rollback()
				return nil, err
			}
			for _, tagID := range tagIDs {
				if err := s.assessmentRepo.CreateAssessmentTagMappingWithParents(tx, asmt.AssessmentSequence, tagID, userId); err != nil {
					tx.Rollback()
					return nil, err
				}
			}
		}
	}

	if err := tx.Commit().Error; err != nil {
		return nil, err
	}

	tags, _ := s.assessmentRepo.GetTagRequestsByAssessmentSequence(asmt.AssessmentSequence)

	response := models.AssessmentResponse{
		AssessmentID:       asmt.AssessmentID,
		AssessmentSequence: asmt.AssessmentSequence,
		Marks:              asmt.Marks,
		AssessmentName:     asmt.AssessmentDesc,
		AssessmentStatus:   createdExt.State,
		Tags:               tags,
	}
	return response, nil
}

func (s *AssessmentServiceImpl) SubmitAssessment(userID string, req models.SubmitUserAssessmentRequest) error {
	session, err := s.assessmentRepo.GetSessionByID(req.SessionID, userID)
	if err != nil {
		return err
	}
	if session == nil {
		return errors.New("invalid session")
	}
	tx := s.db.Begin()
	for _, r := range req.Response {
		err := s.assessmentRepo.SaveAssessmentResponse(tx, session, req.AssessmentSequence, r)
		if err != nil {
			tx.Rollback()
			return err
		}
	}

	_, err = s.assessmentRepo.UpdateAssessmentStatus(tx, session.UserID, session.AssessmentID, "COMPLETED")
	if err != nil {
		tx.Rollback()
		return err
	}

	if err := tx.Commit().Error; err != nil {
		return err
	}
	return nil
}

func (s *AssessmentServiceImpl) DistributeAssessmentUser(assessmentSeq string, userIDs []string) error {
	asmt, err := s.assessmentRepo.GetAssessmentMstByAssmtSeq(assessmentSeq)
	if err != nil {
		return nil
	}
	if asmt == nil {
		return errors.New("assessment not found")
	}
	tx := s.db.Begin()
	for _, uid := range userIDs {
		_, err := s.assessmentRepo.UpdateAssessmentStatus(tx, uid, assessmentSeq, "ASSIGNED")
		if err != nil {
			tx.Rollback()
			return err
		}
	}
	if err := tx.Commit().Error; err != nil {
		return err
	}
	return nil
}

func (s *AssessmentServiceImpl) DistributeAssessmentManager(assessmentSeq string, userIDs []string) error {
	asmt, err := s.assessmentRepo.GetAssessmentMstByAssmtSeq(assessmentSeq)
	if err != nil {
		return nil
	}
	if asmt == nil {
		return errors.New("assessment not found")
	}
	tx := s.db.Begin()
	for _, uid := range userIDs {
		_, err := s.assessmentRepo.AddManagerAssessmentMapping(tx, uid, assessmentSeq)
		if err != nil {
			tx.Rollback()
			return err
		}
	}
	if err := tx.Commit().Error; err != nil {
		return err
	}
	return nil
}

func (s *AssessmentServiceImpl) SaveAssessmentTypingRespone(input *models.AssessmentTypingResult, userId string) error {
	assessmentTypingResult := models.AssessmentTypingResult{
		AssessmentSeq:   input.AssessmentSeq,
		WPM:             input.WPM,
		Accuracy:        input.Accuracy,
		TotalWords:      input.TotalWords,
		CorrectWords:    input.CorrectWords,
		IncorrectWords:  input.IncorrectWords,
		CharactersTyped: input.CharactersTyped,
		UserID:          userId,
		CreatedOn:       time.Now(),
	}
	tx := s.db.Begin()
	_, err := s.assessmentRepo.AddAssessmentTypingResult(tx, assessmentTypingResult)
	if err != nil {
		tx.Rollback()
		return err
	}
	if err := tx.Commit().Error; err != nil {
		return err
	}
	return nil
}

func (s *AssessmentServiceImpl) UpdateAssessmentService(request models.UpdateAssessmentRequest) error {
	assment, err := s.assessmentRepo.GetAssessmentMstByAssmtSeq(request.AssessmentSequence)
	if err != nil {
		return err
	}
	if assment == nil {
		return errors.New("assessment not found")
	}
	tx := s.db.Begin()
	if len(request.DeletedQuestionIDs) > 0 {
		for _, id := range request.DeletedQuestionIDs {
			asmtQtnUpdates := models.AssessmentQuestionMst{
				AssessmentQuestionID: id,
				IsDeleted:            true,
				IsActive:             false,
				ModifiedOn:           time.Now(),
			}
			err := s.assessmentRepo.UpdateAssessmentQuestion(tx, &asmtQtnUpdates)
			if err != nil {
				log.Println("Error while Updating Assessment Question:", id, " ::", err)
				tx.Rollback()
				return err
			}
		}
	}

	assessmentUpdates := &models.AssessmentMst{
		AssessmentID: assment.AssessmentID,
		ModifiedOn:   time.Now(),
	}
	if request.AssessmentDetails.AssessmentDesc != nil {
		assessmentUpdates.AssessmentDesc = *request.AssessmentDetails.AssessmentDesc
	}
	if request.AssessmentDetails.Duration != nil {
		assessmentUpdates.Duration = *request.AssessmentDetails.Duration
	}
	if request.AssessmentDetails.Marks != nil {
		assessmentUpdates.Marks = *request.AssessmentDetails.Marks
	}
	if request.AssessmentDetails.Instruction != nil {
		assessmentUpdates.Instructions = *request.AssessmentDetails.Instruction
	}
	if request.AssessmentDetails.AssessmentType != nil {
		assessmentUpdates.AssessmentType = *request.AssessmentDetails.AssessmentType
	}

	err = s.assessmentRepo.UpdateAssessment(tx, assessmentUpdates)
	if err != nil {
		log.Println("Error while Updating Assessment:", err)
		tx.Rollback()
		return err
	}

	dhlSurveyUpdates := &models.DhlSurveySurveyExt{
		AssessmentSequence: assment.AssessmentSequence,
	}
	if request.AssessmentDetails.Deadline != nil {
		dhlSurveyUpdates.Deadline = *request.AssessmentDetails.Deadline
	}
	if request.AssessmentDetails.State != nil {
		dhlSurveyUpdates.State = *request.AssessmentDetails.State
	}
	if request.AssessmentDetails.TimeLimit != nil {
		dhlSurveyUpdates.TimeLimit = float64(*request.AssessmentDetails.TimeLimit)
	}

	if request.AssessmentDetails.ServiceLineID != nil {
		dhlSurveyUpdates.ServiceLineID = *request.AssessmentDetails.ServiceLineID
	}

	if request.AssessmentDetails.BusinessPartnerID != nil {
		dhlSurveyUpdates.BusinessPartnerID = *request.AssessmentDetails.BusinessPartnerID
	}
	if request.AssessmentDetails.SubBusinessPartnerID != nil {
		dhlSurveyUpdates.SubBusinessPartnerID = *request.AssessmentDetails.SubBusinessPartnerID
	}
	if request.AssessmentDetails.ServiceGroupID != nil {
		dhlSurveyUpdates.ServiceGroupID = *request.AssessmentDetails.ServiceGroupID
	}

	if request.AssessmentDetails.ServiceID != nil {
		dhlSurveyUpdates.ServiceID = *request.AssessmentDetails.ServiceID
	}

	if request.AssessmentDetails.CenterId != nil {
		dhlSurveyUpdates.CenterID = *request.AssessmentDetails.CenterId
	}

	if request.AssessmentDetails.AllowShowResult != nil {
		dhlSurveyUpdates.ShowResult = *request.AssessmentDetails.AllowShowResult
	}

	if request.AssessmentDetails.AllowViewCertificate != nil {
		dhlSurveyUpdates.Certificate = *request.AssessmentDetails.AllowViewCertificate
	}

	err = s.assessmentRepo.UpdateDHLAssessmentExt(tx, dhlSurveyUpdates)
	if err != nil {
		log.Println("Error while Updating DHL Assessment Ext:", err)
		tx.Rollback()
		return err
	}

	for _, question := range request.Questions {
		if question.QuestionID == 0 {
			newQ := models.QuestionMain{
				QuestionTitle: question.Title,
				QuestionType:  question.QuestionType,
			}

			log.Println("Adding new Question")

			nQId, err := s.assessmentRepo.AddNewQuestion(tx, newQ, assment.AssessmentSequence)
			if err != nil {
				log.Println("Error while Add new Question:", err)
				tx.Rollback()
				return err
			}
			for _, option := range question.Options {
				newOpt := models.OptionMain{
					QuestionID:  nQId,
					OptionLabel: option.Label,
					SequenceID:  0,
					AnswerScore: option.Score,
					IsAnswer:    option.IsCorrect,
				}
				err := s.assessmentRepo.AddNewOption(tx, newOpt)
				if err != nil {
					log.Println("Error while Adding new option:", err)
					tx.Rollback()
					return err
				}
			}

			// Handle question tags for new question (with parent-child tag support)
			if len(question.Tags) > 0 {
				for _, tagReq := range question.Tags {
					tagIDs, err := s.assessmentRepo.ProcessTagRequest(tx, tagReq, "system")
					if err != nil {
						log.Println("Error while processing tag request:", err)
						tx.Rollback()
						return err
					}
					for _, tagID := range tagIDs {
						if err := s.assessmentRepo.CreateQuestionTagMappingWithParents(tx, nQId, tagID, "system"); err != nil {
							log.Println("Error while creating question tag mapping:", err)
							tx.Rollback()
							return err
						}
					}
				}
			}
		} else {
			questionUpdates := models.QuestionMain{
				QuestionID:    question.QuestionID,
				QuestionTitle: question.Title,
				QuestionType:  question.QuestionType,
			}
			err := s.assessmentRepo.UpdateQuestion(tx, &questionUpdates)
			if err != nil {
				log.Println("Error while Updating Question:", err)
				tx.Rollback()
				return err
			}

			// Handle question tags for updated question (with parent-child tag support)
			// First, delete existing tag mappings
			if err := tx.Where("question_id = ?", question.QuestionID).Delete(&models.QuestionTagMapping{}).Error; err != nil {
				log.Println("Error while deleting existing question tags:", err)
				tx.Rollback()
				return err
			}
			// Then add new tag mappings
			if len(question.Tags) > 0 {
				for _, tagReq := range question.Tags {
					tagIDs, err := s.assessmentRepo.ProcessTagRequest(tx, tagReq, "system")
					if err != nil {
						log.Println("Error while processing tag request:", err)
						tx.Rollback()
						return err
					}
					for _, tagID := range tagIDs {
						if err := s.assessmentRepo.CreateQuestionTagMappingWithParents(tx, question.QuestionID, tagID, "system"); err != nil {
							log.Println("Error while creating question tag mapping:", err)
							tx.Rollback()
							return err
						}
					}
				}
			}

			for _, option := range question.Options {
				if option.OptionID == 0 {
					newOpt := models.OptionMain{
						QuestionID:  question.QuestionID,
						OptionLabel: option.Label,
						SequenceID:  0,
						AnswerScore: option.Score,
						IsAnswer:    option.IsCorrect,
					}
					err := s.assessmentRepo.AddNewOption(tx, newOpt)
					if err != nil {
						log.Println("Error while Adding new option for question:", question.QuestionID, ":", err)
						tx.Rollback()
						return err
					}
				} else {
					optionUpdates := models.OptionMain{
						OptionID:    option.OptionID,
						AnswerScore: option.Score,
						IsAnswer:    option.IsCorrect,
						OptionLabel: option.Label,
					}
					err := s.assessmentRepo.UpdateOption(tx, &optionUpdates)
					if err != nil {
						log.Println("Error while Updating Option:", err)
						tx.Rollback()
						return err
					}
				}
			}
		}
	}

	if err := tx.Commit().Error; err != nil {
		return err
	}
	return nil
}

func (s *AssessmentServiceImpl) UpdateAssessmentStatusService(req models.UpdateAssessmentStatusRequest) error {
	updates := &models.DhlSurveySurveyExt{
		State:              req.AssessmentStatus,
		AssessmentSequence: req.AssessmentSequence,
	}
	tx := s.db.Begin()
	err := s.assessmentRepo.UpdateDHLAssessmentExt(tx, updates)
	if err != nil {
		tx.Rollback()
		return err
	}
	if err := tx.Commit().Error; err != nil {
		return err
	}
	return nil
}

func (s *AssessmentServiceImpl) GenerateAssessmentExcel(ctx context.Context, filter models.AssessmentReportFilter) ([]byte, error) {
	rows, err := s.assessmentRepo.GetAssessmentReport(ctx, filter)
	if err != nil {
		return nil, err
	}

	f := excelize.NewFile()
	sheet := "Report"
	f.SetSheetName("Sheet1", sheet)

	headers := []string{
		"Employee", "Team Lead", "Team Manager", "Sr Manager",
		"SDL", "SLL", "Skill Set", "Date",
		"Assessment Title", "Status", "Attempts",
		"Assigned", "Passed", "Failed", "Not Started",
		"Assessment Sequence",
	}

	// Write headers
	for i, h := range headers {
		cell, _ := excelize.CoordinatesToCellName(i+1, 1)
		f.SetCellValue(sheet, cell, h)
	}

	// Write data
	for rIndex, row := range rows {
		values := []interface{}{
			row.EmployeeName,
			row.TeamLead,
			row.TeamManager,
			row.SeniorManager,
			row.SDL,
			row.SLL,
			row.SkillSet,
			row.AssessmentDate.Format("2006-01-02 15:04:05"),
			row.AssessmentTitle,
			row.Status,
			row.Attempts,

			row.TotalAssigned,
			row.TotalPassed,
			row.TotalFailed,
			row.TotalNotStarted,

			row.AssessmentSequence,
		}

		for cIndex, v := range values {
			cell, _ := excelize.CoordinatesToCellName(cIndex+1, rIndex+2)
			f.SetCellValue(sheet, cell, v)
		}
	}

	// Export to bytes
	buf, err := f.WriteToBuffer()
	if err != nil {
		return nil, err
	}
	return buf.Bytes(), nil
}

func (s *AssessmentServiceImpl) CreateSessionImage(userID, sessionID string, imageData []byte) (*models.SessionImageResponse, error) {
	log.Printf("=== CreateSessionImage Service - User: %s, Session: %s ===", userID, sessionID)

	if sessionID == "" {
		log.Println("ERROR: session_id is empty")
		return nil, errors.New("session_id is required")
	}
	if len(imageData) == 0 {
		log.Println("ERROR: image data is empty")
		return nil, errors.New("image data is required")
	}
	log.Printf("Image data size: %d bytes", len(imageData))

	// Parse session UUID
	log.Printf("Parsing session UUID: %s", sessionID)
	sessionUUID, err := uuid.Parse(sessionID)
	if err != nil {
		log.Printf("ERROR: Failed to parse session UUID: %v", err)
		return nil, errors.New("invalid session_id format")
	}
	log.Printf("Successfully parsed UUID: %s", sessionUUID.String())

	// Verify session exists and belongs to user
	log.Printf("Verifying session exists for user: %s", userID)
	session, err := s.assessmentRepo.GetSessionByID(sessionID, userID)
	if err != nil {
		log.Printf("ERROR: Database error while fetching session: %v", err)
		return nil, err
	}
	if session == nil {
		log.Printf("ERROR: Session not found or does not belong to user. SessionID: %s, UserID: %s", sessionID, userID)
		return nil, errors.New("session not found or does not belong to user")
	}
	log.Printf("Session verified successfully. AssessmentID: %s", session.AssessmentID)

	// Create session image record
	currentTime := time.Now()
	sessionImage := &models.AssessmentUserSessionImage{
		SessionID:  sessionUUID,
		Image:      imageData,
		CreatedOn:  currentTime,
		CreatedBy:  userID,
		IsActive:   true,
		IsDeleted:  false,
		ModifiedOn: currentTime,
		ModifiedBy: userID,
	}
	log.Printf("Created session image record for session: %s", sessionUUID.String())

	// Save to database in transaction
	log.Println("Starting database transaction to save session image")
	err = s.db.Transaction(func(tx *gorm.DB) error {
		if err := s.assessmentRepo.CreateSessionImage(tx, sessionImage); err != nil {
			log.Printf("ERROR: Transaction failed - %v", err)
			return err
		}
		log.Printf("Session image saved successfully in database with ID: %d", sessionImage.ID)
		return nil
	})

	if err != nil {
		log.Printf("ERROR: Transaction error: %v", err)
		return nil, err
	}

	// Return response
	response := &models.SessionImageResponse{
		ID:        sessionImage.ID,
		SessionID: sessionImage.SessionID,
		CreatedOn: sessionImage.CreatedOn,
	}

	log.Printf("SUCCESS: Returning response with ID: %d", response.ID)
	return response, nil
}
