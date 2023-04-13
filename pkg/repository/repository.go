package repository

import (
	"github.com/jmoiron/sqlx"
	"github.com/underbeers/AdvertService/pkg/models"
)

type AdvertPet interface {
	Create(advertPet models.AdvertPet) error
	GetAll(filter models.AdvertPetFilter) (advertPet []models.AdvertPet, total int64, err error)
	GetFullAdvert(id int) (advert models.FullAdvert, err error)
	ChangeStatus(id int, status string) error
	Delete(id int) error
	Update(id int, input models.UpdateAdvertInput) error
}

type Repository struct {
	AdvertPet
}

func NewRepository(db *sqlx.DB) *Repository {
	return &Repository{
		AdvertPet: NewAdvertPetPostgres(db),
	}
}
