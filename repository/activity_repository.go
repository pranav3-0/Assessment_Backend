package repository

import "gorm.io/gorm"

type ActivityRepository interface {
	LogActivity(activityID int64, logText string) error
}

type activityRepository struct {
	db *gorm.DB
}

func NewActivityRepository(db *gorm.DB) ActivityRepository {
	return &activityRepository{db: db}
}

func (r *activityRepository) LogActivity(activityID int64, logText string) error {
	return r.db.Exec(`
		INSERT INTO activity_log
		(activity_id, log_text, activity_log_time, is_active, is_deleted, created_on)
		VALUES (?, ?, CURRENT_DATE, true, false, NOW())
	`, activityID, logText).Error
}