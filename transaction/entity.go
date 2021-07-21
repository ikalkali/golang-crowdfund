package transaction

import (
	"golang-crowdfunding/campaign"
	"golang-crowdfunding/user"
	"time"
)


type Transaction struct {
	Id         int
	CampaignId int
	UserId     int
	Amount     int
	Status     string
	Code       string
	PaymentURL string
	User user.User
	Campaign campaign.Campaign
	CreatedAt  time.Time
	UpdatedAt time.Time
}