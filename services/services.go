package services

import (
	"github.com/redis/go-redis/v9"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
	"log"
)

type Services struct {
	DB     *gorm.DB
	RdbAcc *redis.Client
	RdbRef *redis.Client
	RdbCac *redis.Client
	RdbLim *redis.Client
	RdbMq  *redis.Client
}

// NewDatabase 初始化数据库连接
func NewDatabase() *Services {
	dsn := "root:@furenjie321@tcp(127.0.0.1:3306)/study?charset=utf8mb4&parseTime=True&loc=Local"
	db, err := gorm.Open(mysql.Open(dsn), &gorm.Config{})
	if err != nil {
		log.Fatal(err)
	}

	//访问令牌的
	rdb := redis.NewClient(&redis.Options{
		Addr:     "127.0.0.1:6379",
		Password: "",
		DB:       0,
	})
	//刷新令牌的
	rdb1 := redis.NewClient(&redis.Options{
		Addr:     "127.0.0.1:6379",
		Password: "",
		DB:       1,
	})
	//实现缓存的
	rdb2 := redis.NewClient(&redis.Options{
		Addr:     "127.0.0.1:6379",
		Password: "",
		DB:       2,
	})
	//限流的
	rdb3 := redis.NewClient(&redis.Options{
		Addr:     "127.0.0.1:6379",
		Password: "",
		DB:       3,
	})
	//消息队列的
	rdb4 := redis.NewClient(&redis.Options{
		Addr:     "127.0.0.1:6379",
		Password: "",
		DB:       4,
	})

	service := Services{
		DB:     db,
		RdbAcc: rdb,
		RdbRef: rdb1,
		RdbCac: rdb2,
		RdbLim: rdb3,
		RdbMq:  rdb4,
	}
	return &service
}
