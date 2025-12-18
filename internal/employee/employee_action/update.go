package employee_action

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

func UpdateLuggageStorage(c *gin.Context, s *services.Services) {
	var luggage models.LuggageStorage
	if err := c.ShouldBind(&luggage); err != nil {

		c.JSON(http.StatusBadRequest, gin.H{
			"success": false,
			"message": "请求数据格式错误",
		})

		logger.Logger.WithFields(logrus.Fields{
			"error": err.Error(),
		}).Error("行李数据绑定错误")
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
	var existingLuggage models.LuggageStorage
	if err := s.DB.First(&existingLuggage, luggage.ID).Error; err != nil {
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

	// 执行更新
	result := s.DB.Model(&models.LuggageStorage{}).Where("id = ?", luggage.ID).Updates(luggage)
	if result.Error != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"success": false,
			"message": "更新行李失败",
		})
		logger.Logger.WithFields(logrus.Fields{
			"error":      result.Error,
			"luggage_id": luggage.ID,
		}).Error("更新行李失败")
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"success": true,
		"message": "行李更新成功",
	})
}

func UpdateLuggage(c *gin.Context, s *services.Services) {
	var luggage models.Luggage
	if err := c.ShouldBind(&luggage); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"success": false,
			"message": "请求数据格式错误",
		})
		return
	}

	var ex models.Luggage
	if err := s.DB.Model(&models.Luggage{}).Where("id = ?", luggage.ID).First(&ex).Error; err != nil {
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

	result := s.DB.Model(&models.Luggage{}).Where("id = ?", luggage.ID).Updates(luggage)
	if result.Error != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"success": false,
			"message": "更新行李失败",
		})
		logger.Logger.WithFields(logrus.Fields{
			"error":      result.Error,
			"luggage_id": luggage.ID,
		}).Error("更新行李失败")
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"success": true,
		"message": "行李更新成功",
	})
}
