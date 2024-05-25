package user

type UserFormatter struct {
	ID         int    `json:"id"`
	Name       string `"json:name"`
	Occupation string `json:"occupation"`
	Email      string `json:"email"`
	Token      string `json:"token"`
}

// Membuat fungsi format aga tidak lagi response json dari struct
// User dari entity.go, karna front-end tidak membutuhkan CreateAt dll.
func FormatUser(user User, token string) UserFormatter {
	formatter := UserFormatter{
		ID:         user.ID,
		Name:       user.Name,
		Occupation: user.Occupation,
		Email:      user.Email,
		Token:      token,
	}
	return formatter
}
