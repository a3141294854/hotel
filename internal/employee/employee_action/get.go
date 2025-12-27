package employee_action

import (
	"github.com/gin-gonic/gin"
	"github.com/sirupsen/logrus"
	"hotel/internal/util"
	"hotel/models"
	"hotel/services"
	"net/http"
)

// GetPickUpCode 通过取件码，获取行李寄存表
func GetPickUpCode(c *gin.Context, s *services.Services) {
	util.Get(c, s.DB, util.RequestList{
		Model:     &models.LuggageStorage{},
		CheckType: "pick_up_code",
		Preloads:  []string{"Guest", "Luggage", "Luggage.Tag", "Luggage.Location", "Photos"},
	})
}

// GetName 通过用户姓名，获取行李寄存表
func GetName(c *gin.Context, s *services.Services) {
	util.Get(c, s.DB, util.RequestList{
		Model:     &models.LuggageStorage{},
		CheckType: "guest_name",
		GetType:   "Guest.Name",
		Preloads:  []string{"Guest", "Luggage", "Luggage.Tag", "Luggage.Location", "Photos"},
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
		Preload("Photos").
		Find(&luggage)
	if result.Error != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"success": false,
			"message": "获取行李失败",
		})
		util.Logger.WithFields(logrus.Fields{
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

// GetGuestID 通过用户id,获取行李寄存表
func GetGuestID(c *gin.Context, s *services.Services) {
	util.Get(c, s.DB, util.RequestList{
		Model:     &models.LuggageStorage{},
		CheckType: "guest_id",
		Preloads:  []string{"Guest", "Luggage", "Luggage.Tag", "Luggage.Location", "Photos"},
	})
}

// GetLocation 根据名字,获取一个地方的行李
func GetLocation(c *gin.Context, s *services.Services) {
	util.Get(c, s.DB, util.RequestList{
		Model:      &models.Location{},
		CheckExist: true,
		CheckType:  "name",
		Preloads:   []string{"Luggage"},
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
		util.Logger.WithFields(logrus.Fields{
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
		util.Logger.WithFields(logrus.Fields{
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
