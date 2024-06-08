package campaign

import (
	"crowdfunding/user"
	"time"
)

type Campaign struct {
	ID               int
	UserID           int
	Name             string
	ShortDescription string
	Description      string
	Perks            string
	BackerCount      int
	GoalAmount       int
	CurrentAmount    int
	Slug             string
	CreatedAt        time.Time
	UpdatedAt        time.Time
	CampaignImages   []CampaignImage `gorm:"foreignKey:CampaignID"`
	User             user.User       `gorm:"foreignKey:UserID"`
}

type CampaignImage struct {
	ID         int `gorm:"primaryKey"`
	CampaignID int `gorm:"column:campaign_id"`
	FileName   string
	IsPrimary  int
	CreatedAt  time.Time
	UpdatedAt  time.Time
}
