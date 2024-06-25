package transaction

import "time"

type CampaignTransactionFormatter struct {
	ID        int       `json:"id"`
	Name      string    `json:"name"`
	Amount    int       `json:"amount"`
	CreatedAt time.Time `json:"created_at"`
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
func FormatCampaignTransactions(transaction []Transaction) []CampaignTransactionFormatter {
	// jika nilai transaksi 0 maka balikkan array kosong {}
	if len(transaction) == 0 {
		return []CampaignTransactionFormatter{}
	}
	// jika nilai ada, akan melooping data transaction
	var transactionsFormatter []CampaignTransactionFormatter
	for _, transaction := range transaction {
		formatter := FormatCampaignTransaction(transaction)
		transactionsFormatter = append(transactionsFormatter, formatter)
	}
	return transactionsFormatter
}
