package auth

import (
	"errors"

	"github.com/dgrijalva/jwt-go"
)

// Membuat service baru
type Service interface {
	GenerateToken(userID int) (string, error)
	// Untuk Validasi token / ttd
	ValidateToken(token string) (*jwt.Token, error)
}

type jwtService struct {
}

// SECRET_KEY ini bersifat sementara ga boleh taro mentah codingan gini
// Ini itu Verify Signature / validasi token ttd online
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

// Fungsi untuk Validasi Token / Signature ttd token
func (s *jwtService) ValidateToken(encondedtoken string) (*jwt.Token, error) {
	token, err := jwt.Parse(encondedtoken, func(token *jwt.Token) (interface{}, error) {
		_, ok := token.Method.(*jwt.SigningMethodHMAC)

		if !ok {
			return nil, errors.New("Invalid token")
		}
		return []byte(SECRET_KEY), nil
	})
	if err != nil {
		return token, err
	}
	return token, nil
}
