package services

import (
	"fmt"
	"github.com/redis/go-redis/v9"
	"github.com/sirupsen/logrus"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
	"hotel/internal/util"
)

// Services 服务结构体
type Services struct {
	DB      *gorm.DB      //数据库连接
	RdbAcc  *redis.Client //访问令牌Redis
	RdbRef  *redis.Client //刷新令牌Redis
	RdbCac  *redis.Client //缓存Redis
	RdbLim  *redis.Client //限流Redis
	RdbRand *redis.Client //随机数Redis
}

// NewDatabase 初始化数据库连接
func NewDatabase(cfg *util.Config) *Services {
	//连接数据库
	dsn := cfg.Database.GetDSN()
	db, err := gorm.Open(mysql.Open(dsn), &gorm.Config{})
	if err != nil {
		util.Logger.WithFields(logrus.Fields{
			"error": err,
		}).Fatal("数据库连接失败")
	}

	//访问令牌Redis
	rdb := redis.NewClient(&redis.Options{
		Addr:     fmt.Sprintf("%s:%d", cfg.Redis.Host, cfg.Redis.Port),
		Password: cfg.Redis.Password,
		DB:       cfg.Redis.Databases.AccessToken,
	})
	//刷新令牌Redis
	rdb1 := redis.NewClient(&redis.Options{
		Addr:     fmt.Sprintf("%s:%d", cfg.Redis.Host, cfg.Redis.Port),
		Password: cfg.Redis.Password,
		DB:       cfg.Redis.Databases.RefreshToken,
	})
	//缓存Redis
	rdb2 := redis.NewClient(&redis.Options{
		Addr:     fmt.Sprintf("%s:%d", cfg.Redis.Host, cfg.Redis.Port),
		Password: cfg.Redis.Password,
		DB:       cfg.Redis.Databases.Cache,
	})
	//限流Redis
	rdb3 := redis.NewClient(&redis.Options{
		Addr:     fmt.Sprintf("%s:%d", cfg.Redis.Host, cfg.Redis.Port),
		Password: cfg.Redis.Password,
		DB:       cfg.Redis.Databases.RateLimit,
	})
	//随机数Redis
	rdb4 := redis.NewClient(&redis.Options{
		Addr:     fmt.Sprintf("%s:%d", cfg.Redis.Host, cfg.Redis.Port),
		Password: cfg.Redis.Password,
		DB:       cfg.Redis.Databases.Random,
	})

	//返回封装的连接
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
