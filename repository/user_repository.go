package repository

import (
	"context"
	"dhl/models"
	"fmt"
	"log"
	"strings"
	"time"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

type UserRepository interface {
	CreateUser(tx *gorm.DB, ctx context.Context, user *models.AssessmentUser, roles []string) error
	FindUserIdBySub(sub string) (string, error)
	FindByUsername(username string) (models.UserWithRoles, error)
	FindByUserId(userID string) (models.AssessmentUser, error)
	UpdateUserProfile(userID uuid.UUID, data models.UserProfileUpdate) error
	AddUserToManagerMapping(tx *gorm.DB, mapping models.UserManagerMapping) (*models.UserManagerMapping, error)

	FetchAllUsersWithExtPaginated(offset, limit int, role *string) ([]models.UserFullData, int64, error)
	FetchUsersManagedBy(offset, limit int, managerID string, role *string) ([]models.UserFullData, int64, error)
	FetchAllUsers(userIds []*string) ([]models.UserWithRoles, error)
}

type UserRepositoryImpl struct {
	db *gorm.DB
}

func NewUserRepository(db *gorm.DB) UserRepository {
	if db == nil {
		panic("database instance is null")
	}
	return &UserRepositoryImpl{db: db}
}

func (r *UserRepositoryImpl) CreateUser(tx *gorm.DB, ctx context.Context, user *models.AssessmentUser, roles []string) error {
	if err := tx.WithContext(ctx).Create(user).Error; err != nil {
		return fmt.Errorf("failed to create user: %w", err)
	}
	userExt := models.DhlAssessmentUserMstExt{
		UserID:    user.UserID.String(),
		CreatedAt: time.Now(),
	}
	if err := tx.WithContext(ctx).Create(&userExt).Error; err != nil {
		return fmt.Errorf("failed to create user ext: %w", err)
	}

	if len(roles) == 0 {
		return nil
	}

	var roleIDs []int64
	if err := tx.WithContext(ctx).
		Table("assessment_user_role_mst").
		Where("role_label IN ?", roles).
		Pluck("role_id", &roleIDs).Error; err != nil {

		return fmt.Errorf("failed to fetch roles: %w", err)
	}

	if len(roleIDs) != len(roles) {
		return fmt.Errorf("some roles do not exist in assessment_user_role_mst")
	}

	mappings := make([]map[string]interface{}, 0)
	for _, roleID := range roleIDs {
		mappings = append(mappings, map[string]interface{}{
			"user_id": user.UserID,
			"role_id": roleID,
		})
	}

	if err := tx.WithContext(ctx).
		Table("assessment_user_role_mapping").
		Create(&mappings).Error; err != nil {

		return fmt.Errorf("failed to assign roles to user: %w", err)
	}

	return nil
}

func (r *UserRepositoryImpl) FindByUsername(username string) (models.UserWithRoles, error) {
	var user models.UserWithRoles

	query := `
SELECT 
	u.user_id,
	u.first_name,
	u.last_name,
	u.email,
	u.phone,
	u.username,
	u.auth_user_id,
	u.user_type,
	COALESCE(ARRAY_AGG(r.role_label ORDER BY r.role_label)
		FILTER (WHERE r.role_label IS NOT NULL), '{}') AS roles


		FROM assessment_user_mst u
		LEFT JOIN assessment_user_role_mapping m 
			ON u.user_id = m.user_id
		LEFT JOIN assessment_user_role_mst r 
			ON m.role_id = r.role_id
		WHERE LOWER(u.username) = LOWER(?)
		GROUP BY 
	u.user_id, u.first_name, u.last_name, u.email, u.phone, 
	u.username, u.auth_user_id, u.user_type;

	`

	err := r.db.Raw(query, username).Scan(&user).Error
	return user, err
}

func (r *UserRepositoryImpl) FindByUserId(userId string) (models.AssessmentUser, error) {
	var user models.AssessmentUser
	err := r.db.Where("user_id = ?", userId).First(&user).Error
	return user, err
}

func (r *UserRepositoryImpl) FindUserIdBySub(sub string) (string, error) {
	var userId string
	query := `
		SELECT u.user_id FROM assessment_user_mst u where u.auth_user_id = ?;
	`
	err := r.db.Raw(query, sub).Scan(&userId).Error
	return userId, err
}

func (r *UserRepositoryImpl) UpdateUser(userID string, updates map[string]interface{}) error {
	if len(updates) == 0 {
		return nil
	}

	return r.db.Model(&models.AssessmentUser{}).
		Where("user_id = ?", userID).
		Updates(updates).
		Error
}

func (r *UserRepositoryImpl) UpdateUserProfile(userID uuid.UUID, data models.UserProfileUpdate) error {
	tx := r.db.Begin()

	// === 1️⃣ Update basic user data ===
	userUpdates := map[string]interface{}{}
	if data.FirstName != nil {
		userUpdates["first_name"] = *data.FirstName
	}
	if data.LastName != nil {
		userUpdates["last_name"] = *data.LastName
	}
	if data.Email != nil {
		userUpdates["email"] = *data.Email
	}
	if data.Phone != nil {
		userUpdates["phone"] = *data.Phone
	}

	if data.AuthUserID != nil {
		userUpdates["auth_user_id"] = *data.AuthUserID
	}
	if data.Password != nil {
		userUpdates["password"] = *data.Password
	}
	if data.NotifyId != nil {
		userUpdates["notify_id"] = *data.NotifyId
	}
	userUpdates["updated_at"] = time.Now()
	if len(userUpdates) > 0 {
		if err := tx.Table("public.assessment_user_mst").
			Where("user_id = ?", userID).
			Updates(userUpdates).Error; err != nil {
			tx.Rollback()
			return err
		}
	}

	// === 2️⃣ Update extended info ===
	extUpdates := map[string]interface{}{}
	if data.CompanyID != nil {
		extUpdates["company_id"] = *data.CompanyID
	}
	if data.Karma != nil {
		extUpdates["karma"] = *data.Karma
	}
	if data.RankID != nil {
		extUpdates["rank_id"] = *data.RankID
	}
	if data.TeamLead != nil {
		extUpdates["team_lead"] = *data.TeamLead
	}
	if data.Manager != nil {
		extUpdates["manager"] = *data.Manager
	}
	if data.SeniorManager != nil {
		extUpdates["senior_manager"] = *data.SeniorManager
	}
	if data.SDL != nil {
		extUpdates["sdl"] = *data.SDL
	}
	if data.SLL != nil {
		extUpdates["sll"] = *data.SLL
	}
	if data.UserMap != nil {
		extUpdates["user_map"] = *data.UserMap
	}
	if data.EmpCode != nil {
		extUpdates["emp_code"] = *data.EmpCode
	}
	if data.Center != nil {
		extUpdates["center"] = *data.Center
	}
	if data.SelectionType != nil {
		extUpdates["selection_type"] = *data.SelectionType
	}

	if len(extUpdates) > 0 {
		var count int64
		if err := tx.Table("public.dhl_assessment_user_mst_ext").
			Where("user_id = ?", userID.String()).
			Count(&count).Error; err != nil {
			tx.Rollback()
			return err
		}
		if count == 0 {
			extUpdates["user_id"] = userID.String()
			if err := tx.Table("public.dhl_assessment_user_mst_ext").
				Create(extUpdates).Error; err != nil {
				tx.Rollback()
				return err
			}
		} else {
			if err := tx.Table("public.dhl_assessment_user_mst_ext").
				Where("user_id = ?", userID.String()).
				Updates(extUpdates).Error; err != nil {
				tx.Rollback()
				return err
			}
		}
	}

	// === 3️⃣ Update roles ===
	if data.Roles != nil {
		// First clear existing roles
		if err := tx.Table("public.assessment_user_role_mapping").
			Where("user_id = ?", userID.String()).
			Delete(nil).Error; err != nil {
			tx.Rollback()
			return err
		}

		// Map role_labels → role_ids
		if len(data.Roles) > 0 {
			var roleIDs []int
			if err := tx.Table("public.assessment_user_role_mst").
				Where("role_label IN ?", data.Roles).
				Pluck("role_id", &roleIDs).Error; err != nil {
				tx.Rollback()
				return err
			}

			// Insert new mappings
			for _, rid := range roleIDs {
				mapping := map[string]interface{}{
					"user_id": userID.String(),
					"role_id": rid,
				}
				if err := tx.Table("public.assessment_user_role_mapping").
					Create(mapping).Error; err != nil {
					tx.Rollback()
					return err
				}
			}
		}
	}

	// === Commit transaction ===
	if err := tx.Commit().Error; err != nil {
		tx.Rollback()
		return err
	}
	return nil
}

func (r *UserRepositoryImpl) AddUserToManagerMapping(tx *gorm.DB, mapping models.UserManagerMapping) (*models.UserManagerMapping, error) {
	err := tx.Create(&mapping).Error
	if err != nil {
		if strings.Contains(err.Error(), "unique_user_manager") {
			updateErr := tx.Model(&models.UserManagerMapping{}).
				Where("user_id = ? AND manager_id = ?", mapping.UserID, mapping.ManagerID).
				Update("is_active", true).Error

			if updateErr != nil {
				return nil, updateErr
			}
			return &mapping, nil
		}
		return nil, err
	}
	return &mapping, nil
}

func (r *UserRepositoryImpl) FetchAllUsers(userIds []*string) ([]models.UserWithRoles, error) {
	var users []models.UserWithRoles

	baseQuery := `
	SELECT 
		u.user_id,
		u.first_name,
		u.last_name,
		u.email,
		u.phone,
		u.username,
		u.is_active,
		ARRAY_AGG(r.role_label ORDER BY r.role_label) AS roles
	FROM 
		public.assessment_user_mst u
	JOIN 
		public.assessment_user_role_mapping m 
			ON u.user_id = m.user_id
	JOIN 
		public.assessment_user_role_mst r 
			ON m.role_id = r.role_id
	`

	if len(userIds) > 0 {
		baseQuery += "WHERE u.user_id IN (?)\n"
	}

	baseQuery += `
	GROUP BY 
		u.user_id, u.first_name, u.last_name, u.email, u.phone, u.username, u.is_active
	`

	if len(userIds) == 0 {
		baseQuery += "LIMIT 50;"
	}

	var err error
	if len(userIds) > 0 {
		err = r.db.Raw(baseQuery, userIds).Scan(&users).Error
	} else {
		err = r.db.Raw(baseQuery).Scan(&users).Error
	}

	if err != nil {
		log.Println(err)
		return nil, err
	}

	return users, nil
}

func (r *UserRepositoryImpl) FetchAllUsersWithExtPaginated(offset, limit int, role *string) ([]models.UserFullData, int64, error) {
	var total int64

	// Step 1: Count total users
	countQuery := `SELECT COUNT(*) FROM public.assessment_user_mst u`
	var countParams []interface{}
	if role != nil && *role != "" {
		countQuery += ` 
		WHERE EXISTS (
			SELECT 1 
			FROM public.assessment_user_role_mapping m
			JOIN public.assessment_user_role_mst r ON m.role_id = r.role_id
			WHERE u.user_id = m.user_id AND r.role_label ILIKE ?
		)`
		countParams = append(countParams, "%"+*role+"%")
	}

	if err := r.db.Raw(countQuery, countParams...).Scan(&total).Error; err != nil {
		return nil, 0, err
	}

	// Step 2: Main query — with role condition if provided
	query := `SELECT u.user_id, u.first_name, u.last_name, u.email, u.phone, u.username,u.user_type, ext.company_id,
	
		ext.karma, ext.rank_id, ext.team_lead, ext.manager, ext.senior_manager, ext.sdl, ext.sll, ext.user_map, ext.emp_code,
		ext.center, ext.selection_type,	COALESCE(string_agg(DISTINCT r.role_label, ','), '') AS roles
	FROM 
		public.assessment_user_mst u
	LEFT JOIN 
		public.dhl_assessment_user_mst_ext ext 
		ON u.user_id::text = ext.user_id
	LEFT JOIN 
		public.assessment_user_role_mapping m 
		ON u.user_id = m.user_id
	LEFT JOIN 
		public.assessment_user_role_mst r 
		ON m.role_id = r.role_id
	`
	// Apply role filter if provided
	if role != nil && *role != "" {
		query += ` WHERE r.role_label ILIKE ?`
	}

	// Group by user fields
	query += `
	GROUP BY 
		u.user_id, u.user_type, ext.company_id, ext.karma, ext.rank_id, ext.team_lead, ext.manager,
		ext.senior_manager, ext.sdl, ext.sll, ext.user_map, ext.emp_code,
		ext.center, ext.selection_type
	ORDER BY u.created_at DESC
	LIMIT ? OFFSET ?;
	`

	// Parameters for query
	params := []interface{}{
		"%" + *role + "%",
		limit,
		offset,
	}

	if role == nil || *role == "" {
		params = []interface{}{
			limit,
			offset,
		}
	}

	// Step 3: Temporary struct for scanning
	type userRow struct {
		UserID        string  `json:"user_id"`
		FirstName     string  `json:"first_name"`
		LastName      string  `json:"last_name"`
		Email         string  `json:"email"`
		Phone         string  `json:"phone"`
		Username      string  `json:"username"`
		CompanyID     *int    `json:"company_id"`
		Karma         *int    `json:"karma"`
		RankID        *int    `json:"rank_id"`
		TeamLead      *string `json:"team_lead"`
		Manager       *string `json:"manager"`
		SeniorManager *string `json:"senior_manager"`
		SDL           *string `json:"sdl"`
		SLL           *string `json:"sll"`
		UserMap       *string `json:"user_map"`
		EmpCode       *string `json:"emp_code"`
		Center        *int    `json:"center"`
		SelectionType *string `json:"selection_type"`
		Roles         string  `json:"roles"`
		UserType      string  `json:"user_type"`
	}

	var rows []userRow
	if err := r.db.Raw(query, params...).Scan(&rows).Error; err != nil {
		log.Println("Error fetching users:", err)
		return nil, 0, err
	}

	// Step 4: Convert to []models.UserFullData
	var users []models.UserFullData
	for _, row := range rows {
		user := models.UserFullData{
			UserID:        uuid.MustParse(row.UserID),
			FirstName:     row.FirstName,
			LastName:      row.LastName,
			Email:         row.Email,
			Phone:         row.Phone,
			Username:      row.Username,
			CompanyID:     row.CompanyID,
			Karma:         row.Karma,
			RankID:        row.RankID,
			TeamLead:      row.TeamLead,
			Manager:       row.Manager,
			SeniorManager: row.SeniorManager,
			SDL:           row.SDL,
			SLL:           row.SLL,
			UserMap:       row.UserMap,
			EmpCode:       row.EmpCode,
			Center:        row.Center,
			SelectionType: row.SelectionType,
			UserType:      row.UserType,
		}

		if row.Roles != "" {
			user.Roles = strings.Split(row.Roles, ",")
		} else {
			user.Roles = []string{}
		}
		users = append(users, user)
	}

	return users, total, nil
}

func (r *UserRepositoryImpl) FetchUsersManagedBy(offset, limit int, managerID string, role *string) ([]models.UserFullData, int64, error) {

	var total int64

	// 1. Count users under this manager
	countQuery := `
		SELECT COUNT(*)
		FROM user_manager_mapping um
		JOIN assessment_user_mst u ON um.user_id = u.user_id::text
		WHERE um.manager_id = ? AND um.is_active = TRUE
	`
	if role != nil && *role != "" {
		countQuery += `
			AND EXISTS (
				SELECT 1 
				FROM assessment_user_role_mapping m
				JOIN assessment_user_role_mst r ON m.role_id = r.role_id
				WHERE u.user_id::text = m.user_id AND r.role_label ILIKE ?
			)
		`
		if err := r.db.Raw(countQuery, managerID, "%"+*role+"%").Scan(&total).Error; err != nil {
			return nil, 0, err
		}
	} else {
		if err := r.db.Raw(countQuery, managerID).Scan(&total).Error; err != nil {
			return nil, 0, err
		}
	}

	// 2. Main query (same fields as admin but with JOIN to manager mapping)
	query := `
	SELECT 
		u.user_id, u.first_name, u.last_name, u.email, u.phone, u.username,
		ext.company_id, ext.karma, ext.rank_id, ext.team_lead, ext.manager,
		ext.senior_manager, ext.sdl, ext.sll, ext.user_map, ext.emp_code,
		ext.center, ext.selection_type,
		COALESCE(string_agg(DISTINCT r.role_label, ','), '') AS roles
	FROM user_manager_mapping um
	JOIN assessment_user_mst u ON um.user_id = u.user_id::text
	LEFT JOIN dhl_assessment_user_mst_ext ext ON u.user_id::text = ext.user_id::text
	LEFT JOIN assessment_user_role_mapping m ON u.user_id::text = m.user_id::text
	LEFT JOIN assessment_user_role_mst r ON m.role_id = r.role_id
	WHERE um.manager_id = ? AND um.is_active = TRUE
	`

	params := []interface{}{managerID}

	if role != nil && *role != "" {
		query += " AND r.role_label ILIKE ? "
		params = append(params, "%"+*role+"%")
	}

	query += `
	GROUP BY 
		u.user_id, ext.company_id, ext.karma, ext.rank_id, ext.team_lead, ext.manager,
		ext.senior_manager, ext.sdl, ext.sll, ext.user_map, ext.emp_code,
		ext.center, ext.selection_type
	ORDER BY u.created_at DESC
	LIMIT ? OFFSET ?
	`

	params = append(params, limit, offset)

	var rows []struct {
		UserID        string
		FirstName     string
		LastName      string
		Email         string
		Phone         string
		Username      string
		CompanyID     *int
		Karma         *int
		RankID        *int
		TeamLead      *string
		Manager       *string
		SeniorManager *string
		SDL           *string
		SLL           *string
		UserMap       *string
		EmpCode       *string
		Center        *int
		SelectionType *string
		Roles         string
	}

	if err := r.db.Raw(query, params...).Scan(&rows).Error; err != nil {
		return nil, 0, err
	}

	// Convert format same as admin
	var users []models.UserFullData
	for _, row := range rows {
		u := models.UserFullData{
			UserID:        uuid.MustParse(row.UserID),
			FirstName:     row.FirstName,
			LastName:      row.LastName,
			Email:         row.Email,
			Phone:         row.Phone,
			Username:      row.Username,
			CompanyID:     row.CompanyID,
			Karma:         row.Karma,
			RankID:        row.RankID,
			TeamLead:      row.TeamLead,
			Manager:       row.Manager,
			SeniorManager: row.SeniorManager,
			SDL:           row.SDL,
			SLL:           row.SLL,
			UserMap:       row.UserMap,
			EmpCode:       row.EmpCode,
			Center:        row.Center,
			SelectionType: row.SelectionType,
		}

		if row.Roles != "" {
			u.Roles = strings.Split(row.Roles, ",")
		} else {
			u.Roles = []string{}
		}

		users = append(users, u)
	}

	return users, total, nil
}
