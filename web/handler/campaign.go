package handler

import (
	"crowdfunding/campaign"
	"net/http"

	"github.com/gin-gonic/gin"
)

type campaignhandler struct {
	campaignService campaign.Service
}

// fungsi service nya
func NewCampaignHandler(campaignService campaign.Service) *campaignhandler {
	return &campaignhandler{campaignService}
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
