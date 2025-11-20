package employee_check

import (
	"errors"
	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
	"hotel/internal/models"
	"hotel/internal/util"
	"log"
	"net/http"
)

func EmployeeRegister(c *gin.Context, db *gorm.DB) {
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
		return
	}

	result := db.Create(&e)
	if result.Error != nil {
		log.Println("插入失败", e.User, result.Error)
	} else {
		c.JSON(http.StatusOK, gin.H{
			"success": true,
			"message": "注册成功",
		})
	}
}

func EmployeeLogin(c *gin.Context, db *gorm.DB) {
	var e struct {
		UserName string `json:"user_name"`
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
	result := db.Select("user_name", "password").
		Where("user_name=?", e.UserName).
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
	accessToken, refreshToken, err := util.GenerateTokenPair(user.ID, user.User)
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
	db.Model(models.RefreshToken{}).Create(&insert)

	c.JSON(http.StatusOK, gin.H{
		"success": true,
		"message": "登录成功",
		"data": gin.H{
			"access_token":  accessToken,
			"refresh_token": refreshToken,
			"expires_in": 900,
			"token_type": "Bearer",
		},
	})
}

func EmployeeLogout(c *gin.Context, db *gorm.DB) {
	claims, ok := c.Get("claims")
	if !ok {
		c.JSON(http.StatusBadRequest, gin.H{
			"success": false,
			"message": "token无效",
		})
		return
	}
	e := claims.(*util.AccessClaims)
	result := db.Model(models.RefreshToken{}).Where("user_id =?", e.UserId).Delete(models.RefreshToken{})
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

func RefreshToken(c *gin.Context, db *gorm.DB) {
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
	result := db.Model(models.RefreshToken{}).
		Where("user_id=?", claims.UserId).First(&employee)
	if result.Error != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"success": false,
			"message": "token无效",
		})
		return
	}
	result2 := db.Model(models.RefreshToken{}).
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
			"expires_in": 900,
			"token_type": "Bearer",
		},
	})

}
