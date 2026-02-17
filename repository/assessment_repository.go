package repository

import (
	"context"
	"dhl/constant"
	"dhl/models"
	"dhl/utils"
	"errors"
	"fmt"
	"log"
	"time"

	"gorm.io/gorm"
)

type AssessmentRepository interface {
	GetAssessmentMstByAssmtSeq(id string) (*models.AssessmentMst, error)
	GetDhlSurveyExtByAssmtSeq(id string) (*models.DhlSurveySurveyExtResponse, error)
	GetOrCreateUserSession(tx *gorm.DB, userID, assessmentID, assessmentType string, partnerID int64) (*models.AssessmentUserSession, error)
	GetUserAssessmentsMap(userID string) ([]models.AssessmentStatus, error)
	GetUserAssessmentStatus(tx *gorm.DB, userID, assessmentID string) (*models.AssessmentStatus, error)
	GetAssessmentQuestions(assessmentSeq string) ([]models.AssessmentQuestionMst, error)
	GetContentByQuestionID(questionID int64) (models.QuestionContentWithType, error)
	GetOptionsByQuestionID(questionID int64) ([]models.OptionMst, error)
	GetContentByID(contentID int64) (models.ContentWithType, error)
	GetSessionByID(sessionID, userID string) (*models.AssessmentUserSession, error)
	GetCertificateDetailsBySessionId(sessionId string) (*models.CertificateDetails, error)
	GetAssessmentsForUserWithPagination(userId string, limit int, offset int) ([]*models.AssessmentListResponse, int64, error)
	GetAssessmentsForManagerWithPagination(managerID string, limit int, offset int) ([]*models.AssessmentListResponse, int64, error)
	GetUserAssessmentSessions(userId, assessmentSequence string) (interface{}, int64, error)
	GetAssessmentsWithPagination(limit int, offset int) ([]*models.AssessmentListResponse, int64, error)
	GetAssessmentsAttendeeInfo(assessmentID string, limit, offset int) ([]*models.AssessmentAttendeesInfo, int64, error)
	GetUserAssessmentResponse(userID *string, assessmentID, sessionId, userType string) (interface{}, error)
	GetAssessmentTypingResult(userId, assessmentSeqence string) (*models.AssessmentTypingResult, error)

	CreateAssessment(ctx context.Context, tx *gorm.DB, assessment *models.AssessmentMst) (*models.AssessmentMst, error)
	CreateAssessmentExt(ctx context.Context, tx *gorm.DB, ext *models.DhlSurveySurveyExt) (*models.DhlSurveySurveyExt, error)
	CreateAssessmentQuestions(ctx context.Context, tx *gorm.DB, assessmentID int64, assessmentSequence string, questionIDs []int64, createdBy string) error
	AddNewQuestion(tx *gorm.DB, question models.QuestionMain, assessmentSequence string) (int64, error)
	AddNewOption(tx *gorm.DB, option models.OptionMain) error
	SaveAssessmentWithQuestions(ctx context.Context, tx *gorm.DB, assessment models.SheetAssessment, userId string) (string, error)
	SaveAssessmentResponse(tx *gorm.DB, session *models.AssessmentUserSession, assessmentSeq string, response models.UserResponse) error
	UpdateAssessmentStatus(tx *gorm.DB, userID, assessmentID, status string) (*models.AssessmentStatus, error)
	AddManagerAssessmentMapping(tx *gorm.DB, userID, assessmentID string) (*models.ManagerAssessmentMapping, error)
	AddAssessmentTypingResult(tx *gorm.DB, input models.AssessmentTypingResult) (*models.AssessmentTypingResult, error)
	GetAssessmentReport(ctx context.Context, filter models.AssessmentReportFilter) ([]models.AssessmentReportRow, error)
	// Questions
	GetPaginatedQuestionsWithOptions(limit, offset int) ([]*models.QuestionDTO, int64, error)

	// Generic Update
	UpdateAssessment(db *gorm.DB, updates *models.AssessmentMst) error
	UpdateDHLAssessmentExt(db *gorm.DB, updates *models.DhlSurveySurveyExt) error
	UpdateAssessmentQuestion(db *gorm.DB, updates *models.AssessmentQuestionMst) error
	UpdateQuestion(db *gorm.DB, updates *models.QuestionMain) error
	UpdateOption(db *gorm.DB, updates *models.OptionMain) error

	// Session Image
	CreateSessionImage(tx *gorm.DB, sessionImage *models.AssessmentUserSessionImage) error
	GetSessionImageBySessionID(sessionID string) ([]byte, error)

	// Tag Operations
	GetOrCreateTagByName(tx *gorm.DB, tagName, createdBy string) (*models.TagMaster, error)
	GetOrCreateTagWithParent(tx *gorm.DB, tagName string, parentTagID *int64, createdBy string) (*models.TagMaster, error)
	ProcessTagRequest(tx *gorm.DB, tagReq models.TagRequest, createdBy string) ([]int64, error)
	GetAllParentTags(tx *gorm.DB, tagID int64) ([]int64, error)
	CreateAssessmentTagMapping(tx *gorm.DB, assessmentSequence string, tagID int64, createdBy string) error
	CreateAssessmentTagMappingWithParents(tx *gorm.DB, assessmentSequence string, tagID int64, createdBy string) error
	GetTagsByAssessmentSequence(assessmentSequence string) ([]string, error)
	GetTagRequestsByAssessmentSequence(assessmentSequence string) ([]models.TagRequest, error)

	// Question Tag Operations
	CreateQuestionTagMapping(tx *gorm.DB, questionID, tagID int64, createdBy string) error
	CreateQuestionTagMappingWithParents(tx *gorm.DB, questionID, tagID int64, createdBy string) error
	GetTagsByQuestionID(questionID int64) ([]string, error)
	GetTagsByQuestionIDs(questionIDs []int64) (map[int64][]string, error)
	GetTagRequestsByQuestionIDs(questionIDs []int64) (map[int64][]models.TagRequest, error)
}

type AssessmentRepositoryImpl struct {
	db *gorm.DB
}

func NewAssessmentRepository(db *gorm.DB) AssessmentRepository {
	return &AssessmentRepositoryImpl{db: db}
}

func (r *AssessmentRepositoryImpl) GetAssessmentMstByAssmtSeq(id string) (*models.AssessmentMst, error) {
	var a models.AssessmentMst
	if err := r.db.Where("assessment_sequence = ?", id).First(&a).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, nil
		}
		return nil, err
	}
	return &a, nil
}

func (r *AssessmentRepositoryImpl) GetDhlSurveyExtByAssmtSeq(id string) (*models.DhlSurveySurveyExtResponse, error) {
	var resp models.DhlSurveySurveyExtResponse

	query := `
        SELECT 
            sse.survey_id,
            sse.state,
            sse.time_limit,
            sse.deadline,
            sse.assessment_sequence,
            sse.attempts_limit,
            sse.certificate,
			sse.is_typing_test,
            dc.center_name,
            dsl.name AS service_line_name,
            dbp.name AS business_partner_name,
            dsbp.name AS sub_business_partner_name,
            dsg.name AS service_group_name,
            ds.service_name

        FROM dhl_survey_survey_ext sse
        LEFT JOIN dhl_center dc 
            ON dc.center_id = sse.center_id

        LEFT JOIN dhl_service_line dsl
            ON dsl.service_line_id = sse.service_line_id

        LEFT JOIN dhl_business_partner dbp
            ON dbp.business_partner_id = sse.business_partner_id

        LEFT JOIN dhl_sub_business_partner dsbp
            ON dsbp.sub_business_partner_id = sse.sub_business_partner_id

        LEFT JOIN dhl_service_group dsg
            ON dsg.service_grp_id = sse.service_group_id

        LEFT JOIN dhl_service ds
            ON ds.service_id = sse.service_id

        WHERE sse.assessment_sequence = ?
    `

	if err := r.db.Raw(query, id).Scan(&resp).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, nil
		}
		return nil, err
	}

	return &resp, nil
}

func (r *AssessmentRepositoryImpl) GetOrCreateUserSession(tx *gorm.DB, userID, assessmentID, assessmentType string, partnerID int64) (*models.AssessmentUserSession, error) {
	var s models.AssessmentUserSession
	createdOn := time.Now()
	err := tx.Where("user_id = ? AND assessment_id = ? AND is_active = true", userID, assessmentID).First(&s).Error
	if errors.Is(err, gorm.ErrRecordNotFound) {
		s = models.AssessmentUserSession{
			SessionID:      models.GenerateUUID(),
			UserID:         userID,
			AssessmentID:   assessmentID,
			AssessmentType: assessmentType,
			PartnerID:      partnerID,
			IsActive:       true,
			IsDeleted:      false,
			AccessTime:     createdOn,
			CreatedOn:      createdOn,
			CreatedBy:      userID,
		}
		if err := tx.Create(&s).Error; err != nil {
			return nil, err
		}
		return &s, nil
	}

	if err == nil {
		s.AccessTime = createdOn
		s.ModifiedOn = createdOn
		s.ModifiedBy = userID
		if err := tx.Save(&s).Error; err != nil {
			return nil, err
		}
		return &s, nil
	}
	return &s, nil
}

