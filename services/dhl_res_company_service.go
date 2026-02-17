package services

import (
	"context"
	"dhl/models"
	"dhl/repository"
)

type DHLResCompanyService interface {
	CreateCompany(ctx context.Context, req models.DHLResCompany) error
	ListCompanies(ctx context.Context) ([]models.DHLResCompany, error)
	UpdateCompany(ctx context.Context, req models.DHLResCompany) error
	DeleteCompany(ctx context.Context, id int64) error
}

type DHLResCompanyServiceImpl struct {
	repo repository.DHLResCompanyRepository
}

func NewDHLResCompanyService(repo repository.DHLResCompanyRepository) DHLResCompanyService {
	return &DHLResCompanyServiceImpl{repo}
}

func (s *DHLResCompanyServiceImpl) CreateCompany(ctx context.Context, req models.DHLResCompany) error {
	return s.repo.Create(ctx, &req)
}

func (s *DHLResCompanyServiceImpl) ListCompanies(ctx context.Context) ([]models.DHLResCompany, error) {
	return s.repo.List(ctx)
}

func (s *DHLResCompanyServiceImpl) UpdateCompany(ctx context.Context, data models.DHLResCompany) error {
	updates := &models.DHLResCompany{
		CompanyID:              data.CompanyID,
		Name:                   data.Name,
		PartnerID:              data.PartnerID,
		CurrencyID:             data.CurrencyID,
		Sequence:               data.Sequence,
		CreatedAt:              data.CreatedAt,
		ParentID:               data.ParentID,
		ReportHeader:           data.ReportHeader,
		ReportFooter:           data.ReportFooter,
		LogoWeb:                data.LogoWeb,
		AccountNo:              data.AccountNo,
		Email:                  data.Email,
		Phone:                  data.Phone,
		CompanyRegistry:        data.CompanyRegistry,
		PaperFormatID:          data.PaperFormatID,
		ExternalReportLayoutID: data.ExternalReportLayoutID,
		BaseOnboardingState:    data.BaseOnboardingState,
		Font:                   data.Font,
		PrimaryColor:           data.PrimaryColor,
		SecondaryColor:         data.SecondaryColor,
		CreatedBy:              data.CreatedBy,
		UpdatedBy:              data.UpdatedBy,
		UpdatedAt:              data.UpdatedAt,
		SocialTwitter:          data.SocialTwitter,
		SocialFacebook:         data.SocialFacebook,
		SocialGithub:           data.SocialGithub,
		SocialLinkedIn:         data.SocialLinkedIn,
		SocialYoutube:          data.SocialYoutube,
		SocialInstagram:        data.SocialInstagram,
		PartnerGID:             data.PartnerGID,
		SnailmailColor:         data.SnailmailColor,
		SnailmailCover:         data.SnailmailCover,
		SnailmailDuplex:        data.SnailmailDuplex,
	}
	return s.repo.Update(ctx, updates)
}

func (s *DHLResCompanyServiceImpl) DeleteCompany(ctx context.Context, id int64) error {
	return s.repo.Delete(ctx, id)
}
