package transaction

import (
	"crowdfunding/campaign"
	"crowdfunding/user"
	"time"
)

type Transaction struct {
	ID         int
	CampaignID int
	UserID     int
	Amount     int
	Status     string
	Code       string
	PaymentURL string
	User       user.User
	Campaign   campaign.Campaign //untuk relasi ke campaign images
	CreatedAt  time.Time
	UpdatedAt  time.Time
}
