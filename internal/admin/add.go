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

func AddPermission(s *services.Services, c *gin.Context) {
	var p models.Permission
	err := c.ShouldBind(&p)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"success": false,
			"message": "请求数据格式错误",
		})
		logger.Logger.WithFields(logrus.Fields{
			"error": err.Error(),
		}).Error("权限数据绑定错误")
		return
	}

	var ex models.Permission
	result := s.DB.Model(models.Permission{}).Where("name = ?", p.Name).First(&ex)
	if result.Error == nil {
		c.JSON(http.StatusConflict, gin.H{
			"success": false,
			"message": "权限已存在",
		})
		return
	}

	result = s.DB.Create(&p)
	if result.Error != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"success": false,
			"message": "内部错误",
		})
		logger.Logger.WithFields(logrus.Fields{
			"error": result.Error,
		}).Error("权限数据库插入错误")
		return
	}
	c.JSON(http.StatusOK, gin.H{
		"success": true,
		"message": "添加成功",
		"data":    p,
	})

}

func AddRole(s *services.Services, c *gin.Context) {
	var r models.Role
	err := c.ShouldBind(&r)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"success": false,
			"message": "请求数据格式错误",
		})
		logger.Logger.WithFields(logrus.Fields{
			"error": err.Error(),
		}).Error("角色数据绑定错误")
		return
	}

	var ex models.Role
	result := s.DB.Model(models.Role{}).Where("name = ?", r.Name).First(&ex)
	if result.Error == nil {
		c.JSON(http.StatusConflict, gin.H{
			"success": false,
			"message": "角色已存在",
		})
		return
	} else if !errors.Is(result.Error, gorm.ErrRecordNotFound) {
		c.JSON(http.StatusInternalServerError, gin.H{
			"success": false,
			"message": "内部错误",
		})
		logger.Logger.WithFields(logrus.Fields{
			"error": result.Error,
		}).Error("角色数据库查询错误")
	}

	result = s.DB.Create(&r)
	if result.Error != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"success": false,
			"message": "内部错误",
		})
		logger.Logger.WithFields(logrus.Fields{
			"error": result.Error,
		}).Error("角色数据库插入错误")
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"success": true,
		"message": "添加成功",
		"data":    r,
	})

}

func AddRolePermission(s *services.Services, c *gin.Context) {

	var req struct {
		RoleID       uint `json:"role_id" binding:"required"`
		PermissionID uint `json:"permission_id" binding:"required"`
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

	var ex1 models.Role
	result := s.DB.Model(models.Role{}).Where("id = ?", req.RoleID).First(&ex1)
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

	var ex2 models.Permission
	result = s.DB.Model(models.Permission{}).Where("id = ?", req.PermissionID).First(&ex2)
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
		logger.Logger.WithFields(logrus.Fields{
			"error":         result.Error,
			"permission_id": req.PermissionID,
		}).Error("权限数据库查询错误")
		return
	}

	role := models.Role{ID: req.RoleID}
	permission := models.Permission{ID: req.PermissionID}
	result1 := s.DB.Model(&role).Association("Permissions").Append(&permission)
	if result1 != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"success": false,
			"message": "内部错误",
		})
		logger.Logger.WithFields(logrus.Fields{
			"error":         result1,
			"role_id":       req.RoleID,
			"permission_id": req.PermissionID,
		}).Error("角色权限数据库插入错误")
		return
	}

}
