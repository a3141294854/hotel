package main

import (
	"github.com/gin-gonic/gin"
	"hotel/internal/util"
	"time"

	"hotel/internal/employee_action"
	"hotel/internal/employee_check"
	"hotel/internal/middleware"
	"hotel/internal/table"
	"hotel/services"
)

func main() {
	service := services.NewDatabase()
	table.Table(service.DB)

	//message_queue.StartTaskProcessor(context.Background(), service)

	util.NewTokenBucketLimiter("local", 10, time.Second, service)

	r := gin.Default()

	r.Use(middleware.RateLimit("local", service))

	r.POST("/employee/register", func(c *gin.Context) {
		employee_check.EmployeeRegister(c, service)
	})
	r.POST("/employee/login", func(c *gin.Context) {
		employee_check.EmployeeLogin(c, service)
	})
	r.POST("/employee/refresh", func(c *gin.Context) {
		employee_check.RefreshToken(c, service)
	})

	e := r.Group("/employee")
	e.Use(middleware.JwtCheck(service))
	{
		e.POST("/logout", middleware.AuthCheck(), func(c *gin.Context) {
			employee_check.EmployeeLogout(c, service)
		})
		e.POST("/add", middleware.AuthCheck(), func(c *gin.Context) {
			employee_action.Add(c, service)
		})
		e.DELETE("/delete", middleware.AuthCheck(), func(c *gin.Context) {
			employee_action.Delete(c, service)
		})
		e.PUT("/update", middleware.AuthCheck(), func(c *gin.Context) {
			employee_action.Update(c, service)
		})

		g := e.Group("/get")
		{
			g.GET("/name", middleware.AuthCheck(), func(c *gin.Context) {
				employee_action.GetName(c, service)
			})
			g.GET("/all", middleware.AuthCheck(), func(c *gin.Context) {
				employee_action.GetAll(c, service)
			})
			g.GET("/guest_id", middleware.AuthCheck(), func(c *gin.Context) {
				employee_action.GetGuestID(c, service)
			})
			g.GET("/location", middleware.AuthCheck(), func(c *gin.Context) {
				employee_action.GetLocation(c, service)
			})
			g.GET("/status", middleware.AuthCheck(), func(c *gin.Context) {
				employee_action.GetStatus(c, service)
			})
			g.POST("/guest_advance", middleware.AuthCheck(), func(c *gin.Context) {
				employee_action.GetAdvance(c, service)
			})
		}

		c := e.Group("/count")
		{
			c.GET("/sum", middleware.AuthCheck(), func(c *gin.Context) {
				employee_action.CountSum(c, service)
			})
			c.GET("/today", middleware.AuthCheck(), func(c *gin.Context) {
				employee_action.CountToday(c, service)
			})
		}

	}
	middleware.FindIp()
	r.Run("0.0.0.0:8080")

}
