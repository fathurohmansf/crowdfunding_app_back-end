package handler

import (
	"crowdfunding/helper"
	"crowdfunding/transaction"
	"net/http"

	"github.com/gin-gonic/gin"
)

// Penerapan Transaction API
// Parameter di uri
// tangkap parameter mapping ke input struct
// panggil service, input struct sebagai parameter
// service, berbekal campaign id bisa panggil repo
// repo mencari data transaction suatu campaign

type transactionHandler struct {
	service transaction.Service
}

func NewTransactionHandler(service transaction.Service) *transactionHandler {
	return &transactionHandler{service}
}

func (h *transactionHandler) GetCampaignTransaction(c *gin.Context) {
	var input transaction.GetCampaignTransactionInput

	err := c.ShouldBindUri(&input)
	if err != nil {
		response := helper.APIResponse("Failed to get campaign's transactions", http.StatusBadRequest, "error", nil)
		c.JSON(http.StatusBadRequest, response)
		return
	}
	transactions, err := h.service.GetTransactionByCampaignID(input)
	if err != nil {
		response := helper.APIResponse("Failed to get campaign's transactions", http.StatusBadRequest, "error", nil)
		c.JSON(http.StatusBadRequest, response)
		return
	}
	response := helper.APIResponse("Campaign's transactions", http.StatusOK, "success", transactions)
	c.JSON(http.StatusOK, response)
}
