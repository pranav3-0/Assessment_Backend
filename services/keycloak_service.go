package services

import (
	"context"
	"dhl/config"
	"dhl/models"
	"encoding/json"
	"fmt"
	"log"
	"os"

	"github.com/Nerzal/gocloak/v13"
)

type KeycloakService struct {
	config map[string]interface{}
}

func NewKeycloakService(configJSON string) *KeycloakService {
	var config map[string]interface{}
	_ = json.Unmarshal([]byte(configJSON), &config)
	return &KeycloakService{config: config}
}

func (k *KeycloakService) LoginUser(ctx context.Context, req models.LoginRequest) (*gocloak.JWT, error) {
	KeycloakRealm := k.config["KEYCLOAK_REALM"].(string)
	KeycloakClientID := k.config["KEYCLOAK_CLIENT_ID"].(string)
	KeycloakClientSecret := k.config["KEYCLOAK_CLIENT_SECRET"].(string)
	KeycloakAuthURL := k.config["KEYCLOAK_AUTH_URL"].(string)

	client := gocloak.NewClient(KeycloakAuthURL)

	token, err := client.Login(ctx, KeycloakClientID, KeycloakClientSecret, KeycloakRealm, req.Username, req.Password)
	if err != nil {
		log.Printf("[ERROR] Keycloak login failed for %s: %v", req.Username, err)
		return nil, fmt.Errorf("invalid credentials")
	}

	if token.AccessToken == "" {
		return nil, fmt.Errorf("access token not found")
	}

	return token, nil
}

func (k *KeycloakService) CreateUser(ctx context.Context, user models.RegisterRequest) (string, error) {
	KeycloakRealm := k.config["KEYCLOAK_REALM"].(string)
	KeycloakClientID := k.config["KEYCLOAK_CLIENT_ID"].(string)
	KeycloakAuthURL := k.config["KEYCLOAK_AUTH_URL"].(string)
	KeycloakClientSecret := k.config["KEYCLOAK_CLIENT_SECRET"].(string)
	client := gocloak.NewClient(KeycloakAuthURL)
	token, err := client.LoginClient(ctx, KeycloakClientID, KeycloakClientSecret, KeycloakRealm)
	if err != nil {
		log.Printf("[ERROR] Failed to login to Keycloak: %v", err)
		return "", err
	}
	newuser := gocloak.User{
		Username:      gocloak.StringP(user.Email),
		Email:         gocloak.StringP(user.Email),
		FirstName:     gocloak.StringP(user.FirstName),
		LastName:      gocloak.StringP(user.LastName),
		Enabled:       gocloak.BoolP(true),
		EmailVerified: gocloak.BoolP(true),
		Attributes: &map[string][]string{
			"phoneNumber": {user.Phone},
		},
		Credentials: &[]gocloak.CredentialRepresentation{
			{
				Type:      gocloak.StringP("password"),
				Value:     gocloak.StringP(user.Password),
				Temporary: gocloak.BoolP(false),
			},
		},
		RealmRoles: &user.Roles,
	}

	userID, err := client.CreateUser(ctx, token.AccessToken, KeycloakRealm, newuser)
	if err != nil {
		return "", fmt.Errorf("user creation failed: %v", err)
	}

	for _, uRole := range user.Roles {
		role, roleErr := client.GetRealmRole(ctx, token.AccessToken, KeycloakRealm, uRole)
		if roleErr != nil {
			log.Printf("[ERROR] Role not found in Keycloak: %v", err)
			return "User role not found at keycloak server", roleErr
		}
		adderr := client.AddRealmRoleToUser(ctx, token.AccessToken, KeycloakRealm, userID, []gocloak.Role{*role})
		if adderr != nil {
			return "Unable to add role to user", adderr
		}
	}

	return userID, nil
}

func (kc *KeycloakService) UpdateUserInKeycloak(keycloakUserID string, userData map[string]interface{}, roles []string) error {
	client := config.Client
	ctx := context.Background()
	username := os.Getenv("KEYCLOAK_ADMIN_USER")
	password := os.Getenv("KEYCLOAK_ADMIN_PASSWORD")
	token, err := client.LoginAdmin(ctx, username, password, "master")
	if err != nil {
		return err
	}

	user := gocloak.User{
		ID:        gocloak.StringP(keycloakUserID),
		FirstName: getStrPtr(userData["firstName"]),
		LastName:  getStrPtr(userData["lastName"]),
		Email:     getStrPtr(userData["email"]),
		Attributes: &map[string][]string{
			"phone": {getStr(userData["phone"])},
		},
	}

	if err := client.UpdateUser(ctx, token.AccessToken, config.KeycloakRealm, user); err != nil {
		return fmt.Errorf("failed to update keycloak user: %w", err)
	}

	// 🔹 If roles provided, update realm-level roles
	if len(roles) > 0 {
		allRoles, err := client.GetRealmRoles(ctx, token.AccessToken, config.KeycloakRealm, gocloak.GetRoleParams{})
		if err != nil {
			return fmt.Errorf("failed to get realm roles: %w", err)
		}

		var assignRoles []gocloak.Role
		for _, roleName := range roles {
			for _, role := range allRoles {
				if *role.Name == roleName {
					assignRoles = append(assignRoles, *role)
				}
			}
		}

		rolesToDelete := make([]gocloak.Role, len(allRoles))
		for i, r := range allRoles {
			rolesToDelete[i] = *r
		}

		if len(assignRoles) > 0 {
			if err := client.DeleteRealmRoleFromUser(ctx, token.AccessToken, config.KeycloakRealm, keycloakUserID, rolesToDelete); err != nil {
				return err
			}
			if err := client.AddRealmRoleToUser(ctx, token.AccessToken, config.KeycloakRealm, keycloakUserID, assignRoles); err != nil {
				return fmt.Errorf("failed to assign roles: %w", err)
			}
		}
	}

	return nil
}

func getStr(v interface{}) string {
	if s, ok := v.(string); ok {
		return s
	}
	return ""
}

func getStrPtr(v interface{}) *string {
	s := getStr(v)
	if s == "" {
		return nil
	}
	return &s
}

func (k *KeycloakService) MigrateRoles(roles []string) {
	ctx := context.Background()
	keycloakAuthURL := os.Getenv("KEYCLOAK_AUTH_URL")
	realm := os.Getenv("KEYCLOAK_REALM")
	client := gocloak.NewClient(keycloakAuthURL)
	token, err := client.LoginAdmin(ctx, os.Getenv("KEYCLOAK_ADMIN_USER"), os.Getenv("KEYCLOAK_ADMIN_PASSWORD"), os.Getenv("KEYCLOAK_MASTER"))
	if err != nil {
		log.Println("Error login Admin:", err)
	}

	for _, role := range roles {
		_, err := client.CreateRealmRole(ctx, token.AccessToken, realm, gocloak.Role{
			Name: gocloak.StringP(role),
		})
		if err != nil {
			log.Println("Error creating role:", role, err)
			continue
		}
	}
}
