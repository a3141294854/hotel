package main

import (
	"github.com/gin-contrib/sessions"
	"github.com/gin-contrib/sessions/cookie"
	"gorm.io/driver/mysql"
	"hotel/internal/table"
	"log"
	"time"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"

	"hotel/internal/employee_action"
	"hotel/internal/employee_check"
	"hotel/internal/function"
)

func main() {
	dsn := "root:@furenjie321@tcp(127.0.0.1:3306)/study?charset=utf8mb4&parseTime=True&loc=Local"
	db, err := gorm.Open(mysql.Open(dsn), &gorm.Config{})
	if err != nil {
		log.Fatal(err)
	}
	table.Table(db)

	store := cookie.NewStore([]byte("secret"))
	limiter := function.NewTokenBucketLimiter(100, time.Second)

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
		g := e.Group("/get")
		{
			g.GET("/name", employee_check.Check(), func(c *gin.Context) {
				employee_action.GetName(c, db)
			})
			g.GET("/all", employee_check.Check(), func(c *gin.Context) {
				employee_action.GetAll(c, db)
			})
			g.GET("/guest_id", employee_check.Check(), func(c *gin.Context) {
				employee_action.GetGuestID(c, db)
			})
			g.GET("/location", employee_check.Check(), func(c *gin.Context) {
				employee_action.GetLocation(c, db)
			})
			g.GET("/status", employee_check.Check(), func(c *gin.Context) {
				employee_action.GetStatus(c, db)
			})
			g.POST("/guest_advance", employee_check.Check(), func(c *gin.Context) {
				employee_action.GetAdvance(c, db)
			})
		}

		c := e.Group("/count")
		{
			c.GET("/sum", employee_check.Check(), func(c *gin.Context) {
				employee_action.CountSum(c, db)
			})
			c.GET("/today", employee_check.Check(), func(c *gin.Context) {
				employee_action.CountToday(c, db)
			})
		}

	}
	r.Run("0.0.0.0:8080")
}
