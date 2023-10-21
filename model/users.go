package model

import(
	"github.com/sirupsen/logrus"
	"gorm.io/gorm"
)

type Users struct{
	gorm.Model
	Name string `json:"name" form:"name" gorm:"type:varchar(255)"`
	Email string `json:"email" form:"email" gorm:"type:varchar(255)"`
	Password string `json:"password" form:"password" gorm:"type:varchar(255)"`
}

type UsersModel struct{
	db *gorm.DB
}

func (um *UsersModel) InitUsers(db *gorm.DB){
	um.db = db
}