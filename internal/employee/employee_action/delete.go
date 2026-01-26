package employee_action

import (
	"errors"
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/sirupsen/logrus"
	"gorm.io/gorm"
	"hotel/internal/util"
	"hotel/models"
	"hotel/services"
	"net/http"
	"time"
)

// DeleteLuggageStorage 删除行李寄存表
func DeleteLuggageStorage(c *gin.Context, s *services.Services) {

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
			return
		} else {
			c.JSON(http.StatusInternalServerError, gin.H{
				"success": false,
				"message": "查询行李失败",
			})
			util.Logger.WithFields(logrus.Fields{
				"error":      err,
				"luggage_id": id,
				"请求id":       c.GetUint("request_id"),
			}).Error("查询行李失败")
		}
		return
	}
	//开启事务
	tx := s.DB.Begin()

	/*删除客户
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
			"请求id":       c.GetUint("request_id"),
		}).Error("删除行李失败")
		return
	}*/

	//删除行李
	for i := range existingLuggage.Luggage {
		luggage := &existingLuggage.Luggage[i]
		result := tx.Model(&models.Luggage{}).Where("id = ?", luggage.ID).Delete(&models.Luggage{})
		if result.Error != nil {
			tx.Rollback()
			c.JSON(http.StatusInternalServerError, gin.H{
				"success": false,
				"message": "删除行李失败",
			})
			util.Logger.WithFields(logrus.Fields{
				"error":      result.Error,
				"luggage_id": id,
				"请求id":       c.GetUint("request_id"),
			}).Error("删除行李失败")
			return
		}
	}
	//删除行李寄存
	result := tx.Model(&models.LuggageStorage{}).Where("id = ?", existingLuggage.ID).Delete(&models.LuggageStorage{})
	if result.Error != nil {
		tx.Rollback()
		c.JSON(http.StatusInternalServerError, gin.H{
			"success": false,
			"message": "删除行李失败",
		})
		util.Logger.WithFields(logrus.Fields{
			"error":      result.Error,
			"luggage_id": id,
			"请求id":       c.GetUint("request_id"),
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
			"请求id":       c.GetUint("request_id"),
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

// CheckOutLuggageStorage 取出行李寄存表
func CheckOutLuggageStorage(c *gin.Context, s *services.Services) {
	code := c.Param("code")

	if code == "" {
		c.JSON(http.StatusBadRequest, gin.H{
			"success": false,
			"message": "取件码不能为空",
		})
		return
	}

	var LuggageStorage models.LuggageStorage
	result := s.DB.Model(&models.LuggageStorage{}).Where("pick_up_code = ?", code).Preload("Luggage").First(&LuggageStorage)
	if result.Error != nil {
		if errors.Is(result.Error, gorm.ErrRecordNotFound) {
			c.JSON(http.StatusNotFound, gin.H{
				"success": false,
				"message": "行李记录不存在",
			})
			return
		}
		util.Logger.WithFields(logrus.Fields{
			"error":      result.Error,
			"luggage_id": LuggageStorage.ID,
			"请求id":       c.GetUint("request_id"),
		})
		c.JSON(http.StatusInternalServerError, gin.H{
			"success": false,
			"message": "查询行李失败",
		})
		return
	}
	now := time.Now()
	LuggageStorage.CheckOutAt = &now

	tx := s.DB.Begin()
	//更新行李寄存
	result = tx.Model(&LuggageStorage).Where("id = ?", LuggageStorage.ID).Updates(LuggageStorage)
	if result.Error != nil {
		tx.Rollback()
		c.JSON(http.StatusInternalServerError, gin.H{
			"success": false,
			"message": "更新行李失败",
		})
		return
	}
	// 更新所有关联的行李
	for i := range LuggageStorage.Luggage {
		luggage := &LuggageStorage.Luggage[i] // ← 使用索引获取指针
		luggage.CheckOutAt = &now

		result = tx.Model(luggage).Where("id = ?", luggage.ID).Updates(luggage)
		if result.Error != nil {
			tx.Rollback()
			c.JSON(http.StatusInternalServerError, gin.H{
				"success": false,
				"message": "更新行李失败",
			})
			util.Logger.WithFields(logrus.Fields{
				"error":      result.Error,
				"luggage_id": luggage.ID,
				"请求id":       c.GetUint("request_id"),
			}).Error("更新行李失败")
			return
		}

		result = tx.Model(luggage).Where("id = ?", luggage.ID).Delete(luggage)
		if result.Error != nil {
			tx.Rollback()
			c.JSON(http.StatusInternalServerError, gin.H{
				"success": false,
				"message": "删除行李失败",
			})
			util.Logger.WithFields(logrus.Fields{
				"error":      result.Error,
				"luggage_id": luggage.ID,
				"请求id":       c.GetUint("request_id"),
			}).Error("删除行李失败")
			return
		}
	}

	//删除行李寄存
	result = tx.Model(&models.LuggageStorage{}).Where("id = ?", LuggageStorage.ID).Delete(&models.LuggageStorage{})
	if result.Error != nil {
		tx.Rollback()
		c.JSON(http.StatusInternalServerError, gin.H{
			"success": false,
			"message": "删除行李失败",
		})
		util.Logger.WithFields(logrus.Fields{
			"error":      result.Error,
			"luggage_id": LuggageStorage.ID,
			"请求id":       c.GetUint("request_id"),
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
			"luggage_id": LuggageStorage.ID,
			"请求id":       c.GetUint("request_id"),
		}).Error("事务提交失败")
		return
	}
	//删除redis缓存，更出新的取件码
	s.RdbRand.Del(c, fmt.Sprintf("%d:%s", LuggageStorage.HotelID, LuggageStorage.PickUpCode))

	c.JSON(http.StatusOK, gin.H{
		"success": true,
		"message": "行李删除成功",
	})
}

// DeleteLuggage 删除行李
func DeleteLuggage(c *gin.Context, s *services.Services) {
	id := c.Param("id")
	//检查是否存在
	var ex models.Luggage
	result := s.DB.Model(&models.Luggage{}).
		Preload("LuggageStorage").
		Where("id = ?", id).First(&ex)
	if result.Error != nil {
		if errors.Is(result.Error, gorm.ErrRecordNotFound) {
			c.JSON(http.StatusNotFound, gin.H{
				"success": false,
				"message": "行李不存在",
			})
			util.Logger.WithFields(logrus.Fields{
				"error":      result.Error,
				"luggage_id": id,
				"请求id":       c.GetUint("request_id"),
			}).Error("行李不存在")
			return
		}
		c.JSON(http.StatusInternalServerError, gin.H{
			"success": false,
			"message": "内部错误",
		})
		util.Logger.WithFields(logrus.Fields{
			"error":      result.Error,
			"luggage_id": id,
			"请求id":       c.GetUint("request_id"),
		}).Error("行李数据库查询错误")
		return
	}
	//开启事务
	tx := s.DB.Begin()
	//删除行李
	result = tx.Model(&models.Luggage{}).Where("id = ?", id).Delete(&models.Luggage{})
	if result.Error != nil {
		tx.Rollback()
		c.JSON(http.StatusInternalServerError, gin.H{
			"success": false,
			"message": "内部错误",
		})
		util.Logger.WithFields(logrus.Fields{
			"error":      result.Error,
			"luggage_id": id,
			"请求id":       c.GetUint("request_id"),
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
			"luggage_id": id,
			"请求id":       c.GetUint("request_id"),
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
				"luggage_id": id,
				"请求id":       c.GetUint("request_id"),
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
			"luggage_id": id,
			"请求id":       c.GetUint("request_id"),
		}).Error("事务提交失败")
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"success": true,
		"message": "行李删除成功",
	})

}

// CheckOutLuggage 取出行李
func CheckOutLuggage(c *gin.Context, s *services.Services) {
	id := c.Param("id")

	//查找有无行李
	var Luggage models.Luggage
	result := s.DB.Model(&models.Luggage{}).Where("id = ?", id).First(&Luggage)
	if result.Error != nil {
		if errors.Is(result.Error, gorm.ErrRecordNotFound) {
			c.JSON(http.StatusNotFound, gin.H{
				"success": false,
				"message": "行李不存在",
			})
			util.Logger.WithFields(logrus.Fields{
				"error":      result.Error,
				"luggage_id": id,
				"请求id":       c.GetUint("request_id"),
			}).Error("行李不存在")
			return
		}
		util.Logger.WithFields(logrus.Fields{
			"error":      result.Error,
			"luggage_id": id,
			"请求id":       c.GetUint("request_id"),
		})
		c.JSON(http.StatusInternalServerError, gin.H{
			"success": false,
			"message": "内部错误",
		})
		return
	}

	tx := s.DB.Begin()
	//修改时间
	now := time.Now()

	result = tx.Model(&models.Luggage{}).Where("id = ?", id).Update("check_out_at", &now)
	if result.Error != nil {
		tx.Rollback()
		util.Logger.WithFields(logrus.Fields{
			"error":      result.Error,
			"luggage_id": id,
			"请求id":       c.GetUint("request_id"),
		})
		c.JSON(http.StatusInternalServerError, gin.H{
			"success": false,
			"message": "内部错误",
		})
		return
	}

	result = tx.Model(&models.Luggage{}).Where("id = ?", id).Delete(&models.Luggage{})
	if result.Error != nil {
		tx.Rollback()
		util.Logger.WithFields(logrus.Fields{
			"error":      result.Error,
			"luggage_id": id,
			"请求id":       c.GetUint("request_id"),
		})
		c.JSON(http.StatusInternalServerError, gin.H{
			"success": false,
			"message": "内部错误",
		})
		return
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
			"luggage_id": id,
			"请求id":       c.GetUint("request_id"),
		}).Error("事务提交失败")
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"success": true,
		"message": "行李取出成功",
	})
}

