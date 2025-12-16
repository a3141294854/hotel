package employee_action

import (
	"github.com/gin-gonic/gin"
	"hotel/models"
	"hotel/services"
	"net/http"
	"time"
)

// CountSum 获取总行李数量
func CountSum(c *gin.Context, s *services.Services) {
	var count int64
	s.DB.Model(&models.LuggageStorage{}).Where("status = ?", "寄存中").Count(&count)

	c.JSON(http.StatusOK, gin.H{
		"success": true,
		"message": "行李数量获取成功",
		"count":   count,
	})
}

// CountToday 获取今日行李数量
func CountToday(c *gin.Context, s *services.Services) {
	today := time.Now().Format("2006-01-02")

	// 统计今天新增的行李
	var todayAdded int64
	s.DB.Model(&models.LuggageStorage{}).
		Where("DATE(created_at) = ?", today).
		Count(&todayAdded)

	// 统计今天取出的行李
	var todayTaken int64
	s.DB.Model(&models.LuggageStorage{}).
		Unscoped().
		Where("status = ? AND DATE(updated_at) = ?", "已取出", today).
		Count(&todayTaken)

	c.JSON(http.StatusOK, gin.H{
		"success": true,
		"message": "获取今日统计成功",
		"data": gin.H{
			"date":        today,
			"today_added": todayAdded,
			"today_taken": todayTaken,
		},
	})

}
