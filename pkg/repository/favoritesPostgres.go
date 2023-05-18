package repository

import (
	"errors"
	"fmt"
	"github.com/google/uuid"
	"github.com/jmoiron/sqlx"
	"github.com/underbeers/AdvertService/pkg/models"
	"strings"
)

type FavoritesPostgres struct {
	db *sqlx.DB
}

func NewFavoritesPostgres(db *sqlx.DB) *FavoritesPostgres {
	return &FavoritesPostgres{db: db}
}

func (f FavoritesPostgres) Create(fav models.Favorites) error {
	tx, err := f.db.Begin()
	if err != nil {
		return err
	}

	createFavoritesQuery := fmt.Sprintf("INSERT INTO %s (%s, user_id) VALUES ($1, $2)", favorites, fav.Type)
	_, err = tx.Exec(createFavoritesQuery, fav.Id, fav.UserId)
	if err != nil {
		tx.Rollback()
		return errors.New(err.Error())
	}

	return tx.Commit()
}

func createGetAllFavoritesQuery(filter models.FavoritesFilter) string {

	query := fmt.Sprintf("SELECT COALESCE(fav.id, 0) as id, COALESCE(fav.advert_id, 0) as advert_id, "+
		"COALESCE(fav.organization_id, 0) as organization_id, COALESCE(fav.specialist_id, 0) as specialist_id, "+
		"COALESCE(fav.event_id, 0) as event_id, fav.user_id FROM %s fav ", favorites)

	if filter.FavoritesId != 0 || filter.UserId != uuid.Nil {

		query += "WHERE "
		setValues := make([]string, 0)

		if filter.FavoritesId != 0 {
			setValues = append(setValues, fmt.Sprintf("fav.id = %d", filter.FavoritesId))
		}

		if filter.UserId != uuid.Nil {
			setValues = append(setValues, fmt.Sprintf("fav.user_id = '%s'", filter.UserId.String()))
		}

		query += strings.Join(setValues, " AND ")

	}

	return query
}

func (f FavoritesPostgres) GetFavorites(filter models.FavoritesFilter) (fav []models.Favorites, err error) {

	query := createGetAllFavoritesQuery(filter)
	err = f.db.Select(&fav, query)

	return fav, err
}

func createFavoritesAdvertsQuery(filter models.FavoritesFilter) string {

	query := fmt.Sprintf("SELECT fav.id as favorites_id, ap.id, ap.pet_card_id, ap.user_id, ap.price, ap.description, ap.city_id, "+
		"ap.district_id, ap.chat, ap.phone, ap.status, ap.publication, pc.pet_name, pc.thumbnail_photo AS main_photo, "+
		"ct.city, ds.district, CASE pc.male WHEN True THEN 'Мальчик' WHEN False THEN 'Девочка' END AS gender, "+
		"pt.pet_type, br.breed_name, pc.birth_date FROM %s fav ",
		favorites)

	query += "INNER JOIN advert_pet ap ON ap.id = fav.advert_id INNER JOIN pet_card pc ON ap.pet_card_id = pc.id INNER JOIN pet_type pt ON pc.pet_type_id = pt.id " +
		"INNER JOIN breed br ON pc.breed_id = br.id INNER JOIN city ct ON ap.city_id = ct.id INNER JOIN district ds " +
		"ON ap.district_id = ds.id "

	if filter.FavoritesId != 0 || filter.UserId != uuid.Nil {

		query += "WHERE "
		setValues := make([]string, 0)

		if filter.FavoritesId != 0 {
			setValues = append(setValues, fmt.Sprintf("fav.id = %d", filter.FavoritesId))
		}

		if filter.UserId != uuid.Nil {
			setValues = append(setValues, fmt.Sprintf("fav.user_id = '%s'", filter.UserId.String()))
		}

		query += strings.Join(setValues, " AND ")

	}

	return query
}

func (f FavoritesPostgres) GetFavoritesAdverts(filter models.FavoritesFilter) (fav []models.AdvertPet, err error) {

	query := createFavoritesAdvertsQuery(filter)
	err = f.db.Select(&fav, query)

	return fav, err
}

func (f FavoritesPostgres) Delete(id int) error {
	query := fmt.Sprintf("DELETE FROM %s fav WHERE fav.id = $1",
		favorites)
	_, err := f.db.Exec(query, id)

	return err
}