func (r *AssessmentRepositoryImpl) GetUserAssessmentsMap(userID string) ([]models.AssessmentStatus, error) {
	var status []models.AssessmentStatus
	err := r.db.Where("user_id = ?", userID).Find(&status).Error
	return status, err
}

func (r *AssessmentRepositoryImpl) GetAssessmentsForUserWithPagination(userId string, limit int, offset int) ([]*models.AssessmentListResponse, int64, error) {
	var assessments []*models.AssessmentListResponse
	var totalRecords int64
	params := []interface{}{userId}
	countQuery := `
         SELECT COUNT(*)
    FROM assessment_mst am
    JOIN dhl_survey_survey_ext sse 
        ON am.assessment_sequence = sse.assessment_sequence
    JOIN assessment_status ast
        ON ast.assessment_id = am.assessment_sequence
    LEFT JOIN job_descriptions jd
        ON jd.job_id = am.job_id
    WHERE ast.user_id = ?
        `

	if err := r.db.Raw(countQuery, userId).Scan(&totalRecords).Error; err != nil {
		return nil, 0, err
	}

	query := `
       SELECT 
    am.assessment_id,
    am.assessment_sequence,
    am.assessment_desc AS assessment_title,
    sse.time_limit,
    am.marks,
    sse.deadline,
    sse.state as assessment_status,
    sse.show_result,
    sse.certificate,
    dc.center_name,
    dsl.name AS service_line_name,
    dbp.name AS business_partner_name,
    dsbp.name AS sub_business_partner_name,
    dsg.name AS service_group_name,
    ds.service_name AS service_name,
    ast.assessment_status AS state,
    jd.title AS job_title
FROM assessment_mst am
JOIN dhl_survey_survey_ext sse 
    ON am.assessment_sequence = sse.assessment_sequence
JOIN assessment_status ast
    ON ast.assessment_id = am.assessment_sequence
LEFT JOIN job_descriptions jd
    ON jd.job_id = am.job_id
LEFT JOIN dhl_center dc
    ON dc.center_id = sse.center_id
LEFT JOIN dhl_service_line dsl
    ON dsl.service_line_id = sse.service_line_id
LEFT JOIN dhl_business_partner dbp
    ON dbp.business_partner_id = sse.business_partner_id
LEFT JOIN dhl_sub_business_partner dsbp
    ON dsbp.sub_business_partner_id = sse.sub_business_partner_id
LEFT JOIN dhl_service_group dsg
    ON dsg.service_grp_id = sse.service_group_id
LEFT JOIN dhl_service ds
    ON ds.service_id = sse.service_id
WHERE ast.user_id = ?
ORDER BY am.assessment_id DESC
LIMIT ? OFFSET ?

    `

	params = append(params, limit, offset)

	if err := r.db.Raw(query, params...).Scan(&assessments).Error; err != nil {
		return nil, 0, err
	}

	return assessments, totalRecords, nil
}

func (r *AssessmentRepositoryImpl) GetAssessmentsForManagerWithPagination(managerID string, limit int, offset int) ([]*models.AssessmentListResponse, int64, error) {
	var assessments []*models.AssessmentListResponse
	var totalRecords int64
	params := []interface{}{managerID}
	countQuery := `
        SELECT COUNT(*)
        FROM assessment_mst am
        JOIN dhl_survey_survey_ext sse 
            ON am.assessment_sequence = sse.assessment_sequence
        JOIN manager_assessment_mapping msm
            ON msm.assessment_sequence = am.assessment_sequence
		WHERE msm.manager_id = ?
        `

	if err := r.db.Raw(countQuery, managerID).Scan(&totalRecords).Error; err != nil {
		return nil, 0, err
	}

	query := `
        SELECT 
            am.assessment_id,
            am.assessment_sequence,
            am.assessment_desc AS assessment_title,
            sse.time_limit,
            am.marks,
            sse.deadline,
            sse.state,
            dc.center_name,
            dsl.name AS service_line_name,
            dbp.name AS business_partner_name,
            dsbp.name AS sub_business_partner_name,
            dsg.name AS service_group_name,
            ds.service_name AS service_name
        FROM assessment_mst am
        JOIN dhl_survey_survey_ext sse 
            ON am.assessment_sequence = sse.assessment_sequence
		JOIN manager_assessment_mapping msm
            ON msm.assessment_sequence = am.assessment_sequence
        LEFT JOIN dhl_center dc
            ON dc.center_id = sse.center_id
        LEFT JOIN dhl_service_line dsl
            ON dsl.service_line_id = sse.service_line_id
        LEFT JOIN dhl_business_partner dbp
            ON dbp.business_partner_id = sse.business_partner_id
        LEFT JOIN dhl_sub_business_partner dsbp
            ON dsbp.sub_business_partner_id = sse.sub_business_partner_id
        LEFT JOIN dhl_service_group dsg
            ON dsg.service_grp_id = sse.service_group_id
        LEFT JOIN dhl_service ds
            ON ds.service_id = sse.service_id
		 WHERE msm.manager_id = ?
        ORDER BY am.assessment_id DESC
        LIMIT ? OFFSET ?
    `

	params = append(params, limit, offset)

	if err := r.db.Raw(query, params...).Scan(&assessments).Error; err != nil {
		return nil, 0, err
	}

	return assessments, totalRecords, nil
}

func (r *AssessmentRepositoryImpl) GetUserAssessmentSessions(userId, assessmentSequence string) (interface{}, int64, error) {
	var totalCount int64
	var userSessions []models.AssessmentUserSession
	countQuery := `
	select count(*) from assessment_user_session where user_id = ? and assessment_id = ?
	`

	if err := r.db.Raw(countQuery, userId, assessmentSequence).Scan(&totalCount).Error; err != nil {
		return nil, 0, err
	}

	if err := r.db.Model(models.AssessmentUserSession{}).
		Where("user_id = ? and assessment_id = ?", userId, assessmentSequence).Scan(&userSessions).Error; err != nil {
		return nil, 0, err
	}

	return userSessions, totalCount, nil

}

func (r *AssessmentRepositoryImpl) GetUserAssessmentStatus(tx *gorm.DB, userID, assessmentID string) (*models.AssessmentStatus, error) {
	var status models.AssessmentStatus
	if err := tx.Where("user_id = ? AND assessment_id = ?", userID, assessmentID).First(&status).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			status = models.AssessmentStatus{
				UserID:           userID,
				AssessmentID:     assessmentID,
				AssessmentStatus: "STARTED",
				IsActive:         true,
				IsDeleted:        false,
				CreatedBy:        userID,
				CreatedOn:        time.Now(),
			}
			if err := tx.Create(&status).Error; err != nil {
				return nil, err
			}
			return &status, nil
		}
		return nil, err
	}
	return &status, nil
}

func (r *AssessmentRepositoryImpl) GetAssessmentQuestions(assessmentSeq string) ([]models.AssessmentQuestionMst, error) {
	var list []models.AssessmentQuestionMst
	err := r.db.Where("assessment_sequence = ? AND is_deleted = false", assessmentSeq).Order("sequence_id").Find(&list).Error
	return list, err
}

func (r *AssessmentRepositoryImpl) GetContentByQuestionID(questionID int64) (models.QuestionContentWithType, error) {
	var list models.QuestionContentWithType
	err := r.db.Raw(`
		SELECT 
			c.content_id,
			ct.content_type,
			c.font,
			c.value,
			q.question_type_id,
			tc.question_type			
		FROM content_mst c
		INNER JOIN content_type_config ct ON ct.content_type_id = c.content_type_id
		INNER JOIN question_mst q ON q.content_id = c.content_id
		LEFT JOIN question_type_config tc ON tc.question_type_id=q.question_type_id
		WHERE q.question_id = ?`, questionID).First(&list).Error
	return list, err
}

func (r *AssessmentRepositoryImpl) GetOptionsByQuestionID(questionID int64) ([]models.OptionMst, error) {
	var list []models.OptionMst
	err := r.db.Where("question_id = ?", questionID).Find(&list).Error
	return list, err
}

func (r *AssessmentRepositoryImpl) GetContentByID(contentID int64) (models.ContentWithType, error) {
	var content models.ContentWithType
	err := r.db.Raw(`
	SELECT c.content_id,ct.content_type,c.font,c.value
	FROM content_mst c
	INNER JOIN content_type_config ct ON ct.content_type_id = c.content_type_id
	WHERE c.content_id = ?
	`, contentID).Find(&content).Error
	return content, err
}

