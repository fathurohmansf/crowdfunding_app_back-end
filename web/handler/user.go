package handler

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

type userHandler struct {
}

func NewUserHandler() *userHandler {
	return &userHandler{}
}

func (h *userHandler) Index(c *gin.Context) {
	c.HTML(http.StatusOK, "user_index.html", nil)
}

// func (h *userHandler) Index(c *gin.Context) {
// 	users, err := h.userService.GetAll()
// 	if err != nil {
// 		c.HTML(http.StatusInternalServerError, "error.html", nil)
// 		return
// 	}
// 	c.HTML(http.StatusOK, "user_index.html", gin.H{
// 		"users":    users,
// 		"testings": "testing string",
// 	})
// 	//Check data
// 	fmt.Printf("%v", users)
// }
