package transaction

import (
	"crowdfunding/campaign"
	"errors"
)

type service struct {
	repository Repository
	// untuk menambahkan campaign Repo untuk kebutuhan Authorization
	campaignRepository campaign.Repository
}

type Service interface {
	GetTransactionByCampaignID(input GetCampaignTransactionInput) ([]Transaction, error)
	// Membuat kontrak untuk User Transaction API
	GetTransactionByID(userID int) ([]Transaction, error)
}

func NewService(repository Repository, campaignRepository campaign.Repository) *service {
	return &service{repository, campaignRepository}
}

func (s *service) GetTransactionByCampaignID(input GetCampaignTransactionInput) ([]Transaction, error) {

	// Menerapkan Authorization (agar user lain tidak bisa liat data transaksi campaign sendiri)
	// 1. Get Campaign dari Repo
	// 2. Check campaign.user.id != user_id yang melakukan request
	campaign, err := s.campaignRepository.FindByID(input.ID)
	// jika ada error maka list kosong
	if err != nil {
		return []Transaction{}, err
	}
	// jika ada error namun ID user salah maka tampilkan bukan user pembuat campaign
	if campaign.UserID != input.User.ID {
		return []Transaction{}, errors.New(" Not an owner of the campaign")
	}
	// jika benar maka akan return transaction nya
	transaction, err := s.repository.GetByCampaignID(input.ID)
	if err != nil {
		return transaction, err
	}
	return transaction, nil
}

// Function untuk User Transaction API
func (s *service) GetTransactionByUserID(userID int) ([]Transaction, error) {
	// panggil repo/data userID
	transaction, err := s.repository.GetByUserID(userID)
	if err != nil {
		return transaction, err
	}
	return transaction, nil
}
