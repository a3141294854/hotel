package employee_check

import (
	"errors"
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/sirupsen/logrus"
	"gorm.io/gorm"
	"hotel/internal/util"
	"hotel/models"
	"hotel/services"
	"math/rand"
	"net/http"
	"time"
)

// EmployeeRegister 员工注册
func EmployeeRegister(c *gin.Context, s *services.Services) {
	var e models.Employee
	//绑定
	err := c.ShouldBind(&e)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"success": false,
			"message": "请求数据格式错误",
		})
		util.Logger.WithFields(logrus.Fields{
			"error": err.Error(),
		}).Error("员工注册失败")
		return
	}
	//检查必要字段
	if e.User == "" || e.Password == "" || e.Name == "" || e.HotelID == 0 {
		c.JSON(http.StatusBadRequest, gin.H{
			"success": false,
			"message": "请求数据格式错误",
		})
		return
	}
	//检查用户名是否存在
	var existingEmployee models.Employee
	result2 := s.DB.Model(models.Employee{}).Where("user=?", e.User).First(&existingEmployee)
	if result2.Error == nil {
		c.JSON(http.StatusConflict, gin.H{
			"success": false,
			"message": "用户名已存在",
		})
		return
	}
	//设置默认值
	if e.RoleID == 0 {
		e.RoleID = 1
	}

	e.Password, err = util.HashPassword(e.Password)

	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"success": false,
			"message": "注册失败",
		})
		util.Logger.WithFields(logrus.Fields{
			"username": e.User,
			"error":    err,
		}).Error("员工密码加密错误")
		return
	}
	e.LastActiveTime = time.Now()
	result := s.DB.Create(&e)

	//插入
	employee := models.Employee{}
	employee.Password = ""
	s.DB.Model(models.Employee{}).Where("user=?", e.User).First(&employee)

	if result.Error != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"success": false,
			"message": "注册失败",
		})
		util.Logger.WithFields(logrus.Fields{
			"username": e.User,
			"error":    result.Error,
		}).Error("插入员工记录失败")
	} else {
		util.Logger.WithFields(logrus.Fields{
			"username": e.User,
		}).Info("员工注册成功")
		c.JSON(http.StatusCreated, gin.H{
			"success": true,
			"message": "注册成功",
			"data":    employee,
		})
	}
}

// EmployeeLogin 员工登录
func EmployeeLogin(c *gin.Context, s *services.Services) {
	var e struct {
		HotelID  int    `json:"hotel_id"`
		User     string `json:"user"`
		Password string `json:"password"`
	}
	//绑定
	if err := c.ShouldBindJSON(&e); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"success": false,
			"message": "请求数据格式错误",
		})
		util.Logger.WithFields(logrus.Fields{
			"error": err.Error(),
		}).Error("员工登录信息绑定错误")
		return
	}
	//检查必要字段
	if e.User == "" || e.Password == "" {
		c.JSON(http.StatusBadRequest, gin.H{
			"success": false,
			"message": "请求数据格式错误",
		})
		return
	}
	//设置默认值
	if e.HotelID == 0 {
		e.HotelID = 1
	}

	//检查用户名是否存在
	var user models.Employee
	result := s.DB.Model(models.Employee{}).
		Where("user=?", e.User).
		Where("hotel_id=?", e.HotelID).
		First(&user)
	if result.Error != nil {
		if errors.Is(result.Error, gorm.ErrRecordNotFound) {
			c.JSON(http.StatusUnauthorized, gin.H{
				"success": false,
				"message": "没有这名员工",
			})
			return
		} else {
			c.JSON(http.StatusUnauthorized, gin.H{
				"success": false,
				"message": "登录失败",
			})
			util.Logger.WithFields(logrus.Fields{
				"error": result.Error,
			}).Error("员工信息搜索错误")
			return
		}
	}

	if !util.CheckPassword(e.Password, user.Password) {
		c.JSON(http.StatusUnauthorized, gin.H{
			"success": false,
			"message": "密码错误",
		})
		return
	}

	//生成令牌
	accessToken, refreshToken, err := util.GenerateTokenPair(user.ID, user.Name, user.HotelID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"success": false,
			"message": "登录失败",
		})
		util.Logger.WithFields(logrus.Fields{
			"error": err,
		}).Error("JWT token生成错误")
		return
	}

	s.RdbAcc.Set(c, fmt.Sprintf("%d", user.ID), accessToken, util.AccessExpireTime+time.Duration(rand.Intn(100))*time.Second)
	s.RdbRef.Set(c, fmt.Sprintf("%d", user.ID), refreshToken, util.RefreshExpireTime+time.Duration(rand.Intn(100))*time.Second)

	util.Logger.WithFields(logrus.Fields{
		"username": user.User,
	}).Info("员工登录成功")

	c.JSON(http.StatusOK, gin.H{
		"success": true,
		"message": "登录成功",
		"data": gin.H{
			"access_token":  accessToken,
			"refresh_token": refreshToken,
			"expires_in":    900,
			"token_type":    "Bearer",
		},
	})
}

