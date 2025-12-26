package admin

import (
	"fmt"
	"hotel/internal/util"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/sirupsen/logrus"
	"hotel/models"
	"hotel/services"
)

// ChangeEmployeeRole 修改员工角色
func ChangeEmployeeRole(s *services.Services, c *gin.Context) {
	var req struct {
		EmployeeID uint `json:"employee_id" binding:"required"`
		RoleID     uint `json:"role_id" binding:"required"`
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
	//检查员工是否存在
	ok, err := util.ExIf(s.DB, "id", &models.Employee{}, fmt.Sprintf("%d", req.EmployeeID))
	if !ok && err == nil {
		c.JSON(http.StatusNotFound, gin.H{
			"success": false,
			"message": "员工不存在",
		})
		return
	}
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"success": false,
			"message": "内部错误",
		})
		util.Logger.WithFields(logrus.Fields{
			"error": err,
		}).Error("员工数据库查询错误")
		return
	}
	//检查角色是否存在
	ok, err = util.ExIf(s.DB, "id", &models.Role{}, fmt.Sprintf("%d", req.RoleID))
	if !ok && err == nil {
		c.JSON(http.StatusNotFound, gin.H{
			"success": false,
			"message": "角色不存在",
		})
		return
	}
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"success": false,
			"message": "内部错误",
		})
		util.Logger.WithFields(logrus.Fields{
			"error": err,
		}).Error("角色数据库查询错误")
		return
	}
	//更新
	result := s.DB.Model(models.Employee{}).Where("id = ?", req.EmployeeID).Update("role_id", req.RoleID)
	if result.Error != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"success": false,
			"message": "内部错误",
		})
		util.Logger.WithFields(logrus.Fields{
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
