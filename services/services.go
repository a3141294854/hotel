package services

import (
	"fmt"
	"github.com/sirupsen/logrus"
	"hotel/internal/config"
	"hotel/internal/util/logger"

	"github.com/redis/go-redis/v9"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
)

type Services struct {
	DB      *gorm.DB
	RdbAcc  *redis.Client
	RdbRef  *redis.Client
	RdbCac  *redis.Client
	RdbLim  *redis.Client
	RdbRand *redis.Client
}

// NewDatabase 初始化数据库连接
func NewDatabase(cfg *config.Config) *Services {
	dsn := cfg.Database.GetDSN()
	db, err := gorm.Open(mysql.Open(dsn), &gorm.Config{})
	if err != nil {
		logger.Logger.WithFields(logrus.Fields{
			"error": err,
		}).Fatal("数据库连接失败")
	}

	//访问令牌的
	rdb := redis.NewClient(&redis.Options{
		Addr:     fmt.Sprintf("%s:%d", cfg.Redis.Host, cfg.Redis.Port),
		Password: cfg.Redis.Password,
		DB:       cfg.Redis.Databases.AccessToken,
	})
	//刷新令牌的
	rdb1 := redis.NewClient(&redis.Options{
		Addr:     fmt.Sprintf("%s:%d", cfg.Redis.Host, cfg.Redis.Port),
		Password: cfg.Redis.Password,
		DB:       cfg.Redis.Databases.RefreshToken,
	})
	//缓存的
	rdb2 := redis.NewClient(&redis.Options{
		Addr:     fmt.Sprintf("%s:%d", cfg.Redis.Host, cfg.Redis.Port),
		Password: cfg.Redis.Password,
		DB:       cfg.Redis.Databases.Cache,
	})
	//限流的
	rdb3 := redis.NewClient(&redis.Options{
		Addr:     fmt.Sprintf("%s:%d", cfg.Redis.Host, cfg.Redis.Port),
		Password: cfg.Redis.Password,
		DB:       cfg.Redis.Databases.RateLimit,
	})
	//随机取数的
	rdb4 := redis.NewClient(&redis.Options{
		Addr:     fmt.Sprintf("%s:%d", cfg.Redis.Host, cfg.Redis.Port),
		Password: cfg.Redis.Password,
		DB:       cfg.Redis.Databases.Random,
	})

	service := Services{
		DB:      db,
		RdbAcc:  rdb,
		RdbRef:  rdb1,
		RdbCac:  rdb2,
		RdbLim:  rdb3,
		RdbRand: rdb4,
	}
	return &service
}
