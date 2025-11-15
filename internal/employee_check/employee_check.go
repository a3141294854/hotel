package employee_check

import (
	"errors"
	"github.com/gin-contrib/sessions"
	"github.com/gin-gonic/gin"
	"github.com/jinzhu/gorm"
	"hotel/internal/models"
	"log"
	"net/http"
)

func EmployeeRegister(c *gin.Context, db *gorm.DB) {
	var e models.Employee
	if err := c.ShouldBind(&e); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"success": false,
		})
		log.Println("注册失败", err.Error())
		return
	}
	result := db.Create(&e)
	if result.Error != nil {
		log.Println("插入失败", e.User, result.Error)
	}
}

func EmployeeLogin(c *gin.Context, db *gorm.DB) {
	var e models.Employee
	if err := c.ShouldBind(&e); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"success": false,
		})
		log.Println("员工登录信息绑定错误", err.Error())
		return
	}

	var user []models.Employee
	result := db.Select("User", "password").Where("user_name=?", e.User).Where("password=?", e.Password).Find(&user)
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
			})
			log.Println("员工信息搜索错误", result.Error)
			return
		}
	}

	session := sessions.Default(c)
	session.Set(e.User, e.Password)
	session.Save()
}

func EmployeeLogout(c *gin.Context) {
	var e models.Employee
	if err := c.ShouldBind(&e); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"success": false,
		})
		log.Println("员工退出信息绑定错误", err.Error())
		return
	}

	session := sessions.Default(c)
	session.Delete(e.User)
	session.Save()
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
