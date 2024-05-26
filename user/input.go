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

// untuk service API Login
type LoginInput struct {
	Email    string `json:"email" binding:"required,email"`
	Password string `json:"password" binding:"required"`
}

// struct untuk mengcekEmail / EmailChecker
type CheckEmailInput struct {
	Email string `json:"email" binding:"required,email"`
}
