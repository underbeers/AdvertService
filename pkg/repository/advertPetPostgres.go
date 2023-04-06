package repository

import (
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

	createAdvertPetQuery := fmt.Sprintf("INSERT INTO %s (pet_card_id, user_id, price, description, region,"+
		"locality, chat, phone, status, publication) VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10)", advertPetTable)
	_, err = tx.Exec(createAdvertPetQuery, advertPet.PetCardId, advertPet.UserId, advertPet.Price, advertPet.Description,
		advertPet.Region, advertPet.Locality, advertPet.Chat, advertPet.Phone, advertPet.Status, time.Now())
	if err != nil {
		err := tx.Rollback()
		if err != nil {
			return err
		}
		return err
	}

	return tx.Commit()
}

func createAdvertPetQuery(filter models.AdvertPetFilter) string {

	query := fmt.Sprintf("SELECT ap.id, ap.pet_card_id, ap.user_id, ap.price, ap.description, ap.region, "+
		"ap.locality, ap.chat, ap.phone, ap.status, ap.publication FROM %s ap ",
		advertPetTable)

	if filter.AdvertPetId != 0 || filter.UserId != uuid.Nil || filter.Region != "" || filter.Locality != "" ||
		filter.Status != "" || filter.MinPrice != 0 || filter.MaxPrice != 0 {

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

		if filter.Region != "" {
			setValues = append(setValues, fmt.Sprintf("ap.region = '%s'", filter.Region))
		}

		if filter.Locality != "" {
			setValues = append(setValues, fmt.Sprintf("ap.locality = '%s'", filter.Locality))
		}

		if filter.Status != "" {
			setValues = append(setValues, fmt.Sprintf("ap.status = '%s'", filter.Status))
		}

		query += strings.Join(setValues, " AND ")
	}

	return query
}

func (a AdvertPetPostgres) GetAll(filter models.AdvertPetFilter) ([]models.AdvertPet, error) {
	var lists []models.AdvertPet

	query := createAdvertPetQuery(filter)
	err := a.db.Select(&lists, query)
	return lists, err
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

	if input.UserId != nil {
		setValues = append(setValues, fmt.Sprintf("user_id=$%d", argId))
		args = append(args, *input.UserId)
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

	if input.Region != nil {
		setValues = append(setValues, fmt.Sprintf("region=$%d", argId))
		args = append(args, *input.Region)
		argId++
	}

	if input.Locality != nil {
		setValues = append(setValues, fmt.Sprintf("locality=$%d", argId))
		args = append(args, *input.Locality)
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
