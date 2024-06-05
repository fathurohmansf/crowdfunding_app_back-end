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

	campaigns, err := h.service.GetCampaigns(userID)
	if err != nil {
		response := helper.APIResponse("Error to get campaigns", http.StatusBadRequest, "error", nil)
		c.JSON(http.StatusBadRequest, response)
		return
	}
	response := helper.APIResponse("list of campaigns", http.StatusOK, "succes", campaigns)
	c.JSON(http.StatusOK, response)
	return
}
