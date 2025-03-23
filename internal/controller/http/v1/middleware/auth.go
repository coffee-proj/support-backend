package middleware

import (
	"slices"
	"strings"

	resp "github.com/coffee/support/internal/controller/response"
	types "github.com/coffee/support/internal/entity/type"
	"github.com/gin-gonic/gin"
)

func (m *Middleware) CheckAccess(roles ...types.Role) gin.HandlerFunc {
	return func(ctx *gin.Context) {
		header := ctx.GetHeader("Authorization")

		if header == "" {
			resp.AbortErrMsg(ctx, foundErr)
			return
		}

		parts := strings.Split(header, " ")
		if len(parts) < 2 {
			resp.AbortErrMsg(ctx, bearerErr)
			return
		}

		bearer := parts[0]
		token := parts[1]

		if bearer != bearerType {
			resp.AbortErrMsg(ctx, bearerErr)
			return
		}

		claims, err := m.auth.ValidateToken(token)
		if err != nil {
			resp.AbortErrMsg(ctx, err)
			return
		}

		if len(roles) != 0 {
			access := false

			for _, role := range claims.Roles {
				if slices.Contains(roles, role) {
					access = true
					break
				}
			}

			if !access {
				resp.AbortErrMsg(ctx, forbiddenErr)
				return
			}
		}

		ctx.Set("userId", claims.Id)

		ctx.Next()
	}
}
