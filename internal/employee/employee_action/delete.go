package employee_action

import (
	"errors"
	"fmt"
	"hotel/internal/util"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/sirupsen/logrus"
	"gorm.io/gorm"
	"hotel/models"
	"hotel/services"
)

// DeleteStorage 删除行李寄存表
func DeleteStorage(c *gin.Context, s *services.Services) {

	id := c.Param("id")

	//检查是否存在
	var existingLuggage models.LuggageStorage
	if err := s.DB.
		Preload("Guest").
		Preload("Luggage").
		Where("id = ?", id).First(&existingLuggage).Error; err != nil {
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
			util.Logger.WithFields(logrus.Fields{
				"error":      err,
				"luggage_id": id,
			}).Error("查询行李失败")
		}
		return
	}
	//开启事务
	tx := s.DB.Begin()

	//删除客户
	result := tx.Model(&models.Guest{}).Where("id = ?", existingLuggage.GuestID).Delete(&models.Guest{})
	if result.Error != nil {
		tx.Rollback()
		c.JSON(http.StatusInternalServerError, gin.H{
			"success": false,
			"message": "删除行李失败",
		})
		util.Logger.WithFields(logrus.Fields{
			"error":      result.Error,
			"luggage_id": id,
		}).Error("删除行李失败")
		return
	}
	//删除行李寄存
	result = tx.Model(&models.LuggageStorage{}).Where("id = ?", existingLuggage.ID).Delete(&models.LuggageStorage{})
	if result.Error != nil {
		tx.Rollback()
		c.JSON(http.StatusInternalServerError, gin.H{
			"success": false,
			"message": "删除行李失败",
		})
		util.Logger.WithFields(logrus.Fields{
			"error":      result.Error,
			"luggage_id": id,
		}).Error("删除行李失败")
		return
	}
	//提交事务
	err := tx.Commit().Error
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"success": false,
			"message": "删除行李失败",
		})
		util.Logger.WithFields(logrus.Fields{
			"error":      err,
			"luggage_id": id,
		}).Error("事务提交失败")
		return
	}
	//删除redis缓存，更出新的取件码
	s.RdbRand.Del(c, fmt.Sprintf("%d:%s", existingLuggage.HotelID, existingLuggage.PickUpCode))

	c.JSON(http.StatusOK, gin.H{
		"success": true,
		"message": "行李删除成功",
	})
}

// DeleteLuggage 删除行李
func DeleteLuggage(c *gin.Context, s *services.Services) {
	var luggage models.Luggage
	//绑定
	if err := c.ShouldBind(&luggage); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"success": false,
			"message": "请求数据格式错误",
		})
		util.Logger.WithFields(logrus.Fields{
			"error": err.Error(),
		}).Error("请求数据格式错误")
		return
	}
	//检查是否存在
	var ex models.Luggage
	result := s.DB.Model(&models.Luggage{}).
		Preload("LuggageStorage").
		Where("id = ?", luggage.ID).First(&ex)
	if result.Error != nil {
		if errors.Is(result.Error, gorm.ErrRecordNotFound) {
			c.JSON(http.StatusNotFound, gin.H{
				"success": false,
				"message": "行李不存在",
			})
			util.Logger.WithFields(logrus.Fields{
				"error":      result.Error,
				"luggage_id": luggage.ID,
			}).Error("行李不存在")
			return
		}
		c.JSON(http.StatusInternalServerError, gin.H{
			"success": false,
			"message": "内部错误",
		})
		util.Logger.WithFields(logrus.Fields{
			"error":      result.Error,
			"luggage_id": luggage.ID,
		}).Error("行李数据库查询错误")
		return
	}
	//开启事务
	tx := s.DB.Begin()
	//删除行李
	result = tx.Model(&models.Luggage{}).Where("id = ?", luggage.ID).Delete(&models.Luggage{})
	if result.Error != nil {
		tx.Rollback()
		c.JSON(http.StatusInternalServerError, gin.H{
			"success": false,
			"message": "内部错误",
		})
		util.Logger.WithFields(logrus.Fields{
			"error":      result.Error,
			"luggage_id": luggage.ID,
		}).Error("行李数据库删除错误")
		return
	}
	//检查行李寄存表是否还有行李
	var luggageStorage models.LuggageStorage
	result = tx.Model(&models.LuggageStorage{}).
		Preload("Luggage").
		Where("id = ?", ex.LuggageStorageID).First(&luggageStorage)
	if result.Error != nil {
		tx.Rollback()
		c.JSON(http.StatusInternalServerError, gin.H{
			"success": false,
			"message": "内部错误",
		})
		util.Logger.WithFields(logrus.Fields{
			"error":      result.Error,
			"luggage_id": luggage.ID,
		}).Error("行李数据库删除错误")
		return
	}

	if len(luggageStorage.Luggage) == 0 {
		result = tx.Model(&models.LuggageStorage{}).
			Delete(&luggageStorage)
		if result.Error != nil {
			tx.Rollback()
			c.JSON(http.StatusInternalServerError, gin.H{
				"success": false,
				"message": "内部错误",
			})
			util.Logger.WithFields(logrus.Fields{
				"error":      result.Error,
				"luggage_id": luggage.ID,
			}).Error("行李数据库删除错误")
			return
		}
	}
	//提交事务
	err := tx.Commit().Error
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"success": false,
			"message": "内部错误",
		})
		util.Logger.WithFields(logrus.Fields{
			"error":      err,
			"luggage_id": luggage.ID,
		}).Error("事务提交失败")
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"success": true,
		"message": "行李删除成功",
	})

}

