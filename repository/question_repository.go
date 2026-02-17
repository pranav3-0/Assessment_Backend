package repository

import (
	"dhl/models"

	"gorm.io/gorm"
)

type QuestionRepository interface {
	CreateQuestion(tx *gorm.DB, question *models.QuestionMst) error
	CreateContent(tx *gorm.DB, content *models.ContentMst) error
	CreateOption(tx *gorm.DB, option *models.OptionMst) error
}

type QuestionRepositoryImpl struct {
	db *gorm.DB
}

func NewQuestionRepository(db *gorm.DB) QuestionRepository {
	return &QuestionRepositoryImpl{db: db}
}

func (r *QuestionRepositoryImpl) CreateContent(tx *gorm.DB, content *models.ContentMst) error {
	return tx.Create(content).Error
}

func (r *QuestionRepositoryImpl) CreateQuestion(tx *gorm.DB, question *models.QuestionMst) error {
	return tx.Create(question).Error
}

func (r *QuestionRepositoryImpl) CreateOption(tx *gorm.DB, option *models.OptionMst) error {
	return tx.Create(option).Error
}