func (r *AssessmentRepositoryImpl) GetSessionByID(sessionID, userID string) (*models.AssessmentUserSession, error) {
	var session models.AssessmentUserSession
	err := r.db.Where("session_id = ? AND user_id= ?", sessionID, userID).First(&session).Error
	if errors.Is(err, gorm.ErrRecordNotFound) {
		return nil, nil
	}
	return &session, err
}

func (r *AssessmentRepositoryImpl) GetAssessmentTypingResult(userId, assessmentSeqence string) (*models.AssessmentTypingResult, error) {
	var typingResult models.AssessmentTypingResult
	err := r.db.Where("assessment_sequence = ? AND user_id = ?", assessmentSeqence, userId).First(&typingResult).Error
	if err != nil {
		return nil, err
	}
	return &typingResult, nil
}

func (r *AssessmentRepositoryImpl) GetCertificateDetailsBySessionId(sessionId string) (*models.CertificateDetails, error) {
	var canGenerate bool
	query := `
		SELECT  u.user_id, CONCAT(u.first_name, ' ', u.last_name) AS user_name,
            a.assessment_sequence, a.assessment_desc AS assessment_title,
            s.session_id,
			(SELECT COALESCE(SUM(point_assigned),0) FROM public.assessment_result 
			WHERE assessment_sequence = a.assessment_sequence AND assessment_session_id = s.session_id::text AND is_deleted = false
			) as marks_obtained,
			a.marks AS total_marks,
			a.passing_score,
			s.created_on AS completed_at
        FROM assessment_user_session s
        JOIN assessment_mst a ON a.assessment_sequence = s.assessment_id
        JOIN assessment_user_mst u ON u.user_id::text = s.user_id
		JOIN assessment_result ar on ar.assessment_session_id::text=s.session_id::text
        WHERE s.session_id = ?
	`
	var details models.CertificateDetails
	if err := r.db.Raw(query, sessionId).Scan(&details).Error; err != nil {
		return nil, err
	}
	log.Println(details)
	if err := r.db.Raw(`select certificate from public.dhl_survey_survey_ext where assessment_sequence= ?`, details.AssessmentSequence).Scan(&canGenerate).Error; err != nil {
		return nil, err
	}
	if canGenerate == false {
		return nil, errors.New("can not generate certificate to this assessment")
	}

	return &details, nil
}

func (r *AssessmentRepositoryImpl) SaveAssessmentResponse(tx *gorm.DB, session *models.AssessmentUserSession, assessmentSeq string, response models.UserResponse) error {
	currentTime := time.Now()
	for _, id := range response.SelectedOptionID {
		var selectedOption models.OptionMst
		if err := tx.Where("option_id = ? AND question_id = ?", id, response.QuestionID).First(&selectedOption).Error; err != nil {
			return err
		}

		result := models.AssessmentResult{
			CreatedOn:           currentTime,
			IsActive:            true,
			IsDeleted:           false,
			ModifiedOn:          currentTime,
			ActivityID:          1,
			AssessmentSequence:  assessmentSeq,
			AssessmentSessionID: session.SessionID.String(),
			AttemptID:           0,
			AttemptOptionID:     id,
			QuestionID:          response.QuestionID,
			AttemptEndTime:      &currentTime,
			AttemptStartTime:    &currentTime,
			PointAssigned:       int64(selectedOption.AnswerScore),
		}
		if err := tx.Create(&result).Error; err != nil {
			return err
		}
	}
	if err := tx.Model(&models.AssessmentUserSession{}).
		Where("session_id = ?", session.SessionID).
		Updates(map[string]interface{}{
			"is_active":   false,
			"modified_on": currentTime,
		}).Error; err != nil {
		return err
	}
	return nil
}

func (r *AssessmentRepositoryImpl) UpdateAssessmentStatus(tx *gorm.DB, userID, assessmentID, status string) (*models.AssessmentStatus, error) {
	var st models.AssessmentStatus

	err := tx.Where("user_id = ? AND assessment_id = ?", userID, assessmentID).
		First(&st).Error

	if errors.Is(err, gorm.ErrRecordNotFound) {
		st = models.AssessmentStatus{
			UserID:           userID,
			AssessmentID:     assessmentID,
			AssessmentStatus: status,
			CreatedOn:        time.Now(),
			CreatedBy:        userID,
			IsActive:         true,
			IsDeleted:        false,
			ModifiedOn:       time.Now(),
			ModifiedBy:       userID,
		}
		return &st, tx.Create(&st).Error
	}
	err = tx.Model(&st).Where("id = ?", st.ID).
		Updates(map[string]interface{}{
			"assessment_status": status,
		}).Error
	if err != nil {
		return nil, err
	}
	return &st, nil
}

func (r *AssessmentRepositoryImpl) AddManagerAssessmentMapping(tx *gorm.DB, userID, assessmentID string) (*models.ManagerAssessmentMapping, error) {
	mapping := models.ManagerAssessmentMapping{
		ManagerID:          userID,
		AssessmentSequence: assessmentID,
		IsActive:           true,
		CreatedOn:          time.Now(),
	}
	return &mapping, tx.Create(&mapping).Error
}

func (r *AssessmentRepositoryImpl) AddAssessmentTypingResult(tx *gorm.DB, input models.AssessmentTypingResult) (*models.AssessmentTypingResult, error) {
	err := tx.Create(&input).Error
	if err != nil {
		return nil, err
	}
	return &input, nil
}

func (r *AssessmentRepositoryImpl) UpdateAssessment(db *gorm.DB, updates *models.AssessmentMst) error {
	if updates.AssessmentID == 0 {
		return errors.New("assessment_id required for update")
	}
	updateMap := utils.BuildUpdateMap(updates)
	return db.Model(&models.AssessmentMst{}).
		Where("assessment_id = ?", updates.AssessmentID).
		Updates(updateMap).Error
}

func (r *AssessmentRepositoryImpl) UpdateDHLAssessmentExt(db *gorm.DB, updates *models.DhlSurveySurveyExt) error {
	if updates.AssessmentSequence == "" {
		return errors.New("assessment_sequence required for update")
	}
	updateMap := utils.BuildUpdateMap(updates)
	return db.Model(&models.DhlSurveySurveyExt{}).
		Where("assessment_sequence = ?", updates.AssessmentSequence).
		Updates(updateMap).Error
}

func (r *AssessmentRepositoryImpl) UpdateAssessmentQuestion(db *gorm.DB, updates *models.AssessmentQuestionMst) error {
	if updates.AssessmentQuestionID == 0 {
		return errors.New("assessment_question_id required for update")
	}
	updateMap := utils.BuildUpdateMap(updates)
	return db.Model(&models.AssessmentQuestionMst{}).
		Where("assessment_question_id = ?", updates.AssessmentQuestionID).
		Updates(updateMap).Error
}

func (r *AssessmentRepositoryImpl) UpdateQuestion(tx *gorm.DB, question *models.QuestionMain) error {
	var qmst models.QuestionMst
	if err := tx.Where("question_id = ?", question.QuestionID).First(&qmst).Error; err != nil {
		return err
	}
	cUpdates := &models.ContentMst{
		Value: question.QuestionTitle,
	}
	cUpdateMap := utils.BuildUpdateMap(cUpdates)

	err := tx.Model(&models.ContentMst{}).Where("content_id = ?", qmst.ContentID).Updates(cUpdateMap).Error
	if err != nil {
		return err
	}

	qUpdates := &models.QuestionMst{
		ModifiedOn: time.Now(),
	}
	qUpdateMap := utils.BuildUpdateMap(qUpdates)

	err = tx.Model(&models.QuestionMst{}).Where("question_id = ?", question.QuestionID).Updates(qUpdateMap).Error
	if err != nil {
		return err
	}
	return nil

}

func (r *AssessmentRepositoryImpl) UpdateOption(tx *gorm.DB, option *models.OptionMain) error {
	var opst models.OptionMst
	if err := tx.Where("option_id = ?", option.OptionID).First(&opst).Error; err != nil {
		return err
	}

	cUpdates := &models.ContentMst{
		Value: option.OptionLabel,
	}
	cUpdateMap := utils.BuildUpdateMap(cUpdates)

	err := tx.Model(&models.ContentMst{}).Where("content_id = ?", opst.ContentID).Updates(cUpdateMap).Error
	if err != nil {
		return err
	}

	oUpdates := &models.OptionMain{
		IsAnswer:    option.IsAnswer,
		AnswerScore: option.AnswerScore,
	}
	updateMap := utils.BuildUpdateMap(oUpdates)
	err = tx.Model(&models.OptionMst{}).
		Where("option_id = ?", option.OptionID).
		Updates(updateMap).Error
	if err != nil {
		return err
	}
	return nil
}

