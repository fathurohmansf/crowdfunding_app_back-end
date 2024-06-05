package campaign

// membuat struct
type CampaignFormatter struct {
	ID               int    `json:"id"`
	UserID           int    `json: "user_id"`
	Name             string `json:"name"`
	ShortDescription string `json: "short_description"`
	ImageURL         string `json: "image_url"`
	GoalAmount       int    `json: "goal_amout"`
	CurrentAmount    int    `json: "current_amount"`
}

// fungsi untuk format struct campaign di entity.go supaya jadi struct CampaignFormatter
func FormatCampaign(campaign Campaign) CampaignFormatter {
	CampaignFormatter := CampaignFormatter{}
	CampaignFormatter.ID = campaign.ID
	CampaignFormatter.UserID = campaign.userID
	CampaignFormatter.Name = campaign.Name
	CampaignFormatter.ShortDescription = campaign.ShortDescription
	CampaignFormatter.GoalAmount = campaign.GoalAmount
	CampaignFormatter.CurrentAmount = campaign.CurrentAmount
	CampaignFormatter.ImageURL = ""

	// untuk pengecek an bahwa punya gambar
	if len(campaign.CampaignImages) > 0 {
		CampaignFormatter.ImageURL = campaign.CampaignImages[0].FileName
	}
	return CampaignFormatter
}
