package util

import (
	"context"
	"errors"
	"fmt"
	"github.com/redis/go-redis/v9"
	"github.com/sirupsen/logrus"
	"gorm.io/gorm"
	"reflect"
	"strconv"
	"strings"

	"math/big"

	"crypto/rand"
)

// GeneratePickUpCode 生成取件码
func GeneratePickUpCode(RdbRand *redis.Client, HotelID uint) (string, error) {
	//加载最大数
	ma := big.NewInt(1000000)
	//随机生成
	now, _ := rand.Int(rand.Reader, ma)
	code := fmt.Sprintf("%06d", now)
	//判断是否重复
	ex, err := RdbRand.Exists(context.Background(), code).Result()
	if err != nil {
		return "", err
	}
	if ex == 1 {
		return GeneratePickUpCode(RdbRand, HotelID)

	} else {
		RdbRand.Set(context.Background(), fmt.Sprintf("%d:%s", HotelID, code), "1", 0)
		return code, nil
	}
}

// ExIf 判断是否存在这个数据 true 存在 false 不存在
// 传入数据库连接，查询的类型，查询的值
func ExIf(db *gorm.DB, ty string, model interface{}, value string) (bool, error) {

	if ty == "id" {
		num, err := strconv.ParseUint(value, 10, 64)
		if err != nil {
			Logger.WithFields(logrus.Fields{
				"error": err,
			}).Error("转换id失败")
			return false, err
		}
		uintnum := uint(num)
		result := db.Model(model).Where("id = ?", uintnum).First(model)
		if result.Error != nil {
			if errors.Is(result.Error, gorm.ErrRecordNotFound) {
				return false, nil
			} else {
				return false, result.Error
			}
		} else {
			return true, nil
		}
	}

	if ty == "name" {
		result := db.Model(model).Where("name = ?", value).First(model)
		if result.Error != nil {
			if errors.Is(result.Error, gorm.ErrRecordNotFound) {
				return false, nil
			} else {
				return false, result.Error
			}
		} else {
			return true, nil
		}
	}

	if ty == "user" {
		result := db.Model(model).Where("user = ?", value).First(model)
		if result.Error != nil {
			if errors.Is(result.Error, gorm.ErrRecordNotFound) {
				return false, nil
			} else {
				return false, result.Error
			}
		} else {
			return true, nil
		}
	}

	return false, errors.New("不支持的类型")
}

// ExIfByField 根据字段名从模型中获取值并检查是否存在
func ExIfByField(db *gorm.DB, ty string, model interface{}) (bool, error) {
	reqValue := reflect.ValueOf(model).Elem()                   //获取指针指向的值
	fieldValue := reqValue.FieldByName(ConvertSnakeToCamel(ty)) //根据字段名获取对应的值

	if !fieldValue.IsValid() {
		return false, errors.New("字段 " + ty + " 不存在")
	}

	var value string
	//返回类型种类
	switch fieldValue.Kind() {
	case reflect.String:
		value = fieldValue.String()
	case reflect.Uint, reflect.Uint8, reflect.Uint16, reflect.Uint32, reflect.Uint64:
		value = fmt.Sprintf("%d", fieldValue.Uint())
	default:
		value = fmt.Sprintf("%v", fieldValue.Interface())
	}

	return ExIf(db, ty, model, value)
}

func ConvertSnakeToCamel(s string) string {
	mappings := map[string]string{
		"id":   "ID",
		"user": "User",
		"name": "Name",
	}

	FieldName, ok := mappings[s]
	if ok {
		return FieldName
	}

	if len(s) == 0 {
		return s
	}
	//如果没有记录，默认开头大写
	return strings.ToUpper(s[:1]) + s[1:]
}
