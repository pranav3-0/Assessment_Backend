package constant

const (
	Version = "/api/version"
)

const (
	Success  = "success"
	Failure  = "failure"
	Active   = "Active"
	InActive = "Inactive"
	CREATE   = "Create"
	UPDATE   = "Update"
	DELETE   = "Delete"
	Running  = "running"
)

const (
	Register                    = "/user-register"
	Login                       = "/user-login"
	Logout                      = "/logout"
	UpdateUser                  = "/update-user"
	Users                       = "/users"
	UserProfile                 = "/user-profile"
	MigrateKeycloak             = "/migrate-to-keycloak"
	MigrateRolesKeycloak        = "/migrate-roles"
	MigrateNotify               = "/migrate-notify"
	ImportAssessment            = "/import-assessment"
	GenerateAssessment          = "/generate-assessment"
	SaveGeneratedAssessment     = "/save-generated-assessment"
	Assessment                  = "/assessment"
	AssessmentSubmit            = "/assessment-submit"
	TypeAssessmentSubmit        = "/typing-assessment-submit"
	Assessments                 = "/assessments"
	ManagerAssessments          = "/manager-assessments"
	AssessmentStatus            = "/assessment-status"
	AssessmentDuplicate         = "/assessment-duplicate"
	AssessmentReport            = "/assessment-report"
	Questions                   = "/questions"
	Question                    = "/question"
	Contact                     = "/contact"
	ContactResponse             = "/contact-response"
	DistributeAssessmentUser    = "/distribute-assessment-user"
	DistributeAssessmentManager = "/distribute-assessment-manger"
	MapUsersToManager           = "/map-users-to-manger"
	DHLBusinessPartner          = "/business-partner"
	Certificate                 = "/certificate"
	RegisterBulk                = "/user-register-bulk"
	SessionImage                = "/session-image"
	JobDescription              = "/job-description"
	JobDescriptions             = "/job-descriptions"
)

type UserRole string

const (
	Admin   UserRole = "admin"
	Manager UserRole = "manager"
	User    UserRole = "user"
)

var ValidUserRoles = []UserRole{
	Admin,
	User,
	Manager,
}

type AssessmentState string

const (
	Draft  AssessmentState = "draft"
	Open   AssessmentState = "open"
	Closed AssessmentState = "closed"
)
