package service

import "itfest-2025/internal/repository"

type Service struct {
}

func NewService(repository *repository.Repository) *Service {
	return &Service{}
}