// Admin Repo functions
func (r *AssessmentRepositoryImpl) GetAssessmentsWithPagination(limit int, offset int) ([]*models.AssessmentListResponse, int64, error) {
	var assessments []*models.AssessmentListResponse
	var totalRecords int64

	countQuery := `
	SELECT COUNT(*)
	FROM assessment_mst am
	JOIN dhl_survey_survey_ext sse 
		ON am.assessment_sequence = sse.assessment_sequence
	`
	if err := r.db.Raw(countQuery).Scan(&totalRecords).Error; err != nil {
		return nil, 0, err
	}

	query := `
SELECT 
    am.assessment_id,
    am.assessment_sequence,
    am.assessment_desc AS assessment_title,
	am.assessment_type,  
    sse.time_limit,
    am.marks,
    am.valid_to AS deadline,
    sse.state,
    jd.title AS job_title,              
    dc.center_name,
    dsl.name AS service_line_name,
    dbp.name AS business_partner_name,
    dsbp.name AS sub_business_partner_name,
    dsg.name AS service_group_name,
    ds.service_name AS service_name
FROM assessment_mst am
JOIN dhl_survey_survey_ext sse 
    ON am.assessment_sequence = sse.assessment_sequence
LEFT JOIN job_descriptions jd              
    ON am.job_id = jd.job_id               
LEFT JOIN dhl_center dc
    ON dc.center_id = sse.center_id
LEFT JOIN dhl_service_line dsl
    ON dsl.service_line_id = sse.service_line_id
LEFT JOIN dhl_business_partner dbp
    ON dbp.business_partner_id = sse.business_partner_id
LEFT JOIN dhl_sub_business_partner dsbp
    ON dsbp.sub_business_partner_id = sse.sub_business_partner_id
LEFT JOIN dhl_service_group dsg
    ON dsg.service_grp_id = sse.service_group_id
LEFT JOIN dhl_service ds
    ON ds.service_id = sse.service_id
ORDER BY am.created_on DESC
LIMIT ? OFFSET ?
`

	if err := r.db.Raw(query, limit, offset).Scan(&assessments).Error; err != nil {
		return nil, 0, err
	}

	return assessments, totalRecords, nil
}


func (r *AssessmentRepositoryImpl) GetAssessmentsAttendeeInfo(assessmentID string, limit, offset int) ([]*models.AssessmentAttendeesInfo, int64, error) {

	var attendees []*models.AssessmentAttendeesInfo
	var totalCount int64

	countQuery := `
		SELECT COUNT(DISTINCT user_id)
FROM assessment_status
WHERE assessment_id = ?;

	`

	if err := r.db.Raw(countQuery, assessmentID).Scan(&totalCount).Error; err != nil {
		return nil, 0, err
	}
	log.Println("AssessmentID received:", assessmentID)

	mainQuery := `
		SELECT 
    ast.assessment_id,
    ast.user_id,
    ast.assessment_status,
    u.first_name,
    u.last_name,
    u.email,
    dc.center_name
FROM assessment_status ast
LEFT JOIN assessment_user_mst u
    ON ast.user_id = u.user_id::text
LEFT JOIN dhl_assessment_user_mst_ext ue
    ON ue.user_id = u.user_id::text
LEFT JOIN dhl_center dc
    ON dc.center_id = ue.center
WHERE ast.assessment_id = ?
ORDER BY ast.modified_on DESC
LIMIT ? OFFSET ?;


	`

	err := r.db.Raw(mainQuery, assessmentID, limit, offset).Scan(&attendees).Error
	if err != nil {
		return nil, 0, err
	}

	return attendees, totalCount, nil
}

func (r *AssessmentRepositoryImpl) GetUserAssessmentResponse(userID *string, assessmentSeq, sessionID, userType string) (interface{}, error) {
	var assessmentDetails models.AssessmentListResponse
	var userResponses []*models.AssessmentUserResponse
	if userID != nil {
		session, _ := r.GetSessionByID(sessionID, *userID)
		if session == nil {
			return nil, errors.New("invalid session")
		}
	}
	// --- 1️⃣ Fetch assessment details (only once)
	assessmentQuery := `
		SELECT DISTINCT
			am.assessment_id,
			am.assessment_sequence,
			am.assessment_desc AS assessment_title,
			am.marks,
			am.passing_score,
			dc.center_name,
			dsl.name AS service_line_name,
			dsp.name AS sub_business_partner_name,
			dsg.name AS service_group_name,
			dse.show_result,
			dse.certificate,
			dse.is_typing_test,
			ds.service_name,
			(
				SELECT COALESCE(SUM(ar.point_assigned), 0)
				FROM assessment_result ar
				WHERE ar.assessment_sequence = am.assessment_sequence
				AND ar.assessment_session_id = ?
				AND ar.is_deleted = false
			) AS marks_obtained,
			ROUND(
					(
						(
							SELECT COALESCE(SUM(ar.point_assigned), 0)
							FROM assessment_result ar
							WHERE ar.assessment_sequence = am.assessment_sequence
							AND ar.assessment_session_id = ?
							AND ar.is_deleted = false
						) * 100.0 / NULLIF(am.marks, 0)
					),
					2
				) AS user_score,
			CASE 
				WHEN (
						(
							SELECT COALESCE(SUM(ar.point_assigned), 0)
							FROM assessment_result ar
							WHERE ar.assessment_sequence = am.assessment_sequence
							AND ar.assessment_session_id = ?
							AND ar.is_deleted = false
						) * 100.0 / NULLIF(am.marks, 0)
					) >= am.passing_score THEN 'Pass'
				ELSE 'Fail'
			END AS result_status
		FROM assessment_mst am
		LEFT JOIN dhl_survey_survey_ext dse ON dse.assessment_sequence = am.assessment_sequence
		LEFT JOIN audit_sub_business_partner_mapping sp_map ON sp_map.odoo_server_id = dse.sub_business_partner_id
		LEFT JOIN dhl_sub_business_partner dsp ON dsp.sub_business_partner_id = sp_map.new_server_id
		LEFT JOIN audit_center_mapping acm ON acm.old_center_id = dse.center_id
		LEFT JOIN dhl_center dc ON dc.center_id = acm.new_center_id
		LEFT JOIN audit_service_line_mapping aslm ON aslm.odoo_server_id = dse.service_line_id
		LEFT JOIN dhl_service_line dsl ON dsl.service_line_id = aslm.new_server_id
		LEFT JOIN audit_service_group_mapping asgm ON asgm.odoo_server_id = dse.service_group_id
		LEFT JOIN dhl_service_group dsg ON dsg.service_grp_id = asgm.new_server_id
		LEFT JOIN audit_service_mapping asm ON asm.odoo_service_id = dse.service_id
		LEFT JOIN dhl_service ds ON ds.service_id = asm.new_service_id
		WHERE am.assessment_sequence = ?
		LIMIT 1;
	`
	if err := r.db.Raw(assessmentQuery, sessionID, sessionID, sessionID, assessmentSeq).Scan(&assessmentDetails).Error; err != nil {
		return nil, err
	}
	if userType == string(constant.User) && assessmentDetails.ShowResult == false {
		return nil, errors.New("results not available to users")
	}

	if assessmentDetails.IsTypingTest == true {
		var userId string
		if userID != nil {
			userId = *userID
		} else {
			err := r.db.Table("assessment_user_session").Select("user_id").Where("session_id = ?", sessionID).Scan(&userId).Error
			if err != nil {
				return nil, fmt.Errorf("error while getting user session %w", err)
			}
		}

		typingResult, _ := r.GetAssessmentTypingResult(userId, assessmentSeq)
		if typingResult != nil {
			assessmentDetails.AttemptedTypingTest = true
			assessmentDetails.TypingResult = typingResult
		}
	}

	// --- 2️⃣ Fetch user responses for given session
	responseQuery := `
		SELECT 
			q.question_id,
			qc.value AS question_text,

			-- Selected option (null if not attempted)
			ur.attempt_option_id AS selected_option_id,
			oc_user.value AS selected_option_text,

			-- Points (null for skipped)
			ur.point_assigned,

			-- Answer Status
			CASE 
				WHEN ur.attempt_option_id IS NULL THEN 'Unattempted'
				WHEN om_user.is_answer = TRUE THEN 'Correct'
				ELSE 'Incorrect'
			END AS answer_status,

			-- Skipped flag
			CASE 
				WHEN ur.attempt_option_id IS NULL THEN TRUE 
				ELSE FALSE 
			END AS skipped

		FROM assessment_question_mst aq
		JOIN question_mst q 
			ON aq.question_id = q.question_id
		JOIN content_mst qc 
			ON q.content_id = qc.content_id

		-- LEFT JOIN to get attempt data (may not exist)
		LEFT JOIN assessment_result ur
			ON ur.question_id = aq.question_id
		AND ur.assessment_sequence = aq.assessment_sequence
		AND ur.assessment_session_id = ?
		AND ur.is_deleted = FALSE

		-- Attempted option join
		LEFT JOIN option_mst om_user 
			ON ur.attempt_option_id = om_user.option_id
		LEFT JOIN content_mst oc_user 
			ON om_user.content_id = oc_user.content_id

		-- Correct options join
		LEFT JOIN option_mst om_correct 
			ON om_correct.question_id = q.question_id
		AND om_correct.is_answer = TRUE
		LEFT JOIN content_mst oc_correct 
			ON om_correct.content_id = oc_correct.content_id

		WHERE 
			aq.assessment_sequence = ?
			AND aq.is_deleted = FALSE

		GROUP BY 
			q.question_id,
			qc.value,
			ur.attempt_option_id,
			oc_user.value,
			ur.point_assigned,
			om_user.is_answer

		ORDER BY q.question_id;
	`

	if err := r.db.Raw(responseQuery, sessionID, assessmentSeq).Scan(&userResponses).Error; err != nil {
		return nil, err
	}
	for _, usrRes := range userResponses {
		options, err := r.GetOptionsByQuestionID(usrRes.QuestionID)
		if err != nil {
			return nil, err
		}

		var answers []models.Answer
		for _, opt := range options {
			optContents, _ := r.GetContentByID(opt.ContentID)
			answers = append(answers, models.Answer{
				AnswerID:      int(opt.OptionID),
				Sequence:      &opt.SequenceID,
				OptionLabel:   optContents.Value,
				CorrectAnswer: opt.IsAnswer,
			})
		}

		usrRes.Answers = answers
	}

	// Fetch session image if exists
	log.Printf("Fetching session image for session: %s", sessionID)
	sessionImage, err := r.GetSessionImageBySessionID(sessionID)
	if err != nil {
		log.Printf("WARNING: Failed to fetch session image: %v", err)
		// Continue without image - don't fail the entire request
	}
	if sessionImage != nil {
		log.Printf("Session image found, size: %d bytes", len(sessionImage))
		assessmentDetails.SessionImage = sessionImage
	} else {
		log.Println("No session image found for this session")
	}

	response := map[string]interface{}{
		"assessmentDetails": &assessmentDetails,
		"userResponses":     userResponses,
	}

	return response, nil
}

