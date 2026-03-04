package services

import (
	"dhl/models"
	"dhl/repository"
)

type DashboardService interface {
	GetDashboard() (*models.DashboardResponse, error)
}

type dashboardService struct {
	dashboardRepo repository.DashboardRepository
}

func NewDashboardService(repo repository.DashboardRepository) DashboardService {
	return &dashboardService{
		dashboardRepo: repo,
	}
}

func (s *dashboardService) GetDashboard() (*models.DashboardResponse, error) {

	totalJDs, totalAssessments, totalUsers, pendingAssignments, err :=
		s.dashboardRepo.GetDashboardCounts()
	if err != nil {
		return nil, err
	}

	recentActivities, err := s.dashboardRepo.GetRecentActivities()
	if err != nil {
		return nil, err
	}

	return &models.DashboardResponse{
		TotalJDs:           totalJDs,
		TotalAssessments:   totalAssessments,
		TotalUsers:         totalUsers,
		PendingAssignments: pendingAssignments,
		RecentActivities:   recentActivities,
	}, nil
}