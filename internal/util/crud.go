package util

import (
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/sirupsen/logrus"
	"gorm.io/gorm"
	"net/http"
	"reflect"
	"strings"
)

type RequestList struct {
	Model      interface{} // 模型
	CheckExist bool        // 是否检查存在
	CheckType  string      //判断是否存在的数据的类型，要小写，查结构体会有函数转换成大写,也是主要的操作对象
	GetType    string      //如果是跨表查询就添加，要大写,包含表名和字段名，如："Guests.Name"
	CheckField []string    // 检查必要的字段名，要大写，只用来查结构体
	Preloads   []string    // 预加载,要大写
}

// Create 通用创建数据
func Create(c *gin.Context, db *gorm.DB, list RequestList) {
	//确定模型
	req := list.Model
	//绑定数据
	err := c.ShouldBind(req)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"success": false,
			"message": "请求数据格式错误",
		})
		return
	}
	//需要检查的字段
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
	//检查是否存在数据
	if list.CheckExist {
		ty := list.CheckType
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
	query := db.Model(req)
	//检查有无要预加载的表
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

// Delete 通用删除数据
func Delete(c *gin.Context, db *gorm.DB, list RequestList) {
	//确定模型
	req := list.Model
	//绑定数据
	err := c.ShouldBind(req)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"success": false,
			"message": "请求数据格式错误",
		})
		return
	}
	//需要检查的字段
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
	//检查是否存在数据
	if list.CheckExist {
		ty := list.CheckType
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
	query := db.Model(req)
	//检查有无要预加载的表
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

// Update 通用更新数据
func Update(c *gin.Context, db *gorm.DB, list RequestList) {
	//确定模型
	req := list.Model
	//绑定数据
	err := c.ShouldBind(req)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"success": false,
			"message": "请求数据格式错误",
		})
		return
	}
	//需要检查的字段
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

	//检查是否存在数据
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
	//检查有无要预加载的表
	if len(list.Preloads) > 0 {
		for _, v := range list.Preloads {
			query = query.Preload(v)
		}
	}

	va := reflect.ValueOf(req).Elem().FieldByName(ConvertSnakeToCamel(list.CheckType)).Interface()
	result := query.Where(fmt.Sprintf("%s = ?", list.CheckType), va).Updates(req)
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
		//"data":    req,
	})
}

// Get 通用get函数
func Get(c *gin.Context, db *gorm.DB, list RequestList) {
	req := list.Model
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

	va := c.Query(list.CheckType)
	if va == "" {
		c.JSON(http.StatusBadRequest, gin.H{
			"success": false,
			"message": "请提供" + list.CheckType + "参数",
		})
		return
	}
	var get string
	//检查有无需要跨表查询
	if list.GetType == "" {
		get = list.CheckType
	} else {
		get = list.GetType
	}
	// 检查是否是关联查询（包含点号）
	if strings.Contains(get, ".") {
		// 关联查询：需要使用 Joins
		parts := strings.Split(get, ".")
		if len(parts) == 2 {
			// parts[0] = "Guest", parts[1] = "Name"
			query = query.Joins(parts[0]) //可以自动把单数结构名转换为复数
		}
	}

	result := query.Where(fmt.Sprintf("%s = ?", get), va).Find(sliceValue.Addr().Interface())

	// 4. 转换为 interface{}
	res := sliceValue.Interface() // []models.LuggageStorage
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
		"count":   result.RowsAffected, //表示受影响的行数，这里就是数量
	})

}
