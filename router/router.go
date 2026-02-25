package router

import (
	"dhl/constant"
	"dhl/controller"
	"dhl/repository"
	"dhl/services"
	"net/http"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

func InitializeRoutes(apiGroup *gin.RouterGroup, db *gorm.DB) {
	var userRepo = repository.NewUserRepository(db)
	var clientRepo = repository.NewClientRepository(db)
	var assessmentRepo = repository.NewAssessmentRepository(db)
	var contactRepo = repository.NewContactRepository(db)
	var dhlBusinessPartnerRepository = repository.NewDHLBusinessPartnerRepository(db)
	var dhlCenterRepository = repository.NewDHLCenterRepository(db)
	var dhlResCompanyRepository = repository.NewDHLResCompanyRepository(db)
	var dhlResPartnerIndustryRepository = repository.NewDHLResPartnerIndustryRepository(db)
	var dhlServiceRepository = repository.NewDHLServiceRepository(db)
	var dhlServiceGroupRepository = repository.NewDHLServiceGroupRepository(db)
	var dhlServiceLineRepository = repository.NewDHLServiceLineRepository(db)
	var dhlSubBusinessPartnerRepository = repository.NewDHLSubBusinessPartnerRepository(db)
	var dhlSubServiceRepository = repository.NewDHLSubServiceRepository(db)

	var jobRepo = repository.NewJobDescriptionRepository(db)
    var jobService = services.NewJobDescriptionService(jobRepo, db)
	var questionRepo = repository.NewQuestionRepository(db)
    var questionService = services.NewQuestionService(questionRepo, db, assessmentRepo)


	var userService = services.NewUserService(userRepo, clientRepo, db)
	var assessmentService = services.NewAssessmentService(assessmentRepo, db)
	var geminiService = services.NewGeminiService()
	var contactService = services.NewContactService(contactRepo)
	var dhlBusinessPartnerService = services.NewDHLBusinessPartnerService(dhlBusinessPartnerRepository)
	var dhlCenterService = services.NewDHLCenterService(dhlCenterRepository)
	var dhlResCompanyService = services.NewDHLResCompanyService(dhlResCompanyRepository)

	var dhlResPartnerIndustryService = services.NewDHLResPartnerIndustryService(dhlResPartnerIndustryRepository)
	var dhlServiceService = services.NewDHLServiceService(dhlServiceRepository)
	var dhlServiceGroupService = services.NewDHLServiceGroupService(dhlServiceGroupRepository)
	var dhlServiceLineService = services.NewDHLServiceLineService(dhlServiceLineRepository)
	var dhlSubBusinessPartnerService = services.NewDHLSubBusinessPartnerService(dhlSubBusinessPartnerRepository)
	var dhlSubServiceService = services.NewDHLSubServiceService(dhlSubServiceRepository)
	var notificationService = services.NewNotificationService(userRepo, assessmentRepo)
	var authService = services.NewAuthService(userRepo, clientRepo, notificationService, db)

	var userController = controller.NewUserController(userService, authService)
	var adminController = controller.NewAdminController(userService, authService, assessmentService, notificationService, contactService,jobService, questionService)
	var assessmentController = controller.NewAssessmentController(assessmentService, userService, geminiService)
	var publicController = controller.NewPublicController(contactService)
	var mastersController = controller.NewMastersController(dhlBusinessPartnerService, dhlCenterService, dhlResCompanyService,
		dhlResPartnerIndustryService, dhlServiceService, dhlServiceGroupService, dhlServiceLineService, dhlSubBusinessPartnerService, dhlSubServiceService)

	UserRoutes(apiGroup, userController)
	AdminRoutes(apiGroup, adminController, assessmentController, mastersController)
	QuestionAuthorRoutes(apiGroup, adminController, assessmentController, mastersController)
	AssessmentRoutes(apiGroup, assessmentController)
	OpenRoutes(apiGroup, publicController, adminController)
}

func getAdminRoutes(adminController *controller.AdminController, assessmentController *controller.AssessmentController, mastersController *controller.MastersController,) Routes {
	return Routes{
		Route{"Admin", http.MethodPost, constant.Assessments, adminController.GetAssessments},
		Route{"Admin", http.MethodPost, constant.ManagerAssessments, adminController.GetManagerAssessments},
		Route{"Admin", http.MethodGet, constant.Assessment, assessmentController.GetAssessment},
		Route{"Admin", http.MethodPost, constant.Users, adminController.GetUsers},
		Route{"Admin", http.MethodPost, constant.UpdateUser, adminController.UpdateUserProfile},
		Route{"Admin", http.MethodPost, constant.ImportAssessment, assessmentController.UploadAssessment},
		Route{"Admin", http.MethodPost, constant.GenerateAssessment, assessmentController.GenerateAssessmentWithAI},
		Route{"Admin", http.MethodPost, constant.SaveGeneratedAssessment, assessmentController.SaveGeneratedAssessment},
		Route{"Admin", http.MethodPost, constant.Assessment, assessmentController.CreateAssessment},
		Route{"Admin", http.MethodPost, constant.AssessmentDuplicate, assessmentController.DuplicateAssessment},
		Route{"Admin", http.MethodPut, constant.AssessmentStatus, adminController.UpdateAssessmentStatusController},
		Route{"Admin", http.MethodPut, constant.Assessment, adminController.UpdateAssessment},
		Route{"Admin", http.MethodPost, constant.Questions, adminController.GetQuestionsController},
		Route{"Admin", http.MethodPost, constant.DistributeAssessmentUser, adminController.DistributeAssessmentToUserController},
		Route{"Admin", http.MethodPost, constant.DistributeAssessmentManager, adminController.DistributeAssessmentToManagerController},
		Route{"Admin", http.MethodPost, constant.MapUsersToManager, adminController.MapUsersToManagerController},
		Route{"Admin", http.MethodPost, constant.AssessmentReport, assessmentController.DownloadAssessmentReport},

		Route{"Admin", http.MethodPost, constant.JobDescription, adminController.CreateJobDescription},
        Route{"Admin", http.MethodPost, constant.JobDescriptions, adminController.GetJobDescriptions},
		Route{"Admin", http.MethodPost, constant.Question, adminController.CreateMultipleQuestions},
		Route{"Admin", http.MethodPost, constant.AssessmentUserResult, adminController.GetAssessmentUserResult},
		Route{"Admin", http.MethodPost, constant.CheckAssessmentAssignment, adminController.CheckAssessmentAssignment},



		Route{"Contact Form", http.MethodPost, constant.ContactResponse, adminController.ListContactRespController},
		// Master Routes
		Route{"CreateDHLBusinessPartner", http.MethodPost, "/business-partner", mastersController.CreateDHLBusinessPartner},
		Route{"ListDHLBusinessPartners", http.MethodGet, "/business-partner", mastersController.ListDHLBusinessPartners},
		Route{"UpdateDHLBusinessPartner", http.MethodPut, "/business-partner", mastersController.UpdateDHLBusinessPartner},
		Route{"DeleteDHLBusinessPartner", http.MethodDelete, "/business-partner/:id", mastersController.DeleteDHLBusinessPartner},

		Route{"CreateDHLCenter", http.MethodPost, "/dhl-center", mastersController.CreateDHLCenter},
		Route{"ListDHLCenters", http.MethodGet, "/dhl-center", mastersController.ListDHLCenters},
		Route{"UpdateDHLCenter", http.MethodPut, "/dhl-center", mastersController.UpdateDHLCenter},
		Route{"DeleteDHLCenter", http.MethodDelete, "/dhl-center/:id", mastersController.DeleteDHLCenter},

		Route{"CreateDHLResCompany", http.MethodPost, "/res-company", mastersController.CreateDHLResCompany},
		Route{"ListDHLResCompanys", http.MethodGet, "/res-company", mastersController.ListDHLResCompanys},
		Route{"UpdateDHLResCompany", http.MethodPut, "/res-company", mastersController.UpdateDHLResCompany},
		Route{"DeleteDHLResCompany", http.MethodDelete, "/res-company/:id", mastersController.DeleteDHLResCompany},

		Route{"CreateDHLResPartnerIndustry", http.MethodPost, "/partner-industry", mastersController.CreateDHLResPartnerIndustry},
		Route{"ListDHLResPartnerIndustry", http.MethodGet, "/partner-industry", mastersController.ListDHLResPartnerIndustry},
		Route{"UpdateDHLResPartnerIndustry", http.MethodPut, "/partner-industry", mastersController.UpdateDHLResPartnerIndustry},
		Route{"DeleteDHLResPartnerIndustry", http.MethodDelete, "/partner-industry/:id", mastersController.DeleteDHLResPartnerIndustry},

		Route{"CreateDHLService", http.MethodPost, "/dhl-service", mastersController.CreateDHLService},
		Route{"ListDHLServices", http.MethodGet, "/dhl-service", mastersController.ListDHLServices},
		Route{"UpdateDHLService", http.MethodPut, "/dhl-service", mastersController.UpdateDHLService},
		Route{"DeleteDHLService", http.MethodDelete, "/dhl-service/:id", mastersController.DeleteDHLService},

		Route{"CreateDHLServiceGroup", http.MethodPost, "/dhl-service-group", mastersController.CreateDHLServiceGroup},
		Route{"ListDHLServiceGroups", http.MethodGet, "/dhl-service-group", mastersController.ListDHLServiceGroups},
		Route{"UpdateDHLServiceGroup", http.MethodPut, "/dhl-service-group", mastersController.UpdateDHLServiceGroup},
		Route{"DeleteDHLServiceGroup", http.MethodDelete, "/dhl-service-group/:id", mastersController.DeleteDHLServiceGroup},

		Route{"CreateDHLServiceLine", http.MethodPost, "/dhl-service-line", mastersController.CreateDHLServiceLine},
		Route{"ListDHLServiceLine", http.MethodGet, "/dhl-service-line", mastersController.ListDHLServiceLine},
		Route{"UpdateDHLServiceLine", http.MethodPut, "/dhl-service-line", mastersController.UpdateDHLServiceLine},
		Route{"DeleteDHLServiceLine", http.MethodDelete, "/dhl-service-line/:id", mastersController.DeleteDHLServiceLine},

		Route{"CreateDHLSubBusinessPartner", http.MethodPost, "/dhl-sub-business-partner", mastersController.CreateDHLSubBusinessPartner},
		Route{"ListDHLSubBusinessPartner", http.MethodGet, "/dhl-sub-business-partner", mastersController.ListDHLSubBusinessPartner},
		Route{"UpdateDHLSubBusinessPartner", http.MethodPut, "/dhl-sub-business-partner", mastersController.UpdateDHLSubBusinessPartner},
		Route{"DeleteDHLSubBusinessPartner", http.MethodDelete, "/dhl-sub-business-partner/:id", mastersController.DeleteDHLSubBusinessPartner},

		Route{"CreateDHLSubService", http.MethodPost, "/dhl-sub-service", mastersController.CreateDHLSubService},
		Route{"ListDHLSubService", http.MethodGet, "/dhl-sub-service", mastersController.ListDHLSubService},
		Route{"UpdateDHLSubService", http.MethodPut, "/dhl-sub-service", mastersController.UpdateDHLSubService},
		Route{"DeleteDHLSubService", http.MethodDelete, "/dhl-sub-service/:id", mastersController.DeleteDHLSubService},
	}
}

func getQuestionAuthorRoutes(adminController *controller.AdminController, assessmentController *controller.AssessmentController, mastersController *controller.MastersController) Routes {
	return Routes{
		Route{"Question Author", http.MethodPost, constant.Assessments, adminController.GetAssessments},
		Route{"Question Author", http.MethodPost, constant.Users, adminController.GetUsers},
		Route{"Question Author", http.MethodPost, constant.ImportAssessment, assessmentController.UploadAssessment},
		Route{"Question Author", http.MethodPost, constant.GenerateAssessment, assessmentController.GenerateAssessmentWithAI},
		Route{"Question Author", http.MethodPost, constant.SaveGeneratedAssessment, assessmentController.SaveGeneratedAssessment},
		Route{"Question Author", http.MethodPost, constant.Assessment, assessmentController.CreateAssessment},
		Route{"Question Author", http.MethodPost, constant.AssessmentDuplicate, assessmentController.DuplicateAssessment},
		Route{"Question Author", http.MethodPut, constant.Assessment, adminController.UpdateAssessment},
		Route{"Question Author", http.MethodPost, constant.Questions, adminController.GetQuestionsController},
		Route{"Question Author", http.MethodPost, constant.DistributeAssessmentManager, adminController.DistributeAssessmentToManagerController},
	}
}

func getUserRoutes(userController *controller.UserController) Routes {
	return Routes{
		Route{"User", http.MethodPost, constant.Register, userController.RegisterUser},
		Route{"User", http.MethodPost, constant.RegisterBulk, userController.RegisterUsersFromCSV},
		Route{"User", http.MethodPost, constant.Login, userController.LoginUser},
		Route{"User", http.MethodPost, constant.Logout, userController.LogoutUser},
		Route{"User", http.MethodPost, constant.UserProfile, userController.GetProfile},
	}
}

func getAssessmentRoutes(assessmentController *controller.AssessmentController) Routes {
	return Routes{
		Route{"Assessment", http.MethodPost, constant.Assessment, assessmentController.GetUserAssessment},
		Route{"Assessment", http.MethodPost, constant.AssessmentSubmit, assessmentController.SubmitUserAssessment},
		Route{"Assessment", http.MethodPost, constant.TypeAssessmentSubmit, assessmentController.SubmitUserAssessmentTypingResult},
		Route{"Assessment", http.MethodPost, constant.Assessments, assessmentController.GetUserAssessments},
		Route{"Assessment", http.MethodPost, constant.Certificate, assessmentController.GetUserAssessmentCerficiate},
		Route{"Assessment", http.MethodPost, constant.SessionImage, assessmentController.CreateSessionImage},
		

		Route{"Assessment", http.MethodPost, constant.AssessmentVerificationPhoto, assessmentController.UploadPhoto},
		Route{"Assessment", http.MethodPost, constant.AssessmentVerificationVoice, assessmentController.UploadVoice},
		Route{"Assessment", http.MethodPost, constant.AssessmentStart, assessmentController.StartAssessment},
		Route{"Assessment", http.MethodPost, constant.AssessmentResultView, assessmentController.GetUserAssessmentResult},

	}
}

func getOpenRoutes(publicController *controller.PublicController, adminController *controller.AdminController) Routes {
	return Routes{
		Route{"Contact Form", http.MethodPost, constant.Contact, publicController.SubmitContactFormController},
		Route{"Admin", http.MethodPost, constant.MigrateKeycloak, adminController.MigrateUsersToKeycloak},
		Route{"Admin", http.MethodPost, constant.MigrateRolesKeycloak, adminController.MigrateRoleToKeycloak},
		Route{"Admin", http.MethodPost, constant.MigrateNotify, adminController.CreateUsersOnNotify},
	}
}
