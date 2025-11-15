package main

import (
	"github.com/gin-contrib/sessions"
	"github.com/gin-contrib/sessions/cookie"
	_ "github.com/go-sql-driver/mysql"
	"hotel/internal/table"
	"log"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/jinzhu/gorm"

	"hotel/internal/employee_action"
	"hotel/internal/employee_check"
	"hotel/internal/function"
)

func main() {
	db, err := gorm.Open("mysql", "root:@furenjie321@(127.0.0.1:3306)/study?charset=utf8mb4&parseTime=True&loc=Local")
	if err != nil {
		log.Fatal(err)
	}
	defer db.Close()
	table.Table(db)

	store := cookie.NewStore([]byte("secret"))
	limiter := function.NewTokenBucketLimiter(10, time.Second)

	r := gin.Default()
	r.Use(sessions.Sessions("session", store))

	// 限流中间件
	r.Use(func(c *gin.Context) {
		if !limiter.Allow() {
			c.JSON(429, gin.H{
				"success": false,
				"message": "请求过于频繁，请稍后再试",
			})
			c.Abort()
			return
		}
		c.Next()
	})

	e := r.Group("/employee")
	{
		e.POST("/register", func(c *gin.Context) {
			employee_check.EmployeeRegister(c, db)
		})
		e.POST("/login", func(c *gin.Context) {
			employee_check.EmployeeLogin(c, db)
		})
		e.POST("/logout", func(c *gin.Context) {
			employee_check.EmployeeLogout(c)
		})
		e.POST("/add", employee_check.Check(), func(c *gin.Context) {
			employee_action.Add(c, db)
		})
		e.DELETE("/delete", employee_check.Check(), func(c *gin.Context) {
			employee_action.Delete(c, db)
		})
		e.PUT("/update", employee_check.Check(), func(c *gin.Context) {
			employee_action.Update(c, db)
		})
		e.GET("/get", employee_check.Check(), func(c *gin.Context) {
			employee_action.Get(c, db)
		})
	}
	r.Run("0.0.0.0:8080")

}
