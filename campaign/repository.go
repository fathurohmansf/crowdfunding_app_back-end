package campaign

import "gorm.io/gorm"

type Repository interface {
	FindAll() ([]Campaign, error)
	FindByUserID(userID int) ([]Campaign, error)
	// GET Campaign ByID
	FindByID(ID int) (Campaign, error)
	// untuk Create Campaign API
	Save(campaign Campaign) (Campaign, error)
	// UPDATE Campaign API untuk update
	Update(campaign Campaign) (Campaign, error)
	// UPLOAD Campaign Image API
	CreateImage(CampaignImage CampaignImage) (CampaignImage, error)
	// ISPRIMARY = TRUE , untuk campaign images
	MarkAllImagesAsNonPrimary(campaignID int) (bool, error)
}

type repository struct {
	db *gorm.DB
}

func NewRepository(db *gorm.DB) *repository {
	return &repository{db}
}

// Kita buat dulu FIndAll
func (r *repository) FindAll() ([]Campaign, error) {
	var campaigns []Campaign

	err := r.db.Find(&campaigns).Error //.Find("CampaignImages", "campaign_images.is_primary = 1")
	if err != nil {
		return campaigns, err
	}
	return campaigns, nil
}

// Fungsi untuk ambil data FindByUserID
func (r *repository) FindByUserID(userID int) ([]Campaign, error) {
	var campaigns []Campaign
	// fungsi Preload ini relasi untuk campaign images
	//err := r.db.Where("user_id = ?", userID).Preload("CampaignImages").Find(&campaigns).Error
	// Di nonaktifkan karna hanya untuk menampilkan primary true saja (1)
	err := r.db.Where("user_id = ?", userID).Preload("CampaignImages", "campaign_images.is_primary = 1").Find(&campaigns).Error
	// Cek error kalo udh ga error apus aja
	//err := r.db.Where("user_id = ?", userID).Find(&campaigns).Error
	if err != nil {
		return campaigns, err
	}
	return campaigns, nil
}

// buat implementasi dari interface FindByID
func (r *repository) FindByID(ID int) (Campaign, error) {
	var campaign Campaign
	err := r.db.Preload("User").Preload("CampaignImages").Where("id = ?", ID).Find(&campaign).Error // .Preload("User")
	if err != nil {
		return campaign, err
	}
	return campaign, nil
}

// buat implementasi fungsi dari Create Campaign API
func (r *repository) Save(campaign Campaign) (Campaign, error) {
	err := r.db.Create(&campaign).Error
	if err != nil {
		return campaign, err
	}
	return campaign, nil
}

// UPDATE implementasi update Campaign API
func (r *repository) Update(campaign Campaign) (Campaign, error) {
	err := r.db.Save(&campaign).Error
	if err != nil {
		return campaign, err
	}
	return campaign, nil
}

// UPLOAD Campaign Image API point 1
func (r *repository) CreateImage(CampaignImage CampaignImage) (CampaignImage, error) {
	err := r.db.Create(&CampaignImage).Error
	if err != nil {
		return CampaignImage, err
	}
	return CampaignImage, nil
}

// Membuat fungsi dari interfaces ISprimary images = true, point 2
func (r *repository) MarkAllImagesAsNonPrimary(campaignID int) (bool, error) {
	// QUERY nya
	// Update campaign_images SET is_primary = false WHERE campaign_id = 1
	err := r.db.Model(&CampaignImage{}).Where("campaign_id = ?", campaignID).Update("is_primary", false).Error

	if err != nil {
		return false, err
	}
	return true, nil
}
