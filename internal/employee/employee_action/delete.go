package employee_action

import (
	"errors"
	"fmt"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/sirupsen/logrus"
	"gorm.io/gorm"
	"hotel/internal/util/logger"
	"hotel/models"
	"hotel/services"
)

// DeleteStorage 删除行李寄存表
func DeleteStorage(c *gin.Context, s *services.Services) {
	var luggage models.LuggageStorage
	if err := c.ShouldBind(&luggage); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"success": false,
			"message": "请求数据格式错误",
		})
		return
	}
	var existingLuggage models.LuggageStorage
	if err := s.DB.
		Preload("Guest").
		Where("status = ?", "寄存中").Where("id = ?", luggage.ID).First(&existingLuggage).Error; err != nil {
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
			logger.Logger.WithFields(logrus.Fields{
				"error":      err,
				"luggage_id": luggage.ID,
			}).Error("查询行李失败")
		}
		return
	}
	tx := s.DB.Begin()

	luggage.Status = "已取出"
	result := tx.Model(&models.LuggageStorage{}).Where("id = ?", luggage.ID).Updates(luggage)
	if result.Error != nil {
		tx.Rollback()
		c.JSON(http.StatusInternalServerError, gin.H{
			"success": false,
			"message": "删除行李失败",
		})
		logger.Logger.WithFields(logrus.Fields{
			"error":      result.Error,
			"luggage_id": luggage.ID,
		}).Error("删除行李失败")
		return
	}

	result2 := tx.Where("id = ?", luggage.ID).Delete(&models.LuggageStorage{})
	if result2.Error != nil {
		tx.Rollback()
		c.JSON(http.StatusInternalServerError, gin.H{
			"success": false,
			"message": "删除行李失败",
		})
		logger.Logger.WithFields(logrus.Fields{
			"error":      result2.Error,
			"luggage_id": luggage.ID,
		}).Error("删除行李失败")
		return
	}

	err := tx.Commit().Error
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"success": false,
			"message": "删除行李失败",
		})
		logger.Logger.WithFields(logrus.Fields{
			"error":      err,
			"luggage_id": luggage.ID,
		}).Error("事务提交失败")
		return
	}
	s.RdbRand.Del(c, fmt.Sprintf("%d:%s", luggage.HotelID, luggage.PickUpCode))
	s.RdbCac.Del(c, existingLuggage.Guest.Name)
	s.RdbCac.Del(c, fmt.Sprintf("%d", luggage.GuestID))

	c.JSON(http.StatusOK, gin.H{
		"success": true,
		"message": "行李删除成功",
	})

}