// DeleteLocation 删除寄存室
func DeleteLocation(c *gin.Context, s *services.Services) {
	ID := c.Param("id")
	//检查是否存
	var ex models.Location
	result := s.DB.Model(&models.Location{}).
		Preload("Luggage").
		Where("id = ?", ID).First(&ex)
	if result.Error != nil {
		if errors.Is(result.Error, gorm.ErrRecordNotFound) {
			c.JSON(http.StatusNotFound, gin.H{
				"success": false,
				"message": "位置不存在",
			})
			util.Logger.WithFields(logrus.Fields{
				"error":       result.Error,
				"location_id": ID,
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
	if len(ex.Luggage) > 0 {
		c.JSON(http.StatusBadRequest, gin.H{
			"success": false,
			"message": "位置下还有行李，无法删除",
		})
		util.Logger.WithFields(logrus.Fields{
			"error":       result.Error,
			"location_id": ID,
		}).Error("位置下还有行李，无法删除")
		return
	}
	//删除位置
	result = s.DB.Model(&models.Location{}).Where("id = ?", ID).Delete(&models.Location{})
	if result.Error != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"success": false,
			"message": "内部错误",
		})
		util.Logger.WithFields(logrus.Fields{
			"error":       result.Error,
			"location_id": ID,
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

	PickUpCode := c.Param("pick_up_code")
	//检查是否存在
	var existingLuggage models.LuggageStorage
	if err := s.DB.
		Preload("Guest").
		Preload("Luggage").
		Where("pick_up_code = ?", PickUpCode).First(&existingLuggage).Error; err != nil {
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
				"luggage_id": existingLuggage.ID,
				"请求id":       c.GetUint("request_id"),
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
			"luggage_id": existingLuggage.ID,
			"请求id":       c.GetUint("request_id"),
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
			"luggage_id": existingLuggage.ID,
			"请求id":       c.GetUint("request_id"),
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
			"luggage_id": existingLuggage.ID,
			"请求id":       c.GetUint("request_id"),
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

func DeleteTag(c *gin.Context, s *services.Services) {
	id := c.Param("id")
	if id == "" {
		c.JSON(http.StatusBadRequest, gin.H{
			"success": false,
			"message": "id不能为空",
		})
	}

	var tag models.Tag
	result := s.DB.Model(&models.Tag{}).Where("id = ?", id).First(&tag)
	if result.Error != nil {
		if errors.Is(result.Error, gorm.ErrRecordNotFound) {
			c.JSON(http.StatusNotFound, gin.H{
				"success": false,
				"message": "标签不存在",
			})
			util.Logger.WithFields(logrus.Fields{
				"error":  result.Error,
				"tag_id": id,
				"请求id":   c.GetUint("request_id"),
			}).Error("标签不存在")
			return
		}
		c.JSON(http.StatusInternalServerError, gin.H{
			"success": false,
			"message": "查询标签失败",
		})
		util.Logger.WithFields(logrus.Fields{
			"error":  result.Error,
			"tag_id": id,
			"请求id":   c.GetUint("request_id"),
		}).Error("查询标签失败")
		return
	}

	result = s.DB.Model(&models.Luggage{}).Where("tag_id = ?", tag.ID).Delete(&models.Luggage{})
	if !errors.Is(result.Error, gorm.ErrRecordNotFound) {
		c.JSON(http.StatusNotFound, gin.H{
			"success": false,
			"message": "标签下有行李",
		})
		util.Logger.WithFields(logrus.Fields{
			"error":  result.Error,
			"tag_id": id,
			"请求id":   c.GetUint("request_id"),
		}).Error("标签下有行李")
		return
	}

	result = s.DB.Model(&models.Tag{}).Where("id = ?", id).Delete(&models.Tag{})
	if result.Error != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"success": false,
			"message": "删除标签失败",
		})
		util.Logger.WithFields(logrus.Fields{
			"error":  result.Error,
			"tag_id": id,
			"请求id":   c.GetUint("request_id"),
		}).Error("删除标签失败")
		return
	}
	c.JSON(http.StatusOK, gin.H{
		"success": true,
		"message": "标签删除成功",
	})
}

func DeleteHotel(c *gin.Context, s *services.Services) {
	id := c.Param("id")
	if id == "" {
		c.JSON(http.StatusBadRequest, gin.H{
			"success": false,
			"message": "id不能为空",
		})
	}

	var hotel models.Hotel
	result := s.DB.Model(&models.Hotel{}).Where("id = ?", id).First(&hotel)
	if result.Error != nil {
		if errors.Is(result.Error, gorm.ErrRecordNotFound) {
			c.JSON(http.StatusNotFound, gin.H{
				"success": false,
				"message": "酒店不存在",
			})
			util.Logger.WithFields(logrus.Fields{
				"error":    result.Error,
				"hotel_id": id,
				"请求id":     c.GetUint("request_id"),
			}).Error("酒店不存在")
			return
		}
		c.JSON(http.StatusInternalServerError, gin.H{
			"success": false,
			"message": "查询酒店失败",
		})
		util.Logger.WithFields(logrus.Fields{
			"error":    result.Error,
			"hotel_id": id,
			"请求id":     c.GetUint("request_id"),
		}).Error("查询酒店失败")
		return
	}

	result = s.DB.Model(&models.Hotel{}).Where("id = ?", id).Delete(&models.Hotel{})
	if result.Error != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"success": false,
			"message": "删除酒店失败",
		})
		util.Logger.WithFields(logrus.Fields{
			"error":    result.Error,
			"hotel_id": id,
			"请求id":     c.GetUint("request_id"),
		}).Error("删除酒店失败")
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"success": true,
		"message": "酒店删除成功",
	})

}