func (r *AssessmentRepositoryImpl) AddNewQuestion(tx *gorm.DB, question models.QuestionMain, assessmentSequence string) (int64, error) {
	content := models.ContentMst{
		ContentTypeID: 1,
		Value:         question.QuestionTitle,
	}
	if err := tx.Create(&content).Error; err != nil {
		return 0, err
	}
	qmst := models.QuestionMst{
		ContentID:  content.ContentID,
		IsActive:   true,
		IsDeleted:  false,
		CreatedOn:  time.Now(),
		ModifiedOn: time.Now(),
	}
	if err := tx.Create(&qmst).Error; err != nil {
		return 0, err
	}
	asmtQ := models.AssessmentQuestionMst{
		QuestionID:         qmst.QuestionID,
		CreatedOn:          time.Now(),
		IsActive:           true,
		IsDeleted:          false,
		ModifiedOn:         time.Now(),
		AssessmentSequence: assessmentSequence,
	}
	log.Println("assessmentSequence:", assessmentSequence)
	log.Println(asmtQ)
	if err := tx.Create(&asmtQ).Error; err != nil {
		return 0, err
	}
	return qmst.QuestionID, nil

}

func (r *AssessmentRepositoryImpl) AddNewOption(tx *gorm.DB, option models.OptionMain) error {
	content := models.ContentMst{
		ContentTypeID: 1,
		Value:         option.OptionLabel,
	}
	if err := tx.Create(&content).Error; err != nil {
		return err
	}
	omst := models.OptionMst{
		ContentID:   content.ContentID,
		IsAnswer:    option.IsAnswer,
		QuestionID:  option.QuestionID,
		SequenceID:  option.SequenceID,
		AnswerScore: option.AnswerScore,
	}
	if err := tx.Create(&omst).Error; err != nil {
		return err
	}

	return nil

}

func (r *AssessmentRepositoryImpl) SaveAssessmentWithQuestions(ctx context.Context, tx *gorm.DB, assessment models.SheetAssessment, userId string) (string, error) {
	createdAt := time.Now()
	// 1️⃣ Insert into assessment_mst using struct
	assessmentMst := models.AssessmentMst{
		AssessmentDesc: assessment.AssessmentName,
		CreatedOn:      createdAt,
		CreatedBy:      userId,
		IsActive:       true,
		IsDeleted:      false,
		ModifiedOn:     createdAt,
		ModifiedBy:     userId,
		Marks:          int64(len(assessment.Questions)),
		AssessmentType: "survey",
	}
	if err := tx.WithContext(ctx).Create(&assessmentMst).Error; err != nil {
		return "", fmt.Errorf("failed to insert assessment_mst: %w", err)
	}

	assessmentID := assessmentMst.AssessmentID

	assessmentSequence := "D0000" + fmt.Sprintf("%d", assessmentID)

	if err := tx.WithContext(ctx).
		Table("assessment_mst").
		Where("assessment_id = ?", assessmentID).
		Update("assessment_sequence", assessmentSequence).Error; err != nil {
		return "", fmt.Errorf("failed to update assessment_sequence: %w", err)
	}

	// 2️⃣ Insert into dhl_survey_survey_ext using struct
	surveyExt := models.DhlSurveySurveyExt{
		SurveyID:           assessmentID,
		State:              string(constant.Draft),
		AssessmentSequence: assessmentSequence,
		Certificate:        true,
		AttemptsLimit:      1,
		ShowResult:         true,
		TimeLimit:          30,
		IsTypingTest:       false,
	}
	if err := tx.WithContext(ctx).Create(&surveyExt).Error; err != nil {
		return "", fmt.Errorf("failed to insert dhl_survey_survey_ext: %w", err)
	}

	// 3️⃣ Loop through questions and collect question IDs
	var questionIDs []int64

	for _, q := range assessment.Questions {
		// Insert question text in content_mst
		contentQ := models.ContentMst{
			ContentTypeID: 1,
			Value:         q.Title,
		}
		if err := tx.WithContext(ctx).Create(&contentQ).Error; err != nil {
			return "", fmt.Errorf("failed to insert question content: %w", err)
		}

		questionContentID := contentQ.ContentID

		// Insert into question_mst
		questionTypeID := 0
		if id, ok := utils.QuestionTypeMap[q.QuestionType]; ok {
			questionTypeID = id
		}

		questionMst := models.QuestionMst{
			ContentID:      questionContentID,
			QuestionTypeID: int64(questionTypeID),
			IsActive:       true,
			IsDeleted:      false,
			CreatedOn:      createdAt,
			ModifiedOn:     createdAt,
		}
		if err := tx.WithContext(ctx).Create(&questionMst).Error; err != nil {
			return "", fmt.Errorf("failed to insert question_mst: %w", err)
		}

		questionID := questionMst.QuestionID
		questionIDs = append(questionIDs, questionID) // Collect question ID

		// Link question to assessment
		// Note: skipping_allowed should be opposite of mandatory_to_answer
		skippingAllowed := !q.MandatoryToAnswer

		assessmentQ := models.AssessmentQuestionMst{
			AssessmentID:       int(assessmentID),
			QuestionID:         questionID,
			AssessmentSequence: assessmentSequence,
			IsActive:           true,
			IsDeleted:          false,
			SkippingAllowed:    skippingAllowed,
			CreatedOn:          createdAt,
			CreatedBy:          userId,
			ModifiedOn:         createdAt,
			ModifiedBy:         userId,
			SequenceID:         0,
			CorrectPoints:      0,
			DurationInSeconds:  0,
			NegativePoints:     0,
			DifficultyLevel:    "",
		}

		if err := tx.WithContext(ctx).Create(&assessmentQ).Error; err != nil {
			return "", fmt.Errorf("failed to insert assessment_question_mst: %w", err)
		}

		// 4️⃣ Loop through options
		for _, opt := range q.Options {
			contentOpt := models.ContentMst{
				ContentTypeID: 1,
				Value:         opt.Label,
			}
			if err := tx.WithContext(ctx).Create(&contentOpt).Error; err != nil {
				return "", fmt.Errorf("failed to insert option content: %w", err)
			}

			optionContentID := contentOpt.ContentID

			optionMst := models.OptionMst{
				ContentID:   optionContentID,
				IsAnswer:    opt.IsCorrect,
				QuestionID:  questionID,
				AnswerScore: opt.Score,
			}
			if err := tx.WithContext(ctx).Create(&optionMst).Error; err != nil {
				return "", fmt.Errorf("failed to insert option_mst: %w", err)
			}
		}

		// Handle question tags if provided (with parent-child tag support)
		if len(q.Tags) > 0 {
			for _, tagReq := range q.Tags {
				tagIDs, err := r.ProcessTagRequest(tx, tagReq, userId)
				if err != nil {
					return "", fmt.Errorf("failed to process tag request for question: %w", err)
				}
				// Create mappings for all tags (parent and children)
				for _, tagID := range tagIDs {
					if err := r.CreateQuestionTagMappingWithParents(tx, questionID, tagID, userId); err != nil {
						return "", fmt.Errorf("failed to create question tag mapping: %w", err)
					}
				}
			}
		}
	}

	if len(assessment.Tags) > 0 {
		for _, tagReq := range assessment.Tags {
			tagIDs, err := r.ProcessTagRequest(tx, tagReq, userId)
			if err != nil {
				return "", fmt.Errorf("failed to process tag request for assessment: %w", err)
			}
			for _, tagID := range tagIDs {
				if err := r.CreateAssessmentTagMappingWithParents(tx, assessmentSequence, tagID, userId); err != nil {
					return "", fmt.Errorf("failed to create assessment tag mapping: %w", err)
				}

				// Map to ALL questions in the assessment
				for _, qID := range questionIDs {
					if err := r.CreateQuestionTagMappingWithParents(tx, qID, tagID, userId); err != nil {
						return "", fmt.Errorf("failed to create question tag mapping for assessment tag: %w", err)
					}
				}
			}
		}
	}

	return assessmentSequence, nil
}

