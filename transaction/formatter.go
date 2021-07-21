package transaction

import "time"

type CampaignTransactionResponse struct {
	Id        int       `json:"id"`
	Name      string    `json:"name"`
	Amount    int       `json:"amount"`
	CreatedAt time.Time `json:"created_at"`
}

type UserTransactionResponse struct {
	Id int `json:"id"`
	Amount int `json:"amount"`
	Status string `json:"status"`
	CreatedAt time.Time `json:"created_at"`
	Campaign UserTransactionResponseCampaign `json:"campaign"`
}

type UserTransactionResponseCampaign struct {
	Name string `json:"name"`
	ImageUrl string `json:"image_url"`
}

type TransactionResponse struct {
	Id int `json:"id"`
	CampaignId int `json:"campaign_id"`
	Amount int `json:"amount"`
	Status string `json:"status"`
	Code string `json:"code"`
	PaymentURL string `json:"payment_url"`
}

func FormatCampaignTransaction(transaction Transaction) CampaignTransactionResponse {
	formatter := CampaignTransactionResponse{}
	formatter.Id = transaction.Id
	formatter.Name = transaction.User.Name
	formatter.Amount = transaction.Amount
	return formatter
}

func FormatTransactionArr(transactions []Transaction) []CampaignTransactionResponse {
	if len(transactions) == 0 {
		return []CampaignTransactionResponse{}
	}

	var transactionsResponse []CampaignTransactionResponse
	for _, transaction := range transactions {
		transactionsResponse = append(transactionsResponse, FormatCampaignTransaction(transaction))
	}

	return transactionsResponse
}

func FormatUserTransaction(transaction Transaction) UserTransactionResponse {
	formatter := UserTransactionResponse{}
	formatter.Id = transaction.Id
	formatter.Amount = transaction.Amount
	formatter.Status = transaction.Status
	formatter.CreatedAt = transaction.CreatedAt
	formatter.Campaign.Name = transaction.Campaign.Name
	formatter.Campaign.ImageUrl = transaction.Campaign.CampaignImages[0].FileName

	return formatter
}

func FormatUserTransactionArr(transactions []Transaction) []UserTransactionResponse {
	if len(transactions) == 0 {
		return []UserTransactionResponse{}
	}

	var transactionResponse []UserTransactionResponse
	for _, transaction := range transactions {
		transactionResponse = append(transactionResponse, FormatUserTransaction(transaction))
	}

	return transactionResponse
}

func FormatTransaction(t Transaction) TransactionResponse {
	var formatter TransactionResponse
	formatter.Id = t.Id
	formatter.CampaignId = t.CampaignId
	formatter.Amount = t.Amount
	formatter.Status = t.Status
	formatter.Code = t.Code
	formatter.PaymentURL = t.PaymentURL

	return formatter
}