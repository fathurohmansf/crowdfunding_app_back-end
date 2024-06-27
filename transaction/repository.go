package transaction

import "gorm.io/gorm"

type repository struct {
	db *gorm.DB
}

type Repository interface {
	GetByCampaignID(campaignID int) ([]Transaction, error)
	// Get userID
	GetByUserID(userID int) ([]Transaction, error)
	// Midtrans
	Save(transaction Transaction) (Transaction, error)
}

func NewRepository(db *gorm.DB) *repository {
	return &repository{db}
}

// implementasi cari Transaksi berdasarkan CampaignID
func (r *repository) GetByCampaignID(campaignID int) ([]Transaction, error) {
	var transaction []Transaction

	// Sebelum Tidak mengurutkan data berdasarkan ID transaction / Created_At
	//err := r.db.Preload("User").Where("campaign_id = ?", campaignID).Find(&transaction).Error
	// Sesudah mengurutkan data berdasarkan ID transaction / Cerated_At
	err := r.db.Preload("User").Where("campaign_id = ?", campaignID).Order("id desc").Find(&transaction).Error
	if err != nil {
		return transaction, err
	}
	return transaction, nil
}

// User Transaction API
func (r *repository) GetByUserID(userID int) ([]Transaction, error) {
	var transactions []Transaction
	// jadi kita men load Campaign yang mempunyai relasi ke campaign images tapi hanya bisa akses/tampil yg is_primary=1
	// karna di transaction tidak punya relasi ke campaign_images
	// baru kita cari userID dan Cari Transaksi berdasarkan dataID
	err := r.db.Preload("Campaign.CampaignImages", "campaign_images.is_primary = 1").Where("user_id = ?", userID).Order("id desc").Find(&transactions).Error
	if err != nil {
		return transactions, err
	}
	return transactions, nil
}

// Membuat fungsi Midtrans payment
func (r *repository) Save(transaction Transaction) (Transaction, error) {
	err := r.db.Create(&transaction).Error
	if err != nil {
		return transaction, err
	}
	return transaction, nil
}
