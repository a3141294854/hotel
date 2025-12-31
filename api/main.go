package main

import (
	"fmt"
	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
	"hotel/internal/table"
	"hotel/internal/util"
	"hotel/services"
	"time"
)

func main() {
	//加载配置
	cfg, err := util.LoadConfig("")
	if err != nil {
		fmt.Println("配置文件加载失败:", err.Error())
		return
	}
	//设置Gin运行模式
	gin.SetMode(cfg.Server.Mode)
	//启动日志
	util.InitLogger(cfg.Log.Level, cfg.Log.Output, cfg.Log.FilePath, cfg.Log.MaxSize, cfg.Log.MaxBackups, cfg.Log.MaxAge)
	//建立数据库连接和表
	service := services.NewDatabase(cfg)
	table.Table(service.DB)

	//创建全局限流
	util.NewTokenBucketLimiter(cfg.RateLimiting.Default.Name, cfg.RateLimiting.Default.Capacity, cfg.RateLimiting.Default.FillRate, service.RdbLim)
	//启动jwt
	util.ConfigJwt(cfg.JWT.AccessTokenDuration, cfg.JWT.RefreshTokenDuration, cfg.JWT.SecretKey)

	//启动路由
	r := gin.New()
	// ====================
	// 添加 CORS 中间件配置
	// ====================

	// 使用 cors.New() 创建一个新的 CORS 中间件实例
	// cors.Config{} 是一个结构体,用于配置 CORS 的各种选项
	r.Use(cors.New(cors.Config{

		// 允许跨域请求的来源(域名/IP)
		// "*" 表示允许所有来源(开发环境使用)
		// 生产环境应该指定具体域名,例如: []string{"https://www.example.com"}
		AllowOrigins: []string{"*"},

		// 允许的 HTTP 请求方法
		// 这些方法是浏览器可以发送给后端的请求类型
		AllowMethods: []string{
			"GET",     // 获取数据
			"POST",    // 提交数据
			"PUT",     // 更新数据
			"DELETE",  // 删除数据
			"OPTIONS", // 预检请求(浏览器在跨域前自动发送)
		},

		// 允许的请求头(Header)
		// 浏览器可以在请求中包含这些 HTTP 头
		AllowHeaders: []string{
			"Origin",           // 请求来源(浏览器自动添加)
			"Content-Type",     // 请求内容类型,如 application/json
			"Authorization",    // 认证信息,如 JWT token
			"X-Requested-With", // 标识是 AJAX 请求
		},

		// 允许前端访问的响应头(Header)
		// 这些响应头可以被前端 JavaScript 读取
		ExposeHeaders: []string{
			"Content-Length", // 响应内容长度
			"Authorization",  // 认证 token
		},

		// 允许发送凭证(Cookie、Authorization 等)
		// true 表示前端可以携带认证信息
		AllowCredentials: true,

		// 预检请求的缓存时间
		// 浏览器会缓存 OPTIONS 请求的结果 12 小时
		// 这段时间内不会再次发送预检请求,提高性能
		MaxAge: 12 * time.Hour,
	}))

	// ====================
	// CORS 配置结束
	// ====================
	open(r, service, cfg)

}
