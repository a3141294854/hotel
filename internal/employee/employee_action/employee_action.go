package employee_action

import (
	"encoding/json"
	"errors"
	"fmt"
	"hotel/internal/util"
	"net/http"
	"strconv"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/sirupsen/logrus"
	"gorm.io/gorm"
	"hotel/internal/util/logger"
	"hotel/models"
	"hotel/services"
)

// AddLuggage 添加行李
func AddLuggage(c *gin.Context, s *services.Services) {

	var req models.LuggageStorage
	if err := c.ShouldBind(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"success": false,
			"message": "请求数据格式错误",
		})
		logger.Logger.WithFields(logrus.Fields{
			"error": err.Error(),
		}).Error("行李数据绑定错误")
		return
	}
	fmt.Println(req)

	if req.GuestName == "" || req.GuestPhone == "" || req.GuestRoom == "" {
		c.JSON(http.StatusBadRequest, gin.H{
			"success": false,
			"message": "请求数据格式错误",
		})
		fmt.Println(req.GuestName, req.GuestPhone, req.GuestRoom)
		return
	}

	// 先创建或获取客户记录
	var guest models.Guest
	result := s.DB.Where("name = ?", req.GuestName).First(&guest)
	if result.Error != nil {
		if errors.Is(result.Error, gorm.ErrRecordNotFound) {
			// 客户不存在，创建新客户
			guest = models.Guest{
				Name: req.GuestName,
			}
			if err := s.DB.Create(&guest).Error; err != nil {
				c.JSON(http.StatusInternalServerError, gin.H{
					"success": false,
					"message": "创建客户记录失败",
				})
				logger.Logger.WithFields(logrus.Fields{
					"error":      err,
					"guest_name": req.GuestName,
				}).Error("创建客户记录失败")
				return
			}
		} else {
			// 其他查询错误
			c.JSON(http.StatusInternalServerError, gin.H{
				"success": false,
				"message": "查询客户失败",
			})
			logger.Logger.WithFields(logrus.Fields{
				"error":      result.Error,
				"guest_name": req.GuestName,
			}).Error("查询客户失败")
			return
		}
	}
	a, _ := c.Get("hotel_id")
	req.HotelID = a.(uint)

	b, _ := c.Get("employee_id")
	req.OperatorID = b.(uint)

	d, _ := c.Get("employee_name")
	req.OperatorName = d.(string)

	req.Status = "寄存中"
	code, err := util.GeneratePickUpCode(s, a.(uint))
	req.PickUpCode = code
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"success": false,
			"message": "生成取件码失败",
		})
		logger.Logger.WithFields(logrus.Fields{
			"error": err.Error(),
		}).Error("生成取件码失败")
		return
	}

	result = s.DB.Create(&req)
	if result.Error != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"success": false,
			"message": "创建行李记录失败",
		})
		logger.Logger.WithFields(logrus.Fields{
			"error":      result.Error,
			"guest_name": req.GuestName,
		}).Error("创建行李记录失败")
		return
	}

	c.JSON(http.StatusCreated, gin.H{
		"success": true,
		"message": "行李添加成功",
		"data":    req,
	})

}

// Delete 删除指定资源
func Delete(c *gin.Context, s *services.Services) {
	var luggage models.Luggage
	if err := c.ShouldBind(&luggage); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"success": false,
			"message": "请求数据格式错误",
		})
		return
	}
	var existingLuggage models.Luggage
	if err := s.DB.Preload("luggage_storage").Where("status = ?", "寄存中").First(&existingLuggage, luggage.ID).Error; err != nil {
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
	result := tx.Model(&models.Luggage{}).Where("id = ?", luggage.ID).Updates(luggage)
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

	result2 := tx.Where("id = ?", luggage.ID).Delete(&models.Luggage{})
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
	s.RdbRand.Del(c, fmt.Sprintf("%d:%s", luggage.HotelID, luggage.LuggageStorage.PickUpCode))

	c.JSON(http.StatusOK, gin.H{
		"success": true,
		"message": "行李删除成功",
	})

}

