package matrepo

import (
	"github.com/syronz/limberr"
	"omono/domain/material/matmodel"
	"omono/internal/core"
	"omono/internal/core/corerr"
	"omono/internal/core/corterm"
	"omono/internal/core/validator"
	"omono/internal/param"
	"omono/internal/types"
	"omono/pkg/dict"
	"omono/pkg/helper"
	"reflect"
	"time"
)

// TagRepo for injecting engine
type TagRepo struct {
	Engine *core.Engine
	Cols   []string
}

// ProvideTagRepo is used in wire and initiate the Cols
func ProvideTagRepo(engine *core.Engine) TagRepo {
	return TagRepo{
		Engine: engine,
		Cols:   helper.TagExtracter(reflect.TypeOf(matmodel.Tag{}), matmodel.TagTable),
	}
}

// FindByID finds the tag via its id
func (p *TagRepo) FindByID(fix types.FixedCol) (tag matmodel.Tag, err error) {
	err = p.Engine.ReadDB.Table(matmodel.TagTable).
		Where("company_id = ? AND node_id = ? AND id = ?", fix.CompanyID, fix.NodeID, fix.ID.ToUint64()).
		First(&tag).Error

	tag.ID = fix.ID
	err = p.dbError(err, "E7173958", tag, corterm.List)

	return
}

// List returns an array of tags
func (p *TagRepo) List(params param.Param) (tags []matmodel.Tag, err error) {
	var colsStr string
	if colsStr, err = validator.CheckColumns(p.Cols, params.Select); err != nil {
		err = limberr.Take(err, "E7177215").Build()
		return
	}

	var whereStr string
	if whereStr, err = params.ParseWhereDelete(p.Cols); err != nil {
		err = limberr.Take(err, "E7110593").Custom(corerr.ValidationFailedErr).Build()
		return
	}

	err = p.Engine.ReadDB.Table(matmodel.TagTable).Select(colsStr).
		Where(whereStr).
		Order(params.Order).
		Limit(params.Limit).
		Offset(params.Offset).
		Find(&tags).Error

	err = p.dbError(err, "E7168596", matmodel.Tag{}, corterm.List)

	return
}

// Count of tags, mainly calls with List
func (p *TagRepo) Count(params param.Param) (count int64, err error) {
	var whereStr string
	if whereStr, err = params.ParseWhere(p.Cols); err != nil {
		err = limberr.Take(err, "E7114485").Custom(corerr.ValidationFailedErr).Build()
		return
	}

	err = p.Engine.ReadDB.Table(matmodel.TagTable).
		Where(whereStr).
		Count(&count).Error

	err = p.dbError(err, "E7114173", matmodel.Tag{}, corterm.List)
	return
}

// Save the tag, in case it is not exist create it
func (p *TagRepo) Save(tag matmodel.Tag) (u matmodel.Tag, err error) {
	if err = p.Engine.DB.Table(matmodel.TagTable).Save(&tag).Error; err != nil {
		err = p.dbError(err, "E7170407", tag, corterm.Updated)
	}

	p.Engine.DB.Table(matmodel.TagTable).Where("id = ?", tag.ID).Find(&u)
	return
}

// Create a tag
func (p *TagRepo) Create(tag matmodel.Tag) (u matmodel.Tag, err error) {
	if err = p.Engine.DB.Table(matmodel.TagTable).Create(&tag).Scan(&u).Error; err != nil {
		err = p.dbError(err, "E7130884", tag, corterm.Created)
	}
	return
}

// Delete the tag
func (p *TagRepo) Delete(tag matmodel.Tag) (err error) {
	now := time.Now()
	tag.DeletedAt = &now
	if err = p.Engine.DB.Table(matmodel.TagTable).Save(&tag).Error; err != nil {
		err = p.dbError(err, "E7118944", tag, corterm.Deleted)
	}
	return
}

// dbError is an internal method for generate proper dataeace error
func (p *TagRepo) dbError(err error, code string, tag matmodel.Tag, action string) error {
	switch corerr.ClearDbErr(err) {
	case corerr.Nil:
		err = nil

	case corerr.NotFoundErr:
		err = corerr.RecordNotFoundHelper(err, code, corterm.ID, tag.ID, corterm.Tags)

	case corerr.ForeignErr:
		err = limberr.Take(err, code).
			Message(corerr.SomeVRelatedToThisVSoItIsNotV, dict.R(corterm.Tags),
				dict.R(corterm.Tag), dict.R(action)).
			Custom(corerr.ForeignErr).Build()

	case corerr.DuplicateErr:
		err = limberr.Take(err, code).
			Message(corerr.VWithValueVAlreadyExist, dict.R(corterm.Tag), tag.Tag).
			Custom(corerr.DuplicateErr).Build()
		err = limberr.AddInvalidParam(err, "name", corerr.VisAlreadyExist, tag.Tag)

	case corerr.ValidationFailedErr:
		err = corerr.ValidationFailedHelper(err, code)

	default:
		err = limberr.Take(err, code).
			Message(corerr.InternalServerError).
			Custom(corerr.InternalServerErr).Build()
	}

	return err
}
