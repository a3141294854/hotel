package admin

import (
	"github.com/gin-gonic/gin"
	"hotel/models"
	"hotel/services"
	"log"
	"net/http"
)

func GetAllRole(s *services.Services, c *gin.Context) {
	var roles []models.Role
	result := s.DB.Model(models.Role{}).Preload("Permissions").Find(&roles)
	if result.Error != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"success": false,
			"message": "内部错误",
		})
		log.Println("角色数据库查询错误:", result.Error)
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
		log.Println("权限数据库查询错误:", result.Error)
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
	result := s.DB.Model(models.Employee{}).Preload("Role").Find(&employees)
	if result.Error != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"success": false,
			"message": "内部错误",
		})
		log.Println("员工数据库查询错误:", result.Error)
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
