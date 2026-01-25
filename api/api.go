package api

import (
	"fmt"
	"hotel/internal/employee/employee_action"
	"hotel/internal/employee/employee_check"
	"hotel/internal/util"

	"github.com/gin-gonic/gin"
	"github.com/sirupsen/logrus"
	"hotel/internal/admin"
	"hotel/internal/middleware"
	"hotel/services"
)

// Init 启动路由配置
func Init(r *gin.Engine, service *services.Services, cfg *util.Config) {
	//中间件配置
	r.Use(middleware.Recovery())
	r.Use(middleware.RequestIDMiddleware())
	r.Use(middleware.LogRequest())
	r.Use(middleware.RateLimit("local", service))

	//公共接口
	r.POST("/login", func(c *gin.Context) {
		employee_check.EmployeeLogin(c, service)
	})
	r.POST("/logout", func(c *gin.Context) {
		employee_check.EmployeeLogout(c, service)
	})
	r.POST("/refresh", func(c *gin.Context) {
		employee_check.RefreshToken(c, service)
	})
	r.POST("/register", middleware.JwtCheck(service), middleware.AuthCheck(service), middleware.CheckAction("管理员"), func(c *gin.Context) {
		employee_check.EmployeeRegister(c, service)
	})

	r.GET("/photos/:filename", middleware.JwtCheck(service), func(c *gin.Context) {
		employee_action.StaticDownloadPhoto(c)
	})

	internal := r.Group("/internal")
	internal.Use(middleware.JwtCheck(service))
	internal.Use(middleware.AuthCheck(service))
	internal.Use(middleware.CheckAction("内部接口"))

	//行李寄存表
	luggageStorage := internal.Group("/luggageStorage")
	{
		luggageStorage.POST("", func(c *gin.Context) {
			employee_action.AddLuggageStorage(c, service)
		})
		luggageStorage.GET("", func(c *gin.Context) {
			employee_action.GetLuggageStorage(c, service)
		})
		luggageStorage.PUT("/:id", func(c *gin.Context) {
			employee_action.UpdateLuggageStorage(c, service)
		})
		luggageStorage.POST("/code/:code", func(c *gin.Context) {
			employee_action.CheckOutLuggageStorage(c, service)
		})
		luggageStorage.DELETE("/:id", func(c *gin.Context) {
			employee_action.DeleteLuggageStorage(c, service)
		})
		luggageStorage.DELETE("/code/:pick_up_code", func(c *gin.Context) {
			employee_action.DeleteLuggageStorageByCode(c, service)
		})
	}
	//行李表
	luggage := internal.Group("/luggage")
	{
		luggage.POST("/:id", func(c *gin.Context) {
			employee_action.CheckOutLuggage(c, service)
		})
		luggage.DELETE("/:id", func(c *gin.Context) {
			employee_action.DeleteLuggage(c, service)
		})
		luggage.PUT("/:id", func(c *gin.Context) {
			employee_action.UpdateLuggage(c, service)
		})
		luggage.GET("", func(c *gin.Context) {
			employee_action.GetLuggage(c, service)
		})
	}
	//标签表
	tag := internal.Group("/tag")
	{
		tag.POST("", func(c *gin.Context) {
			employee_action.AddMac(c, service)
		})
		tag.PUT("/:id", func(c *gin.Context) {
			employee_action.UpdateTag(c, service)
		})
		tag.DELETE("/:id", func(c *gin.Context) {
			employee_action.DeleteTag(c, service)
		})
		tag.GET("", func(c *gin.Context) {
			employee_action.GetTag(c, service)
		})
	}
	//照片表
	photo := internal.Group("/photo")
	{
		photo.POST("", func(c *gin.Context) {
			employee_action.UploadPhoto(c, service)
		})
		photo.GET("/:filename", func(c *gin.Context) {
			employee_action.DownloadPhoto(c)
		})
		photo.POST("/touch", func(c *gin.Context) {
			employee_action.PhotoTouchLuggageStorage(c, service)
		})
	}
	//寄存室表
	location := internal.Group("/location")
	{
		location.POST("", func(c *gin.Context) {
			employee_action.AddLocation(c, service)
		})
		location.PUT("/:id", func(c *gin.Context) {
			employee_action.UpdateLocation(c, service)
		})
		location.DELETE("/:id", func(c *gin.Context) {
			employee_action.DeleteLocation(c, service)
		})
		location.GET("", func(c *gin.Context) {
			employee_action.GetLocation(c, service)
		})
	}
	//客户表
	guest := internal.Group("/guest")
	{
		guest.GET("", func(c *gin.Context) {
			employee_action.GetGuest(c, service)
		})
	}
	//酒店表
	hotel := internal.Group("/hotel")
	{
		hotel.POST("", func(c *gin.Context) {
			admin.AddHotel(c, service)
		})
		hotel.DELETE("/:id", func(c *gin.Context) {
			employee_action.DeleteHotel(c, service)
		})
		hotel.PUT("/:id", func(c *gin.Context) {
			employee_action.UpdateHotel(c, service)
		})
		hotel.GET("", func(c *gin.Context) {
			employee_action.GetHotel(c, service)
		})
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
				admin.AddHotel(c, service)
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
