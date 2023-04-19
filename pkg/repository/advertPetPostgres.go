package repository

import (
	"errors"
	"fmt"
	"github.com/google/uuid"
	"github.com/jmoiron/sqlx"
	"github.com/sirupsen/logrus"
	"github.com/underbeers/AdvertService/pkg/models"
	"strings"
	"time"
)

type AdvertPetPostgres struct {
	db *sqlx.DB
}

func (a AdvertPetPostgres) Create(advertPet models.AdvertPet) error {
	tx, err := a.db.Begin()
	if err != nil {
		return err
	}

	createAdvertPetQuery := fmt.Sprintf("INSERT INTO %s (pet_card_id, user_id, price, description, city_id,"+
		"district_id, chat, phone, status, publication) VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10)", advertPetTable)
	_, err = tx.Exec(createAdvertPetQuery, advertPet.PetCardId, advertPet.UserId, advertPet.Price, advertPet.Description,
		advertPet.CityId, advertPet.DistrictId, advertPet.Chat, advertPet.Phone, advertPet.Status, time.Now())
	if err != nil {
		tx.Rollback()
		return errors.New(err.Error())
	}

	var advert models.AdvertPet
	query := fmt.Sprintf("SELECT pc.user_id FROM pet_card pc WHERE pc.id = %d", advertPet.PetCardId)
	err = a.db.Get(&advert, query)

	if advert.UserId != advertPet.UserId {
		tx.Rollback()
		return errors.New(err.Error())
	}

	return tx.Commit()
}

func createAdvertPetQuery(query string, filter models.AdvertPetFilter) string {

	query += "INNER JOIN pet_card pc ON ap.pet_card_id = pc.id INNER JOIN pet_type pt ON pc.pet_type_id = pt.id " +
		"INNER JOIN breed br ON pc.breed_id = br.id INNER JOIN city ct ON ap.city_id = ct.id INNER JOIN district ds " +
		"ON ap.district_id = ds.id "

	if filter.AdvertPetId != 0 || filter.UserId != uuid.Nil || filter.CityId != 0 || filter.DistrictId != 0 ||
		filter.Status != "" || filter.MinPrice != 0 || filter.MaxPrice != 0 || filter.BreedId != 0 ||
		filter.PetTypeId != 0 || filter.Gender != "" || filter.PetCardId != 0 {

		query += "WHERE "
		setValues := make([]string, 0)

		if filter.AdvertPetId != 0 {
			setValues = append(setValues, fmt.Sprintf("ap.id = %d", filter.AdvertPetId))
		}

		if filter.UserId != uuid.Nil {
			setValues = append(setValues, fmt.Sprintf("ap.user_id = '%s'", filter.UserId.String()))
		}

		if filter.MinPrice != 0 {
			setValues = append(setValues, fmt.Sprintf("ap.price >= %d", filter.MinPrice))
		}

		if filter.MaxPrice != 0 {
			setValues = append(setValues, fmt.Sprintf("ap.price <= %d", filter.MaxPrice))
		}

		if filter.CityId != 0 {
			setValues = append(setValues, fmt.Sprintf("ap.city_id = %d", filter.CityId))
		}

		if filter.DistrictId != 0 {
			setValues = append(setValues, fmt.Sprintf("ap.district_id = %d", filter.DistrictId))
		}

		if filter.Status != "" {
			setValues = append(setValues, fmt.Sprintf("ap.status = '%s'", filter.Status))
		}

		if filter.PetCardId != 0 {
			setValues = append(setValues, fmt.Sprintf("ap.pet_card_id = %d", filter.PetCardId))
		}

		if filter.PetTypeId != 0 {
			setValues = append(setValues, fmt.Sprintf("pc.pet_type_id = %d", filter.PetTypeId))
		}

		if filter.BreedId != 0 {
			setValues = append(setValues, fmt.Sprintf("pc.breed_id = %d", filter.BreedId))
		}

		if filter.Gender != "" {
			if filter.Gender == "male" {
				setValues = append(setValues, fmt.Sprintf("pc.male = True"))
			} else if filter.Gender == "female" {
				setValues = append(setValues, fmt.Sprintf("pc.male = False"))
			}
		}

		query += strings.Join(setValues, " AND ")
	}

	return query
}

