package service

import (
	"github.com/gin-gonic/gin"
	"github.com/underbeers/AdvertService/pkg/models"
	"github.com/underbeers/AdvertService/pkg/repository"
)

type AdvertPet interface {
	Create(advertPet models.AdvertPet) error
	GetAll(filter models.AdvertPetFilter) (advertPet []models.AdvertPet, total int64, err error)
	GetFullAdvert(id int) (advert models.FullAdvert, err error)
	ChangeStatus(id int, status string) error
	Delete(id int) error
	Update(id int, input models.UpdateAdvertInput) error
}

type Service struct {
	AdvertPet
	Config *repository.Config
	Router *gin.Engine
}

func NewService(repos *repository.Repository, cfg *repository.Config) *Service {
	return &Service{
		AdvertPet: NewAdvertPetService(repos.AdvertPet),
		Config:    cfg,
	}
}
