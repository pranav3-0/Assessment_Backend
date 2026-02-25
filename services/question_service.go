package services

import (
	"dhl/models"
	"dhl/repository"
	"time"

	"gorm.io/gorm"
)

type QuestionService interface {
	CreateQuestion(req models.CreateQuestionRequest, createdBy string) (int64, error)

	CreateMultipleQuestions(
		requests []models.CreateQuestionRequest,
		createdBy string,
	) ([]int64, error)
}

type QuestionServiceImpl struct {
	repo           repository.QuestionRepository
	db             *gorm.DB
	assessmentRepo repository.AssessmentRepository
}

func NewQuestionService(repo repository.QuestionRepository, db *gorm.DB, assessmentRepo repository.AssessmentRepository) QuestionService {
	return &QuestionServiceImpl{
		repo:           repo,
		db:             db,
		assessmentRepo: assessmentRepo,
	}
}

func (s *QuestionServiceImpl) CreateQuestion(
	req models.CreateQuestionRequest,
	createdBy string,
) (int64, error) {

	tx := s.db.Begin()
	now := time.Now()

	content := models.ContentMst{
		ContentTypeID: 1,
		Value:         req.Title,
	}

	if err := s.repo.CreateContent(tx, &content); err != nil {
		tx.Rollback()
		return 0, err
	}

	questionTypeID, err := s.repo.GetQuestionTypeID(tx, req.QuestionType)
	if err != nil {
		tx.Rollback()
		return 0, err
	}

	question := models.QuestionMst{
		ContentID:      content.ContentID,
		QuestionTypeID: questionTypeID,
		IsActive:       true,
		IsDeleted:      false,
		CreatedOn:      now,
		ModifiedOn:     now,
	}

	if err := s.repo.CreateQuestion(tx, &question); err != nil {
		tx.Rollback()
		return 0, err
	}

	for _, opt := range req.Options {

		optContent := models.ContentMst{
			ContentTypeID: 1,
			Value:         opt.Label,
		}

		if err := s.repo.CreateContent(tx, &optContent); err != nil {
			tx.Rollback()
			return 0, err
		}

		option := models.OptionMst{
			ContentID:   optContent.ContentID,
			IsAnswer:    opt.IsCorrect,
			QuestionID:  question.QuestionID,
			AnswerScore: int(opt.Score),
		}

		if err := s.repo.CreateOption(tx, &option); err != nil {
			tx.Rollback()
			return 0, err
		}
	}

	// Tags
	for _, tagReq := range req.Tags {
		tagIDs, err := s.assessmentRepo.ProcessTagRequest(tx, tagReq, createdBy)
		if err != nil {
			tx.Rollback()
			return 0, err
		}

		for _, tagID := range tagIDs {
			if err := s.assessmentRepo.CreateQuestionTagMappingWithParents(
				tx,
				question.QuestionID,
				tagID,
				createdBy,
			); err != nil {
				tx.Rollback()
				return 0, err
			}
		}
	}

	if err := tx.Commit().Error; err != nil {
		return 0, err
	}

	return question.QuestionID, nil
}

func (s *QuestionServiceImpl) CreateMultipleQuestions(
	requests []models.CreateQuestionRequest,
	createdBy string,
) ([]int64, error) {

	tx := s.db.Begin()
	now := time.Now()

	var questionIDs []int64

	for _, req := range requests {

		content := models.ContentMst{
			ContentTypeID: 1,
			Value:         req.Title,
		}

		if err := s.repo.CreateContent(tx, &content); err != nil {
			tx.Rollback()
			return nil, err
		}

		questionTypeID, err := s.repo.GetQuestionTypeID(tx, req.QuestionType)
		if err != nil {
			tx.Rollback()
			return nil, err
		}

		question := models.QuestionMst{
			ContentID:      content.ContentID,
			QuestionTypeID: questionTypeID,
			IsActive:       true,
			IsDeleted:      false,
			CreatedOn:      now,
			ModifiedOn:     now,
		}

		if err := s.repo.CreateQuestion(tx, &question); err != nil {
			tx.Rollback()
			return nil, err
		}

		for _, opt := range req.Options {

			optContent := models.ContentMst{
				ContentTypeID: 1,
				Value:         opt.Label,
			}

			if err := s.repo.CreateContent(tx, &optContent); err != nil {
				tx.Rollback()
				return nil, err
			}

			option := models.OptionMst{
				ContentID:   optContent.ContentID,
				IsAnswer:    opt.IsCorrect,
				QuestionID:  question.QuestionID,
				AnswerScore: int(opt.Score),
			}

			if err := s.repo.CreateOption(tx, &option); err != nil {
				tx.Rollback()
				return nil, err
			}
		}

		for _, tagReq := range req.Tags {
			tagIDs, err := s.assessmentRepo.ProcessTagRequest(tx, tagReq, createdBy)
			if err != nil {
				tx.Rollback()
				return nil, err
			}

			for _, tagID := range tagIDs {
				if err := s.assessmentRepo.CreateQuestionTagMappingWithParents(
					tx,
					question.QuestionID,
					tagID,
					createdBy,
				); err != nil {
					tx.Rollback()
					return nil, err
				}
			}
		}

		questionIDs = append(questionIDs, question.QuestionID)
	}

	if err := tx.Commit().Error; err != nil {
		return nil, err
	}

	return questionIDs, nil
}