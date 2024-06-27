package helper

import "github.com/go-playground/validator/v10"

// Membuat struck response untuk ke front end
type Response struct {
	// ini punya meta yg merupakan type data nya meta
	Meta Meta `json:"meta"`
	// punya data yang type nya adalah interface kenapa interface karna
	// isi data nya bisa bebas, berupa json,int,list string dll.
	Data interface{} `json:"data"`
}

// arti `json:"data"` ini jika di ubah menjadi json
// maka field yg akan di munculkan yg di dalam tanda petik "data" ini
type Meta struct {
	Message string `json:"message"`
	Code    int    `json:"code"`
	Status  string `json:"status"`
}

// karna data nya interface jadi kita matikan
// type data struct{

// }

// nama API besar itu untuk bisa akses ke public, dan penamaan interface dengan data
func APIResponse(message string, code int, status string, data interface{}) Response {
	// isi ini membuat object dari struct response di atas
	meta := Meta{
		Message: message,
		Code:    code,
		Status:  status,
	}

	jsonResponse := Response{
		Meta: meta,
		Data: data,
	}
	return jsonResponse
}

func FormatValidationError(err error) []string {
	var errors []string

	for _, e := range err.(validator.ValidationErrors) {
		errors = append(errors, e.Error())
	}
	return errors
}

// DI BAWAH ini untuk front-end yang minta seperti ini
// meta : {
// 	message : 'Your account has been created',
// 	code: 200,
// 	status: 'success',
// },
// data : {
// 	id : 1
// 	name : "Ahmad Fathurohman",
// 	occupation : "conten creator",
// 	email : "com.fathur@gmail.com",
// 	token : "peterpanyangterdalam",
// }
