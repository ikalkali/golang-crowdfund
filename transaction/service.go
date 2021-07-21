package transaction

import (
	"errors"
	"golang-crowdfunding/campaign"
)


type service struct {
	repository         Repository
	campaignRepository campaign.Repository
}

type Service interface {
	GetTransactionsByCampaignId(campaignId int, userId int) ([]Transaction, error)
}

func NewService(r Repository, cr campaign.Repository) *service {
	return &service{r, cr}
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
