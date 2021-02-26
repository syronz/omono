package material

import "omono/internal/types"

// types for material domain
const (
	CreateCompany types.Event = "company-create"
	UpdateCompany types.Event = "company-update"
	DeleteCompany types.Event = "company-delete"
	ListCompany   types.Event = "company-list"
	ViewCompany   types.Event = "company-view"
	ExcelCompany  types.Event = "company-excel"

	CreateColor types.Event = "color-create"
	UpdateColor types.Event = "color-update"
	DeleteColor types.Event = "color-delete"
	ListColor   types.Event = "color-list"
	ViewColor   types.Event = "color-view"
	ExcelColor  types.Event = "color-excel"

	CreateGroup types.Event = "group-create"
	UpdateGroup types.Event = "group-update"
	DeleteGroup types.Event = "group-delete"
	ListGroup   types.Event = "group-list"
	ViewGroup   types.Event = "group-view"
	ExcelGroup  types.Event = "group-excel"

	CreateUnit types.Event = "unit-create"
	UpdateUnit types.Event = "unit-update"
	DeleteUnit types.Event = "unit-delete"
	ListUnit   types.Event = "unit-list"
	ViewUnit   types.Event = "unit-view"
	ExcelUnit  types.Event = "unit-excel"

	CreateTag types.Event = "tag-create"
	UpdateTag types.Event = "tag-update"
	DeleteTag types.Event = "tag-delete"
	ListTag   types.Event = "tag-list"
	ViewTag   types.Event = "tag-view"
	ExcelTag  types.Event = "tag-excel"

	CreateProduct types.Event = "product-create"
	UpdateProduct types.Event = "product-update"
	DeleteProduct types.Event = "product-delete"
	ListProduct   types.Event = "product-list"
	ViewProduct   types.Event = "product-view"
	ExcelProduct  types.Event = "product-excel"
)
