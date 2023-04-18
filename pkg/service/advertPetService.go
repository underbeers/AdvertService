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

func (a *AdvertPetService) Create(advertPet models.AdvertPet) error {
	return a.repo.Create(advertPet)
}

func (a *AdvertPetService) GetAll(filter models.AdvertPetFilter) (advertPet []models.AdvertPet, total int64, err error) {
	return a.repo.GetAll(filter)
}

func (a *AdvertPetService) GetFullAdvert(id int) (advert models.FullAdvert, err error) {
	return a.repo.GetFullAdvert(id)
}

func (a *AdvertPetService) ChangeStatus(id int, status string) error {
	return a.repo.ChangeStatus(id, status)
}

func (a *AdvertPetService) Delete(id int) error {
	return a.repo.Delete(id)
}

func (a *AdvertPetService) Update(id int, input models.UpdateAdvertInput) error {
	if err := input.Validate(); err != nil {
		return err
	}
	return a.repo.Update(id, input)
}
