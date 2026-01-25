package employee_action

import (
	"errors"
	"github.com/gin-gonic/gin"
	"github.com/sirupsen/logrus"
	"gorm.io/gorm"
	"hotel/internal/util"
	"hotel/models"
	"hotel/services"
	"net/http"
)

// GetLuggageStorage 获取行李寄存表
func GetLuggageStorage(c *gin.Context, s *services.Services) {

	query := s.DB.Model(&models.LuggageStorage{})

	if c.Query("guest_name") != "" {
		query = query.Preload("Guest").Where("Guest.name = ?", c.Query("guest_name"))
	}
	if c.Query("pick_up_code") != "" {
		query = query.Where("pick_up_code = ?", c.Query("pick_up_code"))
	}
	if c.Query("status") != "" {
		query = query.Where("status = ?", c.Query("status"))
	}
	if c.Query("operator_name") != "" {
		query = query.Where("operator_name = ?", c.Query("operator_name"))
	}
	if c.Query("id") != "" {
		query = query.Where("id = ?", c.Query("id"))
	}
	if c.Query("unscoped") != "" {
		query = query.Unscoped()
	}

	var luggageStorage []models.LuggageStorage
	result := query.
		Preload("Guest").
		Preload("Luggage").
		Preload("Photos").
		Preload("Luggage.Tag").
		Preload("Luggage.Location").
		Preload("Photos").
		Find(&luggageStorage)

	if result.Error != nil {
		if errors.Is(result.Error, gorm.ErrRecordNotFound) {
			c.JSON(http.StatusNotFound, gin.H{
				"success": false,
				"message": "未找到符合条件的行李",
			})
			return
		}
		c.JSON(http.StatusInternalServerError, gin.H{
			"success": false,
			"message": "获取行李失败",
		})
		util.Logger.WithFields(logrus.Fields{
			"error": result.Error,
			"请求id":  c.GetUint("request_id"),
		}).Error("获取行李失败")
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"success": true,
		"message": "行李获取成功",
		"data":    luggageStorage,
	})

}

func GetLuggage(c *gin.Context, s *services.Services) {
	query := s.DB.Model(&models.Luggage{})
	if c.Query("id") != "" {
		query = query.Where("id = ?", c.Query("id"))
	}
	if c.Query("tag_id") != "" {
		query = query.Where("tag_id = ?", c.Query("tag_id"))
	}
	if c.Query("location_id") != "" {
		query = query.Where("location_id = ?", c.Query("location_id"))
	}
	if c.Query("name") != "" {
		query = query.Where("name = ?", c.Query("name"))
	}
	if c.Query("mac") != "" {
		query = query.Preload("Tag")
		query = query.Where("Tag.mac = ?", c.Query("mac"))
	}
	if c.Query("unscoped") != "" {
		query = query.Unscoped()
	}

	var luggage []models.Luggage
	result := query.
		Preload("Tag").
		Preload("Location").
		Find(&luggage)

	if result.Error != nil {
		if errors.Is(result.Error, gorm.ErrRecordNotFound) {
			c.JSON(http.StatusNotFound, gin.H{
				"success": false,
				"message": "未找到符合条件的行李",
			})
			return
		}
		c.JSON(http.StatusInternalServerError, gin.H{
			"success": false,
			"message": "获取行李失败",
		})
		util.Logger.WithFields(logrus.Fields{
			"error": result.Error,
			"请求id":  c.GetUint("request_id"),
		}).Error("获取行李失败")
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"success": true,
		"message": "行李获取成功",
		"data":    luggage,
	})
}

