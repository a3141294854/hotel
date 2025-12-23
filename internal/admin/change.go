package admin

import (
	"errors"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/sirupsen/logrus"
	"gorm.io/gorm"
	"hotel/internal/util/logger"
	"hotel/models"
	"hotel/services"
)

// ChangeEmployeeRole 修改员工角色
func ChangeEmployeeRole(s *services.Services, c *gin.Context) {
	var req struct {
		EmployeeID uint `json:"employee_id" binding:"required"`
		RoleID     uint `json:"role_id" binding:"required"`
	}
	err := c.ShouldBind(&req)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"success": false,
			"message": "请求数据格式错误",
		})
		logger.Logger.WithFields(logrus.Fields{
			"error": err.Error(),
		}).Error("请求数据格式错误")
		return
	}

	var ex1 models.Employee
	result := s.DB.Model(models.Employee{}).Where("id = ?", req.EmployeeID).First(&ex1)
	if result.Error != nil {
		if errors.Is(result.Error, gorm.ErrRecordNotFound) {
			c.JSON(http.StatusNotFound, gin.H{
				"success": false,
				"message": "员工不存在",
			})
			return
		}
		c.JSON(http.StatusInternalServerError, gin.H{
			"success": false,
			"message": "内部错误",
		})
		logger.Logger.WithFields(logrus.Fields{
			"error":       result.Error,
			"employee_id": req.EmployeeID,
		}).Error("员工数据库查询错误")
		return
	}

	var ex2 models.Role
	result = s.DB.Model(models.Role{}).Where("id = ?", req.RoleID).First(&ex2)
	if result.Error != nil {
		if errors.Is(result.Error, gorm.ErrRecordNotFound) {
			c.JSON(http.StatusNotFound, gin.H{
				"success": false,
				"message": "角色不存在",
			})
			return
		}
		c.JSON(http.StatusInternalServerError, gin.H{
			"success": false,
			"message": "内部错误",
		})
		logger.Logger.WithFields(logrus.Fields{
			"error":   result.Error,
			"role_id": req.RoleID,
		}).Error("角色数据库查询错误")
		return
	}

	result = s.DB.Model(models.Employee{}).Where("id = ?", req.EmployeeID).Update("role_id", req.RoleID)
	if result.Error != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"success": false,
			"message": "内部错误",
		})
		logger.Logger.WithFields(logrus.Fields{
			"error":       result.Error,
			"employee_id": req.EmployeeID,
			"role_id":     req.RoleID,
		}).Error("员工角色数据库更新错误")
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"success": true,
		"message": "角色修改成功",
	})
}
