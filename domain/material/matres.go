package material

import "omono/internal/types"

// list of resources for material domain
const (
	Domain string = "material"

	CompanyWrite types.Resource = "company:write"
	CompanyRead  types.Resource = "company:read"
	CompanyExcel types.Resource = "company:excel"

	ColorWrite types.Resource = "color:write"
	ColorRead  types.Resource = "color:read"
	ColorExcel types.Resource = "color:excel"

	GroupWrite types.Resource = "group:write"
	GroupRead  types.Resource = "group:read"
	GroupExcel types.Resource = "group:excel"

	UnitWrite types.Resource = "unit:write"
	UnitRead  types.Resource = "unit:read"
	UnitExcel types.Resource = "unit:excel"

	TagWrite types.Resource = "tag:write"
	TagRead  types.Resource = "tag:read"
	TagExcel types.Resource = "tag:excel"

	ProductWrite types.Resource = "product:write"
	ProductRead  types.Resource = "product:read"
	ProductExcel types.Resource = "product:excel"
)
