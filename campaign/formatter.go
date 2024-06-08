package campaign

// membuat struct
type CampaignFormatter struct {
	ID               int    `json:"id"`
	UserID           int    `json:"user_id"`
	Name             string `json:"name"`
	ShortDescription string `json:"short_description"`
	ImageURL         string `json:"image_url"`
	GoalAmount       int    `json:"goal_amout"`
	CurrentAmount    int    `json:"current_amount"`
	Slug             string `json:"slug"`
}

// fungsi untuk format struct campaign di entity.go supaya jadi struct CampaignFormatter
func FormatCampaign(campaign Campaign) CampaignFormatter {
	CampaignFormatter := CampaignFormatter{}
	CampaignFormatter.ID = campaign.ID
	CampaignFormatter.UserID = campaign.UserID
	CampaignFormatter.Name = campaign.Name
	CampaignFormatter.ShortDescription = campaign.ShortDescription
	CampaignFormatter.GoalAmount = campaign.GoalAmount
	CampaignFormatter.CurrentAmount = campaign.CurrentAmount
	CampaignFormatter.Slug = campaign.Slug
	CampaignFormatter.ImageURL = ""

	// untuk pengecek an bahwa punya gambar
	if len(campaign.CampaignImages) > 0 {
		CampaignFormatter.ImageURL = campaign.CampaignImages[0].FileName
	}
	return CampaignFormatter
}

// fungsi untuk slice of campaign parameter nya bisa di call
func FormatCampaigns(campaigns []Campaign) []CampaignFormatter {
	// Ini di buat lebih simple null jadi 0 (lebih clean code)
	campaignsFormatter := []CampaignFormatter{}
	// pengecekan jika campaigns 0 itu null jadi 0 di json
	// if len(campaigns) == 0 {
	// 	return []CampaignFormatter{}
	// }
	// buat variable
	// var campaignsFormatter []CampaignFormatter

	// buat perulangan jika banyak campaigns
	for _, campaign := range campaigns {
		// dapatkan single object campaign dulu
		campaignFormatter := FormatCampaign(campaign)
		// baru di append menjadi object campaigns jadi banyakkkk
		campaignsFormatter = append(campaignsFormatter, campaignFormatter)
	}
	return campaignsFormatter
}
