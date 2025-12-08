package util

import (
	"github.com/gin-gonic/gin"
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
			"data":    err.Error(),
		})
		log.Println("权限数据绑定错误:", err.Error())
		return
	}
	s.DB.Create(&p)
}
