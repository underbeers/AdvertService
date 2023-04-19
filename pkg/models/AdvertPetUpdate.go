package models

import (
	"errors"
	"github.com/google/uuid"
	"time"
)

type UpdateAdvertInput struct {
	Id          *int       `json:"id" db:"id"`
	PetCardId   *int       `json:"petCardID" db:"pet_card_id"`
	UserId      *uuid.UUID `json:"userID" db:"user_id"`
	Price       *int       `json:"price" db:"price"`
	Description *string    `json:"description" db:"description"`
	CityId      *int       `json:"cityID" db:"city_id"`
	DistrictId  *int       `json:"districtID" db:"district_id"`
	Chat        *bool      `json:"chat" db:"chat"`
	Phone       *string    `json:"phone" db:"phone"`
	Status      *string    `json:"status" db:"status"`
	Publication *time.Time `json:"publication" db:"publication"`
}

func (i UpdateAdvertInput) Validate() error {
	if i.Id == nil && i.PetCardId == nil && i.Price == nil && i.Description == nil &&
		i.CityId == nil && i.DistrictId == nil && i.Chat == nil && i.Phone == nil && i.Status == nil &&
		i.Publication == nil {
		return errors.New("update structure has no values")
	}
	return nil
}
