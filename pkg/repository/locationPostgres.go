package repository

import (
	"fmt"
	"github.com/jmoiron/sqlx"
	"github.com/underbeers/AdvertService/pkg/models"
	"strings"
)

type LocationPostgres struct {
	db *sqlx.DB
}

func NewLocationPostgres(db *sqlx.DB) *LocationPostgres {
	return &LocationPostgres{db: db}
}

func createCityQuery(filter models.CityFilter) string {
	query := fmt.Sprintf(`SELECT id, city FROM city `)
	if filter.CityId != 0 {
		query += fmt.Sprintf(`WHERE id = %d`, filter.CityId)
	}
	return query
}

func (l LocationPostgres) GetCities(filter models.CityFilter) (city []models.City, err error) {

	query := createCityQuery(filter)
	err = l.db.Select(&city, query)

	return city, err
}

func createDistrictQuery(filter models.DistrictFilter) string {
	query := fmt.Sprintf(`SELECT id, city_id, district FROM district `)

	if filter.DistrictId != 0 || filter.CityId != 0 {

		query += "WHERE "
		setValues := make([]string, 0)

		if filter.DistrictId != 0 {
			setValues = append(setValues, fmt.Sprintf("id = %d", filter.DistrictId))
		}

		if filter.CityId != 0 {
			setValues = append(setValues, fmt.Sprintf("city_id = %d", filter.CityId))
		}

		query += strings.Join(setValues, " AND ")
	}
	return query
}

func (l LocationPostgres) GetDistricts(filter models.DistrictFilter) (district []models.District, err error) {

	query := createDistrictQuery(filter)
	err = l.db.Select(&district, query)

	return district, err
}
