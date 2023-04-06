package repository

import (
	"fmt"
	"github.com/jmoiron/sqlx"
	"github.com/underbeers/AdvertService/pkg/models"
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

	if filter.AdvertPetId != 0 {
		query += fmt.Sprintf("WHERE ap.id = %d", filter.AdvertPetId)
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
	//TODO implement me
	panic("implement me")
}

func NewAdvertPetPostgres(db *sqlx.DB) *AdvertPetPostgres {
	return &AdvertPetPostgres{db: db}
}
