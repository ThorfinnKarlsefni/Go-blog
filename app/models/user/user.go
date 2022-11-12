package user

import (
	"goblog/app/models"
)

type User struct {
	models.BaseModel

	Name     string `gorm:"cloumn:name;type:varchar(255);not null;unique"`
	Email    string `gorm:"cloumn:email;type:varchar(255);default:NULL;unique"`
	Password string `gorm:"cloumn:password;type:varchar(255)"`
}
