package notmodel

import (
	"omono/internal/core/coract"
	"omono/internal/types"
	"time"

	"gorm.io/gorm"
)

// MessageTable is used inside the repo layer
const (
	MessageTable = "not_messages"
)

// Message model
type Message struct {
	gorm.Model  `gorm:"embedded"`
	CreatedBy   *uint      `json:"created_by"`
	RecipientID uint       `json:"recipient_id"`
	Hash        uint64     `gorm:"not null;unique;type:varchar(50)" json:"hash,omitempty"`
	Title       string     `gorm:"type:varchar(200)" json:"title,omitempty"`
	Message     string     `gorm:"not null" json:"message,omitempty"`
	URI         string     `json:"uri"`
	Part        string     `gorm:"type:varchar(50)" json:"part"`
	Status      types.Enum `gorm:"not null;default:'new';type:enum('new','seen')" json:"status"`
	ViewCount   byte       `json:"view_count"`
	ViewedAt    *time.Time `json:"viewed_at"`
}

// Validate check the type of
func (p *Message) Validate(act coract.Action) (err error) {
	return err
}
