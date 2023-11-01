package model

import (
	"time"

	"github.com/sirupsen/logrus"
	"gorm.io/gorm"
)

type TodoInterface interface {
	AddTodo(newTodo Todo) bool
	GetTodos(page, content int, userID uint) []Todo
	GetTodo(id int, userID uint) *Todo
	GetTodoByStatus(page, content int, userID uint, status string) []Todo
	GetTodoByDate(page, content int, userID uint, date string) []Todo
	UpdateTodo(id int, userID uint, todo Todo) bool
	UpdateTodoStatus(id int, UserID uint, status string) bool
	DeleteTodo(id int, userID uint) bool
}

type Todo struct {
	gorm.Model
	Memo     string    `json:"memo" form:"memo" gorm:"type:varchar(255)"`
	DateTime time.Time `json:"date_time" form:"date_time" gorm:"datetime"`
	// Filename   string
	Status     string
	CategoryID uint `json:"category_id" form:"category_id"`
	Category   Category
	UserID     uint
	User       Users
}

type TodoModel struct {
	db *gorm.DB
}

func (tm *TodoModel) InitTodo(db *gorm.DB) {
	tm.db = db
}

func NewTodoModel(db *gorm.DB) TodoInterface {
	return &TodoModel{
		db: db,
	}
}

func (tm *TodoModel) AddTodo(newTodo Todo) bool {
	if err := tm.db.Create(&newTodo).Error; err != nil {
		logrus.Error("Model: Error Saat Input Todo")
		return false
	}
	return true
}

func (tm *TodoModel) GetTodos(page, content int, userID uint) []Todo {
	todo := []Todo{}
	offset := (page - 1) * content
	if err := tm.db.Limit(content).Offset(offset).Where("user_id = ?", userID).Find(&todo).Error; err != nil {
		logrus.Error("Model: Error Mendapatkan Data Todo ", err.Error())
		return nil
	}
	for i := 0; i < len(todo); i++ {
		category := Category{}
		if err := tm.db.First(&category, todo[i].CategoryID).Error; err != nil {
			logrus.Error("Model: Error Mendapatkan Data Category Todo ", err.Error())
			return nil
		}
		todo[i].Category = category
	}

	for i := 0; i < len(todo); i++ {
		user := Users{}
		if err := tm.db.First(&user, userID).Error; err != nil {
			logrus.Error("Model: Error Mendapatkan Data User Todo ", err.Error())
			return nil
		}
		todo[i].Category.User = user
		todo[i].User = user
	}
	return todo
}

func (tm *TodoModel) GetTodo(id int, userID uint) *Todo {
	todo := Todo{}
	if err := tm.db.Where("user_id = ?", userID).First(&todo, id).Error; err != nil {
		logrus.Error("Model: Error Mendapatkan Data Todo ", err.Error())
		return nil
	}
	category := Category{}
	if err := tm.db.Where("id = ?", todo.CategoryID).First(&category).Error; err != nil {
		logrus.Error("Model: Error Mendapatkan Data Category Todo ", err.Error())
		return nil
	}
	todo.Category = category

	user := Users{}
	if err := tm.db.Where("id = ?", userID).First(&user).Error; err != nil {
		logrus.Error("Model: Error Mendapatkan Data User Todo ", err.Error())
		return nil
	}
	todo.Category.User = user
	todo.User = user
	return &todo
}

func (tm *TodoModel) GetTodoByStatus(page, content int, userID uint, status string) []Todo {
	todo := []Todo{}
	offset := (page - 1) * content
	if err := tm.db.Limit(content).Offset(offset).Where("user_id = ? AND status = ?", userID, status).Find(&todo).Error; err != nil {
		logrus.Error("Model: Error Mendapatkan Data Todo Status ", err.Error())
		return nil
	}
	for i := 0; i < len(todo); i++ {
		category := Category{}
		if err := tm.db.Where("id = ?", todo[i].CategoryID).First(&category).Error; err != nil {
			logrus.Error("Model: Error Mendapatkan Data Category Todo ", err.Error())
			return nil
		}
		todo[i].Category = category
	}

	for i := 0; i < len(todo); i++ {
		user := Users{}
		if err := tm.db.Where("id = ?", userID).First(&user).Error; err != nil {
			logrus.Error("Model: Error Mendapatkan Data User Todo ", err.Error())
			return nil
		}
		todo[i].Category.User = user
		todo[i].User = user
	}
	return todo
}
func (tm *TodoModel) GetTodoByDate(page, content int, userID uint, datetime string) []Todo {
	todo := []Todo{}
	offset := (page - 1) * content
	if err := tm.db.Limit(content).Offset(offset).Where("user_id = ? AND DATE(date_time) = ?", userID, datetime).Find(&todo).Error; err != nil {
		logrus.Error("Model: Error Mendapatkan Data Todo Status ", err.Error())
		return nil
	}
	for i := 0; i < len(todo); i++ {
		category := Category{}
		if err := tm.db.Where("id = ?", todo[i].CategoryID).First(&category).Error; err != nil {
			logrus.Error("Model: Error Mendapatkan Data Category Todo ", err.Error())
			return nil
		}
		todo[i].Category = category
	}

	for i := 0; i < len(todo); i++ {
		user := Users{}
		if err := tm.db.Where("id = ?", userID).First(&user).Error; err != nil {
			logrus.Error("Model: Error Mendapatkan Data User Todo ", err.Error())
			return nil
		}
		todo[i].Category.User = user
		todo[i].User = user
	}
	return todo
}

func (tm *TodoModel) UpdateTodo(id int, userID uint, todo Todo) bool {
	data := tm.GetTodo(id, userID)
	if data == nil {
		logrus.Error("Model: Error Update Todo")
		return false
	}
	data.Memo = todo.Memo
	data.DateTime = todo.DateTime
	data.CategoryID = todo.CategoryID
	if err := tm.db.Save(&data).Error; err != nil {
		logrus.Error("Model: Error Update Todo")
		return false
	}
	return true
}

func (tm *TodoModel) UpdateTodoStatus(id int, userID uint, status string) bool {
	data := tm.GetTodo(id, userID)
	if data == nil {
		logrus.Error("Model: Error Update Todo")
		return false
	}
	data.Status = status
	if err := tm.db.Save(&data).Error; err != nil {
		logrus.Error("Model: Error Update Todo")
		return false
	}
	return true
}

func (tm *TodoModel) DeleteTodo(id int, userID uint) bool {
	todo := Todo{}
	data := tm.GetTodo(id, userID)
	if data == nil {
		logrus.Error("Model: Error Delete Todo")
		return false
	}
	if err := tm.db.Where("user_id = ?", userID).Delete(&todo, id).Error; err != nil {
		logrus.Error("Model: Error Delete Todo")
		return false
	}
	return true
}
