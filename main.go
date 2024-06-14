package main

import (
	"crowdfunding/auth"
	"crowdfunding/campaign"
	"crowdfunding/handler"
	"crowdfunding/helper"
	"crowdfunding/user"
	"log"
	"net/http"
	"strings"

	"github.com/dgrijalva/jwt-go"
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

		// 	}
		// 	fmt.Println("Connection to database is good")

		// 	// membuat var users tipe array dengan membaca data struc dari User
		// 	var users []user.User
		// 	// // untuk mengcek nilai array berisi 0
		// 	// length := len(users)
		// 	// fmt.Println(length)

		// 	// untuk mengambil data struc dari User
		// 	db.Find(&users)

		// untuk mengcek nilai array berisi 2 karna ada database usernya = 2
		// length = len(users)
		// fmt.Println(length)

		// 	for _, user := range users {
		// 		fmt.Println(user.Name)
		// 		fmt.Println(user.Email)
		// 		fmt.Println("===========================")
		// 	}

		// udah ga di pakai
		// router := gin.Default()
		// router.GET("/handler", handler)
		// router.Run()
	}
	userRepository := user.NewRepository(db)
	// tambahkan campaigns API ke db
	campaignRepository := campaign.NewRepository(db)
	userService := user.NewService(userRepository)
	campaignService := campaign.NewService(campaignRepository)
	//transactionRepository := transaction.Repository

	authService := auth.NewService()

	// MANUAL input Campaign baru (CreateCampaign API)
	// input := campaign.CreateCampaignInput{}
	// input.Name = "test CreateCampaignAPI"
	// input.ShortDescription = "short"
	// input.Description = "looooooong"
	// input.GoalAmount = 10000000
	// input.Perks = "Hadiah satu, dua, tiga"
	// inputUser, _ := userService.GetUserByID(17) //ambil data userID = 1
	// input.User = inputUser
	// _, err = campaignService.CreateCampaign(input)
	// if err != nil {
	// 	log.Fatal(err.Error())
	// }

	// DI nonaktifkan karna list of campaign akan di buat HANDLER nya di main.go
	// campaigns, _ := campaignService.GetCampaigns(1) // karna 0 berarti ambil semua campaign
	// fmt.Println(len(campaigns))

	// DI nonaktifkan karna ini manual untuk mencari repository campaign
	// ini FindAll
	//campaigns, err := campaignRepository.FindAll()
	// Ini FindByUserID
	// campaigns, err := campaignRepository.FindByUserID(0)
	// fmt.Println("debug")
	// fmt.Println("debuging")
	// fmt.Println(len(campaigns)) // ini cara mengetahui ada berapa data di campaign

	// for _, campaign := range campaigns {
	// 	fmt.Println(campaign.Name)
	// 	if len(campaign.CampaignImages) > 0 {
	// 		fmt.Println("Jumlah gambar")
	// 		fmt.Println(len(campaign.CampaignImages))
	// 		fmt.Println(campaign.CampaignImages[2].FileName)
	// 	}
	// }

	// INI dimatikan karna akan di buat / dipanggil di middleware
	// Membuat Validate Token JWT manual dulu
	// token, err := authService.ValidateToken("eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJ1c2VyX2lkIjo2fQ.GZXFQ5Pf7tzjTlwiBSeqLNTvCQifXYoaIUwATVa1ZP8")
	// if err != nil {
	// 	fmt.Println("ERROR")
	// 	fmt.Println("ERROR")
	// 	fmt.Println("ERROR")
	// }
	// if token.Valid {
	// 	fmt.Println("TOKEN VALID")
	// } else {
	// 	fmt.Println("TOKEN INVALID")
	// }

	// untuk cek di terminal token nya muncul
	// fmt.Println(authService.GenerateToken(1001))

	// Untuk SET avatar manual UPLOAD berdasarkan ID = 1
	//userService.SaveAvatar(1, "images/1-profile.png")

	// DI nonaktifkan karna sudah di buat service login nya di service.go
	// userByEmail, err := userRepository.FindByEmail("postmanlagio@gmail.com")
	// if err != nil {
	// 	fmt.Println(err.Error())
	// }
	// if userByEmail.ID == 0 {
	// 	fmt.Println("user tidak di temukan")
	// } else {
	// 	fmt.Println(userByEmail.Name)
	// }

	// DI nonaktifkan karna sudah di buat di hanler user.go
	// DAN sudah di buat sessions Login nya di bawah
	// input := user.LoginInput{
	// 	Email:    "postmanlagi@gmail.com",
	// 	Password: "password",
	// }
	// user, err := userService.Login(input)
	// if err != nil {
	// 	fmt.Println("Terjadi Kesalahan")
	// 	fmt.Println(err.Error())
	// }
	// fmt.Println(user.Email)
	// fmt.Println(user.Name)

	// HANDLER
	userHandler := handler.NewUserHandler(userService, authService) // tambahkan authService
	campaignHandler := handler.NewCampaignHandler(campaignService)  // tambahkan campaigns

	router := gin.Default()
	// membuat static images supaya bisa di akses routing folder nya langsung di browser
	// cth = http://localhost:8080/images/16-Green%20Ui%20Design%20Letter%20U%20Logo.png
	router.Static("/images", "./images")
	api := router.Group("/api/v1")

	api.POST("/users", userHandler.RegisterUser)
	api.POST("/sessions", userHandler.Login)
	api.POST("/email_checkers", userHandler.CheckEmailAvailability)
	api.POST("/avatars", authMiddleware(authService, userService), userHandler.UploadAvatar) //Middleware karna mau upload itu harus login user dulu ga sembarangan upload
	// Ambil data campaigns get dari server
	api.GET("/campaigns", campaignHandler.GetCampaigns)
	api.GET("/campaigns/:id", campaignHandler.GetCampaign)
	// route untuk CreateCampaign API
	api.POST("/campaigns", authMiddleware(authService, userService), campaignHandler.CreateCampaign) //Middleware karna mau upload itu harus login user dulu ga sembarangan upload
	// route untuk UpdateCampaign API
	api.PUT("/campaigns/:id", authMiddleware(authService, userService), campaignHandler.UpdateCampaign) //Middleware karna mau upload itu harus login user dulu ga sembarangan upload
	// route untuk UPLOAD Campaign Image API
	api.POST("/campaign-images", authMiddleware(authService, userService), campaignHandler.UploadImage) //Middleware karna mau upload itu harus login user dulu ga sembarangan upload

	router.Run()
}

