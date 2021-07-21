package transaction

import (
	"errors"
	"fmt"
	"golang-crowdfunding/campaign"
	"golang-crowdfunding/payment"
	"golang-crowdfunding/user"
	"strconv"
)


type service struct {
	repository         Repository
	campaignRepository campaign.Repository
	paymentService payment.Service
	userService user.Service
}

type Service interface {
	GetTransactionsByCampaignId(campaignId int, userId int) ([]Transaction, error)
	GetTransactionsByUserId(userId int) ([]Transaction, error)
	CreateTransaction(input CreateTransactionInput) (Transaction, error)
	ProcessPayment(input TransactionNotificationInput) error
}

func NewService(r Repository, cr campaign.Repository, p payment.Service, u user.Service) *service {
	return &service{r, cr, p, u}
}

func (s *service) GetTransactionsByCampaignId(campaignId int, userId int) ([]Transaction, error) {
	campaign, err := s.campaignRepository.FindById(campaignId)
	if err != nil {
		return []Transaction{}, err
	}

	if campaign.UserId != userId {
		return []Transaction{}, errors.New("Unauthorized")
	}

	transactions, err := s.repository.GetByCampaignId(campaignId)
	if err != nil {
		return transactions, err
	}

	return transactions, nil
}

func (s *service) GetTransactionsByUserId(userId int) ([]Transaction, error) {
	transactions, err := s.repository.GetByUserId(userId)
	if err != nil {
		return transactions, err
	}

	return transactions, nil
}

func (s *service) CreateTransaction(input CreateTransactionInput) (Transaction, error) {
	var lastTransaction Transaction
	fmt.Println("SERVICE CALLED")
	lastTransaction, _ = s.repository.Last(lastTransaction)

	increment := int(lastTransaction.Id) + 1

	transaction := Transaction{}
	transaction.Id = increment
	transaction.CampaignId = input.CampaignId
	transaction.Amount = input.Amount
	transaction.UserId = input.UserId
	transaction.Status = "pending"
	transaction.Code = strconv.Itoa(increment)
	transaction.PaymentURL = ""

	t, err := s.repository.Save(transaction)
	if err != nil {
		return t, err
	}

	user, err := s.userService.GetUserById(input.UserId)
	if err != nil {
		return t, err
	}

	paymentTransaction := payment.Transaction{
		Id: t.Id,
		Amount: t.Amount,
	}

	paymentURL, err := s.paymentService.GetPaymentURL(paymentTransaction, user)
	if err != nil {
		return t, err
	}

	t.PaymentURL = paymentURL.RedirectURL
	finalTransaction, err := s.repository.Update(t)
	if err != nil {
		return finalTransaction, err
	}

	return finalTransaction, nil
}

func (s *service) ProcessPayment(input TransactionNotificationInput) error {
	
	transaction_Id , _ := strconv.Atoi(input.OrderID)
	transaction, err := s.repository.GetById(transaction_Id)
	if err != nil {
		return err
	}

	if(input.PaymentType == "credit_card" && input.TransactionStatus == "capture" && input.FraudStatus == "accept") {
		transaction.Status = "paid"
	} else if input.TransactionStatus == "settlement" {
		transaction.Status = "paid"
	} else if input.TransactionStatus == "deny" || input.TransactionStatus == "expire" || input.TransactionStatus == "cancel" {
		transaction.Status = "cancelled"
	}

	updatedTransaction, err := s.repository.Update(transaction)
	if err != nil {
		return err
	}

	campaign, err := s.campaignRepository.FindById(updatedTransaction.CampaignId)
	if err != nil {
		return err
	}

	if updatedTransaction.Status == "paid" {
		campaign.BackerCount += 1
		campaign.CurrentAmount += updatedTransaction.Amount

		_, err := s.campaignRepository.Update(campaign)
		if err != nil {
			return err
		}
	}

	return nil
}