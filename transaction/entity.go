package transaction

import (
	"crowdfunding/campaign"
	"crowdfunding/user"
	"time"

	"github.com/leekchan/accounting"
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

// buat func untuk formatting GoalAmount CMS
func (t Transaction) AmountFormatIDR() string {
	ac := accounting.Accounting{Symbol: "Rp. ", Precision: 2, Thousand: ".", Decimal: ","}
	return ac.FormatMoney(t.Amount)
}
