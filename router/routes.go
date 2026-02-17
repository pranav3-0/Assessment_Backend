package router

import (
	"dhl/auth"
	"dhl/constant"
	"dhl/controller"
	"dhl/database"
	"dhl/middleware"
	"dhl/utils"
	"log"
	"net/http"
	"os"
	"strings"

	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
)

type Route struct {
	Name       string
	Method     string
	Path       string
	HandleFunc func(*gin.Context)
}

type routes struct {
	router *gin.Engine
}

type Routes []Route

var ProtectedRoutes = map[string][]string{
	"/v1/user/user-profile": {"admin", "user", "manager"},
	"/v1/assessment/":       {"admin", "user", "manager"},
	"/v1/admin/":            {"admin", "manager"},
	"/v1/question_author/":  {"admin", "questionnaire_author"},
}

func AdminRoutes(g *gin.RouterGroup, adminController *controller.AdminController, assessmentController *controller.AssessmentController, mastersController *controller.MastersController) {
	admin := g.Group("/admin")
	for _, adminRoute := range getAdminRoutes(adminController, assessmentController, mastersController) {
		protectedHandler := auth.Authenticate(admin.BasePath()+adminRoute.Path, ProtectedRoutes, adminRoute.HandleFunc)
		switch adminRoute.Method {
		case http.MethodPost:
			admin.POST(adminRoute.Path, protectedHandler)
		case http.MethodGet:
			admin.GET(adminRoute.Path, protectedHandler)
		case http.MethodPut:
			admin.PUT(adminRoute.Path, protectedHandler)
		case http.MethodDelete:
			admin.DELETE(adminRoute.Path, protectedHandler)
		}
	}
}

func QuestionAuthorRoutes(g *gin.RouterGroup, adminController *controller.AdminController, assessmentController *controller.AssessmentController, mastersController *controller.MastersController) {
	admin := g.Group("/question_author")
	for _, adminRoute := range getQuestionAuthorRoutes(adminController, assessmentController, mastersController) {
		protectedHandler := auth.Authenticate(admin.BasePath()+adminRoute.Path, ProtectedRoutes, adminRoute.HandleFunc)
		switch adminRoute.Method {
		case http.MethodPost:
			admin.POST(adminRoute.Path, protectedHandler)
		case http.MethodGet:
			admin.GET(adminRoute.Path, protectedHandler)
		case http.MethodPut:
			admin.PUT(adminRoute.Path, protectedHandler)
		case http.MethodDelete:
			admin.DELETE(adminRoute.Path, protectedHandler)
		}
	}
}

func UserRoutes(g *gin.RouterGroup, userController *controller.UserController) {
	user := g.Group("/user")
	user.Use(middleware.RateLimit())
	for _, userRoute := range getUserRoutes(userController) {
		protectedHandler := auth.Authenticate(user.BasePath()+userRoute.Path, ProtectedRoutes, userRoute.HandleFunc)
		switch userRoute.Method {
		case http.MethodPost:
			user.POST(userRoute.Path, protectedHandler)
		case http.MethodGet:
			user.GET(userRoute.Path, protectedHandler)
		}
	}
}

func AssessmentRoutes(g *gin.RouterGroup, assessmentController *controller.AssessmentController) {
	assessment := g.Group("/assessment")
	for _, assessmentRoute := range getAssessmentRoutes(assessmentController) {
		protectedHandler := auth.Authenticate(assessment.BasePath()+assessmentRoute.Path, ProtectedRoutes, assessmentRoute.HandleFunc)
		switch assessmentRoute.Method {
		case http.MethodPost:
			assessment.POST(assessmentRoute.Path, protectedHandler)
		case http.MethodGet:
			assessment.GET(assessmentRoute.Path, protectedHandler)
		}
	}
}

func OpenRoutes(g *gin.RouterGroup, publicController *controller.PublicController, adminController *controller.AdminController) {
	public := g.Group("/public")
	public.Use(middleware.RateLimit())
	for _, publicRoute := range getOpenRoutes(publicController, adminController) {
		switch publicRoute.Method {
		case http.MethodPost:
			public.POST(publicRoute.Path, publicRoute.HandleFunc)
		case http.MethodGet:
			public.GET(publicRoute.Path, publicRoute.HandleFunc)
		case http.MethodPut:
			public.PUT(publicRoute.Path, publicRoute.HandleFunc)
		case http.MethodDelete:
			public.DELETE(publicRoute.Path, publicRoute.HandleFunc)
		}
	}
}

func Routing(envFile string) {
	r := routes{
		router: gin.Default(),
	}

	middleware.StartRateLimiterCleanup()

	corsOrigins := strings.Split(os.Getenv("CORS_ORIGINS"), ",")
	r.router.Use(cors.New(cors.Config{
		AllowOrigins:     corsOrigins,
		AllowMethods:     []string{"GET", "POST", "PUT", "DELETE"},
		AllowHeaders:     []string{"Content-Type", "Content-Length", "Accept-Encoding", "Authorization", "Cache-Control"},
		AllowCredentials: true,
	}))
	r.router.GET("/", func(ctx *gin.Context) {
		ctx.JSON(http.StatusOK, gin.H{"message": "DHL server running..."})
	})
	r.router.GET(constant.Version, func(c *gin.Context) { c.JSON(http.StatusOK, gin.H{"version": utils.GetBuildVersion()}) })
	apiGroup := r.router.Group(os.Getenv("ApiVersion"))
	db := database.GetDBConn()
	InitializeRoutes(apiGroup, db)
	if envFile == "dev" {
		r.router.Run(":" + os.Getenv("GO_SERVER_PORT"))
	} else {
		err := r.router.Run(":" + os.Getenv("GO_SERVER_PORT"))
		if err != nil {
			log.Fatal("Failed to start HTTPS server: ", err)
		}
	}
}
