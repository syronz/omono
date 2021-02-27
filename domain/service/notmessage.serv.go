package service

import (
	"crypto/rand"
	"fmt"
	"github.com/syronz/limberr"
	"math/big"
	"omono/domain/notification/enum/messagestatus"
	"omono/domain/notification/notmodel"
	"omono/domain/notification/notrepo"
	"omono/internal/core"
	"omono/internal/core/coract"
	"omono/internal/core/corerr"
	"omono/internal/param"
	"omono/internal/types"
	"omono/pkg/glog"
	"time"
)

// NotMessageServ for injecting auth notrepo
type NotMessageServ struct {
	Repo   notrepo.MessageRepo
	Engine *core.Engine
}

// ProvideNotMessageService for message is used in wire
func ProvideNotMessageService(p notrepo.MessageRepo) NotMessageServ {
	return NotMessageServ{
		Repo:   p,
		Engine: p.Engine,
	}
}

// FindByID for getting message by it's id
func (p *NotMessageServ) FindByID(fix types.FixedCol) (message notmodel.Message, err error) {
	if message, err = p.Repo.FindByID(fix); err != nil {
		err = corerr.Tick(err, "E8218140", "can't fetch the message", fix.ID)
		return
	}

	return
}

// FindByHash is a safe way for view the notification
func (p *NotMessageServ) FindByHash(hash uint64) (message notmodel.Message, err error) {
	params := param.New()
	params.PreCondition = fmt.Sprintf("not_messages.hash = %v", hash)
	params.Order = "not_messages.id ASC"
	var messages []notmodel.Message

	if messages, err = p.Repo.List(params); err != nil {
		glog.CheckError(err, "error in messages list")
		return
	}

	if len(messages) == 0 {
		err = limberr.New("E8218240").Custom(corerr.NotFoundErr).
			Build()
		return
	}

	message = messages[0]

	message.ViewCount += 1
	message.Status = messagestatus.Seen
	now := time.Now()
	message.ViewedAt = &now

	if message, err = p.Save(message); err != nil {
		return
	}

	return message, err
}

// List of messages, it support pagination and search and return back count
func (p *NotMessageServ) List(params param.Param, scope string) (messages []notmodel.Message,
	count int64, err error) {

	// if params.CompanyID != 0 {
	// 	params.PreCondition = fmt.Sprintf(" not_messages.company_id = '%v'  ", params.CompanyID)
	// }

	switch scope {
	case "":
		fallthrough
	case "all":
		params.PreCondition = fmt.Sprintf(" not_messages.recepient_id = '%v' ", params.UserID)
	case "new":
		params.PreCondition = fmt.Sprintf(" not_messages.recepient_id = '%v' AND not_messages.status = '%v' ",
			params.UserID, messagestatus.New)
	case "sent":
		params.PreCondition = fmt.Sprintf(" not_messages.created_by = '%v' ",
			params.UserID)
	}

	if messages, err = p.Repo.List(params); err != nil {
		glog.CheckError(err, "error in messages list")
		return
	}

	if count, err = p.Repo.Count(params); err != nil {
		glog.CheckError(err, "error in messages count")
	}

	return
}

// Create a message
func (p *NotMessageServ) Create(message notmodel.Message) (createdMessage notmodel.Message, err error) {

	if err = message.Validate(coract.Save); err != nil {
		err = corerr.TickValidate(err, "E8257071", "validation failed in creating the message", message)
		return
	}

	ranGen, err := rand.Int(rand.Reader, big.NewInt(1<<50))
	if err != nil {
		err = limberr.New("can not create random generator", "E8275184").
			Custom(corerr.InternalServerErr).Build()
		return
	}
	message.Hash = ranGen.Uint64()
	message.Status = messagestatus.New

	if createdMessage, err = p.Repo.Create(message); err != nil {
		err = corerr.Tick(err, "E8282358", "message not created", message)
		return
	}

	return
}

// Save a message, if it is exist update it, if not create it
func (p *NotMessageServ) Save(message notmodel.Message) (savedMessage notmodel.Message, err error) {
	if err = message.Validate(coract.Save); err != nil {
		err = corerr.TickValidate(err, "E8257759", corerr.ValidationFailed, message)
		return
	}

	if savedMessage, err = p.Repo.Save(message); err != nil {
		err = corerr.Tick(err, "E8288113", "message not saved")
		return
	}

	return
}

// Delete message, it is soft delete
func (p *NotMessageServ) Delete(fix types.FixedCol) (message notmodel.Message, err error) {
	if message, err = p.FindByID(fix); err != nil {
		err = corerr.Tick(err, "E8217029", "message not found for deleting")
		return
	}

	if err = p.Repo.Delete(message); err != nil {
		err = corerr.Tick(err, "E8212469", "message not deleted")
		return
	}

	return
}

// Excel is used for export excel file
func (p *NotMessageServ) Excel(params param.Param) (messages []notmodel.Message, err error) {
	params.Limit = p.Engine.Envs.ToInt(core.ExcelMaxRows)
	params.Offset = 0
	params.Order = fmt.Sprintf("%v.id ASC", notmodel.MessageTable)

	if messages, err = p.Repo.List(params); err != nil {
		err = corerr.Tick(err, "E8229024", "cant generate the excel list for messages")
		return
	}

	return
}
