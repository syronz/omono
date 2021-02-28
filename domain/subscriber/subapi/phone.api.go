package subapi

import (
	"net/http"
	"omono/domain/base/basterm"
	"omono/domain/service"
	"omono/domain/subscriber"
	"omono/domain/subscriber/submodel"
	"omono/internal/core"
	"omono/internal/core/corterm"
	"omono/internal/response"
	"omono/internal/types"
	"omono/pkg/helper/excel"

	"github.com/gin-gonic/gin"
)

// PhoneAPI for injecting phone service
type PhoneAPI struct {
	Service service.SubPhoneServ
	Engine  *core.Engine
}

// ProvidePhoneAPI for phone is used in wire
func ProvidePhoneAPI(c service.SubPhoneServ) PhoneAPI {
	return PhoneAPI{Service: c, Engine: c.Engine}
}

// FindByID is used for fetch a phone by it's id
func (p *PhoneAPI) FindByID(c *gin.Context) {
	resp := response.New(p.Engine, c, subscriber.Domain)
	var err error
	var phone submodel.Phone
	var id uint

	if id, err = types.StrToUint(c.Param("phoneID")); err != nil {
		return
	}

	if phone, err = p.Service.FindByID(id); err != nil {
		resp.Error(err).JSON()
		return
	}

	resp.Record(subscriber.ViewPhone)
	resp.Status(http.StatusOK).
		MessageT(corterm.VInfo, basterm.Phone).
		JSON(phone)
}

// List of phones
func (p *PhoneAPI) List(c *gin.Context) {
	resp, params := response.NewParam(p.Engine, c, submodel.PhoneTable, subscriber.Domain)

	data := make(map[string]interface{})
	var err error

	if data["list"], data["count"], err = p.Service.List(params); err != nil {
		resp.Error(err).JSON()
		return
	}

	resp.Record(subscriber.ListPhone)
	resp.Status(http.StatusOK).
		MessageT(corterm.ListOfV, basterm.Phones).
		JSON(data)
}

// Create phone
func (p *PhoneAPI) Create(c *gin.Context) {
	resp := response.New(p.Engine, c, subscriber.Domain)
	var phone, createdPhone submodel.Phone
	var err error

	if err = resp.Bind(&phone, "E1053717", subscriber.Domain, basterm.Phone); err != nil {
		return
	}

	if createdPhone, err = p.Service.Create(phone); err != nil {
		resp.Error(err).JSON()
		return
	}

	resp.RecordCreate(subscriber.CreatePhone, phone)
	resp.Status(http.StatusOK).
		MessageT(corterm.VCreatedSuccessfully, basterm.Phone).
		JSON(createdPhone)
}

// Update phone
func (p *PhoneAPI) Update(c *gin.Context) {
	resp := response.New(p.Engine, c, subscriber.Domain)
	var err error

	var phone, phoneBefore, phoneUpdated submodel.Phone
	var id uint

	if id, err = types.StrToUint(c.Param("phoneID")); err != nil {
		return
	}

	if err = resp.Bind(&phone, "E1073908", subscriber.Domain, basterm.Phone); err != nil {
		return
	}

	if phoneBefore, err = p.Service.FindByID(id); err != nil {
		resp.Error(err).JSON()
		return
	}

	phone.ID = id
	if phoneUpdated, err = p.Service.Save(phone); err != nil {
		resp.Error(err).JSON()
		return
	}

	resp.Record(subscriber.UpdatePhone, phoneBefore, phone)
	resp.Status(http.StatusOK).
		MessageT(corterm.VUpdatedSuccessfully, basterm.Phone).
		JSON(phoneUpdated)
}

// Delete phone
func (p *PhoneAPI) Delete(c *gin.Context) {
	resp := response.New(p.Engine, c, subscriber.Domain)
	var err error
	var phone submodel.Phone
	var id uint

	if id, err = types.StrToUint(c.Param("phoneID")); err != nil {
		return
	}

	if phone, err = p.Service.Delete(id); err != nil {
		resp.Error(err).JSON()
		return
	}

	resp.Record(subscriber.DeletePhone, phone)
	resp.Status(http.StatusOK).
		MessageT(corterm.VDeletedSuccessfully, basterm.Phone).
		JSON()
}

// Separate phone
func (p *PhoneAPI) Separate(c *gin.Context) {
	resp := response.New(p.Engine, c, subscriber.Domain)
	var err error
	var aPhone submodel.AccountPhone
	var id uint

	if id, err = resp.GetID(c.Param("accountPhoneID"), "E1042479", basterm.Phone); err != nil {
		return
	}

	if aPhone, err = p.Service.Separate(id); err != nil {
		resp.Error(err).JSON()
		return
	}

	resp.Record(subscriber.DeletePhone, aPhone)
	resp.Status(http.StatusOK).
		MessageT(corterm.VDeletedSuccessfully, basterm.Phone).
		JSON()
}

// Excel generate excel files based on search
func (p *PhoneAPI) Excel(c *gin.Context) {
	resp, params := response.NewParam(p.Engine, c, basterm.Phones, subscriber.Domain)
	var err error

	phones, err := p.Service.Excel(params)
	if err != nil {
		resp.Error(err).JSON()
		return
	}

	ex := excel.New("phone")
	ex.AddSheet("Phones").
		AddSheet("Summary").
		Active("Phones").
		SetPageLayout("landscape", "A4").
		SetPageMargins(0.2).
		SetHeaderFooter().
		SetColWidth("B", "G", 15.3).
		SetColWidth("H", "H", 40).
		Active("Summary").
		SetColWidth("A", "D", 20).
		Active("Phones").
		WriteHeader("ID", "Name", "Code", "Type", "Status", "Updated At").
		SetSheetFields("ID", "Name", "Code", "Type", "Status", "UpdatedAt").
		WriteData(phones).
		AddTable()

	buffer, downloadName, err := ex.Generate()
	if err != nil {
		resp.Error(err).JSON()
		return
	}

	resp.Record(subscriber.ExcelPhone)

	c.Header("Content-Description", "File Transfer")
	c.Header("Content-Disposition", "attachment; filename="+downloadName)
	c.Data(http.StatusOK, "application/octet-stream", buffer.Bytes())

}