// EmployeeLogout 员工退出
func EmployeeLogout(c *gin.Context, s *services.Services) {

	claims, ok := c.Get("claims")
	if !ok {
		c.JSON(http.StatusUnauthorized, gin.H{
			"success": false,
			"message": "token无效",
		})
		return
	}
	e := claims.(*util.AccessClaims)

	s.RdbAcc.Del(c, fmt.Sprintf("%d", e.UserId))
	s.RdbRef.Del(c, fmt.Sprintf("%d", e.UserId))

	c.JSON(http.StatusOK, gin.H{
		"success": true,
		"message": "退出成功",
	})
}

// RefreshToken 员工刷新token
func RefreshToken(c *gin.Context, s *services.Services) {
	var e struct {
		RefreshTokens string `json:"refresh_token"`
	}
	//绑定
	if err := c.ShouldBindJSON(&e); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"success": false,
			"message": "请求数据格式错误",
		})
		util.Logger.WithFields(logrus.Fields{
			"error": err.Error(),
		}).Error("员工刷新token信息绑定错误")
		return
	}

	//解析令牌
	claims, err := util.ParseRefreshToken(e.RefreshTokens)
	if err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{
			"success": false,
			"message": "token无效",
		})
		util.Logger.WithFields(logrus.Fields{
			"error": err,
		}).Error("JWT token解析错误")
		return
	}

	//验证令牌
	refreshToken, err := s.RdbRef.Get(c, fmt.Sprintf("%d", claims.UserId)).Result()
	if err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{
			"success": false,
			"message": "token无效",
		})
		return
	}
	//验证令牌是否一致
	if refreshToken != e.RefreshTokens {
		c.JSON(http.StatusUnauthorized, gin.H{
			"success": false,
			"message": "token无效",
		})
		return
	}
	//生成新令牌
	accessToken, refreshToken, err := util.GenerateTokenPair(claims.UserId, claims.UserName, claims.HotelID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"success": false,
			"message": "内部错误",
		})
		util.Logger.WithFields(logrus.Fields{
			"error": err,
		}).Error("JWT token生成错误")
		return
	}
	s.RdbAcc.Set(c, fmt.Sprintf("%d", claims.UserId), accessToken, util.AccessExpireTime+time.Duration(rand.Intn(100))*time.Second)
	s.RdbRef.Set(c, fmt.Sprintf("%d", claims.UserId), refreshToken, util.RefreshExpireTime+time.Duration(rand.Intn(100))*time.Second)

	c.JSON(http.StatusOK, gin.H{
		"success": true,
		"message": "token刷新成功",
		"data": gin.H{
			"access_token":  accessToken,
			"refresh_token": refreshToken,
			"expires_in":    900,
			"token_type":    "Bearer",
		},
	})

}
