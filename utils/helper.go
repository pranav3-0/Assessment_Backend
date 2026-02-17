package utils

import (
	"dhl/models"
	"errors"
	"fmt"
	"math"
	"mime/multipart"
	"os"
	"reflect"
	"strconv"
	"strings"

	"github.com/gin-gonic/gin"
	"github.com/xuri/excelize/v2"
)

func GetBuildVersion() string {
	buildVersion := os.Getenv("BUILD_VERSION")
	return buildVersion
}

func GetPaginationParams(c *gin.Context) (int, int, int) {
	limit, err := strconv.Atoi(c.DefaultQuery("limit", "10"))
	if err != nil || limit <= 0 {
		limit = 10
	}

	page, err := strconv.Atoi(c.DefaultQuery("page", "1"))
	if err != nil || page <= 0 {
		page = 1
	}

	offset := (page - 1) * limit

	return page, limit, offset
}

func GetPagination(limit int, page int, offset int, totalRecords int64) *models.Pagination {
	totalPages := int(math.Ceil(float64(totalRecords) / float64(limit)))

	return &models.Pagination{
		Limit:      limit,
		Page:       page,
		Offset:     offset,
		Total:      totalRecords,
		TotalPages: int64(totalPages),
	}
}

func GetUserIDFromContext(ctx *gin.Context, getUserIdBySubFunc func(string) (string, error)) (string, string, bool, error) {
	sub, subExists := ctx.Get("sub")
	if !subExists {
		return "", "", false, errors.New("user not found")
	}
	roleValues, _ := ctx.Get("userRoles")
	var roles []string
	if roleValues != nil {
		roles = roleValues.([]string)
	}
	highestRole := resolveHighestRole(roles)

	delegateUserID := ctx.GetHeader("X-Delegate-User-Id")
	if delegateUserID != "" {
		return highestRole, delegateUserID, true, nil
	}
	userID, err := getUserIdBySubFunc(sub.(string))
	if err != nil {
		return "", "", false, fmt.Errorf("failed to get user ID by sub: %w", err)
	}
	return highestRole, userID, false, nil
}

func ParseExcelToJSON(file multipart.File, filename string) (*models.SheetAssessment, error) {
	f, err := excelize.OpenReader(file)
	if err != nil {
		return nil, err
	}
	defer f.Close()

	rows, err := f.GetRows(f.GetSheetName(0))
	if err != nil {
		return nil, err
	}

	var assessment models.SheetAssessment
	assessment.AssessmentName = filename

	var questions []models.SheetQuestion
	var currentQuestion *models.SheetQuestion

	for i, row := range rows {
		if i == 0 {
			// Skip header row
			continue
		}

		var (
			questionTitle = strings.TrimSpace(getCell(row, 0)) // A
			isCorrectStr  = strings.TrimSpace(getCell(row, 1)) // B
			scoreStr      = strings.TrimSpace(getCell(row, 2)) // C
			optionLabel   = strings.TrimSpace(getCell(row, 3)) // D
			questionType  = strings.TrimSpace(getCell(row, 4)) // E
			mandatoryStr  = strings.TrimSpace(getCell(row, 5)) // F
		)

		// Detect new question
		if questionTitle != "" {
			if currentQuestion != nil {
				questions = append(questions, *currentQuestion)
			}
			currentQuestion = &models.SheetQuestion{
				Title:             questionTitle,
				MandatoryToAnswer: strings.EqualFold(mandatoryStr, "true"),
				QuestionType:      questionType,
				Options:           []models.SheetOption{},
			}
		}

		if currentQuestion == nil {
			continue
		}

		isCorrect := strings.EqualFold(isCorrectStr, "true") || isCorrectStr == "1"
		score := 0
		fmt.Sscanf(scoreStr, "%d", &score)

		if optionLabel != "" {
			currentQuestion.Options = append(currentQuestion.Options, models.SheetOption{
				Label:     optionLabel,
				IsCorrect: isCorrect,
				Score:     score,
			})
		}
	}

	if currentQuestion != nil {
		questions = append(questions, *currentQuestion)
	}

	assessment.Questions = questions
	return &assessment, nil
}

func getCell(row []string, index int) string {
	if index < len(row) {
		return row[index]
	}
	return ""
}

