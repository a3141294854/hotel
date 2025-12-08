package table

import (
	"gorm.io/gorm"
	"hotel/models"
	"log"
)

func Table(db *gorm.DB) {
	err := db.AutoMigrate(&models.Employee{})
	if err != nil {
		log.Println("创建员工表失败:", err.Error())
	}

	err = db.AutoMigrate(&models.Guest{})
	if err != nil {
		log.Println("创建客户表失败:", err.Error())
	}

	err = db.AutoMigrate(&models.Luggage{})
	if err != nil {
		log.Println("创建行李表失败:", err.Error())
	}

	err = db.AutoMigrate(&models.Role{})
	if err != nil {
		log.Println("创建角色表失败:", err.Error())
	}

	err = db.AutoMigrate(&models.Permission{})
	if err != nil {
		log.Println("创建权限表失败:", err.Error())
	}

	log.Println("创建表成功")
	/*
		p1 := models.Permission{
			Name: "查看行李",
		}
		p2 := models.Permission{
			Name: "创建行李",
		}
		p3 := models.Permission{
			Name: "更新行李",
		}
		p4 := models.Permission{
			Name: "删除行李",
		}
		db.Create(&p1)
		db.Create(&p2)
		db.Create(&p3)
		db.Create(&p4)

		insert := models.Role{
			Name:        "员工",
			Permissions: []models.Permission{p1, p2, p3, p4},
		}
		db.Create(&insert)

	*/
}
