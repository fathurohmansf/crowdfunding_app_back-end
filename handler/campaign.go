package handler

import (
	"crowdfunding/campaign"
	"crowdfunding/helper"
	"crowdfunding/user"
	"fmt"
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

// GETALL Campaign buat handler baru, abis di buat taro handler di main.go jgn lupa
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

// CREATE Campaign API untuk Buat handler baru
// Penerapan Create API
// Tangkap parameter dari user ke input struct
// Ambil current user dari jwt/handler
// nge test manual input repo dan service (main.go)
// Panggil service, parameternya input struct (dan juga buat slug) (service.go dan input.go)
// panggil repository untuk simpan data campaign baru (repository.go)

// Membuat fungsi Create Campaign API
func (h *campaignHandler) CreateCampaign(c *gin.Context) {
	// buat var untuk menangkap parameter yg di input oleh user
	var input campaign.CreateCampaignInput

	err := c.ShouldBindJSON(&input)
	if err != nil {
		// ambil format error dari FormatError
		errors := helper.FormatValidationError(err)
		errorMessage := gin.H{"errors": errors}

		response := helper.APIResponse("Failed to create campaign", http.StatusUnprocessableEntity, "error", errorMessage)
		c.JSON(http.StatusUnprocessableEntity, response)
		return
	}

	// ambil data user dari JWT
	currentUser := c.MustGet("currentUser").(user.User)
	input.User = currentUser
	newCapaign, err := h.service.CreateCampaign(input) // panggil service
	if err != nil {
		response := helper.APIResponse("Failed to create campaign", http.StatusBadRequest, "error", nil)
		c.JSON(http.StatusUnprocessableEntity, response)
		return
	}
	response := helper.APIResponse("Success to create campaign", http.StatusOK, "success", campaign.FormatCampaign(newCapaign))
	c.JSON(http.StatusOK, response)
}

// UPDATE Campaign API handler
// handler ( yg ada 2 tadi diterapkan )
// mapping dari input ke input struct (ada 2 , 1.dari uri byID(api/v1/camapaign/1) 2.dari form input user)
// input dari user, dan juga input yg ada di uri (passing ke service)
// service ( untuk buat logic gimana update ) (find campaign byID uri, lalu tangkap parameter yg ada di input form)
// repository update data campaign

// buat fungsi handler di sini UpdateCampaign API
func (h *campaignHandler) UpdateCampaign(c *gin.Context) {
	// SATU dari Uri ByID, ambil code dari handler GetCampaign
	var inputID campaign.GetCampaignDetailInput
	err := c.ShouldBindUri(&inputID)
	if err != nil {
		response := helper.APIResponse(" Failed to update campaign", http.StatusBadRequest, "Error", nil)
		c.JSON(http.StatusBadRequest, response)
		return
	}
	// DUA dari form input user, ambil code dari handler CreateCampaign
	var inputData campaign.CreateCampaignInput

	err = c.ShouldBindJSON(&inputData)
	if err != nil {
		// ambil format error dari FormatError
		errors := helper.FormatValidationError(err)
		errorMessage := gin.H{"errors": errors}

		response := helper.APIResponse("Failed to update campaign", http.StatusUnprocessableEntity, "error", errorMessage)
		c.JSON(http.StatusUnprocessableEntity, response)
		return
	}
	// ambil data user dari JWT, agar user lain tidak bisa ubah campaign punya pribadi
	currentUser := c.MustGet("currentUser").(user.User)
	inputData.User = currentUser

	// TIGA fungsi jika error update dan update campaign baru
	updatedCampaign, err := h.service.UpdateCampaign(inputID, inputData)
	if err != nil {
		response := helper.APIResponse("Failed to update campaign", http.StatusUnprocessableEntity, "error", nil)
		c.JSON(http.StatusUnprocessableEntity, response)
		return
	}
	response := helper.APIResponse("Success to update campaign", http.StatusOK, "success", campaign.FormatCampaign(updatedCampaign))
	c.JSON(http.StatusOK, response)
}

// UPLOAD Campaign Image API
// hanlder :
// 1. tangkap input dan ubah struct input
// 2. save image campaign ke suatu folder
// Service (kondisi manggil point 2 di repo, panggil repo point 1)
// repository : (repository.go)
// 1. create image/save data image ke dalam tabel campaign_images
// 2. ubah is_primary true ke false, true yg sebelum nya akan jadi false, jadi inti nya hanya 1 gambar aja yang is_primary = true

// implemetasi handler UPLOAD Campaign Image API
func (h *campaignHandler) UploadImage(c *gin.Context) {
	var input campaign.CreateCampaignImageInput
	// karna pakai form jadi ShouldBind aja, bukan ShouldBindJSON
	err := c.ShouldBind(&input)
	if err != nil {
		errors := helper.FormatValidationError(err)
		errorMessage := gin.H{"errors": errors}

		response := helper.APIResponse("Failed to create campaign image", http.StatusUnprocessableEntity, "error", errorMessage)
		c.JSON(http.StatusUnprocessableEntity, response)
		return
	}
	// ambil data user dari JWT, agar user lain tidak bisa ubah campaign punya pribadi
	currentUser := c.MustGet("currentUser").(user.User)
	input.User = currentUser
	userID := currentUser.ID

	file, err := c.FormFile("file")
	if err != nil {
		data := gin.H{"is_uploaded": false}
		response := helper.APIResponse("Failed to upload campaign image", http.StatusBadRequest, "error", data)
		c.JSON(http.StatusBadRequest, response)
		return
	}
	// images/1-namafile.png (ini yang BARU karna ada ID di depan 1)
	path := fmt.Sprintf("images/%d-%s", userID, file.Filename)
	err = c.SaveUploadedFile(file, path)
	if err != nil {
		data := gin.H{"is_uploaded": false}
		response := helper.APIResponse("Failed to upload Campaign Image", http.StatusBadRequest, "error", data)
		c.JSON(http.StatusBadRequest, response)
		return
	}
	_, err = h.service.SaveCampaignImage(input, path)
	if err != nil {
		data := gin.H{"is_uploaded": false}
		response := helper.APIResponse("Failed to upload campaign image", http.StatusBadRequest, "error", data)
		c.JSON(http.StatusBadRequest, response)
		return
	}
	data := gin.H{"is_uploaded": true}
	response := helper.APIResponse("Campaign image successfuly uploaded", http.StatusOK, "success", data)
	c.JSON(http.StatusOK, response)
}
