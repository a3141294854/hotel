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
	cfg, err := config.LoadConfig("")
	if err != nil {
		fmt.Println("配置文件加载失败:", err.Error())
		return
	}
	logger.InitLogger(cfg.Log.Level, cfg.Log.Output, cfg.Log.FilePath, cfg.Log.MaxSize, cfg.Log.MaxBackups, cfg.Log.MaxAge)

	service := services.NewDatabase(cfg)
	table.Table(service.DB)

	util.NewTokenBucketLimiter(cfg.RateLimiting.Default.Name, cfg.RateLimiting.Default.Capacity, cfg.RateLimiting.Default.FillRate, service)
	util.ConfigJwt(cfg)

	r := gin.Default()
	open(r, service, cfg)

}
