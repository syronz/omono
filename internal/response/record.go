package response

import (
	"omono/domain/base"
	"omono/domain/base/basmodel"
	"omono/domain/base/basrepo"
	"omono/domain/service"
	"omono/internal/types"
)

// RecordCreateInstant make it simpler for calling the record
func (r *Response) RecordCreateInstant(ev types.Event, newData interface{}) {
	r.Record(ev, nil, newData)
}

// RecordInstant is used for saving activity
// TODO: deprecated
func (r *Response) RecordInstant(ev types.Event, data ...interface{}) {
	activityServ := service.ProvideBasActivityService(basrepo.ProvideActivityRepo(r.Engine))
	activityServ.Record(r.Context, ev, data...)
}

// Record will send the activity for read/update/delete to the AcitivityCh
func (r *Response) Record(ev types.Event, data ...interface{}) {
	r.initiateRecordCh(ev, data)
}

// RecordCreate will send the activity for creation to the AcitivityCh
func (r *Response) RecordCreate(ev types.Event, newData interface{}) {
	r.initiateRecordCh(ev, nil, newData)
}

func (r *Response) initiateRecordCh(ev types.Event, data ...interface{}) {
	activityServ := service.ProvideBasActivityService(basrepo.ProvideActivityRepo(r.Engine))

	var userID types.RowID
	var companyID, nodeID uint64
	var username string

	recordType := activityServ.FindRecordType(data...)
	before, after := activityServ.FillBeforeAfter(recordType, data...)

	if len(data) > 0 && !r.Engine.Envs.ToBool(base.RecordWrite) {
		return
	}

	if len(data) == 0 && !r.Engine.Envs.ToBool(base.RecordRead) {
		return
	}

	if activityServ.IsRecordSetInEnvironment(recordType) {
		return
	}
	if companyIDtmp, ok := r.Context.Get("COMPANY_ID"); ok {
		companyID = companyIDtmp.(uint64)
	}
	if nodeIDtmp, ok := r.Context.Get("NODE_ID"); ok {
		nodeID = nodeIDtmp.(uint64)
	}
	if userIDtmp, ok := r.Context.Get("USER_ID"); ok {
		userID = userIDtmp.(types.RowID)
	}
	if usernameTmp, ok := r.Context.Get("USERNAME"); ok {
		username = usernameTmp.(string)
	}

	activity := basmodel.Activity{
		FixedCol: types.FixedCol{
			CompanyID: companyID,
			NodeID:    nodeID,
		},
		Event:    ev.String(),
		UserID:   userID,
		Username: username,
		IP:       r.Context.ClientIP(),
		URI:      r.Context.Request.RequestURI,
		Before:   string(before),
		After:    string(after),
	}

	r.Engine.ActivityCh <- activity

	_ = activity
	// activityServ.RecordCh(ac

}
