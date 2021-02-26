package basmid

import (
	"net/http"
	"omono/domain/base"
	"omono/domain/base/basrepo"
	"omono/domain/service"
	"omono/internal/core"
	"omono/internal/types"

	"github.com/gin-gonic/gin"
)

type accessMid struct {
	engine *core.Engine
}

// NewAccessMid is a simpler way for access to the struct
func NewAccessMid(engine *core.Engine) accessMid {
	return accessMid{
		engine: engine,
	}
}

// Check will analyze if the user should have access to special resource or not
func (p *accessMid) Check(resource types.Resource) gin.HandlerFunc {

	return func(c *gin.Context) {

		accessService := service.ProvideBasAccessService(basrepo.ProvideAccessRepo(p.engine))
		accessResult := accessService.CheckAccess(c, resource)

		if c.Query("deleted") == "true" {
			accessResult = accessService.CheckAccess(c, base.ReadDeleted)
		}

		if accessResult {
			//TODO: Implement custom error handling
			c.AbortWithStatusJSON(http.StatusForbidden, gin.H{"msg": "you don't have permission"})
			return
		}

		c.Next()

	}
}
