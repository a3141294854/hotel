package table

import (
	"github.com/sirupsen/logrus"
	"gorm.io/gorm"
	"hotel/internal/util"
	"hotel/models"
	"time"
)

// Table 创建数据库表
func Table(db *gorm.DB) {
	err := db.AutoMigrate(&models.Employee{})
	if err != nil {
		util.Logger.WithFields(logrus.Fields{
			"error": err.Error(),
		}).Error("创建员工表失败")
	}

	err = db.AutoMigrate(&models.Guest{})
	if err != nil {
		util.Logger.WithFields(logrus.Fields{
			"error": err.Error(),
		}).Error("创建客户表失败")
	}

	err = db.AutoMigrate(&models.LuggageStorage{})
	if err != nil {
		util.Logger.WithFields(logrus.Fields{
			"error": err.Error(),
		}).Error("创建行李存储表失败")
	}

	err = db.AutoMigrate(&models.Role{})
	if err != nil {
		util.Logger.WithFields(logrus.Fields{
			"error": err.Error(),
		}).Error("创建角色表失败")
	}

	err = db.AutoMigrate(&models.Permission{})
	if err != nil {
		util.Logger.WithFields(logrus.Fields{
			"error": err.Error(),
		}).Error("创建权限表失败")
	}

	err = db.AutoMigrate(&models.Location{})
	if err != nil {
		util.Logger.WithFields(logrus.Fields{
			"error": err.Error(),
		}).Error("创建位置表失败")
	}

	err = db.AutoMigrate(&models.Hotel{})
	if err != nil {
		util.Logger.WithFields(logrus.Fields{
			"error": err.Error(),
		}).Error("创建酒店表失败")
	}

	err = db.AutoMigrate(&models.Luggage{})
	if err != nil {
		util.Logger.WithFields(logrus.Fields{
			"error": err.Error(),
		}).Error("创建行李表失败")
	}

	err = db.AutoMigrate(&models.Tag{})
	if err != nil {
		util.Logger.WithFields(logrus.Fields{
			"error": err.Error(),
		}).Error("创建标签表失败")
	}

	util.Logger.Info("数据库表创建成功")

	open(db)

}

func open(db *gorm.DB) {
	p1 := models.Permission{
		Name: "查看行李",
	}
	ok, err := util.ExIf(db, "name", &models.Permission{}, "查看行李")
	if !ok {
		db.Create(&p1)
		if err != nil {
			util.Logger.WithFields(logrus.Fields{
				"error": err,
			}).Error("查询权限表失败")
		}
	} else {
		if err != nil {
			util.Logger.WithFields(logrus.Fields{
				"error": err,
			}).Error("查询权限表失败")
		}
	}

	p2 := models.Permission{
		Name: "创建行李",
	}
	ok, err = util.ExIf(db, "name", &models.Permission{}, "创建行李")
	if !ok {
		db.Create(&p2)
		if err != nil {
			util.Logger.WithFields(logrus.Fields{
				"error": err,
			}).Error("查询权限表失败")
		}
	} else {
		if err != nil {
			util.Logger.WithFields(logrus.Fields{
				"error": err,
			}).Error("查询权限表失败")
		}
	}

	p3 := models.Permission{
		Name: "更新行李",
	}
	ok, err = util.ExIf(db, "name", &models.Permission{}, "更新行李")
	if !ok {
		db.Create(&p3)
		if err != nil {
			util.Logger.WithFields(logrus.Fields{
				"error": err,
			}).Error("查询权限表失败")
		}
	} else {
		if err != nil {
			util.Logger.WithFields(logrus.Fields{
				"error": err,
			}).Error("查询权限表失败")
		}
	}

	p4 := models.Permission{
		Name: "删除行李",
	}
	ok, err = util.ExIf(db, "name", &models.Permission{}, "删除行李")
	if !ok {
		db.Create(&p4)
		if err != nil {
			util.Logger.WithFields(logrus.Fields{
				"error": err,
			}).Error("查询权限表失败")
		}
	} else {
		if err != nil {
			util.Logger.WithFields(logrus.Fields{
				"error": err,
			}).Error("查询权限表失败")
		}
	}

	p5 := models.Permission{
		Name: "管理员",
	}
	ok, err = util.ExIf(db, "name", &models.Permission{}, "管理员")
	if !ok {
		db.Create(&p5)
		if err != nil {
			util.Logger.WithFields(logrus.Fields{
				"error": err,
			}).Error("查询权限表失败")
		}
	} else {
		if err != nil {
			util.Logger.WithFields(logrus.Fields{
				"error": err,
			}).Error("查询权限表失败")
		}
	}

	insert := models.Role{
		Name:        "员工",
		Permissions: []models.Permission{p1, p2, p3, p4},
	}
	ok, err = util.ExIf(db, "name", &models.Role{}, "员工")
	if !ok {
		db.Create(&insert)
		if err != nil {
			util.Logger.WithFields(logrus.Fields{
				"error": err,
			}).Error("查询角色表失败")
		}
	} else {
		if err != nil {
			util.Logger.WithFields(logrus.Fields{
				"error": err,
			}).Error("查询角色表失败")
		}
	}

	insert5 := models.Role{
		Name:        "管理员",
		Permissions: []models.Permission{p1, p2, p3, p4, p5},
	}
	ok, err = util.ExIf(db, "name", &models.Role{}, "管理员")
	if !ok {
		db.Create(&insert5)
		if err != nil {
			util.Logger.WithFields(logrus.Fields{
				"error": err,
			}).Error("查询角色表失败")
		}
	} else {
		if err != nil {
			util.Logger.WithFields(logrus.Fields{
				"error": err,
			}).Error("查询角色表失败")
		}
	}

	insert2 := models.Hotel{
		Name: "酒店1",
	}
	ok, err = util.ExIf(db, "name", &models.Hotel{}, "酒店1")
	if !ok {
		db.Create(&insert2)
		if err != nil {
			util.Logger.WithFields(logrus.Fields{
				"error": err,
			}).Error("查询酒店表失败")
		}
	} else {
		if err != nil {
			util.Logger.WithFields(logrus.Fields{
				"error": err,
			}).Error("查询酒店表失败")
		}
	}

	insert3 := models.Employee{
		User:           "admin",
		Password:       "admin123",
		LastActiveTime: time.Now(),
		RoleID:         2,
		HotelID:        1,
	}
	ok, err = util.ExIf(db, "user", &models.Employee{}, "admin")
	if !ok {
		db.Create(&insert3)
		if err != nil {
			util.Logger.WithFields(logrus.Fields{
				"error": err,
			}).Error("查询员工表失败")
		}
	} else {
		if err != nil {
			util.Logger.WithFields(logrus.Fields{
				"error": err,
			}).Error("查询员工表失败")
		}
	}

	insert4 := models.Location{
		Name:    "房间1",
		HotelID: 1,
	}
	ok, err = util.ExIf(db, "name", &models.Location{}, "房间1")
	if !ok {
		db.Create(&insert4)
		if err != nil {
			util.Logger.WithFields(logrus.Fields{
				"error": err,
			}).Error("查询位置表失败")
		}
	} else {
		if err != nil {
			util.Logger.WithFields(logrus.Fields{
				"error": err,
			}).Error("查询位置表失败")
		}
	}

}
