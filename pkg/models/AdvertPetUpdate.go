package models

import (
	"errors"
	"github.com/google/uuid"
	"time"
)

type UpdateAdvertInput struct {
	Id          *int       `json:"id" db:"id"`
	PetCardId   *int       `json:"petCardID" db:"pet_card_id" binding:"required"`
	UserId      *uuid.UUID `json:"userID" db:"user_id" binding:"required"`
	Price       *int       `json:"price" db:"price" binding:"required"`
	Description *string    `json:"description" db:"description"`
	Region      *string    `json:"region" db:"region"`
	Locality    *string    `json:"locality" db:"locality"`
	Chat        *bool      `json:"chat" db:"chat"`
	Phone       *string    `json:"phone" db:"phone"`
	Status      *string    `json:"status" db:"status"`
	Publication *time.Time `json:"publication" db:"publication"`
}

func (i UpdateAdvertInput) Validate() error {
	if i.Id == nil && i.PetCardId == nil && i.UserId == nil && i.Price == nil && i.Description == nil &&
		i.Region == nil && i.Locality == nil && i.Chat == nil && i.Phone == nil && i.Status == nil &&
		i.Publication == nil {
		return errors.New("update structure has no values")
	}
	return nil
}
