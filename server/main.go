package main

import (
	"log"
	"time"

	"github.com/gin-gonic/gin"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"

	"hotel/internal/employee_action"
	"hotel/internal/employee_check"
	"hotel/internal/middleware"
	"hotel/internal/table"
)

func main() {
	dsn := "root:@furenjie321@tcp(127.0.0.1:3306)/study?charset=utf8mb4&parseTime=True&loc=Local"
	db, err := gorm.Open(mysql.Open(dsn), &gorm.Config{})
	if err != nil {
		log.Fatal(err)
	}
	table.Table(db)

	limiter := middleware.NewTokenBucketLimiter(100, time.Second)

	r := gin.Default()

	r.Use(middleware.RateLimit(limiter))

	e := r.Group("/employee")
	{
		e.POST("/register", func(c *gin.Context) {
			employee_check.EmployeeRegister(c, db)
		})
		e.POST("/login", func(c *gin.Context) {
			employee_check.EmployeeLogin(c, db)
		})
		
		// 需要JWT认证的路由组
		auth := e.Group("/")
		auth.Use(middleware.JwtCheck())
		{
			auth.POST("/refresh", func(c *gin.Context) {
				employee_check.RefreshToken(c, db)
			})
			auth.POST("/logout", middleware.AuthCheck(), func(c *gin.Context) {
				employee_check.EmployeeLogout(c, db)
			})
			auth.POST("/add", middleware.AuthCheck(), func(c *gin.Context) {
				employee_action.Add(c, db)
			})
			auth.DELETE("/delete", middleware.AuthCheck(), func(c *gin.Context) {
				employee_action.Delete(c, db)
			})
			auth.PUT("/update", middleware.AuthCheck(), func(c *gin.Context) {
				employee_action.Update(c, db)
			})
			
			g := auth.Group("/get")
			{
				g.GET("/name", middleware.AuthCheck(), func(c *gin.Context) {
					employee_action.GetName(c, db)
				})
				g.GET("/all", middleware.AuthCheck(), func(c *gin.Context) {
					employee_action.GetAll(c, db)
				})
				g.GET("/guest_id", middleware.AuthCheck(), func(c *gin.Context) {
					employee_action.GetGuestID(c, db)
				})
				g.GET("/location", middleware.AuthCheck(), func(c *gin.Context) {
					employee_action.GetLocation(c, db)
				})
				g.GET("/status", middleware.AuthCheck(), func(c *gin.Context) {
					employee_action.GetStatus(c, db)
				})
				g.POST("/guest_advance", middleware.AuthCheck(), func(c *gin.Context) {
					employee_action.GetAdvance(c, db)
				})
			}

			c := auth.Group("/count")
			{
				c.GET("/sum", middleware.AuthCheck(), func(c *gin.Context) {
					employee_action.CountSum(c, db)
				})
				c.GET("/today", middleware.AuthCheck(), func(c *gin.Context) {
					employee_action.CountToday(c, db)
				})
			}
		}
	}
	middleware.FindIp()
	r.Run("0.0.0.0:8080")

}
