package admin

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/sirupsen/logrus"
	"hotel/internal/util/logger"
	"hotel/models"
	"hotel/services"
)

func GetAllRole(s *services.Services, c *gin.Context) {
	var roles []models.Role
	result := s.DB.Model(models.Role{}).Preload("Permissions").Find(&roles)
	if result.Error != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"success": false,
			"message": "内部错误",
		})
		logger.Logger.WithFields(logrus.Fields{
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

func GetAllPermission(s *services.Services, c *gin.Context) {
	var permissions []models.Permission
	result := s.DB.Model(models.Permission{}).Find(&permissions)
	if result.Error != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"success": false,
			"message": "内部错误",
		})
		logger.Logger.WithFields(logrus.Fields{
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

func GetAllEmployee(s *services.Services, c *gin.Context) {
	var employees []models.Employee
	result := s.DB.Model(models.Employee{}).Preload("Role").Preload("Role.Permissions").Find(&employees)
	if result.Error != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"success": false,
			"message": "内部错误",
		})
		logger.Logger.WithFields(logrus.Fields{
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

func GetEmployee() {

}
