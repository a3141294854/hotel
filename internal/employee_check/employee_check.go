package employee_check

import (
	"errors"
	"fmt"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/sirupsen/logrus"
	"gorm.io/gorm"
	"hotel/internal/util"
	"hotel/internal/util/logger"
	"hotel/models"
	"hotel/services"
)

// EmployeeRegister 员工注册
func EmployeeRegister(c *gin.Context, s *services.Services) {
	var e models.Employee
	if err := c.ShouldBind(&e); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"success": false,
			"message": "请求数据格式错误",
		})
		logger.Logger.WithFields(logrus.Fields{
			"error": err.Error(),
		}).Error("员工注册失败")
		return
	}
	if e.User == "" || e.Password == "" || e.Name == "" {
		c.JSON(http.StatusBadRequest, gin.H{
			"success": false,
			"message": "请求数据格式错误",
		})
		return
	}
	var existingEmployee models.Employee
	result2 := s.DB.Model(models.Employee{}).Where("user=?", e.User).First(&existingEmployee)
	if result2.Error == nil { // 如果没有错误，说明找到了用户
		c.JSON(http.StatusConflict, gin.H{
			"success": false,
			"message": "用户名已存在",
		})
		return
	}
	//员工默认角色为1，即普通员工
	e.RoleID = 1
	result := s.DB.Create(&e)
	employee := models.Employee{}
	s.DB.Model(models.Employee{}).Where("user=?", e.User).First(&employee)
	if result.Error != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"success": false,
			"message": "注册失败",
		})
		logger.Logger.WithFields(logrus.Fields{
			"username": e.User,
			"error":    result.Error,
		}).Error("插入员工记录失败")
	} else {
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
		User     string `json:"user"`
		Password string `json:"password"`
	}
	if err := c.ShouldBindJSON(&e); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"success": false,
			"message": "请求数据格式错误",
		})
		logger.Logger.WithFields(logrus.Fields{
			"error": err.Error(),
		}).Error("员工登录信息绑定错误")
		return
	}

	var user models.Employee
	result := s.DB.Model(models.Employee{}).
		Where("user=?", e.User).
		Where("password=?", e.Password).
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
			logger.Logger.WithFields(logrus.Fields{
				"error": result.Error,
			}).Error("员工信息搜索错误")
			return
		}
	}

	accessToken, refreshToken, err := util.GenerateTokenPair(user.ID, user.Name)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"success": false,
			"message": "登录失败",
		})
		logger.Logger.WithFields(logrus.Fields{
			"error": err,
		}).Error("JWT token生成错误")
		return
	}

	s.RdbAcc.Set(c, fmt.Sprintf("%d", user.ID), accessToken, 5*time.Minute)
	s.RdbRef.Set(c, fmt.Sprintf("%d", user.ID), refreshToken, 24*time.Hour)

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
	if err := c.ShouldBindJSON(&e); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"success": false,
			"message": "请求数据格式错误",
		})
		logger.Logger.WithFields(logrus.Fields{
			"error": err.Error(),
		}).Error("员工刷新token信息绑定错误")
		return
	}

	claims, err := util.ParseRefreshToken(e.RefreshTokens)
	if err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{
			"success": false,
			"message": "token无效",
		})
		logger.Logger.WithFields(logrus.Fields{
			"error": err,
		}).Error("JWT token解析错误")
		return
	}

	refreshToken, err := s.RdbRef.Get(c, fmt.Sprintf("%d", claims.UserId)).Result()
	if err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{
			"success": false,
			"message": "token无效",
		})
		return
	}

	if refreshToken != e.RefreshTokens {
		c.JSON(http.StatusUnauthorized, gin.H{
			"success": false,
			"message": "token无效",
		})
		return
	}

	accessToken, refreshToken, err := util.GenerateTokenPair(claims.UserId, claims.UserName)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"success": false,
			"message": "内部错误",
		})
		logger.Logger.WithFields(logrus.Fields{
			"error": err,
		}).Error("JWT token生成错误")
		return
	}
	s.RdbAcc.Set(c, fmt.Sprintf("%d", claims.UserId), accessToken, 5*time.Minute)
	s.RdbRef.Set(c, fmt.Sprintf("%d", claims.UserId), refreshToken, 24*time.Hour)

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
