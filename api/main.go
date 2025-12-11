package main

import (
	"fmt"

	"github.com/gin-gonic/gin"
	"github.com/sirupsen/logrus"
	"hotel/internal/admin"
	"hotel/internal/config"
	"hotel/internal/employee_action"
	"hotel/internal/employee_check"
	"hotel/internal/middleware"
	"hotel/internal/table"
	"hotel/internal/util"
	"hotel/internal/util/logger"
	"hotel/services"
)

func main() {
	cfg, err := config.LoadConfig("")
	if err != nil {
		logger.Logger.WithFields(logrus.Fields{
			"error": err,
		}).Error("加载配置文件失败")
		return
	}
	logger.InitLogger(cfg.Log.Level, cfg.Log.Output, cfg.Log.FilePath, cfg.Log.MaxSize, cfg.Log.MaxBackups, cfg.Log.MaxAge)

	service := services.NewDatabase(cfg)
	table.Table(service.DB)

	//message_queue.StartTaskProcessor(context.Background(), service)
	util.NewTokenBucketLimiter(cfg.RateLimiting.Default.Name, cfg.RateLimiting.Default.Capacity, cfg.RateLimiting.Default.FillRate, service)
	util.ConfigJwt(cfg)

	r := gin.Default()

	r.Use(middleware.RequestIDMiddleware())
	r.Use(middleware.LogRequest())
	r.Use(middleware.RateLimit("local", service))

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
	logger.Logger.WithFields(logrus.Fields{
		"mode": cfg.Server.Mode,
	}).Info("服务器启动")
	err = r.Run(cfg.Server.Host + fmt.Sprintf(":%d", cfg.Server.Port))

	if err != nil {
		logger.Logger.WithFields(logrus.Fields{
			"error": err,
			"mode":  cfg.Server.Mode,
		}).Error("服务器启动失败")
		return
	}

}