// Middleware Authentication
// Cara nya = (di kerjain dari bawah dulu)
// Ambil nilai header Authorization: Bearer tokentokentoken (jadi dari client kirim Auth bearer, lalu ambil header nya)
// Dari header Authorization, kita ambil nilai tokennya aja
// Kita Validasi token (pake service Validatetoken)
// jika token valid,
// kita ambil user_id
// ambil user dari db berdasarkan user_id lewat service (user/service.go)
// kita set context isinya user (context itu tempat untuk menyimpan suatu nilai nanti bisa di GET dari tempat yg lain)

// Kita kerjakan Middleware(1)
func authMiddleware(authService auth.Service, userService user.Service) gin.HandlerFunc {
	// Kita kerjakan Middleware (1)
	return func(c *gin.Context) {
		authHeader := c.GetHeader("Authorization")
		if !strings.Contains(authHeader, "Bearer") {
			response := helper.APIResponse(" Unautorized", http.StatusUnauthorized, "error", nil)
			// pake AbortWithStatusJSON karna middle, kalo proses nya lancar dari user ke middle dulu baru ke UploadAvatar
			c.AbortWithStatusJSON(http.StatusUnauthorized, response)
			return
		}
		// Untuk menghapus Bearer nya, hanya dapet token nya aja gini cara nya
		// Bearer tokentokentoken
		tokenString := ""
		arrayToken := strings.Split(authHeader, " ")
		if len(arrayToken) == 2 {
			tokenString = arrayToken[1]
		}

		token, err := authService.ValidateToken(tokenString)
		if err != nil {
			response := helper.APIResponse(" Unautorized", http.StatusUnauthorized, "error", nil)
			// pake AbortWithStatusJSON karna middle, kalo proses nya lancar dari user ke middle dulu baru ke UploadAvatar
			c.AbortWithStatusJSON(http.StatusUnauthorized, response)
			return
		}

		claim, ok := token.Claims.(jwt.MapClaims)

		if !ok || !token.Valid {
			response := helper.APIResponse(" Unautorized", http.StatusUnauthorized, "error", nil)
			// pake AbortWithStatusJSON karna middle, kalo proses nya lancar dari user ke middle dulu baru ke UploadAvatar
			c.AbortWithStatusJSON(http.StatusUnauthorized, response)
			return
		}
		userID := int(claim["user_id"].(float64))

		user, err := userService.GetUserByID(userID)
		if err != nil {
			response := helper.APIResponse(" Unautorized", http.StatusUnauthorized, "error", nil)
			// pake AbortWithStatusJSON karna middle, kalo proses nya lancar dari user ke middle dulu baru ke UploadAvatar
			c.AbortWithStatusJSON(http.StatusUnauthorized, response)
			return
		}

		c.Set("currentUser", user)
	}

}

// Cara manual RegisterUserInput karna sudah di buat auto oleh handler
// userInput := user.RegisterUserInput{}
// userInput.Name = " Test simpan dari service"
// userInput.Email = "contoh@gmail.com"
// userInput.Occupation = "anak band"
// userInput.Password = "password"

// userService.RegisterUser(userInput)

// beda lagi
// user := user.User{
// 	Name: "Test Simpan",
// }
// userRepository.Save(user)

// NOTED
// Untuk mengakses DataBase itu harus mengunakan function layering
// input dari user
// handler : mapping input dari user -> struct input
// service : melakukan mapping dari struct input ke struct user
// repository
// db

// Noted Routes nya seperti ini :
// input
// handler mapping input ke struct
// handler mapping ke struct User
// repositorty save struct User ke db
// db

// HANDLER TEST untuk mencek dan menampilkan table struck users
// func handler(c *gin.Context) {
// 	dsn := "root:@tcp(127.0.0.1:3306)/crowdfunding?charset=utf8mb4&parseTime=True&loc=Local"
// 	db, err := gorm.Open(mysql.Open(dsn), &gorm.Config{})

// 	if err != nil {
// 		log.Fatal(err.Error())
// 	}

// 	// Ambil Nilai dari table User struc nya
// 	var users []user.User
// 	db.Find(&users)

// 	// lalu hasil nya kita tampilkan menggunakan method c.JSON panggil variable users
// 	c.JSON(http.StatusOK, users)
// 	// Lalu coba run localhost:8080/handler = maka akan muncul list users
// }
