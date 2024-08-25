package handler

import (
	"crowdfunding/user"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
)

type userHandler struct {
	userService user.Service
}

func NewUserHandler(userService user.Service) *userHandler {
	return &userHandler{userService}
}

func (h *userHandler) Index(c *gin.Context) {
	users, err := h.userService.GetAllUsers()
	if err != nil {
		c.HTML(http.StatusInternalServerError, "error.html", nil)
		return
	}

	// Test page Error handling
	// users, _ := h.userService.GetAllUsers()
	// if true {
	// 	//later untuk handling error
	// 	c.HTML(http.StatusInternalServerError, "error.html", nil)
	// 	return
	// }
	c.HTML(http.StatusOK, "user_index.html", gin.H{"users": users}) //pakai map gin.H untuk bisa akses ke var users, dgn key(untuk template) & value(dari service)
}

// func untuk new User
func (h *userHandler) New(c *gin.Context) {
	c.HTML(http.StatusOK, "user_new.html", nil)
}

// func untuk data form user kirim ke database
func (h *userHandler) Create(c *gin.Context) {
	var input user.FormCreateUserInput

	err := c.ShouldBind(&input)
	if err != nil {
		// skip
		input.Error = err
		c.HTML(http.StatusOK, "user_new.html", input)
		return
	}

	RegisterInput := user.RegisterUserInput{}
	RegisterInput.Name = input.Name
	RegisterInput.Email = input.Email
	RegisterInput.Occupation = input.Occupation
	RegisterInput.Password = input.Password

	_, err = h.userService.RegisterUser(RegisterInput)
	if err != nil {
		//skip
		c.HTML(http.StatusInternalServerError, "error.html", nil)
		return
	}
	c.Redirect(http.StatusFound, ("/users"))
}

// Func untuk Edit/Update User di halaman user_edit.html
func (h *userHandler) Edit(c *gin.Context) {
	// Tangkap ID user nya dulu menggunakan IdParam
	idParam := c.Param("id")
	// gunakan servis GetUserByID di user/service.go dan convert ID int to string
	id, _ := strconv.Atoi(idParam)

	registeredUser, err := h.userService.GetUserByID(id)
	if err != nil {
		c.HTML(http.StatusInternalServerError, "error.html", nil)
	}
	c.HTML(http.StatusOK, "user_edit.html", registeredUser) // Mem passing nilai/value yang akan di tampilkan

}

// func Update User
func (h *userHandler) Update(c *gin.Context) {
	// Tangkap ID user nya dulu menggunakan IdParam
	idParam := c.Param("id")
	// gunakan servis GetUserByID di user/service.go dan convert ID int to string
	id, _ := strconv.Atoi(idParam)

	var input user.FormUpdateUserInput

	err := c.ShouldBind(&input)
	if err != nil {
		//skip
	}
	// Bind manual ID nya
	input.ID = id
	_, err = h.userService.UpdateUser(input)
	if err != nil {
		c.HTML(http.StatusInternalServerError, "error.html", nil)
		return
	}
	c.Redirect(http.StatusFound, "/users")
}
