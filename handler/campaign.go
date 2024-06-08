package handler

import (
	"crowdfunding/campaign"
	"crowdfunding/helper"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
)

// Penerapan List Campaign API
// INGAT di kerjakan dari bawah dulu
// tangkap parameter di handler
// handler ke service
// buat formatter sesuai JSON yg di inginkan Front-End(formatter.go)
// service yg menentukan repository mana yg di-call (service.go)
// repository, buat dua : FindAll, FindByUserID(repository.go)
// db

// buat struct
type campaignHandler struct {
	service campaign.Service
}

// fungsi ini untuk membuat object / struct yg nanti di panggil main.go
func NewCampaignHandler(service campaign.Service) *campaignHandler {
	return &campaignHandler{service}
}

// membuat func routing = api/v1/campaigns
func (h *campaignHandler) GetCampaigns(c *gin.Context) {
	// Harus Convert to int
	userID, _ := strconv.Atoi(c.Query("user_id"))

	// BALIKAN dari campaigns ini adalah slice of []campaign
	campaigns, err := h.service.GetCampaigns(userID)
	if err != nil {
		response := helper.APIResponse("Error to get campaigns", http.StatusBadRequest, "error", nil)
		c.JSON(http.StatusBadRequest, response)
		return
	}
	response := helper.APIResponse("list of campaigns", http.StatusOK, "succes", campaign.FormatCampaigns(campaigns))
	c.JSON(http.StatusOK, response)
	return
}

// buat handler baru, abis di buat taro handler di main.go jgn lupa
func (h *campaignHandler) GetCampaign(c *gin.Context) {
	// bentuk nya : api/v1/campaign/1    ini cth nya berdasarkan id
	// handler : mapping id yg di url ke struct input = service, call formatter (campaign.go & main.go)
	// Service : inputnya struct input => untuk menangkap ID  di url, memanggil repo (service.go & input.go)
	// butuh repository : get campaign by ID (repository.go)

	var input campaign.GetCampaignDetailInput
	err := c.ShouldBindUri(&input)
	if err != nil {
		response := helper.APIResponse(" Failed to get of campaign", http.StatusBadRequest, "Error", nil)
		c.JSON(http.StatusBadRequest, response)
		return
	}
	campaignDetail, err := h.service.GetCampaignByID(input)
	if err != nil {
		response := helper.APIResponse(" Failed to get of campaign saat panggil fungsi GetCampaignByID", http.StatusBadRequest, "Error", nil)
		c.JSON(http.StatusBadRequest, response)
		return
	}

	response := helper.APIResponse("Campaign detail", http.StatusOK, "success", campaign.FormatCampaignDetail(campaignDetail))
	c.JSON(http.StatusOK, response)
}
