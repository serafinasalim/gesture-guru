package models

type User struct {
	Id        int    `gorm:"primaryKey" json:"id"`
	FirstName string `gorm:"varchar(50)" json:"firstName"`
	LastName  string `gorm:"varchar(50)" json:"lastName"`
	Email     string `gorm:"varchar(100)" json:"email"`
	Password  string `gorm:"varchar(100)" json:"password"`
	Profile   string `gorm:"varchar(100)" json:"profile"`
}
