package admin

import (
	"errors"
	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
	"hotel/models"
	"hotel/services"
	"log"
	"net/http"
)

func DeleteEmployee(s *services.Services, c *gin.Context) {

	var req struct {
		EmployeeID uint `json:"employee_id" binding:"required"`
	}
	err := c.ShouldBind(&req)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"success": false,
			"message": "请求数据格式错误",
		})
		log.Println("请求数据格式错误:", err.Error())
		return
	}
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
		log.Println("员工数据库查询错误:", result.Error)
		return
	}

	result = s.DB.Delete(&ex)
	if result.Error != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"success": false,
			"message": "内部错误",
		})
		log.Println("员工数据库删除错误:", result.Error)
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"success": true,
		"message": "删除成功",
	})
}

func DeleteRole(s *services.Services, c *gin.Context) {
	var req struct {
		RoleID uint `json:"role_id" binding:"required"`
	}

	err := c.ShouldBind(&req)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"success": false,
			"message": "请求数据格式错误",
		})
		log.Println("请求数据格式错误:", err.Error())
		return
	}

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
		log.Println("角色数据库查询错误", result.Error)
		return
	}

	tx := s.DB.Begin()

	result = tx.Model(models.Employee{}).Where("role_id = ?", req.RoleID).Update("role_id", 3)
	if result.Error != nil {
		tx.Rollback()
		c.JSON(http.StatusInternalServerError, gin.H{
			"success": false,
			"message": "内部错误",
		})
		log.Println("员工数据库更新错误", result.Error)
		return
	}
	result = tx.Model(models.Role{}).Delete(&ex)
	if result.Error != nil {
		tx.Rollback()
		c.JSON(http.StatusInternalServerError, gin.H{
			"success": false,
			"message": "内部错误",
		})
		log.Println("角色数据库删除错误", result.Error)
		return
	}
	err = tx.Commit().Error
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"success": false,
			"message": "内部错误",
		})
		log.Println("事务提交错误", err)
		return
	}
	c.JSON(http.StatusOK, gin.H{
		"success": true,
		"message": "删除成功",
	})
}

func DeletePermission(s *services.Services, c *gin.Context) {
	var req struct {
		PermissionID uint `json:"permission-id" binding:"required"`
	}

	err := c.ShouldBind(&req)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"success": false,
			"message": "请求数据格式错误",
		})
		log.Println("请求数据格式错误:", err.Error())
		return
	}

	var ex models.Permission
	result := s.DB.Model(models.Permission{}).Where("id = ?", req.PermissionID).First(&ex)
	if result != nil {
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
		log.Println("权限数据库查询错误", result.Error)
		return
	}

	result = s.DB.Model(models.Permission{}).Delete(&ex)
	if result != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"success": false,
			"message": "内部错误",
		})
		log.Println("权限数据库删除错误", result.Error)
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"success": true,
		"message": "删除成功",
	})

}
