package service

import (
	"omono/domain/base"
	"omono/domain/base/basrepo"
	"omono/internal/core"
	"omono/internal/types"
	"omono/pkg/glog"
	"strings"

	"github.com/gin-gonic/gin"
)

// BasAccessServ defining auth service
type BasAccessServ struct {
	Repo   basrepo.AccessRepo
	Engine *core.Engine
}

// ProvideBasAccessService for auth is used in wire
func ProvideBasAccessService(p basrepo.AccessRepo) BasAccessServ {
	return BasAccessServ{Repo: p, Engine: p.Engine}
}

var cacheResource map[uint]string

func init() {
	cacheResource = make(map[uint]string)
}

// CheckAccess is used inside each method to findout if user has permission or not
func (p *BasAccessServ) CheckAccess(c *gin.Context, resource types.Resource) bool {
	var userID uint

	if userIDtmp, ok := c.Get("USER_ID"); ok {
		userID = userIDtmp.(uint)
	} else {
		return true
	}

	var resources string
	var ok bool

	if resources, ok = cacheResource[userID]; !ok {
		var err error
		resources, err = p.Repo.GetUserResources(userID)
		glog.CheckError(err, "error in finding the resources for user", userID)
		BasAccessAddToCache(userID, resources)
	}

	return !strings.Contains(resources, string(resource))

}

func IsSuperAdmin(userID uint) bool {
	return strings.Contains(cacheResource[userID], string(base.SuperAccess))
}

// CheckRange is used for checking if user has access to special range of data
func (p *BasAccessServ) CheckRange(companyID, nodeID uint64) bool {
	if companyID > 0 {
		if companyID != 1001 {
			return false
		}
	}

	return true
}

// BasAccessAddToCache add the resources to the cacheResource
func BasAccessAddToCache(userID uint, resources string) {
	cacheResource[userID] = resources
}

func BasAccessDeleteFromCache(userID uint) {
	delete(cacheResource, userID)
}

func BasAccessResetCache(userID uint) {
	cacheResource[userID] = ""
}

func BasAccessResetFullCache() {
	cacheResource = make(map[uint]string)
}
