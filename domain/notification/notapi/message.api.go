package notapi

import (
	"net/http"
	"omono/domain/notification"
	"omono/domain/notification/notmodel"
	"omono/domain/service"
	"omono/internal/core"
	"omono/internal/core/corterm"
	"omono/internal/response"
	"omono/internal/types"
	"omono/pkg/excel"
	"strconv"

	"github.com/gin-gonic/gin"
)

// MessageAPI for injecting message service
type MessageAPI struct {
	Service service.NotMessageServ
	Engine  *core.Engine
}

// ProvideMessageAPI for message is used in wire
func ProvideMessageAPI(c service.NotMessageServ) MessageAPI {
	return MessageAPI{Service: c, Engine: c.Engine}
}

// FindByID is used for fetch a message by it's id
func (p *MessageAPI) FindByID(c *gin.Context) {
	resp := response.New(p.Engine, c, notification.Domain)
	var err error
	var message notmodel.Message
	var fix types.FixedCol

	if fix, err = resp.GetFixedCol(c.Param("messageID"), "E8279820", corterm.Message); err != nil {
		return
	}

	if !resp.CheckRange(fix.CompanyID) {
		return
	}

	if message, err = p.Service.FindByID(fix); err != nil {
		resp.Error(err).JSON()
		return
	}

	resp.Record(notification.ViewMessage)
	resp.Status(http.StatusOK).
		MessageT(corterm.VInfo, corterm.Message).
		JSON(message)
}

// ViewByHash redirect it to the real route while update the status
func (p *MessageAPI) ViewByHash(c *gin.Context) {
	resp := response.New(p.Engine, c, notification.Domain)

	hash, err := strconv.ParseUint(c.Param("hash"), 10, 64)
	if err != nil {
		resp.Error(err).JSON()
		return
	}

	var message notmodel.Message
	if message, err = p.Service.FindByHash(hash); err != nil {
		resp.Error(err).JSON()
		return
	}

	// content := fmt.Sprintf(`<!DOCTYPE html>
	// <html>
	// <head>
	// <title>HTML Meta Tag</title>
	// </head>
	// <body>
	// <p>%v</p>
	// </body>
	// </html>`, message.URI)

	// t := template.Must(template.New("redirect").Parse(content))
	// t.Execute(c.Writer, "")

	// _ = content

	c.Redirect(http.StatusTemporaryRedirect, message.URI)
	c.Abort()
}

// List of messages
func (p *MessageAPI) List(c *gin.Context) {
	resp, params := response.NewParam(p.Engine, c, notmodel.MessageTable, notification.Domain)

	data := make(map[string]interface{})
	var err error

	if params.CompanyID, err = resp.GetCompanyID("E8240730"); err != nil {
		return
	}

	if !resp.CheckRange(params.CompanyID) {
		return
	}

	scope := c.Query("scope")

	if data["list"], data["count"], err = p.Service.List(params, scope); err != nil {
		resp.Error(err).JSON()
		return
	}

	resp.Record(notification.ListMessage)
	resp.Status(http.StatusOK).
		MessageT(corterm.ListOfV, corterm.Messages).
		JSON(data)
}

// Create message
func (p *MessageAPI) Create(c *gin.Context) {
	resp, params := response.NewParam(p.Engine, c, notmodel.MessageTable, notification.Domain)
	var message, createdMessage notmodel.Message
	var err error

	if message.CompanyID, message.NodeID, err = resp.GetCompanyNode("E8259678", notification.Domain); err != nil {
		resp.Error(err).JSON()
		return
	}

	if message.CompanyID, err = resp.GetCompanyID("E8234341"); err != nil {
		return
	}

	if !resp.CheckRange(message.CompanyID) {
		return
	}

	if err = resp.Bind(&message, "E8220672", notification.Domain, corterm.Message); err != nil {
		return
	}

	message.CreatedBy = params.UserID.ToPointer()
	// message.Status =
	if createdMessage, err = p.Service.Create(message); err != nil {
		resp.Error(err).JSON()
		return
	}

	resp.RecordCreate(notification.CreateMessage, message)
	resp.Status(http.StatusOK).
		MessageT(corterm.VCreatedSuccessfully, corterm.Message).
		JSON(createdMessage)
}

// Update message
func (p *MessageAPI) Update(c *gin.Context) {
	resp := response.New(p.Engine, c, notification.Domain)
	var err error

	var message, messageBefore, messageUpdated notmodel.Message
	var fix types.FixedCol

	if fix, err = resp.GetFixedCol(c.Param("messageID"), "E8231892", corterm.Message); err != nil {
		return
	}

	if !resp.CheckRange(fix.CompanyID) {
		return
	}

	if err = resp.Bind(&message, "E8250051", notification.Domain, corterm.Message); err != nil {
		return
	}

	if messageBefore, err = p.Service.FindByID(fix); err != nil {
		resp.Error(err).JSON()
		return
	}

	message.ID = fix.ID
	message.CompanyID = fix.CompanyID
	message.NodeID = fix.NodeID
	message.CreatedAt = messageBefore.CreatedAt
	if messageUpdated, err = p.Service.Save(message); err != nil {
		resp.Error(err).JSON()
		return
	}

	resp.Record(notification.UpdateMessage, messageBefore, message)
	resp.Status(http.StatusOK).
		MessageT(corterm.VUpdatedSuccessfully, corterm.Message).
		JSON(messageUpdated)
}

// Delete message
func (p *MessageAPI) Delete(c *gin.Context) {
	resp := response.New(p.Engine, c, notification.Domain)
	var err error
	var message notmodel.Message
	var fix types.FixedCol

	if fix, err = resp.GetFixedCol(c.Param("messageID"), "E8247907", corterm.Message); err != nil {
		return
	}

	if message, err = p.Service.Delete(fix); err != nil {
		resp.Error(err).JSON()
		return
	}

	resp.Record(notification.DeleteMessage, message)
	resp.Status(http.StatusOK).
		MessageT(corterm.VDeletedSuccessfully, corterm.Message).
		JSON()
}

// Excel generate excel files eaced on search
func (p *MessageAPI) Excel(c *gin.Context) {
	resp, params := response.NewParam(p.Engine, c, corterm.Messages, notification.Domain)
	var err error

	if params.CompanyID, err = resp.GetCompanyID("E8210400"); err != nil {
		return
	}

	if !resp.CheckRange(params.CompanyID) {
		return
	}

	messages, err := p.Service.Excel(params)
	if err != nil {
		resp.Error(err).JSON()
		return
	}

	ex := excel.New("message")
	ex.AddSheet("Messages").
		AddSheet("Summary").
		Active("Messages").
		SetPageLayout("landscape", "A4").
		SetPageMargins(0.2).
		SetHeaderFooter().
		SetColWidth("B", "E", 15.3).
		SetColWidth("F", "F", 40).
		Active("Summary").
		SetColWidth("A", "D", 20).
		Active("Messages").
		WriteHeader("ID", "Company ID", "Node ID", "Name", "Description", "Updated At").
		SetSheetFields("ID", "CompanyID", "NodeID", "Name", "ExDescription", "UpdatedAt").
		WriteData(messages).
		AddTable()

	buffer, downloadName, err := ex.Generate()
	if err != nil {
		resp.Error(err).JSON()
		return
	}

	resp.Record(notification.ExcelMessage)

	c.Header("Content-Description", "File Transfer")
	c.Header("Content-Disposition", "attachment; filename="+downloadName)
	c.Data(http.StatusOK, "application/octet-stream", buffer.Bytes())

}
