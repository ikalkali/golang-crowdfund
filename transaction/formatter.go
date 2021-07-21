package transaction

import "time"

type CampaignTransactionResponse struct {
	Id        int       `json:"id"`
	Name      string    `json:"name"`
	Amount    int       `json:"amount"`
	CreatedAt time.Time `json:"created_at"`
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