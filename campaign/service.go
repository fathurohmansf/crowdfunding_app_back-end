package campaign

import (
	"fmt"

	"github.com/gosimple/slug"
)

type Service interface {
	GetCampaigns(userID int) ([]Campaign, error)
	// GET campaign ByID
	GetCampaignByID(input GetCampaignDetailInput) (Campaign, error)
	// Membuat Create campaign API
	CreateCampaign(input CreateCampaignInput) (Campaign, error)
}

type service struct {
	repository Repository
}

func NewService(repository Repository) *service {
	return &service{repository}
}

func (s *service) GetCampaigns(userID int) ([]Campaign, error) {
	// switching, kalo userID = 0, jika userID ada
	// data nya maka akan menampilkan campaign sesuai user itu yg di buat
	if userID != 0 {
		campaigns, err := s.repository.FindByUserID(userID)
		if err != nil {
			return campaigns, err
		}
		return campaigns, nil
	}
	campaigns, err := s.repository.FindAll()
	if err != nil {
		return campaigns, err
	}
	return campaigns, nil
}

// implementasi interfaces Get Campaign ByID
func (s *service) GetCampaignByID(input GetCampaignDetailInput) (Campaign, error) {
	campaign, err := s.repository.FindByID(input.ID)
	if err != nil {
		return campaign, err
	}
	return campaign, nil
}

// implementasi interfaces CreateCampaign
func (s *service) CreateCampaign(input CreateCampaignInput) (Campaign, error) {
	campaign := Campaign{}
	campaign.Name = input.Name
	campaign.ShortDescription = input.ShortDescription
	campaign.Description = input.Description
	campaign.Perks = input.Perks
	campaign.GoalAmount = input.GoalAmount
	campaign.UserID = input.User.ID
	// SLUG di buat nanti pakai libarary https://github.com/gosimple/slug = ke terminal = go get -u github.com/gosimple/slug
	// %s untuk string, %d untuk int
	slugCandidate := fmt.Sprintf("%s %d", input.Name, input.User.ID)
	campaign.Slug = slug.Make(slugCandidate) //Nama Camapaign 10 => nama-campaign-10

	// save campaign baru
	newCampaign, err := s.repository.Save(campaign)
	if err != nil {
		return newCampaign, err
	}
	return newCampaign, nil
}
