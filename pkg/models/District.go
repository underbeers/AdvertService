package models

type District struct {
	Id       int    `json:"id" db:"id"`
	CityId   int    `json:"cityID" db:"city_id"`
	District string `json:"district" db:"district"`
}
