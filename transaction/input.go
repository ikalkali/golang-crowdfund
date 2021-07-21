package transaction

type CreateTransactionInput struct {
	Amount     int `json:"amount" validate:"required"`
	CampaignId int `json:"campaign_id" validate:"required"`
	UserId     int
}

type TransactionNotificationInput struct {
	TransactionStatus string `json:"transaction_status"`
	OrderID           string `json:"order_id"`
	PaymentType       string `json:"payment_type"`
	FraudStatus       string `json:"fraud_status"`
}