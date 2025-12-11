package table

import (
	"github.com/sirupsen/logrus"
	"gorm.io/gorm"
	"hotel/internal/util/logger"

	"hotel/models"
)

func Table(db *gorm.DB) {
	err := db.AutoMigrate(&models.Employee{})
	if err != nil {
		logger.Logger.WithFields(logrus.Fields{
			"error": err.Error(),
		}).Error("创建员工表失败")
	}

	err = db.AutoMigrate(&models.Guest{})
	if err != nil {
		logger.Logger.WithFields(logrus.Fields{
			"error": err.Error(),
		}).Error("创建客户表失败")
	}

	err = db.AutoMigrate(&models.Luggage{})
	if err != nil {
		logger.Logger.WithFields(logrus.Fields{
			"error": err.Error(),
		}).Error("创建行李表失败")
	}

	err = db.AutoMigrate(&models.Role{})
	if err != nil {
		logger.Logger.WithFields(logrus.Fields{
			"error": err.Error(),
		}).Error("创建角色表失败")
	}

	err = db.AutoMigrate(&models.Permission{})
	if err != nil {
		logger.Logger.WithFields(logrus.Fields{
			"error": err.Error(),
		}).Error("创建权限表失败")
	}

	logger.Logger.Info("数据库表创建成功")
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
