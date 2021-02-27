package submodel

import (
	"omono/domain/subscriber/enum/accounttype"
	"omono/internal/core/coract"
	"omono/internal/core/corerr"
	"omono/internal/core/corterm"
	"omono/internal/types"
	"omono/pkg/helper"

	"github.com/syronz/dict"
	"github.com/syronz/limberr"
)

// AccountTable is used inside the repo layer
const (
	AccountTable = "sub_accounts"
)

// Account model
type Account struct {
	types.FixedCol
	ParentID  *types.RowID `json:"parent_id"`
	Code      string       `gorm:"unique" json:"code"`
	NameEn    string       `gorm:"unique" json:"name_en,omitempty"`
	NameKu    *string      `gorm:"unique" json:"name_ku,omitempty" `
	Type      types.Enum   `json:"type,omitempty"`
	Status    types.Enum   `gorm:"default:'active';type:enum('active','inactive')" json:"status,omitempty"`
	Balance   float64      `json:"balance"`
	ReadOnly  bool         `gorm:"not null;default:0" json:"read_only"`
	Phones    []Phone      `gorm:"-" json:"phones" table:"-"`
	Childrens []Account    `gorm:"-" json:"childrens" table:"-"`
}

// Validate check the type of fields
func (p *Account) Validate(act coract.Action) (err error) {

	// switch act {
	// case coract.Save:

	// 	if len(p.Name) < 5 {
	// 		err = limberr.AddInvalidParam(err, "name",
	// 			corerr.MinimumAcceptedCharacterForVisV,
	// 			dict.R(corterm.Name), 5)
	// 	}

	// 	if len(p.Name) > 255 {
	// 		err = limberr.AddInvalidParam(err, "name",
	// 			corerr.MaximumAcceptedCharacterForVisV,
	// 			dict.R(corterm.Name), 255)
	// 	}

	if p.Code == "" {
		err = limberr.AddInvalidParam(err, "code",
			corerr.VisRequired, dict.R(corterm.Code))
	}

	// 	if len(p.Description) > 255 {
	// 		err = limberr.AddInvalidParam(err, "description",
	// 			corerr.MaximumAcceptedCharacterForVisV,
	// 			dict.R(corterm.Description), 255)
	// 	}
	// }

	// TODO: it should be checked after API has been created
	if ok, _ := helper.Includes(accounttype.List, p.Type); !ok {
		return limberr.AddInvalidParam(err, "type",
			corerr.AcceptedValueForVareV, dict.R(corterm.Type),
			accounttype.Join())
	}

	return err
}
