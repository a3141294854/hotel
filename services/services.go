package services

import (
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
	"log"
)

type Servers struct {
	DB *gorm.DB
}

func NewDatabase() *Servers {
	dsn := "root:@furenjie321@tcp(127.0.0.1:3306)/study?charset=utf8mb4&parseTime=True&loc=Local"
	db, err := gorm.Open(mysql.Open(dsn), &gorm.Config{})
	if err != nil {
		log.Fatal(err)
	}
	service := Servers{
		DB: db,
	}
	return &service
}
