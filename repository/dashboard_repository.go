package repository

import (
	"dhl/models"
	"gorm.io/gorm"
)

type DashboardRepository interface {
	GetDashboardCounts() (totalJDs, totalAssessments, totalUsers, pendingAssignments int64, err error)
	GetRecentActivities() ([]models.RecentActivity, error)
}

type dashboardRepository struct {
	db *gorm.DB
}

func NewDashboardRepository(db *gorm.DB) DashboardRepository {
	return &dashboardRepository{db: db}
}



func (r *dashboardRepository) GetDashboardCounts() (int64, int64, int64, int64, error) {
	var totalJDs int64
	var totalAssessments int64
	var totalUsers int64
	var pendingAssignments int64

	// Job Descriptions
	if err := r.db.
		Table("job_descriptions").
		Where("is_deleted = false").
		Count(&totalJDs).Error; err != nil {
		return 0, 0, 0, 0, err
	}

	// Assessments
	if err := r.db.
		Table("assessment_mst").
		Where("is_deleted = false").
		Count(&totalAssessments).Error; err != nil {
		return 0, 0, 0, 0, err
	}

	// Users
	if err := r.db.
		Table("assessment_user_mst").
		Where("is_active = true").
		Count(&totalUsers).Error; err != nil {
		return 0, 0, 0, 0, err
	}

	// Pending Assignments
	if err := r.db.
		Table("assessment_status").
		Where("assessment_status != ?", "COMPLETED").
		Count(&pendingAssignments).Error; err != nil {
		return 0, 0, 0, 0, err
	}

	return totalJDs, totalAssessments, totalUsers, pendingAssignments, nil
}



func (r *dashboardRepository) GetRecentActivities() ([]models.RecentActivity, error) {
	var activities []models.RecentActivity

	err := r.db.
		Table("activity_log al").
		Select(`
			am.activity_name as activity,
			al.log_text as details,
			al.activity_log_time as date
		`).
		Joins("JOIN activity_mst am ON al.activity_id = am.activity_id").
		Where("al.is_deleted = false").
		Where("am.activity_name NOT IN ?", []string{
			"Assessment Started",
			"Assessment Completed",
		}).
		Order("al.created_on DESC").
		Limit(10).
		Scan(&activities).Error

	return activities, err
}