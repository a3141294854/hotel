package util

import (
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/sirupsen/logrus"
	"gorm.io/gorm"
	"net/http"
	"reflect"
)

type RequestList struct {
	Model      interface{} // 模型
	CheckExist bool        // 是否检查存在
	CheckType  string      //判断是否存在的数据的类型，要小写，查结构体会有函数转换成大写,也是主要的操作对象
	CheckField []string    // 检查必要的字段名，要大写，只用来查结构体
	Preloads   []string    // 预加载
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

	if len(list.CheckField) > 0 {
		for _, v := range list.CheckField {
			if reflect.ValueOf(req).Elem().FieldByName(v).IsZero() {
				c.JSON(http.StatusBadRequest, gin.H{
					"success": false,
					"message": "字段" + v + "不能为空",
				})
				return
			}
		}
	}

	if list.CheckExist {
		ty := list.CheckType
		ok, err := ExIfByField(db, ty, req)
		if !ok && err == nil {
			c.JSON(http.StatusNotFound, gin.H{
				"success": false,
				"message": "数据不存在",
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
	query := db.Model(req)
	if len(list.Preloads) > 0 {
		for _, v := range list.Preloads {
			query = query.Preload(v)
		}

	}

	result := query.Create(req)
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

func Delete(db *gorm.DB, c *gin.Context, list RequestList) {
	req := list.Model
	err := c.ShouldBind(req)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"success": false,
			"message": "请求数据格式错误",
		})
		return
	}

	if len(list.CheckField) > 0 {
		for _, v := range list.CheckField {
			if reflect.ValueOf(req).Elem().FieldByName(v).IsZero() {
				c.JSON(http.StatusBadRequest, gin.H{
					"success": false,
					"message": "字段" + v + "不能为空",
				})
				return
			}
		}
	}

	if list.CheckExist {
		ty := list.CheckType
		ok, err := ExIfByField(db, ty, req)
		if !ok && err == nil {
			c.JSON(http.StatusNotFound, gin.H{
				"success": false,
				"message": "数据不存在",
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
	query := db.Model(req)
	if len(list.Preloads) > 0 {
		for _, v := range list.Preloads {
			query = query.Preload(v)
		}

	}

	result := query.Delete(req)
	if result.Error != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"success": false,
			"message": "删除失败",
		})
		Logger.WithFields(logrus.Fields{
			"error": result.Error,
		}).Error("删除数据失败")
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"success": true,
		"message": "删除成功",
	})
}

func Update(db *gorm.DB, c *gin.Context, list RequestList) {
	req := list.Model
	err := c.ShouldBind(req)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"success": false,
			"message": "请求数据格式错误",
		})
		return
	}

	if len(list.CheckField) > 0 {
		for _, v := range list.CheckField {
			if reflect.ValueOf(req).Elem().FieldByName(v).IsZero() {
				c.JSON(http.StatusBadRequest, gin.H{
					"success": false,
					"message": "字段" + v + "不能为空",
				})
				return
			}
		}
	}

	if list.CheckExist {
		ty := list.CheckType
		ok, err := ExIfByField(db, ty, req)
		if !ok && err == nil {
			c.JSON(http.StatusNotFound, gin.H{
				"success": false,
				"message": "数据不存在",
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
	query := db.Model(req)
	if len(list.Preloads) > 0 {
		for _, v := range list.Preloads {
			query = query.Preload(v)
		}

	}

	va := reflect.ValueOf(req).Elem().FieldByName(ConvertSnakeToCamel(list.CheckType)).Interface()
	result := query.Model(req).Where(fmt.Sprintf("%s = ?", list.CheckType), va).Updates(req)
	if result.Error != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"success": false,
			"message": "更新失败",
		})
		Logger.WithFields(logrus.Fields{
			"error": result.Error,
		}).Error("更新数据失败")
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"success": true,
		"message": "更新成功",
		"data":    req,
	})
}

func Get(db *gorm.DB, c *gin.Context, list RequestList) {
	req := list.Model
	err := c.ShouldBind(req)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"success": false,
			"message": "请求数据格式错误",
		})
		return
	}

	if len(list.CheckField) > 0 {
		for _, v := range list.CheckField {
			if reflect.ValueOf(req).Elem().FieldByName(v).IsZero() {
				c.JSON(http.StatusBadRequest, gin.H{
					"success": false,
					"message": "字段" + v + "不能为空",
				})
				return
			}
		}
	}

	if list.CheckExist {
		ty := list.CheckType
		ok, err := ExIfByField(db, ty, req)
		if !ok && err == nil {
			c.JSON(http.StatusNotFound, gin.H{
				"success": false,
				"message": "数据不存在",
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
	query := db.Model(req)
	if len(list.Preloads) > 0 {
		for _, v := range list.Preloads {
			query = query.Preload(v)
		}

	}

	// 1. 获取元素类型
	elemType := reflect.TypeOf(list.Model).Elem() // models.LuggageStorage 的 reflect.Type

	// 2. 创建切片类型
	sliceType := reflect.SliceOf(elemType) // []models.LuggageStorage 的 reflect.Type

	// 3. 创建切片值（指针）并解引用
	sliceValue := reflect.New(sliceType).Elem() // []models.LuggageStorage 的 reflect.Value

	// 4. 转换为 interface{}
	res := sliceValue.Interface() // []models.LuggageStorage
	va := reflect.ValueOf(req).Elem().FieldByName(ConvertSnakeToCamel(list.CheckType)).Interface()

	result := query.Model(req).Where(fmt.Sprintf("%s = ?", list.CheckType), va).Find(sliceValue.Addr().Interface())
	if result.Error != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"success": false,
			"message": "获取失败",
		})
		Logger.WithFields(logrus.Fields{
			"error": result.Error,
		}).Error("获取数据失败")
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"success": true,
		"message": "获取成功",
		"data":    res,
	})

}
