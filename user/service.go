package user

import (
	"errors"

	"golang.org/x/crypto/bcrypt"
)

type Service interface {
	RegisterUser(input RegisterInput) (User, error)
	Login(input LoginInput) (User, error)
	SaveAvatar(id int, fileLocation string) (User, error)
	GetUserByID(id int) (User, error)
}

type service struct {
	repository Repository
}

func NewService(repository Repository) *service {
	return &service{repository}
}

func (s *service) RegisterUser(input RegisterInput) (User, error) {
	user := User{}

	user.Name = input.Name
	user.Email = input.Email
	user.Role = "user"
	user.Occupation = input.Occupation

	passwordHash, err := bcrypt.GenerateFromPassword([]byte(input.Password), bcrypt.MinCost)

	if err != nil {
		return user, err
	}

	user.PasswordHash = string(passwordHash)

	newUser, err := s.repository.Save(user)

	if err != nil {
		return user, err
	}

	return newUser, nil
}

func (s *service) Login(input LoginInput) (User, error) {
	email := input.Email
	password := input.Password

	userByEmail, err := s.repository.FindByEmail(email)

	if err != nil {
		return userByEmail, err
	}

	if userByEmail.ID == 0 {
		return userByEmail, errors.New("user not found")
	}

	err = bcrypt.CompareHashAndPassword([]byte(userByEmail.PasswordHash), []byte(password))

	if err != nil {
		return userByEmail, err
	}

	return userByEmail, nil
}

func (s *service) SaveAvatar(id int, fileLocation string) (User, error) {
	user, err := s.repository.FindByID(id)

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

func (s *service) GetUserByID(id int) (User, error) {
	user, err := s.repository.FindByID(id)

	if err != nil {
		return user, err
	}

	if user.ID == 0 {
		return user, errors.New("There is no user with that ID")
	}

	return user, nil
}
