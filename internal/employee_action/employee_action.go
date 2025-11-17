package employee_action

import (
	"errors"
	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
	"hotel/internal/models"
	"log"
	"net/http"
	"time"
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

	// 先创建或获取客户记录
	var guest models.Guest
	result := db.Where("guest_name = ?", luggage.GuestName).First(&guest)
	if result.Error != nil {
		if errors.Is(result.Error, gorm.ErrRecordNotFound) {
			// 客户不存在，创建新客户
			guest = models.Guest{
				Name: luggage.GuestName,
			}
			if err := db.Create(&guest).Error; err != nil {
				c.JSON(http.StatusInternalServerError, gin.H{
					"success": false,
					"message": "创建客户记录失败",
				})
				log.Println("创建客户记录失败:", err)
				return
			}
		} else {
			// 其他查询错误
			c.JSON(http.StatusInternalServerError, gin.H{
				"success": false,
				"message": "查询客户失败",
			})
			log.Println("查询客户失败:", result.Error)
			return
		}
	}

	// 设置行李的 GuestID
	luggage.GuestID = guest.ID

	// 创建行李记录
	result = db.Create(&luggage)
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
			"id":         luggage.ID,
			"guest_id":   luggage.GuestID,
			"guest_name": luggage.GuestName,
			"tag":        luggage.Tag,
			"weight":     luggage.Weight,
			"status":     luggage.Status,
			"location":   luggage.Location,
		},
	})

}

// Delete 删除指定资源
func Delete(c *gin.Context, db *gorm.DB) {
	var luggage models.Luggage
	if err := c.ShouldBind(&luggage); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"success": false,
			"message": "请求数据格式错误",
		})
		return
	}
	var existingLuggage models.Luggage
	if err := db.Where("status = ?", "寄存中").First(&existingLuggage, luggage.ID).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			c.JSON(http.StatusNotFound, gin.H{
				"success": false,
				"message": "行李记录不存在",
			})
		} else {
			c.JSON(http.StatusInternalServerError, gin.H{
				"success": false,
				"message": "查询行李失败",
			})
			log.Println("查询行李失败:", err)
		}
		return
	}

	luggage.Status = "已取出"
	luggage.Location = "已取出"
	result := db.Model(&models.Luggage{}).Where("id = ?", luggage.ID).Updates(luggage)
	if result.Error != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"success": false,
			"message": "删除行李失败",
		})
		log.Println("删除行李失败:", result.Error)
		return
	}

	result2 := db.Where("id = ?", luggage.ID).Delete(&models.Luggage{})
	if result2.Error != nil {
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
		c.JSON(http.StatusBadRequest, gin.H{
			"success": false,
			"message": "请求数据格式错误",
		})
		log.Println("行李数据绑定错误:", err.Error())
		return
	}

	// 检查ID是否提供
	if luggage.ID == 0 {
		c.JSON(http.StatusBadRequest, gin.H{
			"success": false,
			"message": "请提供行李ID",
		})
		return
	}

	// 先检查记录是否存在
	var existingLuggage models.Luggage
	if err := db.First(&existingLuggage, luggage.ID).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			c.JSON(http.StatusNotFound, gin.H{
				"success": false,
				"message": "行李记录不存在",
			})
		} else {
			c.JSON(http.StatusInternalServerError, gin.H{
				"success": false,
				"message": "查询行李失败",
			})
			log.Println("查询行李失败:", err)
		}
		return
	}

	// 执行更新
	result := db.Model(&existingLuggage).Updates(luggage)
	if result.Error != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"success": false,
			"message": "更新行李失败",
		})
		log.Println("更新行李失败:", result.Error)
		return
	}

	// 检查是否真的更新了记录
	if result.RowsAffected == 0 {
		c.JSON(http.StatusBadRequest, gin.H{
			"success": false,
			"message": "没有数据需要更新",
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"success": true,
		"message": "行李更新成功",
		"data": gin.H{
			"id":         existingLuggage.ID,
			"guest_id":   existingLuggage.GuestID,
			"guest_name": existingLuggage.GuestName,
			"tag":        existingLuggage.Tag,
			"weight":     existingLuggage.Weight,
			"status":     existingLuggage.Status,
			"location":   existingLuggage.Location,
		},
	})
}

// GetName 获取行李
func GetName(c *gin.Context, db *gorm.DB) {
	var luggage []models.Luggage

	guestName := c.Query("guest_name")
	if guestName == "" {
		c.JSON(http.StatusBadRequest, gin.H{
			"success": false,
			"message": "请输入客户姓名",
		})
		return
	}

	result := db.Select("id", "guest_id", "guest_name", "tag", "weight", "status", "location").
		Where("guest_name = ? AND status = ?", guestName, "寄存中").
		Find(&luggage)
	if result.Error != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"success": false,
			"message": "获取行李失败",
		})
		log.Println("获取行李失败:", result.Error)
		return
	}

	if len(luggage) == 0 {
		c.JSON(http.StatusNotFound, gin.H{
			"success": false,
			"message": "未找到该客户的行李",
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"success": true,
		"message": "行李获取成功",
		"data":    luggage,
		"count":   len(luggage),
	})
}

// GetAll 获取所有行李
func GetAll(c *gin.Context, db *gorm.DB) {
	var luggage []models.Luggage
	result := db.Select("id", "guest_id", "guest_name", "tag", "weight", "status", "location").
		Where("status = ?", "寄存中").
		Find(&luggage)
	if result.Error != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"success": false,
			"message": "获取行李失败",
		})
		log.Println("获取行李失败:", result.Error)
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"success": true,
		"message": "行李获取成功",
		"data":    luggage,
		"count":   len(luggage),
	})

}

