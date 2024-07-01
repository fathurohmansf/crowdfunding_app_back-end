package transaction

import "time"

type CampaignTransactionFormatter struct {
	ID        int       `json:"id"`
	Name      string    `json:"name"`
	Amount    int       `json:"amount"`
	CreatedAt time.Time `json:"created_at"`
}

// Formatter untuk User Transaction API
type UserTransactionFormatter struct {
	ID        int               `json:"id"`
	Amount    int               `json:"amount"`
	Status    string            `json:"status"`
	CreatedAt time.Time         `json:"created_at"`
	Campaign  CampaignFormatter `json:"campaign"`
}
type CampaignFormatter struct {
	Name     string `json:"name"`
	ImageURL string `json:"image_url"`
}

// struct ini untuk Format Create Transaction di Midtrans menggunakan POST (create)
type TransactionFormatter struct {
	ID         int    `json:"id"`
	CampaignID int    `json:"campaign_id"`
	UserID     int    `json:"user_id"`
	Amount     int    `json:"amount"`
	Status     string `json:"status"`
	Code       string `json:"code"`
	PaymentURL string `json:"payment_url"`
}

// fungsi ini untuk format 1 (one) Transaction Campaign
func FormatCampaignTransaction(transaction Transaction) CampaignTransactionFormatter {
	formatter := CampaignTransactionFormatter{}
	formatter.ID = transaction.ID
	formatter.Name = transaction.User.Name
	formatter.Amount = transaction.Amount
	formatter.CreatedAt = transaction.CreatedAt
	return formatter
}

// fungsi ini untuk format List Of Transaction
func FormatCampaignTransactions(transactions []Transaction) []CampaignTransactionFormatter {
	// jika nilai transaksi 0 maka balikkan array kosong {}
	if len(transactions) == 0 {
		return []CampaignTransactionFormatter{}
	}
	// jika nilai ada, akan melooping data transaction
	var transactionsFormatter []CampaignTransactionFormatter
	for _, transaction := range transactions {
		formatter := FormatCampaignTransaction(transaction)
		transactionsFormatter = append(transactionsFormatter, formatter)
	}
	return transactionsFormatter
}

// Fungsi untuk User Transaction API
func FormatUserTransaction(transaction Transaction) UserTransactionFormatter {
	formatter := UserTransactionFormatter{}
	formatter.ID = transaction.ID
	formatter.Amount = transaction.Amount
	formatter.Status = transaction.Status
	formatter.CreatedAt = transaction.CreatedAt

	campaignFormatter := CampaignFormatter{}
	campaignFormatter.Name = transaction.Campaign.Name
	campaignFormatter.ImageURL = ""
	// Ambil nilai yg lebih dari 0 , yaitu 1
	if len(transaction.Campaign.CampaignImages) > 0 {
		campaignFormatter.ImageURL = transaction.Campaign.CampaignImages[0].FileName
	}
	formatter.Campaign = campaignFormatter
	return formatter
}

// fungsi ini untuk format List Of User Transaction API
func FormatUserTransactions(transactions []Transaction) []UserTransactionFormatter {
	// jika nilai transaksi 0 maka balikkan array kosong {}
	if len(transactions) == 0 {
		return []UserTransactionFormatter{}
	}
	// jika nilai ada, akan melooping data transaction
	var transactionsFormatter []UserTransactionFormatter
	for _, transaction := range transactions {
		formatter := FormatUserTransaction(transaction)
		transactionsFormatter = append(transactionsFormatter, formatter)
	}
	// kembalikan nilai nya
	return transactionsFormatter
}

// fungsi ini untuk format 1 (one) Create Transaction pakai Midtrans
func FormatTransaction(transaction Transaction) TransactionFormatter {
	formatter := TransactionFormatter{}
	formatter.ID = transaction.ID
	formatter.CampaignID = transaction.CampaignID
	formatter.UserID = transaction.UserID
	formatter.Amount = transaction.Amount
	formatter.Status = transaction.Status
	formatter.Code = transaction.Code
	formatter.PaymentURL = transaction.PaymentURL
	return formatter
}
