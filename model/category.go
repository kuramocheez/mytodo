package model

import (
	"time"

	"github.com/sirupsen/logrus"
	"gorm.io/gorm"
)

type CategoryInterface interface {
	AddCategory(newCategory Category) bool
	GetCategories(page int, perpage int, id uint) []Category
	GetCategory(id int, idUser uint) *Category
	UpdateCategory(category Category, id int, idUser uint) bool
	DeleteCategory(id int, idUser uint) bool
}

type Category struct {
	ID        uint      `json:"id" form:"id" gorm:"type:primaryKey;autoIncrement:true"`
	Category  string    `json:"category" form:"category" gorm:"type:varchar(255)"`
	Color     string    `json:"color" form:"color" gorm:"type:varchar(255)"`
	CreatedAt time.Time `json:"created_at" form:"created_at" gorm:"type:datetime"`
	UpdatedAt time.Time `json:"updated_at" form:"updated_at" gorm:"type:datetime"`
	DeletedAt time.Time `json:"deleted_at" form:"deleted_at" gorm:"type:datetime"`
	UserID    uint      `json:"user_id" form:"user_id"`
	User      Users     `json:"user" form:"user"`
}

type CategoryModel struct {
	db *gorm.DB
}

func (cm *CategoryModel) InitCategory(db *gorm.DB) {
	cm.db = db
}

func NewCategoryModel(db *gorm.DB) CategoryInterface {
	return &CategoryModel{
		db: db,
	}
}

func (cm *CategoryModel) AddCategory(newCategory Category) bool {
	if err := cm.db.Create(&newCategory).Error; err != nil {
		logrus.Error("Model: Error Saat Input Category ")
		return false
	}
	return true
}

func (cm *CategoryModel) GetCategories(page, perpage int, id uint) []Category {
	categories := []Category{}
	offset := (page - 1) * perpage
	if err := cm.db.Limit(perpage).Offset(offset).Where("user_id = ?", id).Find(&categories).Error; err != nil {
		logrus.Error("Model: Error Mendapatkan Data Category ", err.Error())
		return nil
	}
	for i := 0; i < len(categories); i++ {
		user := Users{}
		if err := cm.db.First(&user, categories[i].UserID).Error; err != nil {
			logrus.Error("Model: Error Mendapatkan User Data Category ", err.Error())
			return nil
		}
		categories[i].User = user
	}
	return categories
}
func (cm *CategoryModel) GetCategory(id int, idUser uint) *Category {
	category := Category{}
	if err := cm.db.Where("user_id = ?", idUser).First(&category, id).Error; err != nil {
		logrus.Error("Model: Data Category Tidak Ditemukan ", err.Error())
		return nil
	}
	user := Users{}
	if err := cm.db.First(&user, idUser).Error; err != nil {
		logrus.Error("Model: Data User Category Tidak Ditemukan ", err.Error())
		return nil
	}
	category.User = user
	return &category
}
func (cm *CategoryModel) UpdateCategory(categoryUp Category, id int, idUser uint) bool {
	data := cm.GetCategory(id, idUser)
	if data == nil {
		logrus.Error("Model: Error Update Data Category")
		return false
	}
	data.Category = categoryUp.Category
	data.Color = categoryUp.Color
	if err := cm.db.Save(&data).Error; err != nil {
		logrus.Error("Model: Error Update Data Category ", err.Error())
		return false
	}
	return true
}
func (cm *CategoryModel) DeleteCategory(id int, idUser uint) bool {
	category := Category{}
	data := cm.GetCategory(id, idUser)
	if data == nil {
		logrus.Error("Model: Error Delete Category")
		return false
	}
	if err := cm.db.Where("user_id = ?", idUser).Delete(&category, id).Error; err != nil {
		logrus.Error("Model: Error Delete Category", err.Error())
		return false
	}
	return true
}
