package app

import (
	"pr_service/internal/repository"
)

type Service struct {
	userRepository *repository.UserRepository
}

func NewService(userRepository *repository.UserRepository) (*Service, error) {
	return &Service{
		userRepository: userRepository,
	}, nil
}