func DeleteLocation(c *gin.Context, s *services.Services) {
	var location models.Location
	//绑定
	if err := c.ShouldBind(&location); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"success": false,
			"message": "请求数据格式错误",
		})
		util.Logger.WithFields(logrus.Fields{
			"error": err.Error(),
		}).Error("请求数据格式错误")
		return
	}
	//检查是否存在
	var ex models.Location
	result := s.DB.Model(&models.Location{}).
		Preload("Luggage").
		Where("id = ?", location.ID).First(&ex)
	if result.Error != nil {
		if errors.Is(result.Error, gorm.ErrRecordNotFound) {
			c.JSON(http.StatusNotFound, gin.H{
				"success": false,
				"message": "位置不存在",
			})
			util.Logger.WithFields(logrus.Fields{
				"error":       result.Error,
				"location_id": location.ID,
			}).Error("位置不存在")
			return
		}
	}
	//检查是否同一个酒店
	UserHotelId := c.GetUint("hotel_id")
	if UserHotelId != ex.HotelID {
		c.JSON(http.StatusBadRequest, gin.H{
			"success": false,
			"message": "位置不属于该酒店",
		})
		util.Logger.WithFields(logrus.Fields{
			"user_hotel_id":     UserHotelId,
			"location_hotel_id": ex.HotelID,
		}).Error("位置不属于该酒店")
		return
	}
	//检查位置下是否还有行李
	if len(location.Luggage) > 0 {
		c.JSON(http.StatusBadRequest, gin.H{
			"success": false,
			"message": "位置下还有行李，无法删除",
		})
		util.Logger.WithFields(logrus.Fields{
			"error":       result.Error,
			"location_id": location.ID,
		}).Error("位置下还有行李，无法删除")
		return
	}
	//删除位置
	result = s.DB.Model(&models.Location{}).Where("id = ?", location.ID).Delete(&models.Location{})
	if result.Error != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"success": false,
			"message": "内部错误",
		})
		util.Logger.WithFields(logrus.Fields{
			"error":       result.Error,
			"location_id": location.ID,
		}).Error("位置数据库删除错误")
		return
	}
	c.JSON(http.StatusOK, gin.H{
		"success": true,
		"message": "位置删除成功",
	})
}

// DeleteLuggageStorageByCode 根据行李寄存码删除行李
func DeleteLuggageStorageByCode(c *gin.Context, s *services.Services) {

	var luggage models.LuggageStorage
	//绑定
	if err := c.ShouldBind(&luggage); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"success": false,
			"message": "请求数据格式错误",
		})
		return
	}
	//检查是否存在
	var existingLuggage models.LuggageStorage
	if err := s.DB.
		Preload("Guest").
		Preload("Luggage").
		Where("pick_up_code = ?", luggage.PickUpCode).First(&existingLuggage).Error; err != nil {
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
			util.Logger.WithFields(logrus.Fields{
				"error":      err,
				"luggage_id": luggage.ID,
			}).Error("查询行李失败")
		}
		return
	}
	//开启事务
	tx := s.DB.Begin()

	//删除客户
	result := tx.Model(&models.Guest{}).Where("id = ?", existingLuggage.GuestID).Delete(&models.Guest{})
	if result.Error != nil {
		tx.Rollback()
		c.JSON(http.StatusInternalServerError, gin.H{
			"success": false,
			"message": "删除行李失败",
		})
		util.Logger.WithFields(logrus.Fields{
			"error":      result.Error,
			"luggage_id": luggage.ID,
		}).Error("删除行李失败")
		return
	}
	//删除行李寄存表
	result = tx.Model(&models.LuggageStorage{}).Where("pick_up_code = ?", existingLuggage.PickUpCode).Delete(&models.LuggageStorage{})
	if result.Error != nil {
		tx.Rollback()
		c.JSON(http.StatusInternalServerError, gin.H{
			"success": false,
			"message": "删除行李失败",
		})
		util.Logger.WithFields(logrus.Fields{
			"error":      result.Error,
			"luggage_id": luggage.ID,
		}).Error("删除行李失败")
		return
	}
	//提交事务
	err := tx.Commit().Error
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"success": false,
			"message": "删除行李失败",
		})
		util.Logger.WithFields(logrus.Fields{
			"error":      err,
			"luggage_id": luggage.ID,
		}).Error("事务提交失败")
		return
	}
	//删除redis缓存，更出新的取件码
	s.RdbRand.Del(c, fmt.Sprintf("%d:%s", luggage.HotelID, luggage.PickUpCode))

	c.JSON(http.StatusOK, gin.H{
		"success": true,
		"message": "行李删除成功",
	})
}
