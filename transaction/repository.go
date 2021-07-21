package transaction

import (
	"fmt"

	"gorm.io/gorm"
)


type repository struct {
	db *gorm.DB
}

type Repository interface {
	GetByCampaignId(campaignId int) ([]Transaction, error)
	GetByUserId(userId int) ([]Transaction, error)
	GetById(id int) (Transaction, error)
	Save(transaction Transaction) (Transaction, error)
	Last(transaction Transaction) (Transaction, error)
	Update(Transaction Transaction) (Transaction, error)
}

func NewRepository(db *gorm.DB) *repository {
	return &repository{db}
}

func (r *repository) GetByCampaignId(campaignId int) ([]Transaction, error) {
	var transactions []Transaction

	err := r.db.Preload("User").Where("campaign_id = ?", campaignId).Order("id desc").Find(&transactions).Error
	if err != nil {
		return transactions, err
	}

	return transactions, nil
}

func (r *repository) GetByUserId(userId int) ([]Transaction, error) {
	var transactions []Transaction

	err := r.db.Preload("Campaign.CampaignImages", "campaign_images.is_primary = 1").Where("user_id = ?", userId).Order("id desc").Find(&transactions).Error
	if err != nil {
		return transactions, err
	}

	return transactions, nil
}

func (r *repository) Save(transaction Transaction) (Transaction, error) {
	err := r.db.Create(&transaction).Error

	fmt.Println("SAVE CALLED")

	if err != nil {
		return transaction, err
	}

	return transaction, nil
}

func (r *repository) Last(t Transaction) (Transaction, error) {
	err := r.db.Last(&t).Error
	if err != nil {
		return t, err
	}
	return t, nil
}

func (r *repository) Update(t Transaction) (Transaction, error) {
	err := r.db.Save(&t).Error
	if err != nil {
		return t, err
	}
	return t, nil
}

func (r *repository) GetById(id int) (Transaction, error) {
	var t Transaction

	err := r.db.Where("id = ?", id).Find(&t).Error
	if err != nil {
		return t, err
	}
	return t, nil
}