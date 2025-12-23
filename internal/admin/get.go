package admin

import (
	"github.com/gin-gonic/gin"
	"github.com/sirupsen/logrus"
	"hotel/internal/util"
	"hotel/models"
	"hotel/services"
	"net/http"
)

// GetAllRole 获取所有角色
func GetAllRole(s *services.Services, c *gin.Context) {
	var roles []models.Role
	result := s.DB.Model(models.Role{}).Preload("Permissions").Find(&roles)
	if result.Error != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"success": false,
			"message": "内部错误",
		})
		util.Logger.WithFields(logrus.Fields{
			"error": result.Error,
		}).Error("角色数据库查询错误")
		return

	}

	c.JSON(http.StatusOK, gin.H{
		"success": true,
		"message": "获取成功",
		"data":    roles,
	})
}

// GetAllPermission 获取所有权限
func GetAllPermission(s *services.Services, c *gin.Context) {
	var permissions []models.Permission
	result := s.DB.Model(models.Permission{}).Find(&permissions)
	if result.Error != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"success": false,
			"message": "内部错误",
		})
		util.Logger.WithFields(logrus.Fields{
			"error": result.Error,
		}).Error("权限数据库查询错误")
		return
	}
	c.JSON(http.StatusOK, gin.H{
		"success": true,
		"message": "获取成功",
		"data":    permissions,
	})
}

// GetAllEmployee 获取所有员工
func GetAllEmployee(s *services.Services, c *gin.Context) {
	var employees []models.Employee
	result := s.DB.Model(models.Employee{}).Preload("Role").Preload("Role.Permissions").Find(&employees)
	if result.Error != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"success": false,
			"message": "内部错误",
		})
		util.Logger.WithFields(logrus.Fields{
			"error": result.Error,
		}).Error("员工数据库查询错误")
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"success": true,
		"message": "获取成功",
		"data":    employees,
	})

}

// GetAllLocation 获取所有位置
func GetAllLocation(s *services.Services, c *gin.Context) {
	v, _ := c.Get("hotel_id")
	id := v.(uint)
	HotelID := id

	var locations []models.Location
	result := s.DB.Model(models.Location{}).Where("hotel_id = ?", HotelID).Find(&locations)
	if result.Error != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"success": false,
			"message": "内部错误",
		})
		util.Logger.WithFields(logrus.Fields{
			"error": result.Error,
		}).Error("位置数据库查询错误")
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"success": true,
		"message": "获取成功",
		"data":    locations,
	})

}
