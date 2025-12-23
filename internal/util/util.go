package util

import (
	"context"
	"errors"
	"fmt"
	"github.com/sirupsen/logrus"
	"gorm.io/gorm"
	"hotel/internal/util/logger"
	"strconv"

	"hotel/services"
	"math/big"

	"crypto/rand"
)

// GeneratePickUpCode 生成取件码
func GeneratePickUpCode(s *services.Services, HotelID uint) (string, error) {
	//加载最大数
	ma := big.NewInt(1000000)
	//随机生成
	now, _ := rand.Int(rand.Reader, ma)
	code := fmt.Sprintf("%06d", now)
	//判断是否重复
	ex, err := s.RdbRand.Exists(context.Background(), code).Result()
	if err != nil {
		return "", err
	}
	if ex == 1 {
		return GeneratePickUpCode(s, HotelID)

	} else {
		s.RdbRand.Set(context.Background(), fmt.Sprintf("%d:%s", HotelID, code), "1", 0)
		return code, nil
	}
}

// ExIf 判断是否存在 true 存在 false 不存在
func ExIf(db *gorm.DB, ty string, model interface{}, value string) (bool, error) {

	if ty == "id" {
		num, err := strconv.ParseUint(value, 10, 64)
		if err != nil {
			logger.Logger.WithFields(logrus.Fields{
				"error": err,
			}).Error("转换id失败")
			return false, err
		}
		uintnum := uint(num)
		result := db.Model(model).Where("id = ?", uintnum).First(model)
		if result.Error != nil {
			if errors.Is(result.Error, gorm.ErrRecordNotFound) {
				return false, nil
			}
			return false, result.Error
		} else {
			return true, nil
		}
	}

	if ty == "name" {
		result := db.Model(model).Where("name = ?", value).First(model)
		if result.Error != nil {
			if errors.Is(result.Error, gorm.ErrRecordNotFound) {
				return false, nil
			}
			return false, result.Error
		} else {
			return true, nil
		}
	}

	if ty == "user" {
		result := db.Model(model).Where("user = ?", value).First(model)
		if result.Error != nil {
			if errors.Is(result.Error, gorm.ErrRecordNotFound) {
				return false, nil
			}
			return false, result.Error
		} else {
			return true, nil
		}
	}

	return false, errors.New("不支持的类型")
}
