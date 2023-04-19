package models

import (
	"github.com/google/uuid"
	"time"
)

type FullAdvert struct {
	Id            int       `json:"id" db:"id"`
	PetCardId     int       `json:"petCardID" db:"pet_card_id"`
	UserId        uuid.UUID `json:"userID" db:"user_id"`
	Price         int       `json:"price" db:"price"`
	Description   string    `json:"description" db:"description"`
	CityId        int       `json:"cityID" db:"city_id"`
	City          string    `json:"city"`
	DistrictId    int       `json:"districtID" db:"district_id"`
	District      string    `json:"district"`
	Chat          bool      `json:"chat" db:"chat"`
	Phone         string    `json:"phone" db:"phone"`
	Status        string    `json:"status" db:"status"`
	Publication   time.Time `json:"publication" db:"publication"`
	PetName       string    `json:"petName" db:"pet_name"`
	PetTypeId     int       `json:"petTypeID" db:"pet_type_id"`
	PetTypeName   string    `json:"petType" db:"pet_type"`
	BreedId       int       `json:"breedID" db:"breed_id"`
	BreedName     string    `json:"breed" db:"breed_name"`
	Photo         string    `json:"photo" db:"photo"`
	BirthDate     time.Time `json:"birthDate" db:"birth_date"`
	Male          bool      `json:"male" db:"male"`
	Gender        string    `json:"gender" db:"gender"`
	Color         string    `json:"color" db:"color"`
	Care          string    `json:"care" db:"care"`
	Character     string    `json:"petCharacter" db:"pet_character"`
	Pedigree      string    `json:"pedigree" db:"pedigree"`
	Sterilization bool      `json:"sterilization" db:"sterilization"`
	Vaccinations  bool      `json:"vaccinations" db:"vaccinations"`
}