func (a AdvertPetPostgres) GetAll(filter models.AdvertPetFilter) (advertPet []models.AdvertPet, total int64, err error) {
	var lists []models.AdvertPet

	countQuery := fmt.Sprintf("SELECT count(*) FROM %s ap ", advertPetTable)
	err = a.db.QueryRow(createAdvertPetQuery(countQuery, filter)).Scan(&total)

	query := fmt.Sprintf("SELECT ap.id, ap.pet_card_id, ap.user_id, ap.price, ap.description, ap.city_id, "+
		"ap.district_id, ap.chat, ap.phone, ap.status, ap.publication, pc.pet_name, pc.photo AS main_photo, "+
		"ct.city, ds.district FROM %s ap ",
		advertPetTable)
	query = createAdvertPetQuery(query, filter)

	if filter.PublicationSort != false {
		query += " ORDER BY ap.publication DESC"
	} else if filter.MinPriceSort != false {
		query += " ORDER BY ap.price"
	} else if filter.MaxPriceSort != false {
		query += " ORDER BY ap.price DESC"
	}

	query += fmt.Sprintf(" LIMIT %d OFFSET %d", filter.PerPage, (filter.Page-1)*filter.PerPage)
	err = a.db.Select(&lists, query)

	return lists, total, err
}

func (a AdvertPetPostgres) GetFullAdvert(id int) (advert models.FullAdvert, err error) {

	query := fmt.Sprintf("SELECT ap.id, ap.pet_card_id, ap.user_id, ap.price, ap.description, ap.city_id, "+
		"ap.district_id, ap.chat, ap.phone, ap.status, ap.publication, pc.pet_name, pc.pet_type_id, pt.pet_type, "+
		"pc.breed_id, br.breed_name, pc.photo, pc.birth_date, pc.male, CASE pc.male WHEN True THEN 'Мальчик' WHEN "+
		"False THEN 'Девочка' END AS gender, pc.color, pc.care, pc.pet_character, pc.pedigree, pc.sterilization, "+
		"pc.vaccinations, ct.city, ds.district FROM %s ap INNER JOIN pet_card pc ON ap.pet_card_id = pc.id INNER "+
		"JOIN pet_type pt ON pc.pet_type_id = pt.id INNER JOIN breed br ON pc.breed_id = br.id INNER JOIN city ct ON "+
		"ap.city_id = ct.id INNER JOIN district ds ON ap.district_id = ds.id WHERE ap.id = %d",
		advertPetTable, id)
	err = a.db.Get(&advert, query)

	return advert, err
}

func (a AdvertPetPostgres) ChangeStatus(id int, status string) error {

	query := fmt.Sprintf("UPDATE %s ap SET status = '%s' WHERE ap.id = %d",
		advertPetTable, status, id)
	_, err := a.db.Exec(query)

	return err
}

func (a AdvertPetPostgres) Delete(id int) error {
	query := fmt.Sprintf("DELETE FROM %s ap WHERE ap.id = $1",
		advertPetTable)
	_, err := a.db.Exec(query, id)

	return err
}

func (a AdvertPetPostgres) Update(id int, input models.UpdateAdvertInput) error {
	setValues := make([]string, 0)
	args := make([]interface{}, 0)
	argId := 1

	if input.PetCardId != nil {
		setValues = append(setValues, fmt.Sprintf("pet_card_id=$%d", argId))
		args = append(args, *input.PetCardId)
		argId++
	}

	if input.Price != nil {
		setValues = append(setValues, fmt.Sprintf("price=$%d", argId))
		args = append(args, *input.Price)
		argId++
	}

	if input.Description != nil {
		setValues = append(setValues, fmt.Sprintf("description=$%d", argId))
		args = append(args, *input.Description)
		argId++
	}

	if input.CityId != nil {
		setValues = append(setValues, fmt.Sprintf("city_id=$%d", argId))
		args = append(args, *input.CityId)
		argId++
	}

	if input.DistrictId != nil {
		setValues = append(setValues, fmt.Sprintf("district_id=$%d", argId))
		args = append(args, *input.DistrictId)
		argId++
	}

	if input.Chat != nil {
		setValues = append(setValues, fmt.Sprintf("chat=$%d", argId))
		args = append(args, *input.Chat)
		argId++
	}

	if input.Phone != nil {
		setValues = append(setValues, fmt.Sprintf("phone=$%d", argId))
		args = append(args, *input.Phone)
		argId++
	}

	if input.Status != nil {
		setValues = append(setValues, fmt.Sprintf("status=$%d", argId))
		args = append(args, *input.Status)
		argId++
	}

	setValues = append(setValues, fmt.Sprintf("publication=$%d", argId))
	args = append(args, time.Now())
	argId++

	setQuery := strings.Join(setValues, ", ")
	query := fmt.Sprintf("UPDATE %s ap SET %s WHERE ap.id = $%d",
		advertPetTable, setQuery, argId)
	args = append(args, id)

	logrus.Debugf("updateQuery: %s", query)
	logrus.Debugf("args: %s", args)

	_, err := a.db.Exec(query, args...)
	return err

}

func NewAdvertPetPostgres(db *sqlx.DB) *AdvertPetPostgres {
	return &AdvertPetPostgres{db: db}
}
