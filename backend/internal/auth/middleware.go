package auth

import (
	apperrors "backend/pkg/errors"
	"strings"

	"github.com/gin-gonic/gin"
)

func Middleware(authService *Service) gin.HandlerFunc {
	return func(c *gin.Context) {
		header := c.GetHeader("Authorization")
		if header == "" {
			apperrors.Handle(c, apperrors.NewUnauthorized("missing authorization header"))
			c.Abort()
			return
		}

		parts := strings.SplitN(header, " ", 2)
		if len(parts) != 2 || !strings.EqualFold(parts[0], "Bearer") {
			apperrors.Handle(c, apperrors.NewUnauthorized("invalid authorization header format"))
			c.Abort()
			return
		}

		claims, err := authService.ValidateToken(parts[1])
		if err != nil {
			apperrors.Handle(c, err)
			c.Abort()
			return
		}

		c.Set("userID", claims.UserID)
		c.Set("email", claims.Email)
		c.Next()
	}
}
