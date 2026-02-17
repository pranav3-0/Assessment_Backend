package auth

import (
	"context"
	"dhl/config"
	"dhl/constant"
	"dhl/models"
	"log"
	"net/http"
	"strings"

	"github.com/gin-gonic/gin"
)

func AuthToken(requiredRoles ...string) gin.HandlerFunc {
	return func(c *gin.Context) {
		log.Printf("=== AUTH: Authenticating request to %s %s ===", c.Request.Method, c.Request.URL.Path)
		log.Printf("Required roles: %v", requiredRoles)

		authHeader := c.GetHeader("Authorization")
		if authHeader == "" {
			log.Println("ERROR AUTH: Authorization header missing")
			models.ErrorResponse(c, constant.Failure, http.StatusUnauthorized, "Authorization header missing", nil, nil)
			c.Abort()
			return
		}

		if !strings.HasPrefix(authHeader, "Bearer ") {
			headerPreview := authHeader
			if len(authHeader) > 20 {
				headerPreview = authHeader[:20]
			}
			log.Printf("ERROR AUTH: Invalid Authorization header format: %s", headerPreview)
			models.ErrorResponse(c, constant.Failure, http.StatusUnauthorized, "Invalid Authorization header format", nil, nil)
			c.Abort()
			return
		}
		headerPreview := authHeader
		if len(authHeader) > 27 {
			headerPreview = authHeader[:27] + "..."
		}
		log.Printf("Authorization header present: %s", headerPreview)

		tokenStr := authHeader[len("Bearer "):]
		client := config.Client
		ctx := context.Background()

		log.Println("Introspecting token with Keycloak...")
		introspection, err := client.RetrospectToken(ctx, tokenStr, config.KeycloakClientID, config.KeycloakClientSecret, config.KeycloakRealm)
		if err != nil {
			log.Printf("ERROR AUTH: Error introspecting token: %v", err)
			models.ErrorResponse(c, constant.Failure, http.StatusUnauthorized, "Invalid token", nil, err)
			c.Abort()
			return
		}
		log.Printf("Token introspection result - Active: %v", *introspection.Active)

		if !*introspection.Active {
			log.Println("ERROR AUTH: Token is not active (expired or invalid)")
			models.ErrorResponse(c, constant.Failure, http.StatusUnauthorized, "Token is expired!", nil, err)
			c.Abort()
			return
		}

		log.Println("Decoding access token...")
		_, claims, err := client.DecodeAccessToken(ctx, tokenStr, config.KeycloakRealm)
		if err != nil {
			log.Printf("ERROR AUTH: Error decoding access token: %v", err)
			c.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid token"})
			c.Abort()
			return
		}
		if claims == nil {
			log.Println("ERROR AUTH: Token claims are nil")
			models.ErrorResponse(c, constant.Failure, http.StatusUnauthorized, "Invalid token claims", nil, err)
			c.Abort()
			return
		}
		log.Println("Token decoded successfully")

		roles, ok := (*claims)["realm_access"].(map[string]interface{})["roles"].([]interface{})
		if !ok {
			log.Println("ERROR AUTH: Failed to extract roles from token claims")
			models.ErrorResponse(c, constant.Failure, http.StatusUnauthorized, "Invalid token claims", nil, err)
			c.Abort()
			return
		}
		log.Printf("User roles from token: %v", roles)

		sub, ok := (*claims)["sub"]
		if !ok {
			log.Println("ERROR AUTH: Failed to extract 'sub' from token claims")
			models.ErrorResponse(c, constant.Failure, http.StatusUnauthorized, "Invalid token claims", nil, err)
			c.Abort()
			return
		}
		log.Printf("User sub: %v", sub)

		var userRoles []string

		hasRequiredRole := false
		for _, requiredRole := range requiredRoles {
			for _, role := range roles {
				userRoles = append(userRoles, role.(string))
				if role == requiredRole {
					hasRequiredRole = true
					break
				}
			}
			if hasRequiredRole {
				break
			}
		}
		log.Printf("User has required role: %v (User roles: %v, Required roles: %v)", hasRequiredRole, userRoles, requiredRoles)

		if !hasRequiredRole {
			log.Printf("ERROR AUTH: Access denied - user doesn't have required role. User roles: %v, Required: %v", userRoles, requiredRoles)
			models.ErrorResponse(c, constant.Failure, http.StatusForbidden, "Access denied", nil, err)
			c.Abort()
			return
		}

		log.Printf("SUCCESS AUTH: User authenticated successfully with role check passed")
		// Store the claims in the context
		c.Set("claims", claims)
		c.Set("sub", sub)
		c.Set("userRoles", userRoles)
		c.Next()
	}
}

func Authenticate(path string, protectedRoutes map[string][]string, handler gin.HandlerFunc) gin.HandlerFunc {
	log.Printf("=== Setting up authentication for path: %s ===", path)
	for protectedPrefix, roles := range protectedRoutes {
		if strings.HasPrefix(path, protectedPrefix) {
			log.Printf("Path %s matches protected prefix %s with required roles: %v", path, protectedPrefix, roles)
			return gin.HandlerFunc(func(c *gin.Context) {
				AuthToken(roles...)(c)
				if c.IsAborted() {
					log.Println("Request aborted by AuthToken middleware")
					return
				}
				handler(c)
			})
		}
	}
	log.Printf("Path %s is not protected, no authentication required", path)
	return handler
}
