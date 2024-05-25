package user

type RegisterUserInput struct {
	Name       string `json:"name" binding:"required"`
	Occupation string `json:"occupation" binding:"required"`
	Email      string `json:"email" binding:"required,email"`
	Password   string `json:"password" binding:"required"`
}

// penjelasan:
// binding "required" ini arti nya harus di isi ga boleh data kosong
// binding:"required,email" = artinya harus di isi dgn format email