func GetTag(c *gin.Context, s *services.Services) {
	query := s.DB.Model(&models.Tag{})
	if c.Query("mac") != "" {
		query = query.Where("mac = ?", c.Query("mac"))
	}
	if c.Query("name") != "" {
		query = query.Where("name = ?", c.Query("name"))
	}
	if c.Query("id") != "" {
		query = query.Where("id = ?", c.Query("id"))
	}
	var tags []models.Tag
	result := query.Find(&tags)
	if result.Error != nil {
		if errors.Is(result.Error, gorm.ErrRecordNotFound) {
			c.JSON(http.StatusNotFound, gin.H{
				"success": false,
				"message": "未找到符合条件的标签",
			})
			return
		}
		c.JSON(http.StatusInternalServerError, gin.H{
			"success": false,
			"message": "获取标签失败",
		})
		util.Logger.WithFields(logrus.Fields{
			"error": result.Error,
			"请求id":  c.GetUint("request_id"),
		}).Error("获取标签失败")
		return
	}
	c.JSON(http.StatusOK, gin.H{
		"success": true,
		"message": "标签获取成功",
		"data":    tags,
	})
}

func GetLocation(c *gin.Context, s *services.Services) {
	query := s.DB.Model(&models.Location{})
	if c.Query("id") != "" {
		query = query.Where("id = ?", c.Query("id"))
	}
	if c.Query("name") != "" {
		query = query.Where("name = ?", c.Query("name"))
	}
	if c.Query("hotel_id") != "" {
		query = query.Where("hotel_id = ?", c.Query("hotel_id"))
	}
	var locations []models.Location
	result := query.Find(&locations)
	if result.Error != nil {
		if errors.Is(result.Error, gorm.ErrRecordNotFound) {
			c.JSON(http.StatusNotFound, gin.H{
				"success": false,
				"message": "未找到符合条件的位置",
			})
			return
		}
		c.JSON(http.StatusInternalServerError, gin.H{
			"success": false,
			"message": "获取位置失败",
		})
		util.Logger.WithFields(logrus.Fields{
			"error": result.Error,
			"请求id":  c.GetUint("request_id"),
		}).Error("获取位置失败")
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"success": true,
		"message": "位置获取成功",
		"data":    locations,
	})
}

func GetGuest(c *gin.Context, s *services.Services) {
	query := s.DB.Model(&models.Guest{})
	if c.Query("id") != "" {
		query = query.Where("id = ?", c.Query("id"))
	}
	if c.Query("name") != "" {
		query = query.Where("name = ?", c.Query("name"))
	}
	if c.Query("hotel_id") != "" {
		query = query.Where("hotel_id = ?", c.Query("hotel_id"))
	}
	var guests []models.Guest
	result := query.Find(&guests)
	if result.Error != nil {
		if errors.Is(result.Error, gorm.ErrRecordNotFound) {
			c.JSON(http.StatusNotFound, gin.H{
				"success": false,
				"message": "未找到符合条件的访客",
			})
			return
		}
		c.JSON(http.StatusInternalServerError, gin.H{
			"success": false,
			"message": "获取访客失败",
		})
		util.Logger.WithFields(logrus.Fields{
			"error": result.Error,
			"请求id":  c.GetUint("request_id"),
		}).Error("获取访客失败")
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"success": true,
		"message": "访客获取成功",
		"data":    guests,
	})
}

func GetHotel(c *gin.Context, s *services.Services) {
	query := s.DB.Model(&models.Hotel{})
	if c.Query("id") != "" {
		query = query.Where("id = ?", c.Query("id"))
	}
	if c.Query("name") != "" {
		query = query.Where("name = ?", c.Query("name"))
	}
	if c.Query("place") != "" {
		query = query.Where("place = ?", c.Query("place"))
	}

	var hotels []models.Hotel
	result := query.Find(&hotels)
	if result.Error != nil {
		if errors.Is(result.Error, gorm.ErrRecordNotFound) {
			c.JSON(http.StatusNotFound, gin.H{
				"success": false,
				"message": "未找到符合条件的酒店",
			})
			return
		}
		c.JSON(http.StatusInternalServerError, gin.H{
			"success": false,
			"message": "获取酒店失败",
		})
		util.Logger.WithFields(logrus.Fields{
			"error": result.Error,
			"请求id":  c.GetUint("request_id"),
		})
	}
	//返回
	c.JSON(http.StatusOK, gin.H{
		"success": true,
		"message": "酒店获取成功",
		"data":    hotels,
	})
}
