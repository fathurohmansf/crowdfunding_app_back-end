package transaction

import "gorm.io/gorm"

type repository struct {
	db *gorm.DB
}

type Repository interface {
	GetByCampaignID(campaignID int) ([]Transaction, error)
}

func NewRepository(db *gorm.DB) *repository {
	return &repository{db}
}

// implementasi cari Transaksi berdasarkan CampaignID
func (r *repository) GetByCampaignID(campaignID int) ([]Transaction, error) {
	var transaction []Transaction

	err := r.db.Preload("User").Where("campaign_id = ?", campaignID).Find(&transaction).Error
	if err != nil {
		return transaction, err
	}
	return transaction, nil
}
