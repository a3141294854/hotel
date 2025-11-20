package middleware

import (
	"fmt"
	"log"
	"net"
	"net/http"
	"strings"
	"sync"
	"time"

	"github.com/gin-gonic/gin"

	"hotel/internal/util"
)

// TokenBucketLimiter 令牌桶限流器
type TokenBucketLimiter struct {
	Capacity     int           // 桶容量
	FillRate     time.Duration // 添加令牌速率，如每10ms加1个令牌
	tokens       int           // 当前令牌数
	lastFillTime time.Time     // 上次添加令牌时间
	mutex        sync.Mutex
}

// NewTokenBucketLimiter 创建新的令牌桶限流器
func NewTokenBucketLimiter(capacity int, fillRate time.Duration) *TokenBucketLimiter {
	return &TokenBucketLimiter{
		Capacity:     capacity,
		FillRate:     fillRate,
		tokens:       capacity,
		lastFillTime: time.Now(),
	}
}

// Allow 检查是否允许请求通过，返回布尔值
func (l *TokenBucketLimiter) Allow() bool {
	l.mutex.Lock()
	defer l.mutex.Unlock()

	now := time.Now()
	// 计算应该添加的令牌数
	elapsed := now.Sub(l.lastFillTime)
	newTokens := int(elapsed / l.FillRate)

	if newTokens > 0 {
		l.tokens += newTokens
		if l.tokens > l.Capacity {
			l.tokens = l.Capacity
		}
		l.lastFillTime = now
	}

	if l.tokens <= 0 {
		return false
	}
	l.tokens--
	return true
}

// RateLimit 限流中间件
func RateLimit(limiter *TokenBucketLimiter) gin.HandlerFunc {
	return func(c *gin.Context) {
		if !limiter.Allow() {
			c.JSON(http.StatusTooManyRequests, gin.H{
				"success": false,
				"message": "请求过于频繁，请稍后再试",
			})
			c.Abort()
			return
		}
		c.Next()
	}
}

// JwtCheck jwt检查中间件
func JwtCheck() gin.HandlerFunc {
	return func(c *gin.Context) {
		authHeader := c.GetHeader("Authorization")
		if authHeader == "" {
			c.JSON(http.StatusUnauthorized, gin.H{
				"success": false,
				"message": "请先登录",
			})
			c.Abort()
			return
		}

		tokenString := strings.TrimPrefix(authHeader, "Bearer ")
		if tokenString == authHeader {
			c.JSON(http.StatusUnauthorized, gin.H{
				"success": false,
				"message": "Authorization格式应为Bearer {token}",
			})
			c.Abort()
			return
		}

		claims, err := util.ParseAccessToken(tokenString)
		if err != nil {
			c.JSON(http.StatusUnauthorized, gin.H{
				"success": false,
				"message": "token无效",
			})
			log.Println("token无效", tokenString)
			c.Abort()
			return
		}
		c.Set("claims", claims)
		c.Next()
	}
}

// AuthCheck 权限检查中间件
func AuthCheck() gin.HandlerFunc {
	return func(c *gin.Context) {

		user, ok := c.Get("claims")
		if !ok {
			c.JSON(http.StatusUnauthorized, gin.H{
				"success": false,
				"message": "请先登录",
			})
			log.Println("没找到声明")
			c.Abort()
			return
		}
		if user.(*util.AccessClaims).UserId == 0 {
			c.JSON(http.StatusUnauthorized, gin.H{
				"success": false,
				"message": "请先登录",
			})
			log.Println("声明id错误")
			c.Abort()
			return
		}
		c.Next()
	}
}

func FindIp() {
	_, err := net.InterfaceAddrs()
	if err != nil {
		fmt.Println("获取网络地址失败:", err)
		return
	}

	// 获取网络接口信息
	interfaces, err := net.Interfaces()
	if err != nil {
		fmt.Println("获取网络接口失败:", err)
		return
	}

	fmt.Println("=== 可用的网络接口和IP地址 ===")

	for _, iface := range interfaces {
		// 跳过回环接口和非活动接口
		if iface.Flags&net.FlagLoopback != 0 || iface.Flags&net.FlagUp == 0 {
			continue
		}

		addrs, err := iface.Addrs()
		if err != nil {
			continue
		}

		for _, addr := range addrs {
			if ipNet, ok := addr.(*net.IPNet); ok && !ipNet.IP.IsLoopback() {
				if ipNet.IP.To4() != nil {
					ip := ipNet.IP.String()
					status := "活动"
					if iface.Flags&net.FlagRunning == 0 {
						status = "非活动"
					}

					fmt.Printf("接口: %-8s | IPv4: %-15s | 状态: %-6s", iface.Name, ip, status)

					// 判断是否为局域网地址
					if isPrivateIP(ipNet.IP) {
						fmt.Println(" | ✓ 局域网地址 - 其他设备可通过此IP访问")
					} else {
						fmt.Println(" | ✗ 非局域网地址")
					}
				}
			}
		}
	}

	fmt.Println("\n=== 访问建议 ===")
	fmt.Println("1. 使用标记为 '✓ 局域网地址' 的IP进行访问")
	fmt.Println("2. 访问格式: http://IP地址:8080")
	fmt.Println("3. 确保防火墙允许8080端口")
	fmt.Println("4. 确保设备在同一局域网内")
}

// 判断是否为私有IP地址
func isPrivateIP(ip net.IP) bool {
	// 10.0.0.0/8
	if ip4 := ip.To4(); ip4 != nil {
		switch {
		case ip4[0] == 10:
			return true
		case ip4[0] == 172 && ip4[1] >= 16 && ip4[1] <= 31:
			return true
		case ip4[0] == 192 && ip4[1] == 168:
			return true
		}
	}
	return false
}
