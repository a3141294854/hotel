package main

import (
	"fmt"
	"hotel/internal/employee/employee_action"
	"hotel/internal/employee/employee_check"
	"hotel/internal/util"

	"github.com/gin-gonic/gin"
	"github.com/sirupsen/logrus"
	"hotel/internal/admin"
	"hotel/internal/config"
	"hotel/internal/middleware"
	"hotel/services"
)

// open 启动路由配置
func open(r *gin.Engine, service *services.Services, cfg *config.Config) {
	//中间件配置
	r.Use(middleware.RequestIDMiddleware())
	r.Use(middleware.LogRequest())
	r.Use(middleware.RateLimit("local", service))

	//公共接口
	r.POST("/employee/login", func(c *gin.Context) {
		employee_check.EmployeeLogin(c, service)
	})
	r.POST("/employee/refresh", func(c *gin.Context) {
		employee_check.RefreshToken(c, service)
	})

	//员工组
	e := r.Group("/employee")
	e.Use(middleware.JwtCheck(service))
	e.Use(middleware.AuthCheck(service))
	{
		//退出登录
		e.POST("/logout", func(c *gin.Context) {
			employee_check.EmployeeLogout(c, service)
		})

		//添加操作
		a := e.Group("/add")
		a.Use(middleware.CheckAction("创建行李"))
		{
			a.POST("/luggageStorage", func(c *gin.Context) {
				employee_action.AddLuggage(c, service)
			})
			a.POST("/mac", func(c *gin.Context) {
				employee_action.AddMac(c, service)
			})
		}

		//删除操作
		d := e.Group("/delete")
		d.Use(middleware.CheckAction("删除行李"))
		{
			d.POST("/luggageStorage", func(c *gin.Context) {
				employee_action.DeleteStorage(c, service)
			})
			d.POST("/luggage", func(c *gin.Context) {
				employee_action.DeleteLuggage(c, service)
			})
		}

		//更新操作
		u := e.Group("/update")
		u.Use(middleware.CheckAction("更新行李"))
		{
			u.PUT("/luggageStorage", func(c *gin.Context) {
				employee_action.UpdateLuggageStorage(c, service)
			})
			u.PUT("/luggage", func(c *gin.Context) {
				employee_action.UpdateLuggage(c, service)
			})
		}

		//查询操作
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
			g.POST("/guest_advance", func(c *gin.Context) {
				employee_action.GetAdvance(c, service)
			})
			g.GET("/pick_up_code", func(c *gin.Context) {
				employee_action.GetPickUpCode(c, service)
			})
		}

		//统计操作
		c := e.Group("/count")
		{
			c.GET("/sum", func(c *gin.Context) {
				employee_action.CountSum(c, service)
			})
			c.GET("/today", func(c *gin.Context) {
				employee_action.CountToday(c, service)
			})
		}

		//图片的操作
		p := e.Group("/photo")
		{
			p.POST("/upload", func(c *gin.Context) {
				employee_action.UploadPhoto(c, service)
			})

			p.GET("/download/:filename", func(c *gin.Context) {
				employee_action.DownloadPhoto(c)
			})
			p.POST("/touch", func(c *gin.Context) {
				employee_action.PhotoTouchLuggageStorage(c, service)
			})
			p.GET("/get_all", func(c *gin.Context) {
				employee_action.GetAllPhoto(c, service)
			})
		}

	}

	//管理员组
	t := r.Group("/tool")
	t.Use(middleware.JwtCheck(service))
	t.Use(middleware.AuthCheck(service))
	t.Use(middleware.CheckAction("管理员"))
	{
		//添加操作
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
			a.POST("/hotel", func(c *gin.Context) {
				admin.AddHotel(service, c)
			})
		}

		//查询操作
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
			g.GET("/location", func(c *gin.Context) {
				admin.GetAllLocation(service, c)
			})
		}

		//修改操作
		c := t.Group("/change")
		{
			c.POST("/employee_role", func(c *gin.Context) {
				admin.ChangeEmployeeRole(service, c)
			})
		}

		//删除操作
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

	//显示IP地址
	//middleware.FindIp()
	util.Logger.WithFields(logrus.Fields{
		"mode": cfg.Server.Mode,
	}).Info("服务器启动")
	err := r.Run(cfg.Server.Host + fmt.Sprintf(":%d", cfg.Server.Port))

	if err != nil {
		util.Logger.WithFields(logrus.Fields{
			"error": err,
			"mode":  cfg.Server.Mode,
		}).Error("服务器启动失败")
		return
	}

}
