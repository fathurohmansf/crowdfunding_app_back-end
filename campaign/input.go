package campaign

import "crowdfunding/user"

type GetCampaignDetailInput struct {
	// pakai uri
	ID int `uri:"id" binding:"required"`
}

// Membaut struct baru untuk Create Campaign API Sevice.go
type CreateCampaignInput struct {
	Name             string `json:"name"`
	ShortDescription string `json:"short_description"`
	Description      string `json:"description"`
	GoalAmount       int    `json:"goal_amount"`
	Perks            string `json:"perks"`
	User             user.User
}
