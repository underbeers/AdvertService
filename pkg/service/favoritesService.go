package service

import (
	"github.com/underbeers/AdvertService/pkg/models"
	"github.com/underbeers/AdvertService/pkg/repository"
)

type FavoritesService struct {
	repo repository.Favorites
}

func NewFavoritesService(repo repository.Favorites) *FavoritesService {
	return &FavoritesService{repo: repo}
}

func (f *FavoritesService) Create(fav models.Favorites) error {
	return f.repo.Create(fav)
}

func (f *FavoritesService) GetFavorites(filter models.FavoritesFilter) (fav []models.Favorites, err error) {
	return f.repo.GetFavorites(filter)
}

func (f *FavoritesService) GetFavoritesAdverts(filter models.FavoritesFilter) (fav []models.AdvertPet, err error) {
	return f.repo.GetFavoritesAdverts(filter)
}

func (f *FavoritesService) Delete(id int) error {
	return f.repo.Delete(id)
}
