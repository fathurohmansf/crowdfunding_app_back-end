package transaction

import (
	"crowdfunding/campaign"
	"crowdfunding/payment"
	"errors"
	"fmt"
	"strconv"
	"strings"
)

type service struct {
	repository Repository

	campaignRepository campaign.Repository

	paymentService payment.Service
}

type Service interface {
	GetTransactionByCampaignID(input GetCampaignTransactionInput) ([]Transaction, error)

	GetTransactionByUserID(userID int) ([]Transaction, error)

	CreateTransaction(input CreateTransactionInput) (Transaction, error)

	ProcessPayment(input TransactionNotificationInput) error

	GetAllTransactions() ([]Transaction, error)
}

func NewService(repository Repository, campaignRepository campaign.Repository, paymentService payment.Service) *service {
	return &service{repository, campaignRepository, paymentService}
}

func (s *service) GetTransactionByCampaignID(input GetCampaignTransactionInput) ([]Transaction, error) {

	campaign, err := s.campaignRepository.FindByID(input.ID)

	if err != nil {
		return []Transaction{}, err
	}

	if campaign.UserID != input.User.ID {
		return []Transaction{}, errors.New(" Not an owner of the campaign")
	}

	transaction, err := s.repository.GetByCampaignID(input.ID)
	if err != nil {
		return transaction, err
	}
	return transaction, nil
}

func (s *service) GetTransactionByUserID(userID int) ([]Transaction, error) {

	transaction, err := s.repository.GetByUserID(userID)
	if err != nil {
		return transaction, err
	}
	return transaction, nil
}

func (s *service) CreateTransaction(input CreateTransactionInput) (Transaction, error) {
	transaction := Transaction{}
	transaction.CampaignID = input.CampaignID
	transaction.Amount = input.Amount
	transaction.UserID = input.User.ID
	transaction.Status = "panding"
	// transaction.Code = "ORDER-001"

	// Logic to generate Transaction Code
	lastTransaction, err := s.repository.GetLastTransaction(transaction.ID)
	if err != nil {
		return transaction, fmt.Errorf("Gagal mendapatkan ID last Transaction: %v", err)
	}

	if lastTransaction.ID == 0 {
		// define awal Code transaction
		transaction.Code = "ORDER-001"
	} else {
		// increment the last transaction number Code
		parts := strings.Split(lastTransaction.Code, "-")
		if len(parts) == 2 {
			lastNumber, err := strconv.Atoi(parts[1])
			if err != nil {
				return transaction, fmt.Errorf("failed to parse last transaction code: %v", err)
			}
			transaction.Code = fmt.Sprintf("ORDER-%03d", lastNumber+1)
		} else {
			return transaction, fmt.Errorf("invalid transaction code format in last transaction: %s", lastTransaction.Code)
		}
	}

	newTranscation, err := s.repository.Save(transaction)
	if err != nil {
		return newTranscation, err
	}

	paymentTransaction := payment.Transaction{
		ID:     newTranscation.ID,
		Amount: newTranscation.Amount,
	}

	paymentURL, err := s.paymentService.GetPaymentURL(paymentTransaction, input.User)
	if err != nil {
		return newTranscation, err
	}
	newTranscation.PaymentURL = paymentURL
	newTranscation, err = s.repository.Update(newTranscation)
	if err != nil {
		return newTranscation, err
	}
	return newTranscation, nil
}

func (s *service) ProcessPayment(input TransactionNotificationInput) error {
	transaction_id, _ := strconv.Atoi(input.OrderID)

	transaction, err := s.repository.GetByID(transaction_id)
	if err != nil {
		return err
	}
	if input.PaymentType == "credit_card" && input.TransactionStatus == "capture" && input.FraudStatus == "accept" {
		transaction.Status = "paid"
	} else if input.TransactionStatus == "sattlement" {
		transaction.Status = "paid"
	} else if input.TransactionStatus == "deny" || input.TransactionStatus == "expire" || input.TransactionStatus == "cancel" {
		transaction.Status = "cancelled"
	}

	updatedTransaction, err := s.repository.Update(transaction)
	if err != nil {
		return err
	}

	campaign, err := s.campaignRepository.FindByID(updatedTransaction.CampaignID)
	if err != nil {
		return err
	}
	if updatedTransaction.Status == "paid" {
		campaign.BackerCount = campaign.BackerCount + 1
		campaign.CurrentAmount = campaign.CurrentAmount + updatedTransaction.Amount

		_, err := s.campaignRepository.Update(campaign)
		if err != nil {
			return err
		}
	}
	return nil
}

func (s *service) GetAllTransactions() ([]Transaction, error) {
	transactions, err := s.repository.FindAll()
	if err != nil {
		return transactions, err
	}
	return transactions, nil
}
