package models

type City struct {
	Id   int    `json:"id" db:"id"`
	City string `json:"city" db:"city"`
}
