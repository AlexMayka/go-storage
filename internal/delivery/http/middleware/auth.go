package middleware

import (
	"context"
	"github.com/gin-gonic/gin"
	"go-storage/internal/config"
	"go-storage/internal/domain"
	"go-storage/pkg/errors"
	"go-storage/pkg/jwt"
	"go-storage/pkg/logger"
	"strings"
	"sync"
	"time"
)

type UseCaseAuthInterface interface {
	GetRolePermissionsByRoleId(ctx context.Context, roleId string) (*[]domain.Permission, error)
}

type AuthMiddleware struct {
	m    sync.RWMutex
	uc   UseCaseAuthInterface
	cash map[string]AuthCash
}

type AuthCash struct {
	per []domain.Permission
	ext time.Time
}

func NewAuthMiddleware(uc UseCaseAuthInterface) *AuthMiddleware {
	am := &AuthMiddleware{
		uc:   uc,
		cash: make(map[string]AuthCash),
	}

	go am.startCleanup()

	return am
}

func (a *AuthMiddleware) RequireAuth() gin.HandlerFunc {
	return func(c *gin.Context) {
		log := logger.FromContext(c.Request.Context())
		cnfPtr := config.FromContext(c.Request.Context())
		if cnfPtr == nil {
			log.Error("config not found in context")
			errors.HandleError(c, errors.InternalServer("configuration error"))
			c.Abort()
			return
		}
		cnf := *cnfPtr

		header := c.Request.Header.Get("Authorization")
		if header == "" {
			log.Error("missing authorization header")
			errors.HandleError(c, errors.Unauthorized("authorization header required"))
			c.Abort()
			return
		}

		if !strings.HasPrefix(header, "Bearer ") {
			log.Error("invalid authorization", "authHeader", header)
			errors.HandleError(c, errors.Unauthorized("invalid authorization"))
			c.Abort()
			return
		}

		tokenStr := strings.TrimPrefix(header, "Bearer ")
		tokenJWT, err := jwt.ParseToken(tokenStr, []byte(cnf.App.JwtSecret))
		if err != nil {
			log.Error("invalid JWT token", "error", err.Error())
			errors.HandleError(c, errors.Unauthorized(err.Error()))
			c.Abort()
			return
		}

		c.Set("user_id", tokenJWT.UserID)
		c.Set("role_id", tokenJWT.RoleID)
		c.Set("company_id", tokenJWT.CompanyID)

		c.Next()
	}
}

func (a *AuthMiddleware) RequireAnyPermission(permissions []string) gin.HandlerFunc {
	return func(c *gin.Context) {
		log := logger.FromContext(c.Request.Context())
		roleId := c.GetString("role_id")

		a.m.RLock()
		cashPermissions, ok := a.cash[roleId]
		a.m.RUnlock()

		var userPermissions *[]domain.Permission

		if !ok || cashPermissions.ext.Before(time.Now()) {
			dbPermissions, err := a.uc.GetRolePermissionsByRoleId(c.Request.Context(), roleId)
			if err != nil {
				log.Error("get role userPermissions error")
				errors.HandleError(c, errors.Unauthorized(err.Error()))
				c.Abort()
				return
			}

			userPermissions = dbPermissions
			a.m.Lock()
			a.cash[roleId] = AuthCash{
				per: *dbPermissions,
				ext: time.Now().Add(time.Minute * 10),
			}
			a.m.Unlock()
		} else {
			permsCopy := make([]domain.Permission, len(cashPermissions.per))
			copy(permsCopy, cashPermissions.per)
			userPermissions = &permsCopy
		}

		check := CheckAnyPermission(userPermissions, permissions)
		if !check {
			log.Error("insufficient authority", "per")
			errors.HandleError(c, errors.Forbidden("insufficient permissions"))
			c.Abort()
			return
		}

		c.Next()
	}
}

func (a *AuthMiddleware) RequireAllPermission(permissions []string) gin.HandlerFunc {
	return func(c *gin.Context) {
		log := logger.FromContext(c.Request.Context())
		roleId := c.GetString("role_id")

		a.m.RLock()
		cashPermissions, ok := a.cash[roleId]
		a.m.RUnlock()

		var userPermissions *[]domain.Permission

		if !ok || cashPermissions.ext.Before(time.Now()) {
			dbPermissions, err := a.uc.GetRolePermissionsByRoleId(c.Request.Context(), roleId)
			if err != nil {
				log.Error("get role userPermissions error")
				errors.HandleError(c, errors.Unauthorized(err.Error()))
				c.Abort()
				return
			}

			userPermissions = dbPermissions
			a.m.Lock()
			a.cash[roleId] = AuthCash{
				per: *dbPermissions,
				ext: time.Now().Add(time.Minute * 10),
			}
			a.m.Unlock()
		} else {
			permsCopy := make([]domain.Permission, len(cashPermissions.per))
			copy(permsCopy, cashPermissions.per)
			userPermissions = &permsCopy
		}

		check := CheckAllPermission(userPermissions, permissions)
		if !check {
			log.Error("insufficient authority", "per")
			errors.HandleError(c, errors.Forbidden("insufficient permissions"))
			c.Abort()
			return
		}

		c.Next()
	}
}

func CheckAnyPermission(userPermissions *[]domain.Permission, permissions []string) bool {
	userPermsMap := make(map[string]interface{}, len(*userPermissions))
	for _, perm := range *userPermissions {
		userPermsMap[strings.ToLower(perm.Name)] = struct{}{}
	}

	for _, requiredPerm := range permissions {
		if _, ok := userPermsMap[strings.ToLower(requiredPerm)]; ok {
			return true
		}
	}

	return false
}

func CheckAllPermission(userPermissions *[]domain.Permission, permissions []string) bool {
	userPermsMap := make(map[string]interface{}, len(*userPermissions))
	for _, perm := range *userPermissions {
		userPermsMap[strings.ToLower(perm.Name)] = struct{}{}
	}

	verifiedRights := make([]string, 0, len(permissions))
	for _, requiredPerm := range permissions {
		if _, ok := userPermsMap[strings.ToLower(requiredPerm)]; ok {
			verifiedRights = append(verifiedRights, requiredPerm)
		}
	}

	return len(verifiedRights) == len(permissions)
}

func (a *AuthMiddleware) startCleanup() {
	ticker := time.NewTicker(30 * time.Minute)
	defer ticker.Stop()

	for {
		select {
		case <-ticker.C:
			a.cleanupExpiredCache()
		}
	}
}

func (a *AuthMiddleware) cleanupExpiredCache() {
	now := time.Now()
	a.m.Lock()
	defer a.m.Unlock()

	for roleId, cached := range a.cash {
		if cached.ext.Before(now) {
			delete(a.cash, roleId)
		}
	}
}
