package models

import "github.com/google/uuid"

type AdvertPetFilter struct {
	AdvertPetId     int
	PetCardId       int
	UserId          uuid.UUID
	MinPrice        int
	MaxPrice        int
	CityId          int
	DistrictId      int
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

type CityFilter struct {
	CityId int
}

type DistrictFilter struct {
	DistrictId int
	CityId     int
}
