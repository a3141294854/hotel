package main

import (
	"fmt"
	"github.com/gin-gonic/gin"
	"hotel/internal/config"
	"hotel/internal/table"
	"hotel/internal/util"
	"hotel/internal/util/logger"
	"hotel/services"
)

func main() {
	//加载配置
	cfg, err := config.LoadConfig("")
	if err != nil {
		fmt.Println("配置文件加载失败:", err.Error())
		return
	}
	//启动日志
	logger.InitLogger(cfg.Log.Level, cfg.Log.Output, cfg.Log.FilePath, cfg.Log.MaxSize, cfg.Log.MaxBackups, cfg.Log.MaxAge)

	//建立数据库连接和表
	service := services.NewDatabase(cfg)
	table.Table(service.DB)

	//创建全局限流
	util.NewTokenBucketLimiter(cfg.RateLimiting.Default.Name, cfg.RateLimiting.Default.Capacity, cfg.RateLimiting.Default.FillRate, service)
	//启动jwt
	util.ConfigJwt(cfg)

	//启动路由
	r := gin.Default()
	open(r, service, cfg)

}
