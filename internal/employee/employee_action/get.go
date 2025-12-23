package employee_action

import (
	"encoding/json"
	"errors"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/sirupsen/logrus"
	"gorm.io/gorm"
	"hotel/internal/util/logger"
	"hotel/models"
	"hotel/services"
)

// GetPickUpCode 获取行李寄存表
func GetPickUpCode(c *gin.Context, s *services.Services) {
	var luggageStorage models.LuggageStorage
	//绑定
	err := c.ShouldBindJSON(&luggageStorage)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"success": false,
			"message": "行李数据绑定错误",
		})
		return
	}

	//查询行李
	result := s.DB.Model(&models.LuggageStorage{}).
		Preload("Guest").
		Preload("Luggage").
		Preload("Luggage.Tag").
		Preload("Luggage.Location").
		Where("pick_up_code = ?", luggageStorage.PickUpCode).First(&luggageStorage)
	if result.Error != nil {
		if errors.Is(result.Error, gorm.ErrRecordNotFound) {
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
				"error":      result.Error,
				"luggage_id": luggageStorage.ID,
			}).Error("查询行李失败")
		}
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"success": true,
		"message": "行李获取成功",
		"data":    luggageStorage,
	})

}

// GetName 获取行李寄存表
func GetName(c *gin.Context, s *services.Services) {
	//获取客户姓名
	guestName := c.Query("guest_name")
	if guestName == "" {
		c.JSON(http.StatusBadRequest, gin.H{
			"success": false,
			"message": "请输入客户姓名",
		})
		return
	}

	var guest models.Guest
	//查询客户
	result := s.DB.Model(&models.Guest{}).
		Preload("LuggageStorage").
		Preload("LuggageStorage.Luggage").
		Preload("LuggageStorage.Luggage.Tag").
		Preload("LuggageStorage.Luggage.Location").
		Where("name = ?", guestName).
		First(&guest)
	if result.Error != nil {
		if errors.Is(result.Error, gorm.ErrRecordNotFound) {
			c.JSON(http.StatusNotFound, gin.H{
				"success": false,
				"message": "客户不存在",
			})
		} else {
			c.JSON(http.StatusInternalServerError, gin.H{
				"success": false,
				"message": "查询客户失败",
			})
			logger.Logger.WithFields(logrus.Fields{
				"error": result.Error,
			}).Error("查询客户失败")
		}
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"success": true,
		"message": "行李获取成功",
		"data":    guest.LuggageStorage,
		"count":   len(guest.LuggageStorage),
	})
}

// GetAll 获取所有行李寄存表
func GetAll(c *gin.Context, s *services.Services) {
	var luggage []models.LuggageStorage
	result := s.DB.
		Preload("Guest").
		Preload("Luggage").
		Preload("Luggage.Tag").
		Preload("Luggage.Location").
		Find(&luggage)
	if result.Error != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"success": false,
			"message": "获取行李失败",
		})
		logger.Logger.WithFields(logrus.Fields{
			"error": result.Error,
		}).Error("获取行李失败")
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"success": true,
		"message": "行李获取成功",
		"data":    luggage,
		"count":   len(luggage),
	})

}

// GetGuestID 获取行李寄存表
func GetGuestID(c *gin.Context, s *services.Services) {
	guestID := c.Query("guest_id")
	if guestID == "" {
		c.JSON(http.StatusBadRequest, gin.H{
			"success": false,
			"message": "请输入客户ID",
		})
		return
	}

	if val, err := s.RdbCac.Get(c, guestID).Result(); err == nil {
		var luggage []models.LuggageStorage
		val := json.Unmarshal([]byte(val), &luggage)
		if val == nil {
			c.JSON(http.StatusOK, gin.H{
				"success": true,
				"message": "获取行李成功",
				"data":    luggage,
				"count":   len(luggage),
			})
			return
		}
	}

	var luggage []models.LuggageStorage
	result := s.DB.
		Preload("Guest").
		Preload("Luggage").
		Preload("Luggage.Tag").
		Preload("Luggage.Location").
		Where("guest_id = ? AND status = ?", guestID, "寄存中").
		Find(&luggage)
	if result.Error != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"success": false,
			"message": "获取行李失败",
		})
		logger.Logger.WithFields(logrus.Fields{
			"error": result.Error,
		}).Error("获取行李失败")
		return
	}

	if len(luggage) == 0 {
		c.JSON(http.StatusNotFound, gin.H{
			"success": false,
			"message": "未找到该客户的行李",
		})
		return
	}

	val, err := json.Marshal(luggage)
	if err != nil {
		logger.Logger.WithFields(logrus.Fields{
			"error":    err,
			"guest_id": guestID,
		}).Error("JSON序列化失败")
	} else {
		s.RdbCac.Set(c, guestID, val, time.Minute*15)
	}

	c.JSON(http.StatusOK, gin.H{
		"success": true,
		"message": "行李获取成功",
		"data":    luggage,
		"count":   len(luggage),
	})
}

// GetLocation 获取一个地方的行李
func GetLocation(c *gin.Context, s *services.Services) {
	location := c.Query("location")
	if location == "" {
		c.JSON(http.StatusBadRequest, gin.H{
			"success": false,
			"message": "请输入存放地点",
		})
		return
	}

	var loc models.Location
	result := s.DB.Model(models.Location{}).
		Preload("Luggage").
		Where("name = ?", location).First(&loc)
	if result.Error != nil {
		if errors.Is(result.Error, gorm.ErrRecordNotFound) {
			c.JSON(http.StatusNotFound, gin.H{
				"success": false,
				"message": "存放地点不存在",
			})
			return
		}
		c.JSON(http.StatusInternalServerError, gin.H{
			"success": false,
			"message": "获取存放地点失败",
		})
		logger.Logger.WithFields(logrus.Fields{
			"error": result.Error,
		}).Error("获取存放地点失败")
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"success": true,
		"message": "获取行李成功",
		"data":    loc.Luggage,
		"count":   len(loc.Luggage),
	})
}

// GetAdvance 高级查询行李寄存表
func GetAdvance(c *gin.Context, s *services.Services) {

	type AdvanceRequest struct {
		OperatorName string `json:"operator_name"`
		GuestName    string `json:"guest_name"`
		Status       string `json:"status"`
	}

	var req AdvanceRequest
	err := c.ShouldBind(&req)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"success": false,
			"message": "请求数据格式错误",
		})
		logger.Logger.WithFields(logrus.Fields{
			"error": err.Error(),
		}).Error("获取行李信息绑定错误")
		return
	}

	var luggage []models.LuggageStorage

	query := s.DB.Model(&models.LuggageStorage{})

	if req.GuestName != "" {
		query = query.Preload("Guest").Where("guests.name = ?", req.GuestName)
	}
	if req.OperatorName != "" {
		query = query.Where("operator_name = ?", req.OperatorName)
	}

	if req.Status != "" {
		query = query.Where("status = ?", req.Status)
	} else {
		query = query.Where("status = ?", "寄存中")
	}

	result := query.
		Find(&luggage)
	if result.Error != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"success": false,
			"message": "获取行李失败",
		})
		logger.Logger.WithFields(logrus.Fields{
			"error": result.Error,
		}).Error("获取行李失败")
		return
	}

	if len(luggage) == 0 {
		c.JSON(http.StatusNotFound, gin.H{
			"success": false,
			"message": "未找到符合条件的行李",
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
