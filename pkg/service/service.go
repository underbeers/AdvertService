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

type Location interface {
	GetCities(filter models.CityFilter) (city []models.City, err error)
	GetDistricts(filter models.DistrictFilter) (district []models.District, err error)
}

type Favorites interface {
	Create(fav models.Favorites) error
	GetFavorites(filter models.FavoritesFilter) (fav []models.Favorites, err error)
	GetFavoritesAdverts(filter models.FavoritesFilter) (fav []models.AdvertPet, err error)
	Delete(id int) error
}

type Service struct {
	AdvertPet
	Location
	Favorites
	Config *repository.Config
	Router *gin.Engine
}

func NewService(repos *repository.Repository, cfg *repository.Config) *Service {
	return &Service{
		AdvertPet: NewAdvertPetService(repos.AdvertPet),
		Location:  NewLocationService(repos.Location),
		Favorites: NewFavoritesService(repos.Favorites),
		Config:    cfg,
	}
}
