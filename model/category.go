package model

import (
	"github.com/sirupsen/logrus"
	"gorm.io/gorm"
)

type CategoryInterface interface {
	AddCategory(newCategory Category) bool
	GetCategories() []Category
	GetCategory(id int) *Category
	UpdateCategory(category Category, id int) bool
	DeleteCategory(id int) bool
}

type Category struct {
	gorm.Model
	Category string `json:"category" form:"category" gorm:"type:varchar(255)"`
	Color    string `json:"color" form:"color" gorm:"type:varchar(255)"`
	UserID   uint
	User     Users
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
		logrus.Error("Model: Error Saat Input Category")
		return false
	}
	return true
}

func (cm *CategoryModel) GetCategories() []Category {
	categories := []Category{}
	if err := cm.db.Find(&categories).Error; err != nil {
		logrus.Error("Model: Error saat mengambil data buku", err.Error())
		return nil
	}
	return categories
}
func (cm *CategoryModel) GetCategory(id int) *Category {
	category := Category{}
	if err := cm.db.First(&category, id).Error; err != nil {
		logrus.Error("Model: Data Category Tidak Ditemukan", err.Error())
		return nil
	}
	return &category
}
func (cm *CategoryModel) UpdateCategory(category Category, id int) bool {
	// category := Category{}
	// data := cm.GetCategory(id)
	// if data != nil{
	// 	logrus.Error("Model: Data Category Tidak Ditemukan", err.Error())
	// 	return false
	// }

	if err := cm.db.Save(&category).Error; err != nil {
		logrus.Error("Model: Gagal mengubah data category", err.Error())
		return false
	}
	return true
}
func (cm *CategoryModel) DeleteCategory(id int) bool {
	data := cm.GetCategory(id)
	if data != nil {
		logrus.Error("Model: Data Category Tidak Ditemukan")
		return false
	}
	if err := cm.db.Delete(&data).Error; err != nil {
		logrus.Error("Model: Data Category Tidak dapat dihapus", err.Error())
		return false
	}
	return true
}
