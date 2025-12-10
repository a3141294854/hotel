package main

import (
	"github.com/gin-gonic/gin"
	"hotel/internal/admin"
	"hotel/internal/util"
	"log"
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

	/*r.POST("/employee/register", func(c *gin.Context) {
		employee_check.EmployeeRegister(c, service)
	})*/

	r.POST("/employee/login", func(c *gin.Context) {
		employee_check.EmployeeLogin(c, service)
	})
	r.POST("/employee/refresh", func(c *gin.Context) {
		employee_check.RefreshToken(c, service)
	})

	e := r.Group("/employee")
	e.Use(middleware.JwtCheck(service))
	e.Use(middleware.AuthCheck(service))
	{
		e.POST("/logout", func(c *gin.Context) {
			employee_check.EmployeeLogout(c, service)
		})

		a := e.Group("/add")
		a.Use(middleware.CheckAction("创建行李"))
		{
			a.POST("/luggage", func(c *gin.Context) {
				employee_action.AddLuggage(c, service)
			})
		}

		d := e.Group("/delete")
		d.Use(middleware.CheckAction("删除行李"))
		{
			d.POST("/luggage", func(c *gin.Context) {
				employee_action.Delete(c, service)
			})
		}

		u := e.Group("/update")
		u.Use(middleware.CheckAction("更新行李"))
		{
			u.POST("/luggage", func(c *gin.Context) {
				employee_action.Update(c, service)
			})
		}

		g := e.Group("/get")
		g.Use(middleware.CheckAction("查看行李"))
		{
			g.GET("/name", func(c *gin.Context) {
				employee_action.GetName(c, service)
			})
			g.GET("/all", func(c *gin.Context) {
				employee_action.GetAll(c, service)
			})
			g.GET("/guest_id", func(c *gin.Context) {
				employee_action.GetGuestID(c, service)
			})
			g.GET("/location", func(c *gin.Context) {
				employee_action.GetLocation(c, service)
			})
			g.GET("/status", func(c *gin.Context) {
				employee_action.GetStatus(c, service)
			})
			g.POST("/guest_advance", func(c *gin.Context) {
				employee_action.GetAdvance(c, service)
			})
		}

		c := e.Group("/count")
		{
			c.GET("/sum", func(c *gin.Context) {
				employee_action.CountSum(c, service)
			})
			c.GET("/today", func(c *gin.Context) {
				employee_action.CountToday(c, service)
			})
		}

	}

	t := r.Group("/tool")
	t.Use(middleware.JwtCheck(service))
	t.Use(middleware.AuthCheck(service))
	t.Use(middleware.CheckAction("管理员"))
	{
		a := t.Group("/add")
		{
			a.POST("/permission", func(c *gin.Context) {
				admin.AddPermission(service, c)
			})
			a.POST("/role", func(c *gin.Context) {
				admin.AddRole(service, c)
			})
			a.POST("/employee", func(c *gin.Context) {
				employee_check.EmployeeRegister(c, service)
			})
			a.POST("/role_permission", func(c *gin.Context) {
				admin.AddRolePermission(service, c)
			})
		}

		g := t.Group("/get")
		{
			g.GET("/employee", func(c *gin.Context) {
				admin.GetAllEmployee(service, c)
			})
			g.GET("/permission", func(c *gin.Context) {
				admin.GetAllPermission(service, c)
			})
			g.GET("/role", func(c *gin.Context) {
				admin.GetAllRole(service, c)
			})
		}

		c := t.Group("/change")
		{
			c.POST("/employee_role", func(c *gin.Context) {
				admin.ChangeEmployeeRole(service, c)
			})
		}

		d := t.Group("/delete")
		{
			d.POST("/employee", func(c *gin.Context) {
				admin.DeleteEmployee(service, c)
			})
			d.POST("/role", func(c *gin.Context) {
				admin.DeleteRole(service, c)
			})
			d.POST("/permission", func(c *gin.Context) {
				admin.DeletePermission(service, c)
			})
		}

	}

	middleware.FindIp()
	err := r.Run("0.0.0.0:8080")
	if err != nil {
		log.Println("服务器启动失败:", err)
		return
	}

}
