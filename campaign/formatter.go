package campaign

import "strings"

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
	campaignFormatter := CampaignFormatter{}
	campaignFormatter.ID = campaign.ID
	campaignFormatter.UserID = campaign.UserID
	campaignFormatter.Name = campaign.Name
	campaignFormatter.ShortDescription = campaign.ShortDescription
	campaignFormatter.GoalAmount = campaign.GoalAmount
	campaignFormatter.CurrentAmount = campaign.CurrentAmount
	campaignFormatter.Slug = campaign.Slug
	campaignFormatter.ImageURL = ""

	// untuk pengecek an bahwa punya gambar
	if len(campaign.CampaignImages) > 0 {
		campaignFormatter.ImageURL = campaign.CampaignImages[0].FileName
	}
	return campaignFormatter
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

// Buat Struct untuk CampaignDetailByID
// localhost:8080/api/v1/campaigns/1
type CampaignDetailFormatter struct {
	ID               int      `json:"id"`
	Name             string   `json:"name"`
	ShortDescription string   `json:"short_description"`
	Description      string   `json:"description"`
	ImageURL         string   `json:"image_url"`
	GoalAmount       int      `json:"goal_amout"`
	CurrentAmount    int      `json:"current_amout"`
	UserID           int      `json:"user_id"`
	Slug             string   `json:"slug"`
	Perks            []string `json:"perks"`
	// Membuat struct User di dalam struct CampaignDetailFormatter
	User CampaignUserFormatter `json:"user"`
	// Membuat struct images di dlm struct CampaignDetailFormatter
	Images []CampaignImageFormatter `json:"images"`
}
type CampaignUserFormatter struct {
	Name     string `json:"name"`
	ImageURL string `json:"image_url"`
}
type CampaignImageFormatter struct {
	ImageURL  string `json:"image_url"`
	IsPrimary bool   `json:"is_primary"`
}

// Fungsi Format campaigndetailByID
func FormatCampaignDetail(campaign Campaign) CampaignDetailFormatter {
	campaignDetailFormatter := CampaignDetailFormatter{}
	campaignDetailFormatter.ID = campaign.ID
	campaignDetailFormatter.Name = campaign.Name
	campaignDetailFormatter.ShortDescription = campaign.ShortDescription
	campaignDetailFormatter.Description = campaign.Description
	campaignDetailFormatter.GoalAmount = campaign.GoalAmount
	campaignDetailFormatter.CurrentAmount = campaign.CurrentAmount
	campaignDetailFormatter.Slug = campaign.Slug
	campaignDetailFormatter.ImageURL = ""

	if len(campaign.CampaignImages) > 0 {
		campaignDetailFormatter.ImageURL = campaign.CampaignImages[0].FileName
	}

	var perks []string
	// SPLIT ini fungsi nya untuk dapet perulangan string perks satu, perks dua
	for _, perk := range strings.Split(campaign.Perks, ",") {
		perks = append(perks, strings.TrimSpace(perk))
	}

	campaignDetailFormatter.Perks = perks

	// cara panggil user di formatter data nya bakal muncul
	user := campaign.User
	campaignUserFormatter := CampaignUserFormatter{}
	campaignUserFormatter.Name = user.Name
	campaignUserFormatter.ImageURL = user.AvatarFileName

	campaignDetailFormatter.User = campaignUserFormatter

	// Cara panggil Images di formatter data nya bakal muncul
	images := []CampaignImageFormatter{} // karna bentuk nya slice [] harus di append

	for _, image := range campaign.CampaignImages {
		campaignImageFormatter := CampaignImageFormatter{}
		campaignImageFormatter.ImageURL = image.FileName

		isPrimary := false

		if image.IsPrimary == 1 {
			isPrimary = true
		}

		campaignImageFormatter.IsPrimary = isPrimary
		// ini di append nya
		images = append(images, campaignImageFormatter)
	}
	campaignDetailFormatter.Images = images

	return campaignDetailFormatter
}
