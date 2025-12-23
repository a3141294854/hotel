package util

import (
	"github.com/gin-gonic/gin"
	"github.com/sirupsen/logrus"
	"gorm.io/gorm"
	"net/http"
)

type RequestList struct {
	Model          interface{} // 模型
	CheckExist     bool        // 是否检查存在
	CheckExistType string      //判断是否存档的数据的类型
	CheckField     []string    // 检查必要的字段名
	Preloads       []string    // 预加载
}

func Create(db *gorm.DB, c *gin.Context, list RequestList) {
	req := list.Model
	err := c.ShouldBind(req)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"success": false,
			"message": "请求数据格式错误",
		})
		return
	}

	if list.CheckExist {
		ty := list.CheckExistType
		ok, err := ExIfByField(db, ty, req)

		if ok && err == nil {
			c.JSON(http.StatusConflict, gin.H{
				"success": false,
				"message": "数据已存在",
			})
			return
		}

		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{
				"success": false,
				"message": "内部错误",
			})
			Logger.WithFields(logrus.Fields{
				"error": err,
			}).Error("检查数据失败")
			return
		}

	}

	result := db.Create(req)
	if result.Error != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"success": false,
			"message": "创建失败",
		})
		Logger.WithFields(logrus.Fields{
			"error": result.Error,
		}).Error("创建数据失败")
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"success": true,
		"message": "创建成功",
		"data":    req,
	})

}
