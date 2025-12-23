package middleware

import (
	"fmt"
	"github.com/sirupsen/logrus"
	"net"
	"net/http"
	"strings"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"hotel/internal/util"
	"hotel/models"
	"hotel/services"
)

// CheckAction 检查权限中间件 查询上下文中是否有相应权限
func CheckAction(name string) gin.HandlerFunc {
	return func(c *gin.Context) {
		_, ok := c.Get(name)
		if !ok {
			c.JSON(http.StatusUnauthorized, gin.H{
				"success": false,
				"message": "没有权限",
			})
			c.Abort()
			return
		}
		c.Next()
	}
}

// RateLimit 限流中间件
func RateLimit(name string, s *services.Services) gin.HandlerFunc {
	return func(c *gin.Context) {
		if !util.LimiterAllow(name, s.RdbLim, c) {
			c.JSON(http.StatusTooManyRequests, gin.H{
				"success": false,
				"message": "请求过于频繁，请稍后再试",
			})
			util.Logger.WithFields(logrus.Fields{
				"client_ip": c.ClientIP(),
			}).Warn("请求过于频繁")
			c.Abort()
			return
		}
		c.Next()
	}
}

// JwtCheck jwt检查中间件
func JwtCheck(s *services.Services) gin.HandlerFunc {
	return func(c *gin.Context) {
		//检查请求头
		authHeader := c.GetHeader("Authorization")
		if authHeader == "" {
			c.JSON(http.StatusUnauthorized, gin.H{
				"success": false,
				"message": "请先登录",
			})
			c.Abort()
			return
		}
		//提取令牌
		tokenString := strings.TrimPrefix(authHeader, "Bearer ")
		if tokenString == authHeader {
			c.JSON(http.StatusUnauthorized, gin.H{
				"success": false,
				"message": "Authorization格式应为Bearer {token}",
			})
			c.Abort()
			return
		}
		//解析令牌
		claims, err := util.ParseAccessToken(tokenString)
		if err != nil {

			if strings.Contains(err.Error(), "token is expired") {
				c.JSON(http.StatusUnauthorized, gin.H{
					"success": false,
					"message": "访问令牌已过期",
				})
			} else {
				c.JSON(http.StatusUnauthorized, gin.H{
					"success": false,
					"message": "token无效",
				})
			}
			c.Abort()
			return
		}
		//验证令牌有效性
		a := tokenString
		b := s.RdbAcc.Get(c, fmt.Sprintf("%d", claims.UserId)).Val()
		if a != b {
			c.JSON(http.StatusUnauthorized, gin.H{
				"success": false,
				"message": "token无效",
			})
			c.Abort()
			return
		}

		c.Set("claims", claims)
		c.Next()
	}
}

// AuthCheck 解析JWT并设置权限
func AuthCheck(s *services.Services) gin.HandlerFunc {
	return func(c *gin.Context) {
		//获取JWT声明
		claims, ok := c.Get("claims")
		if !ok {
			c.JSON(http.StatusUnauthorized, gin.H{
				"success": false,
				"message": "请先登录",
			})
			util.Logger.WithFields(logrus.Fields{
				"client_ip": c.ClientIP(),
			}).Warn("JWT令牌中未找到声明")
			c.Abort()
			return
		}
		e := claims.(*util.AccessClaims)

		//设置上下文
		c.Set("employee_id", e.UserId)
		c.Set("employee_name", e.UserName)
		c.Set("hotel_id", e.HotelId)
		//更新最后活动时间
		insert := models.Employee{
			LastActiveTime: time.Now(),
		}
		s.DB.Model(&models.Employee{}).Where("id = ?", e.UserId).Updates(insert)
		//查询角色权限
		var employee models.Employee
		s.DB.Model(&models.Employee{}).Preload("Role").Where("id = ?", e.UserId).First(&employee)
		var role models.Role
		s.DB.Model(&models.Role{}).Preload("Permissions").Where("id = ?", employee.RoleID).First(&role)
		for _, v := range role.Permissions {
			//设置权限
			c.Set(v.Name, 1)
		}
		c.Next()
	}
}

// RequestIDMiddleware 请求ID中间件
func RequestIDMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		//检查请求ID
		requestID := c.GetHeader("X-Request-ID")
		if requestID == "" {
			requestID = generateRequestID()
		}

		c.Set("request_id", requestID)

		c.Header("X-Request-ID", requestID)

		c.Next()
	}
}

// generateRequestID 生成请求ID
func generateRequestID() string {
	// 生成UUID v4
	return uuid.New().String()
}

// LogRequest 记录请求日志
func LogRequest() gin.HandlerFunc {
	return func(c *gin.Context) {
		start := time.Now()

		c.Next()

		//记录请求信息
		requestID, _ := c.Get("request_id")
		duration := time.Since(start)
		util.Logger.WithFields(logrus.Fields{
			"request_id": requestID,
			"method":     c.Request.Method,
			"path":       c.Request.URL.Path,
			"status":     c.Writer.Status(),
			"duration":   duration.String(),
		}).Info("请求处理完成")

	}
}

// FindIp 查找并显示IP地址
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

// isPrivateIP 判断是否为私有IP地址
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
