package repository

import (
	"dhl/models"
     "fmt"
	"gorm.io/gorm"
)

type QuestionRepository interface {
	CreateQuestion(tx *gorm.DB, question *models.QuestionMst) error
	CreateContent(tx *gorm.DB, content *models.ContentMst) error
	CreateOption(tx *gorm.DB, option *models.OptionMst) error

	GetQuestionTypeID(tx *gorm.DB, typeValue string) (int64, error) 
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


func (r *QuestionRepositoryImpl) GetQuestionTypeID(tx *gorm.DB, typeValue string) (int64, error){
	var id int64

	err := tx.
		Table("question_type_config").
		Select("question_type_id").
		Where("question_type = ?", typeValue).
		Scan(&id).Error

	if err != nil {
		return 0, err
	}

	if id == 0 {
		return 0, fmt.Errorf("invalid question_type: %s", typeValue)
	}

	return id, nil
}