package table

import (
	"github.com/jinzhu/gorm"
	"hotel/internal/models"
	"log"
)

func Table(db *gorm.DB) {
	err := db.AutoMigrate(&models.Employee{})
	if err != nil {
		log.Println("创建员工表失败:", err.Error)
	}

	err = db.AutoMigrate(&models.Luggage{})
	if err != nil {
		log.Println("创建行李表失败:", err.Error)
	}

	err = db.AutoMigrate(&models.Guest{})
	if err != nil {
		log.Println("创建客户表失败:", err.Error)
	}

	log.Println("创建表成功")
}