func (r *AssessmentRepositoryImpl) CreateAssessment(ctx context.Context, tx *gorm.DB, assessment *models.AssessmentMst) (*models.AssessmentMst, error) {
	if err := tx.WithContext(ctx).Table("assessment_mst").Create(assessment).Error; err != nil {
		return nil, fmt.Errorf("failed to insert assessment_mst: %w", err)
	}
	assessmentID := assessment.AssessmentID
	if assessmentID == 0 {
		return nil, fmt.Errorf("failed to fetch assessment_id (unexpected zero)")
	}
	assessmentSequence := fmt.Sprintf("D%05d", assessmentID)
	if err := tx.WithContext(ctx).
		Table("assessment_mst").
		Where("assessment_id = ?", assessmentID).
		Update("assessment_sequence", assessmentSequence).Error; err != nil {
		return nil, fmt.Errorf("failed to update assessment_sequence: %w", err)
	}
	assessment.AssessmentSequence = assessmentSequence
	return assessment, nil
}

func (r *AssessmentRepositoryImpl) CreateAssessmentExt(ctx context.Context, tx *gorm.DB, ext *models.DhlSurveySurveyExt) (*models.DhlSurveySurveyExt, error) {
	if err := tx.WithContext(ctx).Table("dhl_survey_survey_ext").
		Create(ext).Error; err != nil {

		return nil, fmt.Errorf("failed to insert survey_ext: %w", err)
	}

	return ext, nil
}

func (r *AssessmentRepositoryImpl) CreateAssessmentQuestions(ctx context.Context, tx *gorm.DB, assessmentID int64, assessmentSequence string, questionIDs []int64, createdBy string) error {

	for idx, qID := range questionIDs {

		row := models.AssessmentQuestionMst{
			CreatedOn:          time.Now(),
			CreatedBy:          createdBy,
			IsActive:           true,
			IsDeleted:          false,
			ModifiedOn:         time.Now(),
			ModifiedBy:         createdBy,
			AssessmentSequence: assessmentSequence,
			CorrectPoints:      1,
			DurationInSeconds:  0,
			NegativePoints:     0,
			QuestionID:         qID,
			SequenceID:         int64(idx + 1),
			AssessmentID:       int(assessmentID),
		}

		if err := tx.WithContext(ctx).
			Table("assessment_question_mst").
			Create(&row).Error; err != nil {
			return fmt.Errorf("failed to map question %d: %w", qID, err)
		}
	}

	return nil
}

func (r *AssessmentRepositoryImpl) GetPaginatedQuestionsWithOptions(limit, offset int) ([]*models.QuestionDTO, int64, error) {

	// 1. Count total questions
	var total int64
	if err := r.db.
		Table("question_mst").
		Count(&total).Error; err != nil {
		return nil, 0, err
	}

	type rowData struct {
		QuestionID     int64   `gorm:"column:question_id"`
		Title          string  `gorm:"column:question_title"`
		QuestionTypeID string  `gorm:"column:question_type_id"`
		OptionID       *int64  `gorm:"column:option_id"`
		OptionLabel    *string `gorm:"column:option_label"`
		IsCorrect      *bool   `gorm:"column:is_correct"`
		Score          *int    `gorm:"column:score"`
	}

	sql := `
		WITH qids AS (
			SELECT q.question_id
			FROM question_mst q
			ORDER BY q.question_id
			LIMIT ? OFFSET ?
		)
		SELECT 
			q.question_id,
			qc.value AS question_title,
			q.question_type_id,
			o.option_id,
			oc.value AS option_label,
			o.is_answer AS is_correct,
			COALESCE(o.answer_score, 0) AS score
		FROM qids
		JOIN question_mst q ON q.question_id = qids.question_id
		JOIN content_mst qc ON qc.content_id = q.content_id
		LEFT JOIN option_mst o ON o.question_id = q.question_id
		LEFT JOIN content_mst oc ON oc.content_id = o.content_id
		ORDER BY q.question_id;
	`

	var rows []rowData

	if err := r.db.Raw(sql, limit, offset).Scan(&rows).Error; err != nil {
		return nil, 0, err
	}

	// 3. Grouping logic
	questionsMap := make(map[int64]*models.QuestionDTO)

	for _, row := range rows {

		// if question not yet created
		if _, exists := questionsMap[row.QuestionID]; !exists {
			questionsMap[row.QuestionID] = &models.QuestionDTO{
				QuestionID:   row.QuestionID,
				Title:        row.Title,
				QuestionType: row.QuestionTypeID,
				Options:      []models.OptionDTO{},
			}
		}

		// add option if exists
		if row.OptionID != nil {
			opt := models.OptionDTO{
				OptionID:  *row.OptionID,
				Label:     *row.OptionLabel,
				IsCorrect: *row.IsCorrect,
				Score:     *row.Score,
			}
			questionsMap[row.QuestionID].Options = append(questionsMap[row.QuestionID].Options, opt)
		}
	}

	// Convert map → slice
	questions := make([]*models.QuestionDTO, 0, len(questionsMap))
	for _, q := range questionsMap {
		questions = append(questions, q)
	}

	return questions, total, nil
}

func (r *AssessmentRepositoryImpl) GetAssessmentReport(ctx context.Context, filter models.AssessmentReportFilter) ([]models.AssessmentReportRow, error) {

	var rows []models.AssessmentReportRow

	query := `
select 
    CONCAT(u.first_name, ' ', u.last_name) AS employee_name,
    ue.team_lead,
    ue.manager as team_manager,
    ue.senior_manager,
    ue.sdl,
    ue.sll,
    '' as skill_set,
    ast.created_on as assessment_date,
    am.assessment_desc as assessment_title,
    ast.assessment_status as status,
    ae.attempts_limit as attempts,
    (select count(*) from assessment_status where assessment_id = ast.assessment_id) as total_assigned,

    -- total passed
    (select count(*) from assessment_status 
        where assessment_id = ast.assessment_id and assessment_status = 'Passed') as total_passed,

    -- total failed
    (select count(*) from assessment_status 
        where assessment_id = ast.assessment_id and assessment_status = 'Failed') as total_failed,

    -- total not started
    (select count(*) from assessment_status 
        where assessment_id = ast.assessment_id and assessment_status = 'Not Started') as total_not_started,

    ast.assessment_id as assessment_sequence,

    -- marks obtained
    (select sum(point_assigned) 
        from assessment_result 
        where assessment_session_id = (
            select session_id::text 
            from assessment_user_session 
            where assessment_id = ast.assessment_id 
              and user_id = ast.user_id
            order by created_on desc
            limit 1
        )
    ) as marks_obtained,

    am.marks as total_marks,
    ast.user_id

from assessment_status ast
LEFT JOIN assessment_user_mst u 
    on u.user_id::text = ast.user_id
LEFT JOIN dhl_assessment_user_mst_ext ue
    on ue.user_id = ast.user_id
LEFT JOIN assessment_mst am
    on am.assessment_sequence = ast.assessment_id
LEFT JOIN dhl_survey_survey_ext ae
    on ae.assessment_sequence = ast.assessment_id
WHERE 1=1 
    `
	params := []interface{}{}

	if filter.FromDate != nil && filter.ToDate != nil {
		query += " AND ast.created_on >= ? AND ast.created_on < ?"
		params = append(params, filter.FromDate, filter.ToDate)
	}
	if filter.EmployeeName != "" {
		query += " AND CONCAT(u.first_name, ' ', u.last_name) ILIKE ?"
		params = append(params, "%"+filter.EmployeeName+"%")
	}
	if filter.AssessmentID != "" {
		query += " AND ast.assessment_id = ?"
		params = append(params, filter.AssessmentID)
	}
	if filter.CenterID != nil {
		query += " AND ae.center_id = ?"
		params = append(params, *filter.CenterID)
	}
	if filter.Status != "" {
		query += " AND ast.assessment_status ILIKE ?"
		params = append(params, filter.Status)
	}
	if filter.QuestionnaireTitle != "" {
		query += " AND am.assessment_desc ILIKE ?"
		params = append(params, "%"+filter.QuestionnaireTitle+"%")
	}

	query += " ORDER BY ast.created_on DESC"

	if err := r.db.WithContext(ctx).Raw(query, params...).Scan(&rows).Error; err != nil {
		return nil, err
	}
	return rows, nil
}

