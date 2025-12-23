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
		Where("id = ?", luggage.ID).First(&existingLuggage).Error; err != nil {
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
	luggage.Status = "已取出"
	//更新行李寄存表
	result := tx.Model(&models.LuggageStorage{}).Where("id = ?", luggage.ID).Updates(luggage)
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
	//删除行李
	result2 := tx.Where("id = ?", luggage.ID).Delete(&models.LuggageStorage{})
	if result2.Error != nil {
		tx.Rollback()
		c.JSON(http.StatusInternalServerError, gin.H{
			"success": false,
			"message": "删除行李失败",
		})
		util.Logger.WithFields(logrus.Fields{
			"error":      result2.Error,
			"luggage_id": luggage.ID,
		}).Error("删除行李失败")
		return
	}
	//删除客户
	result = tx.Model(&models.Guest{}).Where("id = ?", existingLuggage.GuestID).Delete(&models.Guest{})
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
	s.RdbRand.Del(c, fmt.Sprintf("%d:%s", luggage.HotelID, luggage.PickUpCode))

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
