package util

import (
	"context"
	"encoding/json"
	"github.com/sirupsen/logrus"
	"time"

	"github.com/gin-gonic/gin"

	"hotel/models"

	"hotel/internal/util/logger"
	"hotel/services"
)

func NewTokenBucketLimiter(name string, capacity int, fillRate time.Duration, s *services.Services) {
	limiter := models.TokenBucketLimiter{
		Capacity:     capacity,
		FillRate:     fillRate,
		Tokens:       capacity,
		LastFillTime: time.Now().UnixNano(),
	}
	insert, err := json.Marshal(limiter)
	if err != nil {
		logger.Logger.WithFields(logrus.Fields{
			"name":  name,
			"error": err,
		}).Error("限流器序列化失败")
		return
	}
	s.RdbLim.Set(context.Background(), name, string(insert), 0)
}

func LimiterAllow(Name string, s *services.Services, c *gin.Context) bool {
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
limiter.LastFillTime = tonumber(ARGV[1])
--检查
local flag = 0
if limiter.Tokens > 0 then
	limiter.Tokens = limiter.Tokens - 1
	flag = 1
end
--序列化
local insert = cjson.encode(limiter)
--更新
redis.call('SET', KEYS[1], insert)
return {flag, limiter.Tokens}
 `
	result, err := s.RdbLim.Eval(c, luaScipt, []string{Name}, time.Now().UnixNano()).Result()
	if err != nil {
		logger.Logger.WithFields(logrus.Fields{
			"name":  Name,
			"error": err,
		}).Error("限流器执行脚本失败")
		return false
	}
	results := result.([]interface{})
	flag := results[0].(int64)
	tokens := results[1].(int64)
	logger.Logger.WithFields(logrus.Fields{
		"name":   Name,
		"tokens": tokens,
	}).Debug("限流器状态")

	if flag == 1 {
		return true
	} else {
		return false
	}
}
