package models

import (
	"gorm.io/gorm"
)

type Tasks struct {
	ID uint `gorm:"primary key;autoIncrement" json:"id"`

	Title *string `json:"title"`

	Description *string `json:"description"`
	Due_date    *string `json:"date"`
	// Age    *uint   `json:"age"`
	Status *string `json:"status"`
}

func MigrateTask(db *gorm.DB) error {
	err := db.AutoMigrate(&Tasks{})
	return err
}
