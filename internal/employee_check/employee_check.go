package employee_check

import (
	"errors"
	"github.com/gin-contrib/sessions"
	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
	"hotel/internal/models"
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
	var e models.Employee
	if err := c.ShouldBind(&e); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"success": false,
			"message": "请求数据格式错误",
		})
		log.Println("员工登录信息绑定错误", err.Error())
		return
	}

	var user []models.Employee
	result := db.Select("user_name", "password").Where("user_name=?", e.User).Where("password=?", e.Password).Find(&user)
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

	session := sessions.Default(c)
	session.Set("user", e.User)
	session.Save()
	c.JSON(http.StatusOK, gin.H{
		"success": true,
		"message": "登录成功",
	})
}

func EmployeeLogout(c *gin.Context) {
	var e models.Employee
	if err := c.ShouldBind(&e); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"success": false,
			"message": "请求数据格式错误",
		})
		log.Println("员工退出信息绑定错误", err.Error())
		return
	}

	session := sessions.Default(c)
	session.Delete("user")
	session.Save()
	c.JSON(http.StatusOK, gin.H{
		"success": true,
		"message": "退出成功",
	})
}

func Check() gin.HandlerFunc {
	return func(c *gin.Context) {
		session := sessions.Default(c)
		user := session.Get("user")
		if user == nil {
			c.JSON(http.StatusUnauthorized, gin.H{
				"success": false,
				"message": "请先登录",
			})
			c.Abort()
			return
		}
		c.Next()
	}
}
