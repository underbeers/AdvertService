package models

import "github.com/google/uuid"

type AdvertPetFilter struct {
	AdvertPetId int
	UserId      uuid.UUID
	MinPrice    int
	MaxPrice    int
	Region      string
	Locality    string
	Status      string
}
