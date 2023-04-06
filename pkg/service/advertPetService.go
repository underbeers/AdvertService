package service

import (
	"github.com/underbeers/AdvertService/pkg/models"
	"github.com/underbeers/AdvertService/pkg/repository"
)

type AdvertPetService struct {
	repo repository.AdvertPet
}

func NewAdvertPetService(repo repository.AdvertPet) *AdvertPetService {
	return &AdvertPetService{repo: repo}
}

func (s *AdvertPetService) Create(advertPet models.AdvertPet) error {
	return s.repo.Create(advertPet)
}

func (s *AdvertPetService) GetAll(filter models.AdvertPetFilter) ([]models.AdvertPet, error) {
	return s.repo.GetAll(filter)
}

func (s *AdvertPetService) ChangeStatus(id int, status string) error {
	return s.repo.ChangeStatus(id, status)
}

func (s *AdvertPetService) Delete(id int) error {
	return s.repo.Delete(id)
}

func (s *AdvertPetService) Update(id int, input models.UpdateAdvertInput) error {
	if err := input.Validate(); err != nil {
		return err
	}
	return s.repo.Update(id, input)
}
