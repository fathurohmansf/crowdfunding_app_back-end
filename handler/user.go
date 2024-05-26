package handler

import (
	"crowdfunding/helper"
	"crowdfunding/user"
	"net/http"

	"github.com/gin-gonic/gin"
)

type userHandler struct {
	userService user.Service
}

func NewUserHandler(userService user.Service) *userHandler {
	return &userHandler{userService}
}

func (h *userHandler) RegisterUser(c *gin.Context) {
	// tangkap input dari user
	// map input dari user ke struct RegisterUserInput
	// struct di atas kita passing sebagai paramter service

	var input user.RegisterUserInput

	err := c.ShouldBindJSON(&input)
	if err != nil {
		// buat format error menjadi array of string
		// var errors []string

		// for _, e := range err.(validator.ValidationErrors) {
		// 	errors = append(errors, e.Error())
		// }

		// ambil format error dari FormatError
		errors := helper.FormatValidationError(err)

		errorMessage := gin.H{"errors": errors}

		response := helper.APIResponse("Register account failed", http.StatusUnprocessableEntity, "error", errorMessage)
		c.JSON(http.StatusUnprocessableEntity, response)
		return
		// ketika ada error kita return supaya ga eksekusi di bawah nya
	}

	newUser, err := h.userService.RegisterUser(input)

	if err != nil {
		response := helper.APIResponse("Register account failed", http.StatusBadRequest, "error", nil)
		c.JSON(http.StatusBadRequest, response)
		return
	}

	// token, err := h.jwtService.GenerateToken()
	// token itu akan di buat dulu service nya
	formatter := user.FormatUser(newUser, "tokentokentoken")

	response := helper.APIResponse("Account has been registered", http.StatusOK, "success", formatter)
	c.JSON(http.StatusOK, response)
}

func (h *userHandler) Login(c *gin.Context) {
	// user memasukkan input  berupa (email & password)
	// input ditangkap handler
	// mapping dari input user ke input struct
	// input struct kita passing ke service
	// di service mencari dgn bantuan repository user dengan email x
	// mencocokkan password

	// langkah nya kita buat dimulai dari bawa yaitu repository
	var input user.LoginInput

	if err := c.ShouldBindJSON(&input); err != nil {
		// ambil format error dari FormatError
		errors := helper.FormatValidationError(err)
		errorMessage := gin.H{"errors": errors}

		response := helper.APIResponse("Login failed", http.StatusUnprocessableEntity, "error", errorMessage)
		c.JSON(http.StatusUnprocessableEntity, response)
		return
	}

	loggedinUser, err := h.userService.Login(input)

	// jika useri id == 0 maka kesalahan itu di buat
	if err != nil {
		errorMessage := gin.H{"errors": err.Error()}

		response := helper.APIResponse("Login failed", http.StatusUnprocessableEntity, "error", errorMessage)
		c.JSON(http.StatusUnprocessableEntity, response)
		return
	}

	formatter := user.FormatUser(loggedinUser, "tokentokentoken")

	response := helper.APIResponse("Successfuly Login", http.StatusOK, "success", formatter)
	c.JSON(http.StatusOK, response)
}

// fungsi baru handle untuk Email checker
func (h *userHandler) CheckEmailAvailability(c *gin.Context) {
	// ada input email dari user
	// input email di mapping ke struct input
	// struct input di passing ke service
	// service akan memanggil repository - email sudah ada atau belum (FindByEmail di repository.go)
	// repository akan melalkukan query ke - db

	var input user.CheckEmailInput
	if err := c.ShouldBindJSON(&input); err != nil {
		// ambil format error dari FormatError
		errors := helper.FormatValidationError(err)
		errorMessage := gin.H{"errors": errors}

		response := helper.APIResponse("Email checking failed", http.StatusUnprocessableEntity, "error", errorMessage)
		c.JSON(http.StatusUnprocessableEntity, response)
		return
	}

	isEmailAvailable, err := h.userService.IsEmailAvailable(input)
	if err != nil {
		errorMessage := gin.H{"error": "Server error"}
		response := helper.APIResponse("Email checking failed", http.StatusUnprocessableEntity, "error", errorMessage)
		c.JSON(http.StatusUnprocessableEntity, response)
		return
	}

	data := gin.H{
		"is_available": isEmailAvailable,
	}

	metaMessage := "Email has been registered"

	if isEmailAvailable {
		metaMessage = "Email is availabe"
	}

	response := helper.APIResponse(metaMessage, http.StatusOK, "succes", data)
	c.JSON(http.StatusOK, response)
}
