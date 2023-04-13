package models

import "github.com/google/uuid"

type AdvertPetFilter struct {
	AdvertPetId     int
	PetCardId       int
	UserId          uuid.UUID
	MinPrice        int
	MaxPrice        int
	Region          string
	Locality        string
	Status          string
	MinPriceSort    bool
	MaxPriceSort    bool
	PublicationSort bool
	Page            int
	PerPage         int
	PetTypeId       int
	BreedId         int
	Gender          string
}
