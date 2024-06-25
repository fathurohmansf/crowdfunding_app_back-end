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
}

func NewService(repository Repository, campaignRepository campaign.Repository) *service {
	return &service{repository, campaignRepository}
}

func (s *service) GetTransactionByCampaignID(input GetCampaignTransactionInput) ([]Transaction, error) {

	// Menerapkan Authorization (agar user lain tidak bisa liat data transaksi campaign sendiri)
	// 1. Get Campaign
	// 2. Check campaign.user.id != user_id yang melakukan request
	campaign, err := s.campaignRepository.FindByID(input.ID)
	if err != nil {
		return []Transaction{}, err
	}
	if campaign.UserID != input.User.ID {
		return []Transaction{}, errors.New("Not an owner of the campaign")
	}
	transaction, err := s.repository.GetByCampaignID(input.ID)
	if err != nil {
		return transaction, err
	}
	return transaction, nil
}
