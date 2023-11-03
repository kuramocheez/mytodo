package model

import (
	"github.com/sirupsen/logrus"
	"gorm.io/gorm"
)

type UsersInterface interface {
	Register(newUser Users) *Users
	Login(login Login) *Users
}

type Users struct {
	// ID        uint      `json:"id" form:"id" gorm:"type:primaryKey;autoIncrement:true"`
	gorm.Model
	Name     string `json:"name" form:"name" gorm:"type:varchar(255)"`
	Email    string `json:"email" form:"email" gorm:"type:varchar(255);uniqueIndex"`
	Password string `json:"password" form:"password" gorm:"type:varchar(255)"`
	// CreatedAt time.Time `json:"created_at" form:"created_at" gorm:"type:datetime"`
	// UpdatedAt time.Time `json:"updated_at" form:"updated_at" gorm:"type:datetime"`
	// DeletedAt time.Time `json:"deleted_at" form:"deleted_at" gorm:"type:datetime"`
}

type Login struct {
	Email    string `json:"email" form:"email" gorm:"type:varchar(255)"`
	Password string `json:"password" form:"password" gorm:"type:varchar(255)"`
}

type UsersModel struct {
	db *gorm.DB
}

func (um *UsersModel) InitUsers(db *gorm.DB) {
	um.db = db
}

func NewUsersModel(db *gorm.DB) UsersInterface {
	return &UsersModel{
		db: db,
	}
}

func (um *UsersModel) Register(newUser Users) *Users {
	if err := um.db.Create(&newUser).Error; err != nil {
		logrus.Error("Model: Error Saat Input Data User", err.Error())
		return nil
	}
	category := Category{
		Category: "Kegiatan Saya",
		Color:    "#3f48cc",
		UserID:   newUser.ID,
	}
	if err := um.db.Create(&category).Error; err != nil {
		logrus.Error("Model: Error Saat Input Data User", err.Error())
		return nil
	}
	return &newUser
}

func (um *UsersModel) Login(login Login) *Users {
	users := Users{}
	if err := um.db.Where("email = ? AND password = ?", login.Email, login.Password).First(&users).Error; err != nil {
		logrus.Error("Model: User Tidak Ditemukan ", err.Error())
		return nil
	}
	return &users
}