func BuildUpdateMap(structPtr interface{}) map[string]interface{} {
	v := reflect.ValueOf(structPtr).Elem()
	t := v.Type()

	updateMap := make(map[string]interface{})

	for i := 0; i < v.NumField(); i++ {
		fieldValue := v.Field(i)
		fieldType := t.Field(i)

		dbTag := fieldType.Tag.Get("gorm")
		column := parseGormColumnName(dbTag)
		if column == "" {
			continue
		}

		if !fieldValue.IsZero() {
			updateMap[column] = fieldValue.Interface()
		}
	}

	return updateMap
}

func parseGormColumnName(tag string) string {
	if tag == "" {
		return ""
	}
	for _, part := range reflect.StructTag(tag).Get("column") {
		if part == ' ' {
			break
		}
	}
	for _, segment := range splitGormTag(tag) {
		if len(segment) > 7 && segment[:7] == "column:" {
			return segment[7:]
		}
	}
	return ""
}

func splitGormTag(tag string) []string {
	var parts []string
	start := 0
	for i := 0; i < len(tag); i++ {
		if tag[i] == ';' {
			parts = append(parts, tag[start:i])
			start = i + 1
		}
	}
	if start < len(tag) {
		parts = append(parts, tag[start:])
	}
	return parts
}

func ParseQuestionnaireExcelToJSON(file multipart.File, filename string) (*models.SheetAssessment, error) {
	// Load Excel file
	f, err := excelize.OpenReader(file)
	if err != nil {
		return nil, fmt.Errorf("cannot open excel: %w", err)
	}

	rows, err := f.GetRows("Sheet1") // change if your sheet name differs
	if err != nil {
		return nil, fmt.Errorf("cannot read sheet: %w", err)
	}

	assessment := &models.SheetAssessment{
		AssessmentName: "",
		Questions:      []models.SheetQuestion{},
	}

	var currentQ *models.SheetQuestion

	for i, row := range rows {
		if i == 0 {
			continue // skip header
		}

		// Ensure minimum columns exist
		for len(row) < 9 {
			row = append(row, "")
		}

		questionType := strings.TrimSpace(row[1])
		assessmentName := strings.TrimSpace(row[2])
		title := strings.TrimSpace(row[3])
		optLabel := strings.TrimSpace(row[4])
		scoreStr := strings.TrimSpace(row[5])
		correctStr := strings.TrimSpace(row[6])
		mandatoryStr := strings.TrimSpace(row[7]) // mandatory answer

		// Set assessment name (once)
		if assessment.AssessmentName == "" && assessmentName != "" {
			assessment.AssessmentName = assessmentName
		}

		// If title cell contains question → start new question
		if title != "" {
			if currentQ != nil {
				assessment.Questions = append(assessment.Questions, *currentQ)
			}

			mandatory := false
			if mandatoryStr == "1" || strings.ToUpper(mandatoryStr) == "TRUE" {
				mandatory = true
			}

			currentQ = &models.SheetQuestion{
				Title:             title,
				MandatoryToAnswer: mandatory,
				QuestionType:      questionType,
				Options:           []models.SheetOption{},
			}
		}

		// Create option only if label exists
		if optLabel != "" && currentQ != nil {
			score, _ := strconv.Atoi(scoreStr)
			isCorrect := strings.ToUpper(correctStr) == "TRUE"

			opt := models.SheetOption{
				Label:     optLabel,
				IsCorrect: isCorrect,
				Score:     score,
			}

			currentQ.Options = append(currentQ.Options, opt)
		}
	}

	// Append the final question
	if currentQ != nil {
		assessment.Questions = append(assessment.Questions, *currentQ)
	}

	return assessment, nil
}

var QuestionTypeMap = map[string]int{
	"simple_choice":   1,
	"multiple_choice": 2,
	"free_text":       3,
	"date":            4,
	"numerical_box":   5,
	"typing_test":     6,
	"textbox":         7,
	"matrix":          8,
}

func resolveHighestRole(roles []string) string {
	priority := map[string]int{
		"admin":   3,
		"manager": 2,
		"user":    1,
	}
	highestRole := ""
	highestVal := 0
	for _, r := range roles {
		if val, exists := priority[strings.ToLower(r)]; exists {
			if val > highestVal {
				highestVal = val
				highestRole = r
			}
		}
	}
	return highestRole
}
