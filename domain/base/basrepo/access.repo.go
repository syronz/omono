package basrepo

import (
	"omono/domain/base/basmodel"
	"omono/internal/core"
)

// AccessRepo for injecting engine
type AccessRepo struct {
	Engine *core.Engine
}

// ProvideAccessRepo is used in wire
func ProvideAccessRepo(engine *core.Engine) AccessRepo {
	return AccessRepo{Engine: engine}
}

// GetUserResources is used for finding all resources
func (p *AccessRepo) GetUserResources(userID uint) (result string, err error) {
	resources := struct {
		Resources string
	}{}

	err = p.Engine.ReadDB.Table(basmodel.UserTable).Select("bas_roles.resources").
		Joins("INNER JOIN bas_roles ON bas_users.role_id = bas_roles.id").
		Where("bas_users.id = ?", userID).Scan(&resources).Error

	result = resources.Resources

	return
}
