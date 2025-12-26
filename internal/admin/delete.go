package admin

import (
	"errors"
	"hotel/internal/util"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/sirupsen/logrus"
	"gorm.io/gorm"
	"hotel/models"
	"hotel/services"
)

// DeleteEmployee 删除员工
func DeleteEmployee(s *services.Services, c *gin.Context) {
	var req struct {
		EmployeeID uint `json:"employee_id" binding:"required"`
	}
	//绑定
	err := c.ShouldBind(&req)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"success": false,
			"message": "请求数据格式错误",
		})
		util.Logger.WithFields(logrus.Fields{
			"error": err.Error(),
		}).Error("请求数据格式错误")
		return
	}
	//查找对应员工
	var ex models.Employee
	result := s.DB.Model(models.Employee{}).Where("id = ?", req.EmployeeID).First(&ex)
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
		util.Logger.WithFields(logrus.Fields{
			"error":       result.Error,
			"employee_id": req.EmployeeID,
		}).Error("员工数据库查询错误")
		return
	}
	//删除员工
	result = s.DB.Delete(&ex)
	if result.Error != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"success": false,
			"message": "内部错误",
		})
		util.Logger.WithFields(logrus.Fields{
			"error":       result.Error,
			"employee_id": req.EmployeeID,
		}).Error("员工数据库删除错误")
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"success": true,
		"message": "删除成功",
	})
}

// DeleteRole 删除角色
func DeleteRole(s *services.Services, c *gin.Context) {
	var req struct {
		RoleID uint `json:"role_id" binding:"required"`
	}
	//绑定
	err := c.ShouldBind(&req)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"success": false,
			"message": "请求数据格式错误",
		})
		util.Logger.WithFields(logrus.Fields{
			"error": err.Error(),
		}).Error("请求数据格式错误")
		return
	}
	//查找对应角色
	var ex models.Role
	result := s.DB.Model(models.Role{}).Where("id = ?", req.RoleID).First(&ex)
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
		util.Logger.WithFields(logrus.Fields{
			"error":   result.Error,
			"role_id": req.RoleID,
		}).Error("角色数据库查询错误")
		return
	}

	tx := s.DB.Begin()
	//删除角色权限关联
	err = tx.Model(&ex).Association("Permissions").Clear()
	if err != nil {
		tx.Rollback()
		c.JSON(http.StatusInternalServerError, gin.H{
			"success": false,
			"message": "内部错误",
		})
		util.Logger.WithFields(logrus.Fields{
			"error":   err,
			"role_id": req.RoleID,
		}).Error("员工数据库更新错误")
		return
	}
	//更新员工角色为默认角色，要自己手动赋角色1
	result = tx.Model(models.Employee{}).Where("role_id = ?", req.RoleID).Update("role_id", 3)
	if result.Error != nil {
		tx.Rollback()
		c.JSON(http.StatusInternalServerError, gin.H{
			"success": false,
			"message": "内部错误",
		})
		util.Logger.WithFields(logrus.Fields{
			"error":   result.Error,
			"role_id": req.RoleID,
		}).Error("员工数据库更新错误")
		return
	}
	//删除角色
	result = tx.Model(models.Role{}).Delete(&ex)
	if result.Error != nil {
		tx.Rollback()
		c.JSON(http.StatusInternalServerError, gin.H{
			"success": false,
			"message": "内部错误",
		})
		util.Logger.WithFields(logrus.Fields{
			"error":   result.Error,
			"role_id": req.RoleID,
		}).Error("角色数据库删除错误")
		return
	}
	//提交事务
	err = tx.Commit().Error
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"success": false,
			"message": "内部错误",
		})
		util.Logger.WithFields(logrus.Fields{
			"error": err,
		}).Error("事务提交错误")
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"success": true,
		"message": "删除成功",
	})
}

// DeletePermission 删除权限
func DeletePermission(s *services.Services, c *gin.Context) {
	var req struct {
		PermissionID uint `json:"permission_id" binding:"required"`
	}
	//绑定
	err := c.ShouldBind(&req)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"success": false,
			"message": "请求数据格式错误",
		})
		util.Logger.WithFields(logrus.Fields{
			"error": err.Error(),
		}).Error("请求数据格式错误")
		return
	}
	//查找对应权限
	var ex models.Permission
	result := s.DB.Model(models.Permission{}).Where("id = ?", req.PermissionID).First(&ex)
	if result.Error != nil {
		if errors.Is(result.Error, gorm.ErrRecordNotFound) {
			c.JSON(http.StatusNotFound, gin.H{
				"success": false,
				"message": "权限不存在",
			})
			return
		}
		c.JSON(http.StatusInternalServerError, gin.H{
			"success": false,
			"message": "内部错误",
		})
		util.Logger.WithFields(logrus.Fields{
			"error":         result.Error,
			"permission_id": req.PermissionID,
		}).Error("权限数据库查询错误")
		return
	}

	tx := s.DB.Begin()
	//删除角色权限关联
	err = tx.Model(&ex).Association("Roles").Clear()
	if err != nil {
		tx.Rollback()
		c.JSON(http.StatusInternalServerError, gin.H{
			"success": false,
			"message": "内部错误",
		})
		util.Logger.WithFields(logrus.Fields{
			"error":         err,
			"permission_id": req.PermissionID,
		}).Error("员工数据库更新错误")
		return
	}
	//删除权限
	result = tx.Model(models.Permission{}).Delete(&ex)
	if result.Error != nil {
		tx.Rollback()
		c.JSON(http.StatusInternalServerError, gin.H{
			"success": false,
			"message": "内部错误",
		})

		util.Logger.WithFields(logrus.Fields{
			"error":         result.Error,
			"permission_id": req.PermissionID,
		}).Error("权限数据库删除错误")
		return
	}
	//提交事务
	err = tx.Commit().Error
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"success": false,
			"message": "内部错误",
		})
		util.Logger.WithFields(logrus.Fields{
			"error": err,
		}).Error("事务提交错误")
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"success": true,
		"message": "删除成功",
	})

}
