package handler

import (
	"crowdfunding/campaign"
	"crowdfunding/user"
	"net/http"

	"github.com/gin-gonic/gin"
)

type campaignhandler struct {
	campaignService campaign.Service
	userService     user.Service
}

// fungsi service nya
func NewCampaignHandler(campaignService campaign.Service, userService user.Service) *campaignhandler {
	return &campaignhandler{campaignService, userService}
}

// funsi ambil all data campaign
func (h *campaignhandler) Index(c *gin.Context) {
	campaigns, err := h.campaignService.GetCampaigns(0)
	if err != nil {
		c.HTML(http.StatusInternalServerError, "error.html", nil)
		return
	}

	// jika sukses maka tampilkan ke halaman html
	// dan parsing data gin.H = .campaigns
	c.HTML(http.StatusOK, "campaign_index.html", gin.H{"campaigns": campaigns})
}

// fungsi untuk new campaign dan load html
func (h *campaignhandler) New(c *gin.Context) {
	users, err := h.userService.GetAllUsers()
	if err != nil {
		c.HTML(http.StatusInternalServerError, "error.html", nil)
		return
	}

	input := campaign.FormCreateCampaignInput{}
	input.Users = users
	c.HTML(http.StatusOK, "campaign_new.html", input)
}

// fungsi Submit / Create  campaign di CMS
func (h *campaignhandler) Create(c *gin.Context) {
	var input campaign.FormCreateCampaignInput

	err := c.ShouldBind(&input)
	if err != nil {
		users, e := h.userService.GetAllUsers()
		// users, e := h.userService.GetAllUsers()
		if e != nil {
			c.HTML(http.StatusInternalServerError, "error.html", nil)
			return
		}
		input.Users = users
		input.Error = err

		c.HTML(http.StatusOK, "campaign_new.html", input)
		return
	}

	// dapatkan data user dulu
	user, err := h.userService.GetUserByID(input.UserID)
	if err != nil {
		c.HTML(http.StatusInternalServerError, "error.html", nil)
		return
	}
	// parsing data createCampaignInput
	createCampaignInput := campaign.CreateCampaignInput{}
	createCampaignInput.Name = input.Name
	createCampaignInput.ShortDescription = input.ShortDescription
	createCampaignInput.Description = input.Description
	createCampaignInput.GoalAmount = input.GoalAmount
	createCampaignInput.Perks = input.Perks
	createCampaignInput.User = user

	_, err = h.campaignService.CreateCampaign(createCampaignInput)
	if err != nil {
		c.HTML(http.StatusInternalServerError, "error.html", nil)
		return
	}
	c.Redirect(http.StatusFound, "/campaigns")
}
