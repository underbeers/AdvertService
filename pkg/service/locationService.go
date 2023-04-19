package service

import (
	"github.com/underbeers/AdvertService/pkg/models"
	"github.com/underbeers/AdvertService/pkg/repository"
)

type LocationService struct {
	repo repository.Location
}

func NewLocationService(repo repository.Location) *LocationService {
	return &LocationService{repo: repo}
}

func (l *LocationService) GetCities(filter models.CityFilter) (city []models.City, err error) {
	return l.repo.GetCities(filter)
}
func (l *LocationService) GetDistricts(filter models.DistrictFilter) (district []models.District, err error) {
	return l.repo.GetDistricts(filter)
}
