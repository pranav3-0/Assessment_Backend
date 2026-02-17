package services

import (
	"dhl/models"
	"dhl/repository"
)

type ContactService interface {
	Submit(response *models.ContactUsResponse) error
	List() ([]models.ContactUsResponse, error)
}

type ContactServiceImpl struct {
	repo repository.ContactRepository
}

func NewContactService(repo repository.ContactRepository) ContactService {
	return &ContactServiceImpl{repo}
}

func (s *ContactServiceImpl) Submit(response *models.ContactUsResponse) error {
	return s.repo.Save(response)
}

func (s *ContactServiceImpl) List() ([]models.ContactUsResponse, error) {
	return s.repo.GetAll()
}
