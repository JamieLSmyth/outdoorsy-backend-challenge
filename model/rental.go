package model

type Price struct {
	Day int `json:"day"`
}

type Location struct {
	City string `json:"city"`
	State string `json:"state"`
	Zip string  `json:"zip"`
	Country string  `json:"country"`
	Lat float64  `json:"lat"`
	Lng float64 `json:"lng"`
}

type Rental struct{
	Id int `json:"id"`
	Name string `json:"name"`
	Description string `json:"description"`
	Type string	`json:"type"`
	Make string	`json:"make"`
	Model string `json:"model"`
	Year int `json:"year"`
	Length float64 `json:"length"`
	Sleeps int  `json:"sleeps"`
	PrimaryImageUrl string  `json:"primary_image_url"`
	Price Price `gorm:"embedded" json:"price"`
	Location Location `gorm:"embedded" json:"location"`
	UserId int `json:"-"` //TODO fiqure out how to tell GROM to not worry about this field and still allow the join
	User User `json:"user"`
}
