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

type Repository struct {
	AdvertPet
	Location
	Favorites
}

func NewRepository(db *sqlx.DB) *Repository {
	return &Repository{
		AdvertPet: NewAdvertPetPostgres(db),
		Location:  NewLocationPostgres(db),
		Favorites: NewFavoritesPostgres(db),
	}
}
