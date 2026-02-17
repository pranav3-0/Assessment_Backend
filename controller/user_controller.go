package controller

import (
	"context"
	"dhl/config"
	"dhl/constant"
	"dhl/models"
	"dhl/services"
	"dhl/utils"
	"encoding/csv"
	"io"
	"log"
	"net/http"
	"strings"

	"github.com/gin-gonic/gin"
)

type UserController struct {
	userService services.UserService
	authService services.AuthService
}

func NewUserController(userService services.UserService, authService services.AuthService) *UserController {
	return &UserController{userService: userService, authService: authService}
}

func (uc *UserController) RegisterUser(ctx *gin.Context) {
	var req models.RegisterRequest
	if err := ctx.ShouldBindJSON(&req); err != nil {
		models.ErrorResponse(ctx, constant.Failure, http.StatusBadRequest, "Invalid request", nil, err)
		return
	}

	user, err := uc.authService.RegisterUser(ctx, req)
	if err != nil {
		models.ErrorResponse(ctx, constant.Failure, http.StatusInternalServerError, "Failed to register user", nil, err)
		return
	}

	models.SuccessResponse(ctx, constant.Success, http.StatusOK, "User registered successfully", user, nil, nil)
	return
}

type csvRowResult struct {
	Row    int         `json:"row"`
	Email  string      `json:"email,omitempty"`
	Phone  string      `json:"phone,omitempty"`
	Status string      `json:"status"`          // success | failed
	Error  string      `json:"error,omitempty"` // error message
	User   interface{} `json:"user,omitempty"`  // created user (optional)
}

func (uc *UserController) RegisterUsersFromCSV(ctx *gin.Context) {
	fileHeader, err := ctx.FormFile("file")
	if err != nil {
		models.ErrorResponse(ctx, constant.Failure, http.StatusBadRequest, "CSV file is required (form-data key: file)", nil, err)
		return
	}

	f, err := fileHeader.Open()
	if err != nil {
		models.ErrorResponse(ctx, constant.Failure, http.StatusInternalServerError, "Unable to open uploaded file", nil, err)
		return
	}
	defer f.Close()

	r := csv.NewReader(f)
	r.TrimLeadingSpace = true

	header, err := r.Read()
	if err != nil {
		models.ErrorResponse(ctx, constant.Failure, http.StatusBadRequest, "Invalid CSV (cannot read header)", nil, err)
		return
	}

	// header -> index map
	col := map[string]int{}
	for i, h := range header {
		key := strings.ToLower(strings.TrimSpace(h))
		col[key] = i
	}

	get := func(row []string, key string) string {
		idx, ok := col[key]
		if !ok || idx < 0 || idx >= len(row) {
			return ""
		}
		return strings.TrimSpace(row[idx])
	}

	parseRoles := func(s string) []string {
		s = strings.TrimSpace(s)
		s = strings.Trim(s, `"'`) // remove quotes if any
		if s == "" {
			return nil
		}
		parts := strings.Split(s, ",")
		out := make([]string, 0, len(parts))
		for _, p := range parts {
			role := strings.TrimSpace(p)
			if role != "" {
				out = append(out, role)
			}
		}
		return out
	}

	var (
		results   []csvRowResult
		successes int
		failures  int
		rowNum    = 1 // header is row 1
	)

	for {
		record, err := r.Read()
		if err == io.EOF {
			break
		}
		rowNum++

		if err != nil {
			failures++
			results = append(results, csvRowResult{
				Row:    rowNum,
				Status: "failed",
				Error:  "csv read error: " + err.Error(),
			})
			continue
		}

		// Skip fully empty lines
		empty := true
		for _, v := range record {
			if strings.TrimSpace(v) != "" {
				empty = false
				break
			}
		}
		if empty {
			continue
		}
		roles := parseRoles(get(record, "roles"))
		if len(roles) == 0 {
			roles = []string{"user"}
		}

		req := models.RegisterRequest{
			ClientName: get(record, "client_name"),
			FirstName:  get(record, "first_name"),
			LastName:   get(record, "last_name"),
			Email:      get(record, "email"),
			Phone:      get(record, "phone"),
			Roles:      roles,
			Password:   "password",
		}
		log.Println("Bulk user data", req.ClientName)

		// Optional: quick validation before service call
		// (binding tags won't run here because we aren't using ShouldBindJSON)
		if req.ClientName == "" || req.FirstName == "" || req.Email == "" || req.Phone == "" {
			failures++
			results = append(results, csvRowResult{
				Row:    rowNum,
				Email:  req.Email,
				Phone:  req.Phone,
				Status: "failed",
				Error:  "missing required fields (client_name, first_name, email, phone)",
			})
			continue
		}

		user, regErr := uc.authService.RegisterUser(ctx, req)
		if regErr != nil {
			failures++
			results = append(results, csvRowResult{
				Row:    rowNum,
				Email:  req.Email,
				Phone:  req.Phone,
				Status: "failed",
				Error:  regErr.Error(),
			})
			continue
		}

		successes++
		results = append(results, csvRowResult{
			Row:    rowNum,
			Email:  req.Email,
			Phone:  req.Phone,
			Status: "success",
			User:   user,
		})
	}

	payload := gin.H{
		"totalProcessed": successes + failures,
		"successCount":   successes,
		"failureCount":   failures,
		"results":        results,
	}

	models.SuccessResponse(ctx, constant.Success, http.StatusOK, "Users registered from CSV successfully", payload, nil, nil)
}

func (uc *UserController) LoginUser(ctx *gin.Context) {
	var req models.LoginRequest
	if err := ctx.ShouldBindJSON(&req); err != nil {
		models.ErrorResponse(ctx, constant.Failure, http.StatusBadRequest, "Invalid request", nil, err)
		return
	}

	user, err := uc.authService.LoginUser(ctx, req)
	if err != nil {
		models.ErrorResponse(ctx, constant.Failure, http.StatusInternalServerError, "Failed to login user", nil, err)
		return
	}

	models.SuccessResponse(ctx, constant.Success, http.StatusOK, "success", user, nil, nil)
	return
}

func (uc *UserController) LogoutUser(c *gin.Context) {
	var input struct {
		RefreshToken string `json:"refresh_token"`
	}
	if err := c.ShouldBindJSON(&input); err != nil {
		models.ErrorResponse(c, constant.Failure, http.StatusBadRequest, "Invalid refresh token", nil, err)
		return
	}

	client := config.Client
	ctx := context.Background()

	err := client.Logout(ctx, config.KeycloakClientID, config.KeycloakClientSecret, config.KeycloakRealm, input.RefreshToken)
	if err != nil {
		log.Println("Error logging out from Keycloak:", err)
		models.ErrorResponse(c, constant.Failure, http.StatusUnauthorized, "Error while user logout", nil, err)
		return
	}
	models.SuccessResponse(c, constant.Success, http.StatusOK, "User logged out successfully", nil, nil, nil)
}

func (uc *UserController) GetProfile(ctx *gin.Context) {
	_, userId, _, err := utils.GetUserIDFromContext(ctx, uc.userService.FindUserIdBySub)
	if err != nil {
		models.ErrorResponse(ctx, constant.Failure, http.StatusUnauthorized, err.Error(), nil, err)
		return
	}
	user, err := uc.userService.GetUserProfile(userId)
	if err != nil {
		models.ErrorResponse(ctx, constant.Failure, http.StatusInternalServerError, err.Error(), nil, err)
		return
	}

	models.SuccessResponse(ctx, constant.Success, http.StatusOK, "success", user, nil, nil)
	return
}
