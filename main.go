package main

import (
	"crowdfunding/auth"
	"crowdfunding/campaign"
	"crowdfunding/handler"
	"crowdfunding/helper"
	"crowdfunding/payment"
	"crowdfunding/transaction"
	"crowdfunding/user"
	"log"
	"net/http"
	"path/filepath"
	"strings"

	webHandler "crowdfunding/web/handler"

	"github.com/dgrijalva/jwt-go"
	"github.com/gin-contrib/cors"
	"github.com/gin-contrib/multitemplate"
	"github.com/gin-gonic/gin"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
)

func main() {
	// koneksi handler
	dsn := "root:@tcp(127.0.0.1:3306)/crowdfunding?charset=utf8mb4&parseTime=True&loc=Local"
	db, err := gorm.Open(mysql.Open(dsn), &gorm.Config{})

	if err != nil {
		log.Fatal(err.Error())
	}

	// REPOSITORY
	userRepository := user.NewRepository(db)
	// tambahkan campaigns API ke db
	campaignRepository := campaign.NewRepository(db)
	transactionRepository := transaction.NewRepository(db)

	// SERVICE
	authService := auth.NewService()
	userService := user.NewService(userRepository)
	campaignService := campaign.NewService(campaignRepository)
	// Panggil payment Midtrans service
	paymentService := payment.NewService()
	// Panggil Transaction Service , panggil juga campaignrepository
	transactionService := transaction.NewService(transactionRepository, campaignRepository, paymentService)

	// HANDLER
	userHandler := handler.NewUserHandler(userService, authService) // tambahkan authService
	campaignHandler := handler.NewCampaignHandler(campaignService)  // tambahkan campaigns
	transactionHandler := handler.NewTransactionHandler(transactionService)
	userWebHandler := webHandler.NewUserHandler(userService) // untuk handler web admin

	// ROUTE
	router := gin.Default()
	// CoRS for client
	router.Use(cors.Default())
	// add layoutRendering in untuk web admin
	router.HTMLRender = loadTemplates("./web/templates")
	// Static Image
	router.Static("/images", "./images")
	// Load Images CMS
	router.Static("/css", "./web/assets/css")
	router.Static("/js", "./web/assets/js")
	router.Static("/webfonts", "./web/assets/webfonts")
	api := router.Group("/api/v1")

	api.POST("/users", userHandler.RegisterUser)
	api.POST("/sessions", userHandler.Login)
	api.POST("/email_checkers", userHandler.CheckEmailAvailability)
	api.POST("/avatars", authMiddleware(authService, userService), userHandler.UploadAvatar) //Middleware karna mau upload itu harus login user dulu ga sembarangan upload
	// Fetch User API
	api.GET("/users/fetch", authMiddleware(authService, userService), userHandler.FetchUser)

	// Ambil data campaigns get dari server
	api.GET("/campaigns", campaignHandler.GetCampaigns)
	api.GET("/campaigns/:id", campaignHandler.GetCampaign)
	// route untuk CreateCampaign API
	api.POST("/campaigns", authMiddleware(authService, userService), campaignHandler.CreateCampaign) //Middleware karna mau upload itu harus login user dulu ga sembarangan upload
	// route untuk UpdateCampaign API
	api.PUT("/campaigns/:id", authMiddleware(authService, userService), campaignHandler.UpdateCampaign) //Middleware karna mau upload itu harus login user dulu ga sembarangan upload
	// route untuk UPLOAD Campaign Image API
	api.POST("/campaign-images", authMiddleware(authService, userService), campaignHandler.UploadImage) //Middleware karna mau upload itu harus login user dulu ga sembarangan upload
	// route untuk Campaign Transaction API
	api.GET("/campaigns/:id/transactions", authMiddleware(authService, userService), transactionHandler.GetCampaignTransaction)
	// route untuk User Transaction API
	api.GET("/transactions", authMiddleware(authService, userService), transactionHandler.GetUserTranactions)
	// route untuk Create Transaction Midtrans
	api.POST("/transactions", authMiddleware(authService, userService), transactionHandler.CreateTransaction)
	// route untuk Transaction Notification Midtrans
	api.POST("/transactions/notification", transactionHandler.GetNotification)

	// Route for CMS admin
	// jika ada request ke /users maka akan di arahkan ke webhandler.index
	router.GET("/users", userWebHandler.Index)
	// route untuk newUser di CMS admin
	router.GET("/users/new", userWebHandler.New)

	router.Run()
}

// Kita kerjakan Middleware(1)
func authMiddleware(authService auth.Service, userService user.Service) gin.HandlerFunc {

	// Kita kerjakan Middleware (1)
	return func(c *gin.Context) {
		authHeader := c.GetHeader("Authorization")
		if !strings.Contains(authHeader, "Bearer") {
			response := helper.APIResponse(" Unautorized", http.StatusUnauthorized, "error", nil)
			c.AbortWithStatusJSON(http.StatusUnauthorized, response)
			return
		}

		tokenString := ""
		arrayToken := strings.Split(authHeader, " ")
		if len(arrayToken) == 2 {
			tokenString = arrayToken[1]
		}

		token, err := authService.ValidateToken(tokenString)
		if err != nil {
			response := helper.APIResponse(" Unautorized", http.StatusUnauthorized, "error", nil)
			c.AbortWithStatusJSON(http.StatusUnauthorized, response)
			return
		}

		claim, ok := token.Claims.(jwt.MapClaims)

		if !ok || !token.Valid {
			response := helper.APIResponse(" Unautorized", http.StatusUnauthorized, "error", nil)
			c.AbortWithStatusJSON(http.StatusUnauthorized, response)
			return
		}
		userID := int(claim["user_id"].(float64))

		user, err := userService.GetUserByID(userID)
		if err != nil {
			response := helper.APIResponse(" Unautorized", http.StatusUnauthorized, "error", nil)
			c.AbortWithStatusJSON(http.StatusUnauthorized, response)
			return
		}

		c.Set("currentUser", user)
	}
}

func loadTemplates(templatesDir string) multitemplate.Renderer {
	r := multitemplate.NewRenderer()

	layouts, err := filepath.Glob(templatesDir + "/layouts/*")
	if err != nil {
		panic(err.Error())
	}

	includes, err := filepath.Glob(templatesDir + "/**/*")
	if err != nil {
		panic(err.Error())
	}

	// Generate our templates map from our layouts/ and includes/ directories
	for _, include := range includes {
		layoutCopy := make([]string, len(layouts))
		copy(layoutCopy, layouts)
		files := append(layoutCopy, include)
		r.AddFromFiles(filepath.Base(include), files...)
	}
	return r
}
