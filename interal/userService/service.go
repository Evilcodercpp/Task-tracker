package userservice

import (
	"errors"
	"log"
)

type UserService struct {
	repo UserRepository
}

func NewUserService(repo UserRepository) *UserService {
	return &UserService{repo: repo}
}

func (s *UserService) CreateUser(user *User) error {
	if user.Email == "" {
		return errors.New("email cannot be empty")
	}
	if user.Password == "" {
		return errors.New("password cannot be empty")
	}
	return s.repo.CreateUser(user)
}

func (s *UserService) GetAllUsers() ([]User, error) {
	users, err := s.repo.GetAllUsers()
	if err != nil {
		return nil, err
	}
	if len(users) == 0 {
		log.Println("no users found")
	}
	return users, nil
}

func (s *UserService) GetUserByID(id string) (*User, error) {
	if id == "" {
		return nil, errors.New("id cannot be empty")
	}
	return s.repo.GetUserByID(id)
}

func (s *UserService) UpdateUser(user *User) error {
	if user.Email == "" {
		return errors.New("email cannot be empty")
	}
	existing, err := s.repo.GetUserByID(user.ID)
	if err != nil {
		return err
	}
	existing.Email = user.Email
	if user.Password != "" {
		existing.Password = user.Password
	}
	return s.repo.UpdateUser(existing)
}

func (s *UserService) DeleteUser(id string) error {
	if id == "" {
		return errors.New("id cannot be empty")
	}
	err := s.repo.DeleteUser(id)
	if err != nil {
		log.Printf("failed to delete user %s: %v", id, err)
		return err
	}
	log.Printf("user %s deleted successfully", id)
	return nil
}
