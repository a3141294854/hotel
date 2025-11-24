package models

import (
	"gorm.io/gorm"
	"time"
)

type Employee struct {
	ID       uint   `gorm:"primaryKey;autoIncrement;type:int unsigned"`
	Name     string `json:"name" gorm:"column:name"`
	User     string `json:"user" gorm:"column:user"`
	Password string `json:"password" gorm:"column:password"`
}

type Luggage struct {
	ID        uint           `gorm:"primaryKey;autoIncrement"`
	GuestID   uint           `json:"guest_id" gorm:"column:guest_id"`
	Tag       string         `json:"tag" gorm:"column:tag;size:50"`
	Weight    float32        `json:"weight" gorm:"column:weight"`
	Status    string         `json:"status" gorm:"column:status;size:20"`
	Location  string         `json:"location" gorm:"column:location;size:100"`
	CreatedAt time.Time      `json:"-"`
	UpdatedAt time.Time      `json:"-"`
	DeletedAt gorm.DeletedAt `gorm:"index" json:"-"`

	Guest Guest `gorm:"foreignKey:GuestID"`
}
type Guest struct {
	ID   uint   `gorm:"primaryKey;autoIncrement;type:int unsigned"`
	Name string `json:"name" gorm:"column:guest_name"`
}
