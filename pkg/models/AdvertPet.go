package models

import (
	"github.com/google/uuid"
	"time"
)

type AdvertPet struct {
	Id          int       `json:"id" db:"id"`
	PetCardId   int       `json:"petCardID" db:"pet_card_id" binding:"required"`
	UserId      uuid.UUID `json:"userID" db:"user_id"`
	Price       int       `json:"price" db:"price" binding:"required"`
	Description string    `json:"description" db:"description" binding:"required"`
	CityId      int       `json:"cityID" db:"city_id" binding:"required"`
	City        string    `json:"city"`
	DistrictId  int       `json:"districtID" db:"district_id" binding:"required"`
	District    string    `json:"district"`
	Chat        bool      `json:"chat" db:"chat"`
	Phone       string    `json:"phone" db:"phone"`
	Status      string    `json:"status" db:"status"`
	Publication time.Time `json:"publication" db:"publication"`
	PetName     string    `json:"petName" db:"pet_name"`
	MainPhoto   string    `json:"mainPhoto" db:"main_photo"`
}
