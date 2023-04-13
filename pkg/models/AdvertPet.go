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
	Description string    `json:"description" db:"description"`
	Region      string    `json:"region" db:"region"`
	Locality    string    `json:"locality" db:"locality"`
	Chat        bool      `json:"chat" db:"chat"`
	Phone       string    `json:"phone" db:"phone"`
	Status      string    `json:"status" db:"status"`
	Publication time.Time `json:"publication" db:"publication"`
	PetName     string    `json:"petName" db:"pet_name"`
	MainPhoto   string    `json:"mainPhoto" db:"main_photo"`
}
