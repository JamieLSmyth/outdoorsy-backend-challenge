package model

type User struct {
	Id int `gorm:"primaryKey" json:"id"`
	FirstName string `json:"first_name"`
	LastName string `json:"last_name"`
}