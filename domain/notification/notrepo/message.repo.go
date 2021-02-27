package notrepo

import (
	"github.com/syronz/dict"
	"github.com/syronz/limberr"
	"omono/domain/base/message/basterm"
	"omono/domain/notification/notmodel"
	"omono/internal/core"
	"omono/internal/core/corerr"
	"omono/internal/core/corterm"
	"omono/internal/core/validator"
	"omono/internal/param"
	"omono/internal/types"
	"omono/pkg/helper"
	"reflect"
)

// MessageRepo for injecting engine
type MessageRepo struct {
	Engine *core.Engine
	Cols   []string
}

// ProvideMessageRepo is used in wire and initiate the Cols
func ProvideMessageRepo(engine *core.Engine) MessageRepo {
	return MessageRepo{
		Engine: engine,
		Cols:   helper.TagExtracter(reflect.TypeOf(notmodel.Message{}), notmodel.MessageTable),
	}
}

// FindByID finds the message via its id
func (p *MessageRepo) FindByID(fix types.FixedCol) (message notmodel.Message, err error) {
	err = p.Engine.ReadDB.Table(notmodel.MessageTable).
		Where("id = ?", fix.ID.ToUint64()).
		First(&message).Error

	message.ID = fix.ID
	err = p.dbError(err, "E8270758", message, corterm.List)

	return
}

// List returns an array of messages
func (p *MessageRepo) List(params param.Param) (messages []notmodel.Message, err error) {
	var colsStr string
	if colsStr, err = validator.CheckColumns(p.Cols, params.Select); err != nil {
		err = limberr.Take(err, "E8229822").Build()
		return
	}

	var whereStr string
	if whereStr, err = params.ParseWhere(p.Cols); err != nil {
		err = limberr.Take(err, "E8255021").Custom(corerr.ValidationFailedErr).Build()
		return
	}

	err = p.Engine.ReadDB.Table(notmodel.MessageTable).Select(colsStr).
		Where(whereStr).
		Order(params.Order).
		Limit(params.Limit).
		Offset(params.Offset).
		Find(&messages).Error

	err = p.dbError(err, "E8232329", notmodel.Message{}, corterm.List)

	return
}

// Count of messages, mainly calls with List
func (p *MessageRepo) Count(params param.Param) (count int64, err error) {
	var whereStr string
	if whereStr, err = params.ParseWhere(p.Cols); err != nil {
		err = limberr.Take(err, "E8292390").Custom(corerr.ValidationFailedErr).Build()
		return
	}

	err = p.Engine.ReadDB.Table(notmodel.MessageTable).
		Where(whereStr).
		Count(&count).Error

	err = p.dbError(err, "E8271205", notmodel.Message{}, corterm.List)
	return
}

// Save the message, in case it is not exist create it
func (p *MessageRepo) Save(message notmodel.Message) (u notmodel.Message, err error) {
	if err = p.Engine.DB.Table(notmodel.MessageTable).Save(&message).Error; err != nil {
		err = p.dbError(err, "E8275412", message, corterm.Updated)
	}

	p.Engine.DB.Table(notmodel.MessageTable).Where("id = ?", message.ID).Find(&u)
	return
}

// Create a message
func (p *MessageRepo) Create(message notmodel.Message) (u notmodel.Message, err error) {
	if err = p.Engine.DB.Table(notmodel.MessageTable).Create(&message).Scan(&u).Error; err != nil {
		err = p.dbError(err, "E8235444", message, corterm.Created)
	}
	return
}

// Delete the message
func (p *MessageRepo) Delete(message notmodel.Message) (err error) {
	if err = p.Engine.DB.Table(notmodel.MessageTable).Delete(&message).Error; err != nil {
		err = p.dbError(err, "E8220697", message, corterm.Deleted)
	}
	return
}

// dbError is an internal method for generate proper dataeace error
func (p *MessageRepo) dbError(err error, code string, message notmodel.Message, action string) error {
	switch corerr.ClearDbErr(err) {
	case corerr.Nil:
		err = nil

	case corerr.NotFoundErr:
		err = corerr.RecordNotFoundHelper(err, code, corterm.ID, message.ID, corterm.Messages)

	case corerr.ForeignErr:
		err = limberr.Take(err, code).
			Message(corerr.SomeVRelatedToThisVSoItIsNotV, dict.R(basterm.Users),
				dict.R(corterm.Message), dict.R(action)).
			Custom(corerr.ForeignErr).Build()

	case corerr.ValidationFailedErr:
		err = corerr.ValidationFailedHelper(err, code)

	default:
		err = limberr.Take(err, code).
			Message(corerr.InternalServerError).
			Custom(corerr.InternalServerErr).Build()
	}

	return err
}
