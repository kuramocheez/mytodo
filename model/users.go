package model

import(
	"github.com/sirupsen/logrus"
	"gorm.io/gorm"
)
type UsersInterface interface{
	Register(newUser Users) bool
	Login(email, password string) *Users
}

type Users struct{
	gorm.Model
	Name string `json:"name" form:"name" gorm:"type:varchar(255)"`
	Email string `json:"email" form:"email" gorm:"type:varchar(255)"`
	Password string `json:"password" form:"password" gorm:"type:varchar(255)"`
}

type Login struct{
	Email string `json:"email" form:"email" gorm:"type:varchar(255)"`
	Password string `json:"password" form:"password" gorm:"type:varchar(255)"`
}

type UsersModel struct{
	db *gorm.DB
}

func (um *UsersModel) InitUsers(db *gorm.DB){
	um.db = db
}

func NewUsersModel(db *gorm.DB)
UsersInterface{
	return &UsersModel{
		db: db,
	}
}

func (um *UsersModel) Register(newUser Users) bool{
	if err := um.db.Create(&newUser).Error; err != nil{
		logrus.Error("Model: Error Saat Input Data User", err.Error())
		return false
	
	}
	return true
}

func (um *UsersModel) Login(email, password string) *Users{
	users := Users{}
	if err := um.db.Where("email = ? AND password = ?", mail, password).First(&users).Error; err != nil{
		logrus.Error("Model: User Tidak Ditemukan", err.Error())
	}
	return &users
}