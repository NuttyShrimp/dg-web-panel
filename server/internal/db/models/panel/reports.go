package panel_models

import (
	"degrens/panel/models"
	"errors"

	"gorm.io/gorm"
)

type MsgType string

const (
	Text  MsgType = "text"
	Image         = "image"
)

type Report struct {
	BaseModel
	Title   string      `json:"title"`
	Creator string      `json:"creator"`
	Open    bool        `json:"open"`
	Tags    []ReportTag `gorm:"many2many:report_tags_link" json:"tags"`
	// Registered members, other than staff
	Members  []ReportMember  `json:"members"`
	Messages []ReportMessage `json:"messages"`
}

type ReportMessage struct {
	BaseModel
	Message  string              `json:"message"`
	Type     MsgType             `json:"type"`
	ReportID uint                `json:"-"`
	UserID   *uint               `json:"-"`
	User     User                `json:"-"`
	MemberID *uint               `json:"-"`
	Member   ReportMember        `json:"-"`
	Sender   ReportMessageSender `gorm:"-" json:"sender"`
}

type ReportMessageSender struct {
	models.UserInfo
	SteamId string `json:"steamId"`
}

type ReportMember struct {
	ID       uint   `gorm:"primaryKey" json:"-"`
	Name     string `json:"name"`
	SteamID  string `json:"steamId"`
	ReportID uint   `json:"-"`
}

type ReportTag struct {
	Name    string   `gorm:"primaryKey" json:"name"`
	Color   string   `json:"color"`
	Reports []Report `gorm:"many2many:report_tags_link" json:"-"`
}

func (rm *ReportMessage) BeforeCreate(tx *gorm.DB) error {
	if rm.UserID == nil && rm.MemberID == nil {
		return errors.New("message should be assigned to a user or a report member")
	}
	if rm.UserID != nil && rm.MemberID != nil {
		return errors.New("message cannot be assigned to a user and a report member")
	}
	return nil
}
