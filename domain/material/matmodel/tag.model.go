package matmodel

import (
	"github.com/syronz/dict"
	"github.com/syronz/limberr"
	"omono/internal/core/coract"
	"omono/internal/core/corerr"
	"omono/internal/core/corterm"
	"omono/internal/types"
)

// TagTable is a global instance for working with tag
const (
	TagTable = "mat_tags"
)

// Tag model
type Tag struct {
	types.FixedCol
	ParentID  *types.RowID `json:"parent_id,omitempty"`
	Tag       string       `gorm:"not null;unique;type:varchar(50)" json:"tag,omitempty"`
	Type      string       `gorm:"type:varchar(50)" json:"type,omitempty"`
	KeywordEn string       `gorm:"type:varchar(50)" json:"keyword_en,omitempty"`
	KeywordKu string       `gorm:"type:varchar(50)" json:"keyword_Ku,omitempty"`
	KeywordAr string       `gorm:"type:varchar(50)" json:"keyword_Ar,omitempty"`
}

// Validate check the type of fields
func (p *Tag) Validate(act coract.Action) (err error) {

	switch act {
	case coract.Save:

		if p.Tag == "" {
			err = limberr.AddInvalidParam(err, "tag",
				corerr.VisRequired,
				dict.R(corterm.Tag))
		}

		if len(p.Tag) > 255 {
			err = limberr.AddInvalidParam(err, "tag",
				corerr.MaximumAcceptedCharacterForVisV,
				dict.R(corterm.Tag), 255)
		}

	}

	return err
}
