package service

import (
	"encoding/json"
	"fmt"
	"omono/domain/base"
	"omono/domain/base/basmodel"
	"omono/domain/base/basrepo"
	"omono/internal/core"
	"omono/internal/param"
	"omono/internal/types"
	"omono/pkg/glog"
	"time"

	"github.com/gin-gonic/gin"
)

// RecordType is and int used as an enum
type RecordType int

const (
	read RecordType = iota
	writeBefore
	writeAfter
	writeBoth
)

// BasActivityServ for injecting auth basrepo
type BasActivityServ struct {
	Repo   basrepo.ActivityRepo
	Engine *core.Engine
}

// ProvideBasActivityService for activity is used in wire
func ProvideBasActivityService(p basrepo.ActivityRepo) BasActivityServ {
	return BasActivityServ{Repo: p, Engine: p.Engine}
}

// Save activity
func (p *BasActivityServ) Save(activity basmodel.Activity) (createdActivity basmodel.Activity, err error) {
	createdActivity, err = p.Repo.Create(activity)

	// p.Engine.CheckInfo(err, fmt.Sprintf("Failed in saving activity for %+v", activity))

	return
}

// ActivityWatcher is used for watching activity channel
func (p *BasActivityServ) ActivityWatcher() {
	var arr []basmodel.Activity
	counter := 0
	var activity basmodel.Activity

	tickTimer := time.Tick(p.Engine.Envs.ToDuration(base.ActivityTickTimer) * time.Second)

	for {
		select {
		case activity = <-p.Engine.ActivityCh:
			counter++
			arr = append(arr, activity)
			if counter > p.Engine.Envs.ToInt(base.ActivityFileCounter) {
				p.Repo.CreateBatch(arr)
				counter = 0
				arr = []basmodel.Activity{}
			}
		case <-tickTimer:
			if len(arr) > 0 {
				p.Repo.CreateBatch(arr)
				counter = 0
				arr = []basmodel.Activity{}
			}
		}
	}
}

// Record will save the activity
// TODO: Record is deprecated we should go with channels
func (p *BasActivityServ) Record(c *gin.Context, ev types.Event, data ...interface{}) {
	var userID types.RowID
	var username string

	recordType := p.FindRecordType(data...)
	before, after := p.FillBeforeAfter(recordType, data...)

	if len(data) > 0 && !p.Engine.Envs.ToBool(base.RecordWrite) {
		return
	}

	if len(data) == 0 && !p.Engine.Envs.ToBool(base.RecordRead) {
		return
	}

	if p.IsRecordSetInEnvironment(recordType) {
		return
	}
	if userIDtmp, ok := c.Get("USER_ID"); ok {
		userID = userIDtmp.(types.RowID)
	}
	if usernameTmp, ok := c.Get("USERNAME"); ok {
		username = usernameTmp.(string)
	}

	activity := basmodel.Activity{
		Event:    ev.String(),
		UserID:   userID,
		Username: username,
		IP:       c.ClientIP(),
		URI:      c.Request.RequestURI,
		Before:   string(before),
		After:    string(after),
	}

	_, err := p.Repo.Create(activity)
	glog.CheckError(err, fmt.Sprintf("Failed in saving activity for %+v", activity))
}

// FillBeforeAfter check if there is a need for entering before data or not
func (p *BasActivityServ) FillBeforeAfter(recordType RecordType, data ...interface{}) (before, after []byte) {
	var err error
	if recordType == writeBefore || recordType == writeBoth {
		before, err = json.Marshal(data[0])
		glog.CheckError(err, "error in encoding data to before-json")
	}
	if recordType == writeAfter || recordType == writeBoth {
		after, err = json.Marshal(data[1])
		glog.CheckError(err, "error in encoding data to after-json")
	}

	return
}

// FindRecordType is helper function for finding the best way for recording data
func (p *BasActivityServ) FindRecordType(data ...interface{}) RecordType {
	switch len(data) {
	case 0:
		return read
	case 2:
		return writeBoth
	default:
		if data[0] == nil {
			return writeAfter
		}
	}

	return writeBefore
}

// IsRecordSetInEnvironment check if in the env file record activated or not
func (p *BasActivityServ) IsRecordSetInEnvironment(recordType RecordType) bool {
	switch recordType {
	case read:
		if !p.Engine.Envs.ToBool(base.RecordRead) {
			return true
		}
	default:
		if !p.Engine.Envs.ToBool(base.RecordWrite) {
			return true
		}
	}
	return false
}

// List of activities, it support pagination and search and return back count
func (p *BasActivityServ) List(params param.Param) (data map[string]interface{}, err error) {

	data = make(map[string]interface{})

	data["list"], err = p.Repo.List(params)
	glog.CheckError(err, "activities list")
	if err != nil {
		return
	}

	data["count"], err = p.Repo.Count(params)
	glog.CheckError(err, "activities count")

	return
}
