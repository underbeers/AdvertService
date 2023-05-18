package models

import "github.com/google/uuid"

type Favorites struct {
	Type           string    `json:"type" binding:"required"`
	Id             int       `json:"id" db:"id"`
	UserId         uuid.UUID `json:"userID" db:"user_id"`
	OrganizationId int       `json:"organizationID" db:"organization_id"`
	SpecialistId   int       `json:"specialistID" db:"specialist_id"`
	EventId        int       `json:"eventID" db:"event_id"`
	AdvertId       int       `json:"advertID" db:"advert_id"`
}
