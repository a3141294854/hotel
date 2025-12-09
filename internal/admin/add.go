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

func AddPermission(s *services.Services, c *gin.Context) {
	var p models.Permission
	err := c.ShouldBind(&p)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"success": false,
			"message": "请求数据格式错误",
		})
		log.Println("权限数据绑定错误:", err.Error())
		return
	}

	var ex models.Permission
	result := s.DB.Model(models.Permission{}).Where("name = ?", p.Name).First(&ex)
	if result.Error != nil {
		if errors.Is(result.Error, gorm.ErrRecordNotFound) {
			c.JSON(http.StatusOK, gin.H{
				"success": false,
				"message": "权限已存在",
			})
			return
		} else {
			c.JSON(http.StatusInternalServerError, gin.H{
				"success": false,
				"message": "内部错误",
			})
			log.Println("权限数据库查询错误:", result.Error)
			return
		}
	}

	s.DB.Create(&p)
}

func AddRole(s *services.Services, c *gin.Context) {
	var r models.Role
	err := c.ShouldBind(&r)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"success": false,
			"message": "请求数据格式错误",
		})
		log.Println("角色数据绑定错误:", err.Error())
		return
	}

	var ex models.Role
	result := s.DB.Model(models.Role{}).Where("name = ?", r.Name).First(&ex)
	if result.Error != nil {
		if errors.Is(result.Error, gorm.ErrRecordNotFound) {
			c.JSON(http.StatusOK, gin.H{
				"success": false,
				"message": "角色已存在",
			})
			return
		} else {
			c.JSON(http.StatusInternalServerError, gin.H{
				"success": false,
				"message": "内部错误",
			})
		}
	}

	s.DB.Create(&r)
}
