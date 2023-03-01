package firebase

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

func Allow(endpointRoles []string) gin.HandlerFunc {
	return func(ctx *gin.Context) {
		memberRoles := ctx.GetStringMap("memberRoles")

		allowed := false
		for _, endpointRole := range endpointRoles {
			for memberRole := range memberRoles {
				if !allowed {
					allowed = endpointRole == memberRole
				}
			}
		}

		if !allowed {
			ctx.JSON(http.StatusUnauthorized, gin.H{"error": "Unauthorized"})
			ctx.Abort()
			return
		}

		ctx.Next()
	}
}
