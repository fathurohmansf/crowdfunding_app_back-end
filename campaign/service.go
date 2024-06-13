package campaign

import (
	"errors"
	"fmt"

	"github.com/gosimple/slug"
)

type Service interface {
	GetCampaigns(userID int) ([]Campaign, error)
	// GET campaign ByID
	GetCampaignByID(input GetCampaignDetailInput) (Campaign, error)
	// Membuat Create campaign API
	CreateCampaign(input CreateCampaignInput) (Campaign, error)
	// UPDATE buat update campaign API
	UpdateCampaign(inputID GetCampaignDetailInput, inputData CreateCampaignInput) (Campaign, error)
	// UPLOAD Campaign Image API
	SaveCampaignImage(input CreateCampaignImageInput, fileLocation string) (CampaignImage, error)
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

// UPDATE implementasi interfaces Update
func (s *service) UpdateCampaign(inputID GetCampaignDetailInput, inputData CreateCampaignInput) (Campaign, error) {
	campaign, err := s.repository.FindByID(inputID.ID)
	if err != nil {
		return campaign, err
	}
	// Penbuatan fungsi untuk saat update campaign itu harus User punya campaign ga boleh user lain
	if campaign.UserID != inputData.User.ID {
		// UserID itu yang punya data/campaign , User.ID user yang punya request untuk update data/campaign
		// jika == maka allow, jika != maka muncul error.New di bawah ini
		return campaign, errors.New("Not an owner of the campaign")
	}

	// fungsi di bawah ini untuk update data dari data lama ke data baru untuk di update
	campaign.Name = inputData.Name
	campaign.ShortDescription = inputData.ShortDescription
	campaign.Description = inputData.Description
	campaign.Perks = inputData.Perks
	campaign.GoalAmount = inputData.GoalAmount

	// ini untuk memberikan data / simpan ke repository
	updatedCampaign, err := s.repository.Update(campaign)
	if err != nil {
		return updatedCampaign, err
	}
	return updatedCampaign, nil
}

// UPLOAD Implementasi interface UPLOAD Campaign Image API
func (s *service) SaveCampaignImage(input CreateCampaignImageInput, fileLocation string) (CampaignImage, error) {
	// Membuat definisi awal isPrimary = 0/false
	isPrimary := 0
	// Jika saat user klik/input maka nilai nya true
	if input.IsPrimary {
		isPrimary = 1
	}
	// Pengecekan is_primary nya true
	if input.IsPrimary {
		_, err := s.repository.MarkAllImagesAsNonPrimary(input.CampaignID)
		if err != nil {
			return CampaignImage{}, err
		}
	}
	campaignImage := CampaignImage{}
	campaignImage.CampaignID = input.CampaignID
	campaignImage.IsPrimary = isPrimary
	campaignImage.FileName = fileLocation

	newCampaignImage, err := s.repository.CreateImage(campaignImage)
	if err != nil {
		return newCampaignImage, err
	}
	return newCampaignImage, nil
}
