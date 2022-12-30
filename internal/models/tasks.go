package models

import (
	"gorm.io/gorm"
	"log"
)

type Tasks struct {
	Id      int64  `json:"id"`
	Title   string `json:"title"`
	Descr   string `json:"descr"`
	EndDate string `json:"end_Date"`
	Userid  int64  `json:"userId"`
}

type TaskModel struct {
	Db *gorm.DB
}

func (m *TaskModel) CreateTask(task Tasks) error {
	result := m.Db.Create(&task)
	return result.Error
}

func (m *TaskModel) ShowTaskDb(userId int64) (error, []Tasks) {
	tasks := []Tasks{}

	result := m.Db.Where("userId = ?", userId).Find(&tasks)

	if result.Error != nil {
		log.Print("Ошибка")
	}

	return nil, tasks
}

func (m *TaskModel) DeleteTask(id int64) error {

	result := m.Db.Table("tasks").Where("id=?", id).Delete(&Tasks{})

	if result.Error != nil {
		log.Fatalf("ошибка удаления", result.Error)
	}

	return result.Error

}