func Update(c *gin.Context, s *services.Services) {
	var luggage models.Luggage
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
	var existingLuggage models.Luggage
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
	result := s.DB.Model(&existingLuggage).Updates(luggage)
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

	// 检查是否真的更新了记录
	if result.RowsAffected == 0 {
		c.JSON(http.StatusBadRequest, gin.H{
			"success": false,
			"message": "没有数据需要更新",
		})
		return
	}

	// 获取客户信息
	var guest models.Guest
	s.DB.First(&guest, existingLuggage.GuestID)

	s.RdbCac.Set(c, strconv.Itoa(int(existingLuggage.ID)), guest.Name, time.Minute*15)

	c.JSON(http.StatusOK, gin.H{
		"success": true,
		"message": "行李更新成功",
		"data": gin.H{
			"id":         existingLuggage.ID,
			"guest_id":   existingLuggage.GuestID,
			"guest_name": guest.Name,
			"tag":        existingLuggage.Tag,
			"status":     existingLuggage.Status,
			"location":   existingLuggage.Location,
		},
	})
}

// GetName 获取行李
func GetName(c *gin.Context, s *services.Services) {
	var luggage []models.Luggage

	guestName := c.Query("guest_name")
	if guestName == "" {
		c.JSON(http.StatusBadRequest, gin.H{
			"success": false,
			"message": "请输入客户姓名",
		})
		return
	}

	if val, err := s.RdbCac.Get(c, guestName).Result(); err == nil {
		var luggage []models.Luggage
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

	result := s.DB.Preload("Guest").
		Joins("JOIN guests ON luggages.guest_id = guests.id").
		Where("guests.guest_name = ? AND luggages.status = ?", guestName, "寄存中").
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
			"error":      err,
			"guest_name": guestName,
		}).Error("JSON序列化失败")
	} else {
		s.RdbCac.Set(c, guestName, val, time.Minute*15)
	}

	c.JSON(http.StatusOK, gin.H{
		"success": true,
		"message": "行李获取成功",
		"data":    luggage,
		"count":   len(luggage),
	})
}

// GetAll 获取所有行李
func GetAll(c *gin.Context, s *services.Services) {
	var luggage []models.Luggage
	result := s.DB.
		Preload("Guest").
		Where("status = ?", "寄存中").
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

// GetGuestID 获取行李
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
		var luggage []models.Luggage
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

	var luggage []models.Luggage
	result := s.DB.
		Preload("Guest").
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

// GetLocation 获取行李
func GetLocation(c *gin.Context, s *services.Services) {
	location := c.Query("location")
	if location == "" {
		c.JSON(http.StatusBadRequest, gin.H{
			"success": false,
			"message": "请输入存放地点",
		})
		return
	}

	var luggage []models.Luggage
	result := s.DB.
		Preload("Guest"). // 预加载Guest信息
		Where("location = ? AND status = ?", location, "寄存中").
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
			"message": "未找到该存放地点的行李",
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

// GetStatus 获取行李
func GetStatus(c *gin.Context, s *services.Services) {
	status := c.Query("status")
	if status == "" {
		c.JSON(http.StatusBadRequest, gin.H{
			"success": false,
			"message": "请输入状态",
		})
		return
	}

	var luggage []models.Luggage
	result := s.DB.
		Preload("Guest").
		Where("status = ?", status).
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
			"message": "未找到该状态的行李",
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

// GetAdvance 高级查询
func GetAdvance(c *gin.Context, s *services.Services) {

	type AdvanceRequest struct {
		GuestName string  `json:"guest_name"`
		GuestID   uint    `json:"guest_id"`
		Location  string  `json:"location"`
		Status    string  `json:"status"`
		Tag       string  `json:"tag"`
		Weight    float32 `json:"weight"`
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

	var luggage []models.Luggage

	query := s.DB.Model(&models.Luggage{})

	if req.GuestName != "" {
		query = query.Joins("JOIN guests ON luggages.guest_id = guests.id").
			Where("guests.guest_name = ?", req.GuestName)
	}

	if req.GuestID != 0 {
		query = query.Where("guest_id = ?", req.GuestID)
	}

	if req.Location != "" {
		query = query.Where("location = ?", req.Location)
	}

	if req.Status != "" {
		query = query.Where("status = ?", req.Status)
	} else {
		query = query.Where("status = ?", "寄存中")
	}

	if req.Tag != "" {
		query = query.Where("tag = ?", req.Tag)
	}

	if req.Weight != 0 {
		query = query.Where("weight = ?", req.Weight)
	}

	result := query.Find(&luggage)
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

// CountSum 获取总行李数量
func CountSum(c *gin.Context, s *services.Services) {
	var count int64
	s.DB.Model(&models.Luggage{}).Where("status = ?", "寄存中").Count(&count)

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
	s.DB.Model(&models.Luggage{}).
		Where("DATE(created_at) = ?", today).
		Count(&todayAdded)

	// 统计今天取出的行李
	var todayTaken int64
	s.DB.Model(&models.Luggage{}).
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
