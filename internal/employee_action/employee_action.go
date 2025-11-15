package employee_action

import (
	"github.com/gin-gonic/gin"
	"github.com/jinzhu/gorm"
	"hotel/internal/models"
	"log"
	"net/http"
)

// Add 添加行李
func Add(c *gin.Context, db *gorm.DB) {
	var luggage models.Luggage

	// 绑定请求数据到行李结构体
	if err := c.ShouldBind(&luggage); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"success": false,
			"message": "请求数据格式错误",
		})
		log.Println("行李数据绑定错误:", err.Error())
		return
	}

	// 验证必要字段
	if luggage.GuestID == 0 {
		c.JSON(http.StatusBadRequest, gin.H{
			"success": false,
			"message": "客人ID不能为空",
		})
		return
	}

	if luggage.Tag == "" {
		c.JSON(http.StatusBadRequest, gin.H{
			"success": false,
			"message": "行李标签不能为空",
		})
		return
	}

	// 设置默认状态和位置
	if luggage.Status == "" {
		luggage.Status = "寄存中"
	}

	if luggage.Location == "" {
		luggage.Location = "前台"
	}

	// 创建行李记录
	result := db.Create(&luggage)
	if result.Error != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"success": false,
			"message": "添加行李失败",
		})
		log.Println("添加行李失败:", result.Error)
		return
	}

	// 返回成功响应
	c.JSON(http.StatusOK, gin.H{
		"success": true,
		"message": "行李添加成功",
		"data": gin.H{
			"id":       luggage.ID,
			"guest_id": luggage.GuestID,
			"tag":      luggage.Tag,
			"weight":   luggage.Weight,
			"status":   luggage.Status,
			"location": luggage.Location,
		},
	})
}

// Delete 删除指定资源
func Delete(c *gin.Context, db *gorm.DB) {
	var luggage models.Luggage
	if err := c.ShouldBind(&luggage); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"success": false,
			"message": "删除行李失败",
		})
		return
	}
	result := db.Delete(&luggage)
	if result.Error != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"success": false,
			"message": "删除行李失败",
		})
		return
	}
	c.JSON(http.StatusOK, gin.H{
		"success": true,
		"message": "行李删除成功",
	})
}

func Update(c *gin.Context, db *gorm.DB) {
	var luggage models.Luggage
	if err := c.ShouldBind(&luggage); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"success": false,
			"message": "更新行李失败",
		})
		return
	}
	result := db.Model(&luggage).Updates(luggage)
	if result.Error != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"success": false,
			"message": "更新行李失败",
		})
		return
	}
	c.JSON(http.StatusOK, gin.H{
		"success": true,
		"message": "行李更新成功",
	})
}

// Get 获取行李
func Get(c *gin.Context, db *gorm.DB) {
	var luggage models.Luggage
	if err := c.ShouldBind(&luggage); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"success": false,
			"message": "获取行李失败",
		})
		return
	}
	result := db.Where("")
	if result.Error != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"success": false,
			"message": "获取行李失败",
		})
		return
	}
	c.JSON(http.StatusOK, gin.H{
		"success": true,
		"message": "行李获取成功",
		"data":    luggage,
	})
}
