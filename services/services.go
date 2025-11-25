package services

import (
	"github.com/redis/go-redis/v9"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
	"log"
)

type Services struct {
	DB   *gorm.DB
	RDB  *redis.Client
	RDB1 *redis.Client
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
		Addr:     "localhost:6379",
		Password: "",
		DB:       0,
	})
	//刷新令牌的
	rdb1 := redis.NewClient(&redis.Options{
		Addr:     "localhost:6379",
		Password: "",
		DB:       1,
	})

	service := Services{
		DB:   db,
		RDB:  rdb,
		RDB1: rdb1,
	}
	return &service
}
