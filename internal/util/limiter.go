package util

import (
	"context"
	"encoding/json"
	"github.com/gin-gonic/gin"
	"hotel/models"
	"hotel/services"
	"log"
	"time"
)

func NewTokenBucketLimiter(name string, capacity int, fillRate time.Duration, s *services.Services) {
	limiter := models.TokenBucketLimiter{
		Capacity:     capacity,
		FillRate:     fillRate,
		Tokens:       capacity,
		LastFillTime: time.Now(),
	}
	insert, err := json.Marshal(limiter)
	if err != nil {
		log.Println(name, "限流器序列化失败", err)
		return
	}
	s.RdbLim.Set(context.Background(), name, string(insert), 0)
}

func CreatLock(name string, s *services.Services) bool {
	lock, err := s.RdbLim.SetNX(context.Background(), name+"locked", true, 100*time.Millisecond).Result()
	if err != nil {
		log.Println(name, "锁创建失败", err)
		return false
	}
	if lock {
		return true
	} else {
		return false
	}
}

// LimiterAllow 检查是否允许请求通过，返回布尔值
func LimiterAllow(Name string, s *services.Services, c *gin.Context) bool {
	result := CreatLock(Name, s)
	defer s.RdbLim.Del(c, Name+"locked")
	if !result {
		//log.Println(Name, "有锁")
		return false
	}

	temp := s.RdbLim.Get(c, Name).Val()
	var limiter models.TokenBucketLimiter
	err := json.Unmarshal([]byte(temp), &limiter)
	if err != nil {
		log.Println(Name, "限流器反序列化失败", err)
		return false
	}
	now := time.Now()
	count := int(now.Sub(limiter.LastFillTime) / limiter.FillRate)
	limiter.Tokens = min(limiter.Tokens+count, limiter.Capacity)
	limiter.LastFillTime = now
	flag := true
	if limiter.Tokens > 0 {
		limiter.Tokens--
		flag = true
	} else {
		flag = false
	}

	insert, err := json.Marshal(limiter)
	if err != nil {
		log.Println(Name, "限流器序列化失败", err)
		return false
	}
	s.RdbLim.Set(c, Name, string(insert), 0)
	log.Println(Name, "限流器", limiter.Tokens)
	return flag
}