// GetGuestID 获取行李
func GetGuestID(c *gin.Context, db *gorm.DB) {
	guestID := c.Query("guest_id")
	if guestID == "" {
		c.JSON(http.StatusBadRequest, gin.H{
			"success": false,
			"message": "请输入客户ID",
		})
		return
	}
	var luggage []models.Luggage
	result := db.Select("id", "guest_id", "guest_name", "tag", "weight", "status", "location").
		Where("guest_id = ? AND status = ?", guestID, "寄存中").
		Find(&luggage)
	if result.Error != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"success": false,
			"message": "获取行李失败",
		})
		log.Println("获取行李失败:", result.Error)
		return
	}

	if len(luggage) == 0 {
		c.JSON(http.StatusNotFound, gin.H{
			"success": false,
			"message": "未找到该客户的行李",
		})
		return
	}
}

// GetLocation 获取行李
func GetLocation(c *gin.Context, db *gorm.DB) {
	location := c.Query("location")
	if location == "" {
		c.JSON(http.StatusBadRequest, gin.H{
			"success": false,
			"message": "请输入存放地点",
		})
		return
	}

	var luggage []models.Luggage
	result := db.Select("id", "guest_id", "guest_name", "tag", "weight", "status", "location").
		Where("location = ? AND status = ?", location, "寄存中").
		Find(&luggage)
	if result.Error != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"success": false,
			"message": "获取行李失败",
		})
		log.Println("获取行李失败:", result.Error)
		return
	}

	if len(luggage) == 0 {
		c.JSON(http.StatusNotFound, gin.H{
			"success": false,
			"message": "未找到该存放地点的行李",
		})
		return
	}

}

// GetStatus 获取行李
func GetStatus(c *gin.Context, db *gorm.DB) {
	status := c.Query("status")
	if status == "" {
		c.JSON(http.StatusBadRequest, gin.H{
			"success": false,
			"message": "请输入状态",
		})
		return
	}

	var luggage []models.Luggage
	result := db.Select("id", "guest_id", "guest_name", "tag", "weight", "status", "location").
		Where("status = ?", status).
		Find(&luggage)
	if result.Error != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"success": false,
			"message": "获取行李失败",
		})
		log.Println("获取行李失败:", result.Error)
		return
	}

	if len(luggage) == 0 {
		c.JSON(http.StatusNotFound, gin.H{
			"success": false,
			"message": "未找到该状态的行李",
		})
		return
	}
}

func GetAdvance(c *gin.Context, db *gorm.DB) {
	var find models.Luggage
	err := c.ShouldBind(&find)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"success": false,
			"message": "请求数据格式错误",
		})
		log.Println("获取行李信息绑定错误", err.Error())
		return
	}

	guestName := find.GuestName
	guestID := find.GuestID
	location := find.Location
	status := find.Status
	tag := find.Tag
	weight := find.Weight

	var luggage []models.Luggage
	// 构建基础查询
	query := db.Model(&models.Luggage{})

	// 只有当参数不为空时才添加条件
	if guestName != "" {
		query = query.Where("guest_name = ?", guestName)
	}

	if guestID != 0 {
		query = query.Where("guest_id = ?", guestID)
	}

	if location != "" {
		query = query.Where("location = ?", location)
	}

	if status != "" {
		query = query.Where("status = ?", status)
	} else {
		query = query.Where("status = ?", "寄存中")
	}

	if tag != "" {
		query = query.Where("tag = ?", tag)
	}

	if weight != 0 {
		query = query.Where("weight = ?", weight)
	}

	result := query.Find(&luggage)
	if result.Error != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"success": false,
			"message": "获取行李失败",
		})
		log.Println("获取行李失败:", result.Error)
		return
	}

	if len(luggage) == 0 {
		c.JSON(http.StatusNotFound, gin.H{
			"success": false,
			"message": "未找到该状态的行李",
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"success": true,
		"message": "获取行李成功",
		"data":    luggage,
		"count":   len(luggage),
	})

}

// CountSum 获取总行李数量
func CountSum(c *gin.Context, db *gorm.DB) {
	var count int64
	db.Model(&models.Luggage{}).Where("status = ?", "寄存中").Count(&count)

	c.JSON(http.StatusOK, gin.H{
		"success": true,
		"message": "行李数量获取成功",
		"count":   count,
	})
}

// CountToday 获取今日行李数量
func CountToday(c *gin.Context, db *gorm.DB) {
	today := time.Now().Format("2006-01-02")

	// 统计今天新增的行李（进入）
	var todayAdded int64
	db.Model(&models.Luggage{}).
		Where("DATE(created_at) = ?", today).
		Count(&todayAdded)

	// 统计今天状态变为"已取出"的行李（出去）
	var todayTaken int64
	db.Model(&models.Luggage{}).
		Unscoped().
		Where("status = ? AND DATE(updated_at) = ?", "已取出", today).
		Count(&todayTaken)

	c.JSON(http.StatusOK, gin.H{
		"success": true,
		"message": "获取今日统计成功",
		"data": gin.H{
			"date":        today,
			"today_added": todayAdded, // 今日进入
			"today_taken": todayTaken, // 今日出去
		},
	})

}
