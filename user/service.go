package user

import (
	"errors"

	"golang.org/x/crypto/bcrypt"
)

type Service interface {
	RegisterUser(input RegisterUserInput) (User, error)
	// bikin kontrak di Service untuk Login
	Login(input LoginInput) (User, error)
}

type service struct {
	repository Repository
}

func NewService(repository Repository) *service {
	return &service{repository}
}

func (s *service) RegisterUser(input RegisterUserInput) (User, error) {
	user := User{}
	user.Name = input.Name
	user.Email = input.Email
	user.Occupation = input.Occupation
	passwordHash, err := bcrypt.GenerateFromPassword([]byte(input.Password), bcrypt.MinCost)
	if err != nil {
		return user, err
	}
	user.PasswordHash = string(passwordHash)
	user.Role = "user"

	newUser, err := s.repository.Save(user)
	if err != nil {
		return newUser, err
	}
	return newUser, nil
}

// mapping struct input ke struct User
// simpan struct User melalui repository

// Membuat fungsi Service Login dgn parameter input
func (s *service) Login(input LoginInput) (User, error) {
	email := input.Email
	password := input.Password

	// mencari user yang memiliki email dari repository
	user, err := s.repository.FindByEmail(email)
	if err != nil {
		return user, err
	}
	// buat pengecekan jika ada error
	if user.ID == 0 {
		return user, errors.New("No user found on that email")
	}
	// untuk pengecekan password dari database
	err = bcrypt.CompareHashAndPassword([]byte(user.PasswordHash), []byte(password))
	// pengecekan error password
	if err != nil {
		return user, err
	}
	return user, nil
}
