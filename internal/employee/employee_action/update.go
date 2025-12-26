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
		CheckField: []string{"ID"},
	})
}

// UpdateLuggage 根据id，更新行李
func UpdateLuggage(c *gin.Context, s *services.Services) {
	util.Update(c, s.DB, util.RequestList{
		Model:      &models.Luggage{},
		CheckExist: true,
		CheckType:  "id",
		CheckField: []string{"ID"},
	})
}
