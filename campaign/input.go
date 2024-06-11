package campaign

import "crowdfunding/user"

type GetCampaignDetailInput struct {
	// pakai uri
	ID int `uri:"id" binding:"required"`
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
