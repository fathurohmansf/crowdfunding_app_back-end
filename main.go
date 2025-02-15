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
	"github.com/gin-contrib/sessions"
	"github.com/gin-contrib/sessions/cookie"
	"github.com/gin-gonic/gin"
	"github.com/joho/godotenv"
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
	userWebHandler := webHandler.NewUserHandler(userService)                          // untuk handler web admin CMS
	campaignWebHandler := webHandler.NewCampaignHandler(campaignService, userService) // untuk handler web admin CMS
	transactionsWebHandler := webHandler.NewTransactionHandler(transactionService)    // untuk handler web admin CMS
	sessionsWebHandler := webHandler.NewSessionHandler(userService)                   // untuk handler web admin CMS

	// ROUTE
	router := gin.Default()
	// CoRS for client
	router.Use(cors.Default())

	// session login CMS middleware
	cookieStore := cookie.NewStore([]byte(auth.SECRET_KEY))
	router.Use(sessions.Sessions("senabung", cookieStore))

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
	router.GET("/users", authAdminMiddleware(), userWebHandler.Index)
	// route untuk newUser di CMS admin
	router.GET("/users/new", authAdminMiddleware(), userWebHandler.New)
	// route untuk New user kirim data dari form ke database
	router.POST("/users", authAdminMiddleware(), userWebHandler.Create)
	// route untuk ID Param / ID User nanti buat Edit/Update User
	router.GET("/users/edit/:id", authAdminMiddleware(), userWebHandler.Edit)
	// route untuk update user di CMS
	router.POST("/users/update/:id", authAdminMiddleware(), userWebHandler.Update)
	// route untuk upload avatar di CMS
	router.GET("/users/avatar/:id", authAdminMiddleware(), userWebHandler.NewAvatar)
	// route untuk crate avatar upload di CMS
	router.POST("/users/avatar/:id", authAdminMiddleware(), userWebHandler.CreateAvatar)
	// route untuk all campaign di CMS
	router.GET("/campaigns", authAdminMiddleware(), campaignWebHandler.Index)
	// route untuk new campaign di CMS
	router.GET("/campaigns/new", authAdminMiddleware(), campaignWebHandler.New)
	// route untuk submit new campaign di CMS POST
	router.POST("/campaigns", authAdminMiddleware(), campaignWebHandler.Create)
	// route untuk upload new image campaign di CMS
	router.GET("/campaigns/image/:id", authAdminMiddleware(), campaignWebHandler.NewImage)
	// route untuk submit new image camoaign di CMS
	router.POST("/campaigns/image/:id", authAdminMiddleware(), campaignWebHandler.CreateImage)
	// route untuk edit campaign dari params di CMS
	router.GET("/campaigns/edit/:id", authAdminMiddleware(), campaignWebHandler.Edit)
	// route untuk update campaign di CMS
	router.POST("/campaigns/update/:id", authAdminMiddleware(), campaignWebHandler.Update)
	// route untuk show detail campaign di CMS
	router.GET("/campaigns/show/:id", authAdminMiddleware(), campaignWebHandler.Show)
	// route untuk show all transactions
	router.GET("/transactions", authAdminMiddleware(), transactionsWebHandler.Index)
	// route untuk login CMS
	router.GET("/login", sessionsWebHandler.New)
	// route untuk submit login CMS
	router.POST("/sessions", sessionsWebHandler.Create)

	// Find & load .env file
	err = godotenv.Load(".env")
	if err != nil {
		log.Fatal("Error loading .env file")
	}

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

// middleware for CMS
func authAdminMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		session := sessions.Default(c)
		// gunakan key userID
		userIDSession := session.Get("userID")
		if userIDSession == nil {
			c.Redirect(http.StatusFound, "/login")
			return
		}

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
