package handler

import (
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

func (h *userHandler) Newuser(c *gin.Context) {
	c.HTML(http.StatusOK, "user_new.html", nil)
}
