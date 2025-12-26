package util

import (
	"context"
	"encoding/json"
	"github.com/redis/go-redis/v9"
	"github.com/sirupsen/logrus"
	"time"

	"github.com/gin-gonic/gin"
)

// NewTokenBucketLimiter 创建令牌桶限流器
func NewTokenBucketLimiter(name string, capacity int, fillRate time.Duration, RdbLim *redis.Client) {
	limiter := TokenBucketLimiter{
		Capacity:     capacity, //容量
		FillRate:     fillRate, //填充速率
		Tokens:       capacity, //令牌数
		LastFillTime: time.Now().UnixNano(),
	}
	//序列化
	insert, err := json.Marshal(limiter)
	if err != nil {
		Logger.WithFields(logrus.Fields{
			"name":  name,
			"error": err,
		}).Error("限流器序列化失败")
		return
	}
	//设置
	RdbLim.Set(context.Background(), name, string(insert), 0)
}

// LimiterAllow 限流检查
func LimiterAllow(Name string, RdbLim *redis.Client, c *gin.Context) bool {
	luaScipt := `
--获取数据
local data = redis.call('GET', KEYS[1])

if not data then
	return {0, 0}
end

--解析
local limiter = cjson.decode(data)
--计算
local elapsed = tonumber(ARGV[1]) - limiter.LastFillTime 
local count = math.floor(elapsed / limiter.FillRate)
limiter.Tokens = math.min(limiter.Tokens + count, limiter.Capacity)

--检查
local flag = 0
if limiter.Tokens > 0 then
    limiter.LastFillTime = tonumber(ARGV[1])
	limiter.Tokens = limiter.Tokens - 1
	flag = 1
end
--序列化
local insert = cjson.encode(limiter)
--更新
redis.call('SET', KEYS[1], insert)
return {flag, limiter.Tokens}
 `
	//运行限流脚本
	result, err := RdbLim.Eval(c, luaScipt, []string{Name}, time.Now().UnixNano()).Result()
	if err != nil {
		Logger.WithFields(logrus.Fields{
			"name":  Name,
			"error": err,
		}).Error("限流器执行脚本失败")
		return false
	}
	//获取数据
	results := result.([]interface{})
	flag := results[0].(int64)
	tokens := results[1].(int64)
	Logger.WithFields(logrus.Fields{
		"name":   Name,
		"tokens": tokens,
	}).Debug("限流器状态")

	if flag == 1 {
		return true
	} else {
		return false
	}
}

type TokenBucketLimiter struct {
	Capacity     int           // 桶容量
	FillRate     time.Duration // 添加令牌速率，如每10ms加1个令牌
	Tokens       int           // 当前令牌数
	LastFillTime int64         // 上次添加令牌时间
}
