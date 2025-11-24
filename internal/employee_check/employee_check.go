package employee_check

import (
	"errors"
	"fmt"
	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
	"hotel/internal/util"
	"hotel/models"
	"hotel/services"
	"log"
	"net/http"
)

func EmployeeRegister(c *gin.Context, s *services.Servers) {
	var e models.Employee
	if err := c.ShouldBind(&e); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"success": false,
			"message": "请求数据格式错误",
		})
		log.Println("注册失败", err.Error())
		return
	}
	if e.User == "" || e.Password == "" || e.Name == "" {
		c.JSON(http.StatusBadRequest, gin.H{
			"success": false,
			"message": "请求数据格式错误",
		})
		fmt.Println(e.User, e.Password, e.Name)
		return
	}
	var existingEmployee models.Employee
	result2 := s.DB.Model(models.Employee{}).Where("user=?", e.User).First(&existingEmployee)
	if result2.Error == nil { // 如果没有错误，说明找到了用户
		c.JSON(http.StatusBadRequest, gin.H{
			"success": false,
			"message": "用户名已存在",
		})
		return
	}

	result := s.DB.Create(&e)
	employee := models.Employee{}
	s.DB.Model(models.Employee{}).Where("user=?", e.User).First(&employee)
	if result.Error != nil {
		log.Println("插入失败", e.User, result.Error)
	} else {
		c.JSON(http.StatusOK, gin.H{
			"success": true,
			"message": "注册成功",
			"data":    employee,
		})
	}
}

func EmployeeLogin(c *gin.Context, s *services.Servers) {
	var e struct {
		User     string `json:"user"`
		Password string `json:"password"`
	}
	if err := c.ShouldBindJSON(&e); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"success": false,
			"message": "请求数据格式错误",
		})
		log.Println("员工登录信息绑定错误", err.Error())
		return
	}

	var user models.Employee
	result := s.DB.Select("user", "password").
		Where("user=?", e.User).
		Where("password=?", e.Password).
		First(&user)
	if result.Error != nil {
		if errors.Is(result.Error, gorm.ErrRecordNotFound) {
			c.JSON(http.StatusBadRequest, gin.H{
				"success": false,
				"message": "没有这名员工",
			})
			return
		} else {
			c.JSON(http.StatusBadRequest, gin.H{
				"success": false,
				"message": "登录失败",
			})
			log.Println("员工信息搜索错误", result.Error)
			return
		}
	}
	accessToken, refreshToken, err := util.GenerateTokenPair(user.ID, user.Name)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"success": false,
			"message": "登录失败",
		})
		log.Println("token生成错误", err)
		return
	}

	insert := models.RefreshToken{
		UserID:   user.ID,
		UserName: user.Name,
		Token:    refreshToken,
	}
	s.DB.Model(models.RefreshToken{}).Create(&insert)

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

func EmployeeLogout(c *gin.Context, s *services.Servers) {
	claims, ok := c.Get("claims")
	if !ok {
		c.JSON(http.StatusBadRequest, gin.H{
			"success": false,
			"message": "token无效",
		})
		return
	}
	e := claims.(*util.AccessClaims)
	result := s.DB.Model(models.RefreshToken{}).Where("user_id =?", e.UserId).Delete(models.RefreshToken{})
	if result.Error != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"success": false,
			"message": "退出失败",
		})
		log.Println("退出失败", result.Error)
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"success": true,
		"message": "退出成功",
	})
}

func RefreshToken(c *gin.Context, s *services.Servers) {
	var e struct {
		RefreshToken string `json:"refresh_token"`
	}
	if err := c.ShouldBindJSON(&e); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"success": false,
			"message": "请求数据格式错误",
		})
		log.Println("员工刷新token信息绑定错误", err.Error())
		return
	}

	claims, err := util.ParseRefreshToken(e.RefreshToken)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"success": false,
			"message": "token无效",
		})
		return
	}

	var employee models.RefreshToken
	result := s.DB.Model(models.RefreshToken{}).
		Where("user_id=?", claims.UserId).First(&employee)
	if result.Error != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"success": false,
			"message": "token无效",
		})
		return
	}
	result2 := s.DB.Model(models.RefreshToken{}).
		Where("user_id=?", claims.UserId).
		Where("token=?", e.RefreshToken).
		Delete(models.RefreshToken{})
	if result2.Error != nil {
		log.Println("token删除错误", result2.Error)
		return
	}

	accessToken, refreshToken, err := util.GenerateTokenPair(employee.UserID, employee.UserName)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"success": false,
			"message": "token无效",
		})
		return
	}

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
