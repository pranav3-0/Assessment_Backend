package models

import "time"

type RecentActivity struct {
	Activity string    `json:"activity"`
	Details  string    `json:"details"`
	Date     time.Time `json:"date"`
}

type DashboardResponse struct {
	TotalJDs            int64            `json:"total_jds"`
	TotalAssessments    int64            `json:"total_assessments"`
	TotalUsers          int64            `json:"total_users"`
	PendingAssignments  int64            `json:"pending_assignments"`
	RecentActivities    []RecentActivity `json:"recent_activity"`
}