package employee_action

import (
	"github.com/gin-gonic/gin"
	"hotel/internal/util"
	"hotel/models"
	"hotel/services"
)

// UpdateLuggageStorage 根据id，更新行李寄存表
func UpdateLuggageStorage(c *gin.Context, s *services.Services) {
	util.Update(c, s.DB, util.RequestList{
		Model:      &models.LuggageStorage{},
		CheckExist: true,
		CheckType:  "id",
	})
}

// UpdateLuggage 根据id，更新行李
func UpdateLuggage(c *gin.Context, s *services.Services) {
	util.Update(c, s.DB, util.RequestList{
		Model:      &models.Luggage{},
		CheckExist: true,
		CheckType:  "id",
	})
}
func UpdateTag(c *gin.Context, s *services.Services) {
	util.Update(c, s.DB, util.RequestList{
		Model:      &models.Tag{},
		CheckExist: true,
		CheckType:  "id",
	})
}

func UpdateLocation(c *gin.Context, s *services.Services) {
	util.Update(c, s.DB, util.RequestList{
		Model:      &models.Location{},
		CheckExist: true,
		CheckType:  "id",
	})
}

func UpdateHotel(c *gin.Context, s *services.Services) {
	util.Update(c, s.DB, util.RequestList{
		Model:      &models.Hotel{},
		CheckExist: true,
		CheckType:  "id",
	})
}
