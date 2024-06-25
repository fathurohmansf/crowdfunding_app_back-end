package transaction

import "crowdfunding/user"

type GetCampaignTransactionInput struct {
	// pakai uri
	ID   int `uri:"id" binding:"required"`
	User user.User
}

// Membaut struct baru untuk Create Campaign API Sevice.go
type CreateCampaignInput struct {
	Name             string `json:"name" binding:"required"`
	ShortDescription string `json:"short_description" binding:"required"`
	Description      string `json:"description" binding:"required"`
	GoalAmount       int    `json:"goal_amount" binding:"required"`
	Perks            string `json:"perks" binding:"required"`
	User             user.User
}

// Membuat struct baru untuk UPDATE is_primary menjadi 0 dari 1 (untuk kebutuhan service.go)
type CreateCampaignImageInput struct {
	// pakai form karena user tidak dalam bentuk json seperti di atas
	CampaignID int `form:"campaign_id" binding:"required"`
	// IsPrimary ini ga pakai binding, karna tidak harus true yg harus di upload
	IsPrimary bool `form:"is_primary"`
	User      user.User
}
