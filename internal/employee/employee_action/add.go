package employee_action

import (
	"errors"
	"hotel/internal/util"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/sirupsen/logrus"
	"gorm.io/gorm"
	"hotel/internal/util/logger"
	"hotel/models"
	"hotel/services"
)

// AddLuggage 添加行李
func AddLuggage(c *gin.Context, s *services.Services) {

	var req struct {
		BagCount      int `json:"bag_count"`
		BackpackCount int `json:"backpack_count"`
		BoxCount      int `json:"box_count"`
		OtherCount    int `json:"other_count"`

		GuestName  string `json:"guest_name"`
		GuestPhone string `json:"guest_phone"`
		GuestRoom  string `json:"guest_room"`

		LocationID int `json:"location_id"`

		Macs []string `json:"macs"`

		Status string `json:"status"`
		Remark string `json:"remark"`
	}

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

	if req.GuestName == "" || req.GuestPhone == "" || req.GuestRoom == "" {
		c.JSON(http.StatusBadRequest, gin.H{
			"success": false,
			"message": "客户信息不能为空",
		})
		return
	}
	if req.BagCount == 0 && req.BackpackCount == 0 && req.BoxCount == 0 && req.OtherCount == 0 {
		c.JSON(http.StatusBadRequest, gin.H{
			"success": false,
			"message": "行李数量不能为0",
		})
		return
	}

	// 先创建或获取客户记录
	var guest models.Guest
	result := s.DB.Where("name = ?", req.GuestName).First(&guest)
	if result.Error != nil {
		if errors.Is(result.Error, gorm.ErrRecordNotFound) {
			// 客户不存在，创建新客户
			guest = models.Guest{
				Name:  req.GuestName,
				Phone: req.GuestPhone,
				Room:  req.GuestRoom,
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

	//创建行李表记录
	insert := models.LuggageStorage{
		GuestID:       guest.ID,
		BagCount:      req.BagCount,
		BackpackCount: req.BackpackCount,
		BoxCount:      req.BoxCount,
		Status:        req.Status,
		Remark:        req.Remark,
	}
	//存默认值和获取
	a, _ := c.Get("hotel_id")
	insert.HotelID = a.(uint)

	b, _ := c.Get("employee_id")
	insert.OperatorID = b.(uint)

	d, _ := c.Get("employee_name")
	insert.OperatorName = d.(string)

	insert.Status = "寄存中"
	//默认设置为前台
	if req.LocationID == 0 {
		insert.LocationID = 1
	} else {
		insert.LocationID = uint(req.LocationID)
	}

	code, err := util.GeneratePickUpCode(s, a.(uint))
	insert.PickUpCode = code
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

	result = s.DB.Model(&models.LuggageStorage{}).Create(&insert)
	if result.Error != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"success": false,
			"message": "创建行李记录失败",
		})
		logger.Logger.WithFields(logrus.Fields{
			"error": result.Error,
		}).Error("创建行李记录失败")
		return
	}

	//创建行李mac表记录
	for _, mac := range req.Macs {
		var tag models.Tag
		result = s.DB.Model(&models.Tag{}).Where("mac = ?", mac).First(&tag)
		if result.Error != nil {
			if errors.Is(result.Error, gorm.ErrRecordNotFound) {
				c.JSON(http.StatusNotFound, gin.H{
					"success": false,
					"message": "行李mac不存在",
				})
				return
			} else {
				c.JSON(http.StatusInternalServerError, gin.H{
					"success": false,
					"message": "查询行李mac失败",
				})
				logger.Logger.WithFields(logrus.Fields{
					"error": result.Error,
				}).Error("查询行李mac失败")
				return
			}
		}

		luggage := models.Luggage{
			LuggageStorageID: insert.ID,
			TagID:            tag.ID,
			Status:           "寄存中",
		}

		result = s.DB.Model(&models.Luggage{}).Create(&luggage)
		if result.Error != nil {
			c.JSON(http.StatusInternalServerError, gin.H{
				"success": false,
				"message": "创建行李mac记录失败",
			})
			logger.Logger.WithFields(logrus.Fields{
				"error": result.Error,
			}).Error("创建行李mac记录失败")
			return
		}
	}

	c.JSON(http.StatusCreated, gin.H{
		"success": true,
		"message": "行李添加成功",
		"data":    insert,
	})

}

// AddMac 添加mac
func AddMac(c *gin.Context, s *services.Services) {
	var req models.Tag
	if err := c.ShouldBind(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"success": false,
			"message": "请求数据格式错误",
		})
		return
	}
	if req.Mac == "" {
		c.JSON(http.StatusBadRequest, gin.H{
			"success": false,
			"message": "mac不能为空",
		})
		return
	}

	var ex models.Tag
	result := s.DB.Model(&models.Tag{}).Where("mac = ?", req.Mac).First(&ex)
	if result.Error == nil {
		c.JSON(http.StatusConflict, gin.H{
			"success": false,
			"message": "mac已存在",
		})
		return
	}

	result = s.DB.Model(&models.Tag{}).Create(&req)
	if result.Error != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"success": false,
			"message": "创建行李mac记录失败",
		})
		logger.Logger.WithFields(logrus.Fields{
			"error": result.Error,
		}).Error("创建行李mac记录失败")
		return
	}
	c.JSON(http.StatusCreated, gin.H{
		"success": true,
		"message": "行李mac添加成功",
		"data":    req,
	})

}
