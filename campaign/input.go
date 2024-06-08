package campaign

type GetCampaignDetailInput struct {
	// pakai uri
	ID int `uri:"id" binding:"required"`
}