func (r *AssessmentRepositoryImpl) CreateSessionImage(tx *gorm.DB, sessionImage *models.AssessmentUserSessionImage) error {
	log.Printf("=== Repository: CreateSessionImage - Session: %s ===", sessionImage.SessionID.String())
	log.Printf("Image size: %d bytes, CreatedBy: %s", len(sessionImage.Image), sessionImage.CreatedBy)

	if err := tx.Create(sessionImage).Error; err != nil {
		log.Printf("ERROR: Failed to insert session image into database: %v", err)
		return err
	}

	log.Printf("SUCCESS: Session image inserted with ID: %d", sessionImage.ID)
	return nil
}

func (r *AssessmentRepositoryImpl) GetSessionImageBySessionID(sessionID string) ([]byte, error) {
	log.Printf("=== Repository: GetSessionImageBySessionID - Session: %s ===", sessionID)

	var sessionImage models.AssessmentUserSessionImage
	err := r.db.Where("session_id = ? AND is_deleted = false AND is_active = true", sessionID).
		Order("created_on DESC").
		First(&sessionImage).Error

	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			log.Printf("No image found for session: %s", sessionID)
			return nil, nil
		}
		log.Printf("ERROR: Failed to fetch session image: %v", err)
		return nil, err
	}

	log.Printf("SUCCESS: Found session image with ID: %d, size: %d bytes", sessionImage.ID, len(sessionImage.Image))
	return sessionImage.Image, nil
}

// Tag Operations
func (r *AssessmentRepositoryImpl) GetOrCreateTagByName(tx *gorm.DB, tagName, createdBy string) (*models.TagMaster, error) {
	return r.GetOrCreateTagWithParent(tx, tagName, nil, createdBy)
}

// GetOrCreateTagWithParent creates or retrieves a tag with optional parent
func (r *AssessmentRepositoryImpl) GetOrCreateTagWithParent(tx *gorm.DB, tagName string, parentTagID *int64, createdBy string) (*models.TagMaster, error) {
	var tag models.TagMaster
	currentTime := time.Now()

	// Build query to find existing tag
	query := tx.Where("tag = ? AND is_deleted = false", tagName)
	if parentTagID != nil {
		query = query.Where("parent_tag_id = ?", *parentTagID)
	} else {
		query = query.Where("parent_tag_id IS NULL")
	}

	err := query.First(&tag).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			// Create new tag
			tag = models.TagMaster{
				TagName:     tagName,
				ParentTagID: parentTagID,
				CreatedOn:   currentTime,
				CreatedBy:   createdBy,
				IsActive:    true,
				IsDeleted:   false,
				ModifiedOn:  currentTime,
				ModifiedBy:  createdBy,
			}
			if err := tx.Create(&tag).Error; err != nil {
				return nil, err
			}
			if parentTagID != nil {
				log.Printf("Created new child tag: %s (parent_id=%d) with ID: %d", tagName, *parentTagID, tag.TagID)
			} else {
				log.Printf("Created new parent tag: %s with ID: %d", tagName, tag.TagID)
			}
			return &tag, nil
		}
		return nil, err
	}

	if parentTagID != nil {
		log.Printf("Found existing child tag: %s (parent_id=%d) with ID: %d", tagName, *parentTagID, tag.TagID)
	} else {
		log.Printf("Found existing parent tag: %s with ID: %d", tagName, tag.TagID)
	}
	return &tag, nil
}

// ProcessTagRequest processes a TagRequest and returns all tag IDs created/found
// Logic:
// 1. If ParentTag is empty and ChildTags is empty -> error
// 2. If ParentTag is provided but not in DB -> create it as parent (parent_tag_id = NULL)
// 3. If ParentTag is provided and exists -> use it
// 4. If ChildTags are provided -> create them with parent_tag_id pointing to ParentTag
// 5. If only ParentTag is provided (no children) -> just return the parent tag ID
func (r *AssessmentRepositoryImpl) ProcessTagRequest(tx *gorm.DB, tagReq models.TagRequest, createdBy string) ([]int64, error) {
	var tagIDs []int64

	// Case 1: If both parent and children are empty, return error
	if tagReq.ParentTag == "" && len(tagReq.ChildTags) == 0 {
		return nil, fmt.Errorf("tag request must have either parent_tag or child_tags")
	}

	// Case 2: Only child tags provided (no parent) - create/get them as standalone parent tags
	if tagReq.ParentTag == "" && len(tagReq.ChildTags) > 0 {
		for _, childName := range tagReq.ChildTags {
			childTag, err := r.GetOrCreateTagWithParent(tx, childName, nil, createdBy)
			if err != nil {
				return nil, fmt.Errorf("failed to create/get tag '%s': %w", childName, err)
			}
			tagIDs = append(tagIDs, childTag.TagID)
		}
		return tagIDs, nil
	}

	// Case 3: Parent tag is provided
	// First, create/get the parent tag (with parent_tag_id = NULL)
	parentTag, err := r.GetOrCreateTagWithParent(tx, tagReq.ParentTag, nil, createdBy)
	if err != nil {
		return nil, fmt.Errorf("failed to create/get parent tag '%s': %w", tagReq.ParentTag, err)
	}

	// If no child tags, just return the parent tag ID
	if len(tagReq.ChildTags) == 0 {
		tagIDs = append(tagIDs, parentTag.TagID)
		return tagIDs, nil
	}

	// Case 4: Parent and children both provided
	// Create/get child tags with parent_tag_id pointing to parent
	for _, childName := range tagReq.ChildTags {
		childTag, err := r.GetOrCreateTagWithParent(tx, childName, &parentTag.TagID, createdBy)
		if err != nil {
			return nil, fmt.Errorf("failed to create/get child tag '%s': %w", childName, err)
		}
		tagIDs = append(tagIDs, childTag.TagID)
	}

	return tagIDs, nil
}

func (r *AssessmentRepositoryImpl) CreateAssessmentTagMapping(tx *gorm.DB, assessmentSequence string, tagID int64, createdBy string) error {
	currentTime := time.Now()

	var existing models.AssessmentTagMapping
	err := tx.Where("assessment_sequence = ? AND tag_id = ? AND is_deleted = false", assessmentSequence, tagID).First(&existing).Error
	if err == nil {
		log.Printf("Assessment tag mapping already exists for assessment: %s, tag: %d", assessmentSequence, tagID)
		return nil
	}
	if !errors.Is(err, gorm.ErrRecordNotFound) {
		return err
	}

	mapping := models.AssessmentTagMapping{
		AssessmentSequence: assessmentSequence,
		TagID:              tagID,
		CreatedOn:          currentTime,
		CreatedBy:          createdBy,
		IsActive:           true,
		IsDeleted:          false,
		ModifiedOn:         currentTime,
		ModifiedBy:         createdBy,
	}

	if err := tx.Create(&mapping).Error; err != nil {
		return err
	}
	log.Printf("Created assessment tag mapping: assessment=%s, tag=%d", assessmentSequence, tagID)
	return nil
}

func (r *AssessmentRepositoryImpl) GetTagsByAssessmentSequence(assessmentSequence string) ([]string, error) {
	var tags []string
	query := `
		SELECT tm.tag
		FROM assessment_tag_mapping atm
		INNER JOIN tag_mst tm ON atm.tag_id = tm.tag_id
		WHERE atm.assessment_sequence = ?
		AND atm.is_deleted = false
		AND atm.is_active = true
		AND tm.is_deleted = false
		AND tm.is_active = true
		ORDER BY tm.tag
	`
	if err := r.db.Raw(query, assessmentSequence).Scan(&tags).Error; err != nil {
		return nil, err
	}
	return tags, nil
}

// GetAllParentTags retrieves all parent tags recursively for a given tag
func (r *AssessmentRepositoryImpl) GetAllParentTags(tx *gorm.DB, tagID int64) ([]int64, error) {
	var parentIDs []int64
	currentTagID := tagID

	// Recursively get parent tags (max 10 levels to prevent infinite loops)
	for i := 0; i < 10; i++ {
		var tag models.TagMaster
		if err := tx.Where("tag_id = ? AND is_deleted = false", currentTagID).First(&tag).Error; err != nil {
			if errors.Is(err, gorm.ErrRecordNotFound) {
				break
			}
			return nil, err
		}

		if tag.ParentTagID == nil {
			break
		}

		parentIDs = append(parentIDs, *tag.ParentTagID)
		currentTagID = *tag.ParentTagID
	}

	return parentIDs, nil
}

