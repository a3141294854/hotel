package table

import (
	"fmt"
	"reflect"
	"strings"
	"time"

	"github.com/sirupsen/logrus"
	"gorm.io/gorm"
	"hotel/internal/util"
	"hotel/models"
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

	err = db.AutoMigrate(&models.Photo{})
	if err != nil {
		util.Logger.WithFields(logrus.Fields{
			"error": err.Error(),
		}).Error("创建照片表失败")
	}

	util.Logger.Info("数据库表创建成功")

	open(db)

}

func open(db *gorm.DB) {
	// 创建权限
	permissions := []models.Permission{
		{Name: "查看行李"},
		{Name: "创建行李"},
		{Name: "更新行李"},
		{Name: "删除行李"},
		{Name: "管理员"},
	}

	for _, p := range permissions {
		createIfNotExists(db, "name", &p, "权限表")
	}

	// 获取权限对象
	p1, _ := getPermissionByName(db, "查看行李")
	p2, _ := getPermissionByName(db, "创建行李")
	p3, _ := getPermissionByName(db, "更新行李")
	p4, _ := getPermissionByName(db, "删除行李")
	p5, _ := getPermissionByName(db, "管理员")

	// 创建角色
	roles := []models.Role{
		{Name: "员工", Permissions: []models.Permission{p1, p2, p3, p4}},
		{Name: "管理员", Permissions: []models.Permission{p1, p2, p3, p4, p5}},
	}

	for _, r := range roles {
		createIfNotExists(db, "name", &r, "角色表")
	}

	// 创建酒店
	hotel := models.Hotel{Name: "酒店1"}
	createIfNotExists(db, "name", &hotel, "酒店表")

	// 创建管理员员工
	admin := models.Employee{
		User:           "admin",
		Password:       "admin123",
		LastActiveTime: time.Now(),
		RoleID:         2,
		HotelID:        1,
	}
	createIfNotExists(db, "user", &admin, "员工表")

	// 创建位置
	location := models.Location{
		Name:    "房间1",
		HotelID: 1,
	}
	createIfNotExists(db, "name", &location, "位置表")
}

// createIfNotExists 通用函数：检查数据是否存在，不存在则创建
func createIfNotExists(db *gorm.DB, fieldName string, model interface{}, tableName string) {
	fieldValue := getFieldByName(model, fieldName)
	ok, err := util.ExIf(db, fieldName, model, fieldValue)
	if !ok {
		if err := db.Create(model).Error; err != nil {
			util.Logger.WithFields(logrus.Fields{
				"error": err,
			}).Errorf("创建%s失败", tableName)
		}
	} else if err != nil {
		util.Logger.WithFields(logrus.Fields{
			"error": err,
		}).Errorf("查询%s失败", tableName)
	}
}

// getFieldByName 通过反射获取结构体字段的值
func getFieldByName(model interface{}, fieldName string) string {
	v := reflect.ValueOf(model)
	if v.Kind() == reflect.Ptr {
		v = v.Elem()
	}
	// 将字段名首字母大写以匹配结构体字段
	fieldNameCamel := fieldName
	if len(fieldName) > 0 {
		fieldNameCamel = strings.ToUpper(fieldName[:1]) + fieldName[1:]
	}
	field := v.FieldByName(fieldNameCamel)
	if field.IsValid() {
		return fmt.Sprintf("%v", field.Interface())
	}
	return ""
}

// getPermissionByName 根据名称获取权限对象
func getPermissionByName(db *gorm.DB, name string) (models.Permission, error) {
	var permission models.Permission
	err := db.Where("name = ?", name).First(&permission).Error
	return permission, err
}
