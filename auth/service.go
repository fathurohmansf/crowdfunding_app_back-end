package auth

import (
	"github.com/dgrijalva/jwt-go"
)

// Membuat service baru
type Service interface {
	GenerateToken(userID int) (string, error)
}

type jwtService struct {
}

// SECRET_KEY ini bersifat sementara ga boleh taro mentah codingan gini
var SECRET_KEY = []byte("BWASTARTUP_s3cr3T_k3y")

// aktifkan fungsi service jwt
func NewService() *jwtService {
	return &jwtService{}
}

// Membuat funsi jwt untuk signature di browser sesuai jwt.io
func (s *jwtService) GenerateToken(userID int) (string, error) {
	claim := jwt.MapClaims{}
	// json yg akan di tampilkan key nya user_id
	claim["user_id"] = userID

	// Buat token menggunakan algoritma yg di gunakan HS256
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claim)
	// Untuk Faktor keamanan pakai SECRET_KEY
	signToken, err := token.SignedString(SECRET_KEY)
	if err != nil {
		return signToken, err
	}
	return signToken, nil
}
