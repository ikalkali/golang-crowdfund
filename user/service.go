package user

import (
	"errors"

	"golang.org/x/crypto/bcrypt"
)


type Service interface {
	RegisterUser(input RegisterUserInput) (User, error)
	Login(input LoginInput) (User, error)
	IsEmailAvailable(input CheckEmailInput) (bool, error)
	SaveAvatar(id int, fileLocation string) (User, error)
	GetUserById(id int) (User, error)
}

type service struct {
	repository Repository
}

func NewService(repository Repository) *service {
	return &service{repository}
}

func (s *service) RegisterUser(input RegisterUserInput) (User, error) {
	// get last user for auto increment
	var lastUser User
	lastUser, _ = s.repository.Last(lastUser)

	user := User{}
	user.Id = int(lastUser.Id) + 1
	user.Name = input.Name
	user.Email = input.Email
	user.Occupation = input.Occupation
	passwordHash, err := bcrypt.GenerateFromPassword([]byte(input.Password), bcrypt.MinCost)
	if err != nil {
		return user, err
	}
	user.PasswordHash = string(passwordHash)
	user.Role = "user"

	newUser, err := s.repository.Save(user)
	if err != nil {
		return newUser, err
	}

	return newUser, nil

}

func (s *service) Login(input LoginInput) (User, error) {
	email := input.Email
	password := input.Password

	user,err := s.repository.FindByEmail(email)
	if err != nil {
		return user, err
	}

	if user.Id == 0 {
		return user, errors.New("User tidak ditemukan")
	}

	err = bcrypt.CompareHashAndPassword([]byte(user.PasswordHash), []byte(password))
	if err != nil {
		return user, err
	}
	return user, nil
}

func (s *service) IsEmailAvailable(input CheckEmailInput) (bool, error) {
	email := input.Email

	user, err := s.repository.FindByEmail(email)
	if err != nil {
		return false, err
	}
	if user.Id == 0 {
		return true, nil
	}
	return false, nil
}

func (s *service) SaveAvatar(id int, fileLocation string) (User, error) {
	user, err := s.repository.FindById(id)
	if err != nil {
		return user, err
	}
	
	user.AvatarFileName = fileLocation

	updatedUser, err := s.repository.Update(user)
	if err != nil {
		return updatedUser, err
	}

	return updatedUser, nil
}

func (s *service) GetUserById(id int) (User,error) {
	user, err := s.repository.FindById(id)
	if err != nil {
		return user, err
	}

	if user.Id == 0 {
		return user, errors.New("no user found with that ID")
	}

	return user, nil
}