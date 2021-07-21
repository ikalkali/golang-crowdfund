package transaction

import (
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
	User user.User
	CreatedAt  time.Time
	UpdatedAt time.Time
}