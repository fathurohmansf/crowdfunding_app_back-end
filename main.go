package main

import (
	"crowdfunding/handler"
	"crowdfunding/user"
	"log"

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
	userService := user.NewService(userRepository)
	userHandler := handler.NewUserHandler(userService)

	router := gin.Default()
	api := router.Group("/api/v1")

	api.POST("/users", userHandler.RegisterUser)

	router.Run()

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

}

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
