package user

import "gorm.io/gorm"

type Repository interface {
	Save(user User) (User, error)
	Last(user User) (User, error)
}

type repository struct {
	db *gorm.DB
}

func NewRepository(db *gorm.DB) *repository {
	return &repository{db}
}

func (r *repository) Save(user User) (User, error) {
	err := r.db.Create(&user).Error
	if err != nil {
		return user, err
	}

	return user, nil
}

func (r *repository) Last(user User) (User, error) {
	err := r.db.Last(&user).Error
	if err != nil {
		return user, err
	}
	return user, nil
}