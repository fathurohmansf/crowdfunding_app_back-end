package campaign

type Service interface {
	GetCampaigns(userID int) ([]Campaign, error)
	// GET campaign ByID
	//GetCampaignByID(input GetCampaignDetailInput) (Campaign, error)
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
// func (s *service) GetCampaignByID(input GetCampaignDetailInput) (Campaign, error) {
// 	campaign, err := s.repository.FindByID(input.ID)
// 	if err != nil {
// 		return campaign, err
// 	}
// 	return campaign, nil
// }