// CreateAssessmentTagMappingWithParents creates tag mapping for assessment and all parent tags
func (r *AssessmentRepositoryImpl) CreateAssessmentTagMappingWithParents(tx *gorm.DB, assessmentSequence string, tagID int64, createdBy string) error {
	// Create mapping for the tag itself
	if err := r.CreateAssessmentTagMapping(tx, assessmentSequence, tagID, createdBy); err != nil {
		return err
	}

	// Get all parent tags
	parentTagIDs, err := r.GetAllParentTags(tx, tagID)
	if err != nil {
		return fmt.Errorf("failed to get parent tags: %w", err)
	}

	// Create mappings for all parent tags
	for _, parentTagID := range parentTagIDs {
		if err := r.CreateAssessmentTagMapping(tx, assessmentSequence, parentTagID, createdBy); err != nil {
			return fmt.Errorf("failed to create parent tag mapping: %w", err)
		}
	}

	return nil
}

// Question Tag Operations
func (r *AssessmentRepositoryImpl) CreateQuestionTagMapping(tx *gorm.DB, questionID, tagID int64, createdBy string) error {
	currentTime := time.Now()

	var existing models.QuestionTagMapping
	err := tx.Where("question_id = ? AND tag_id = ? AND is_deleted = false", questionID, tagID).First(&existing).Error
	if err == nil {
		log.Printf("Question tag mapping already exists for question: %d, tag: %d", questionID, tagID)
		return nil
	}
	if !errors.Is(err, gorm.ErrRecordNotFound) {
		return err
	}

	mapping := models.QuestionTagMapping{
		QuestionID: questionID,
		TagID:      tagID,
		CreatedOn:  currentTime,
		CreatedBy:  createdBy,
		IsActive:   true,
		IsDeleted:  false,
		ModifiedOn: currentTime,
		ModifiedBy: createdBy,
	}

	if err := tx.Create(&mapping).Error; err != nil {
		return err
	}
	log.Printf("Created question tag mapping: question=%d, tag=%d", questionID, tagID)
	return nil
}

func (r *AssessmentRepositoryImpl) GetTagsByQuestionID(questionID int64) ([]string, error) {
	var tags []string
	query := `
		SELECT tm.tag
		FROM tag_question_mapping qtm
		INNER JOIN tag_mst tm ON qtm.tag_id = tm.tag_id
		WHERE qtm.question_id = ?
		AND qtm.is_deleted = false
		AND qtm.is_active = true
		AND tm.is_deleted = false
		AND tm.is_active = true
		ORDER BY tm.tag
	`
	if err := r.db.Raw(query, questionID).Scan(&tags).Error; err != nil {
		return nil, err
	}
	return tags, nil
}

func (r *AssessmentRepositoryImpl) GetTagsByQuestionIDs(questionIDs []int64) (map[int64][]string, error) {
	if len(questionIDs) == 0 {
		return make(map[int64][]string), nil
	}

	type QuestionTag struct {
		QuestionID int64  `gorm:"column:question_id"`
		TagName    string `gorm:"column:tag"`
	}

	var results []QuestionTag
	query := `
		SELECT qtm.question_id, tm.tag
		FROM tag_question_mapping qtm
		INNER JOIN tag_mst tm ON qtm.tag_id = tm.tag_id
		WHERE qtm.question_id IN ?
		AND qtm.is_deleted = false
		AND qtm.is_active = true
		AND tm.is_deleted = false
		AND tm.is_active = true
		ORDER BY qtm.question_id, tm.tag
	`
	if err := r.db.Raw(query, questionIDs).Scan(&results).Error; err != nil {
		return nil, err
	}

	tagMap := make(map[int64][]string)
	for _, result := range results {
		tagMap[result.QuestionID] = append(tagMap[result.QuestionID], result.TagName)
	}

	return tagMap, nil
}

func (r *AssessmentRepositoryImpl) GetTagRequestsByAssessmentSequence(assessmentSequence string) ([]models.TagRequest, error) {
	type TagWithParent struct {
		TagID       int64   `gorm:"column:tag_id"`
		TagName     string  `gorm:"column:tag"`
		ParentTagID *int64  `gorm:"column:parent_tag_id"`
		ParentName  *string `gorm:"column:parent_name"`
	}

	var results []TagWithParent
	query := `
		SELECT
			tm.tag_id,
			tm.tag,
			tm.parent_tag_id,
			pt.tag as parent_name
		FROM assessment_tag_mapping atm
		INNER JOIN tag_mst tm ON atm.tag_id = tm.tag_id
		LEFT JOIN tag_mst pt ON tm.parent_tag_id = pt.tag_id
		WHERE atm.assessment_sequence = ?
		AND atm.is_deleted = false
		AND atm.is_active = true
		AND tm.is_deleted = false
		AND tm.is_active = true
		ORDER BY COALESCE(pt.tag, tm.tag), tm.tag
	`
	if err := r.db.Raw(query, assessmentSequence).Scan(&results).Error; err != nil {
		return nil, err
	}

	parentMap := make(map[string][]string)
	standaloneParents := make(map[string]bool)

	for _, result := range results {
		if result.ParentTagID != nil && result.ParentName != nil {
			parentMap[*result.ParentName] = append(parentMap[*result.ParentName], result.TagName)
		} else {
			standaloneParents[result.TagName] = true
		}
	}

	var tagRequests []models.TagRequest
	processed := make(map[string]bool)

	for _, result := range results {
		if result.ParentTagID == nil {
			if !processed[result.TagName] {
				processed[result.TagName] = true
				tagReq := models.TagRequest{
					ParentTag: result.TagName,
					ChildTags: parentMap[result.TagName],
				}
				tagRequests = append(tagRequests, tagReq)
			}
		}
	}

	return tagRequests, nil
}

func (r *AssessmentRepositoryImpl) GetTagRequestsByQuestionIDs(questionIDs []int64) (map[int64][]models.TagRequest, error) {
	if len(questionIDs) == 0 {
		return make(map[int64][]models.TagRequest), nil
	}

	type QuestionTagWithParent struct {
		QuestionID  int64   `gorm:"column:question_id"`
		TagID       int64   `gorm:"column:tag_id"`
		TagName     string  `gorm:"column:tag"`
		ParentTagID *int64  `gorm:"column:parent_tag_id"`
		ParentName  *string `gorm:"column:parent_name"`
	}

	var results []QuestionTagWithParent
	query := `
		SELECT
			qtm.question_id,
			tm.tag_id,
			tm.tag,
			tm.parent_tag_id,
			pt.tag as parent_name
		FROM tag_question_mapping qtm
		INNER JOIN tag_mst tm ON qtm.tag_id = tm.tag_id
		LEFT JOIN tag_mst pt ON tm.parent_tag_id = pt.tag_id
		WHERE qtm.question_id IN ?
		AND qtm.is_deleted = false
		AND qtm.is_active = true
		AND tm.is_deleted = false
		AND tm.is_active = true
		ORDER BY qtm.question_id, COALESCE(pt.tag, tm.tag), tm.tag
	`
	if err := r.db.Raw(query, questionIDs).Scan(&results).Error; err != nil {
		return nil, err
	}

	questionTagMap := make(map[int64][]models.TagRequest)

	for _, qID := range questionIDs {
		parentMap := make(map[string][]string)
		standaloneParents := make(map[string]bool)

		for _, result := range results {
			if result.QuestionID == qID {
				if result.ParentTagID != nil && result.ParentName != nil {
					parentMap[*result.ParentName] = append(parentMap[*result.ParentName], result.TagName)
				} else {
					standaloneParents[result.TagName] = true
				}
			}
		}

		processed := make(map[string]bool)
		var tagRequests []models.TagRequest

		for _, result := range results {
			if result.QuestionID == qID && result.ParentTagID == nil {
				if !processed[result.TagName] {
					processed[result.TagName] = true
					tagReq := models.TagRequest{
						ParentTag: result.TagName,
						ChildTags: parentMap[result.TagName],
					}
					tagRequests = append(tagRequests, tagReq)
				}
			}
		}

		questionTagMap[qID] = tagRequests
	}

	return questionTagMap, nil
}

func (r *AssessmentRepositoryImpl) CreateQuestionTagMappingWithParents(tx *gorm.DB, questionID, tagID int64, createdBy string) error {
	if err := r.CreateQuestionTagMapping(tx, questionID, tagID, createdBy); err != nil {
		return err
	}

	parentTagIDs, err := r.GetAllParentTags(tx, tagID)
	if err != nil {
		return fmt.Errorf("failed to get parent tags: %w", err)
	}

	for _, parentTagID := range parentTagIDs {
		if err := r.CreateQuestionTagMapping(tx, questionID, parentTagID, createdBy); err != nil {
			return fmt.Errorf("failed to create parent tag mapping: %w", err)
		}
	}

	return nil
}

