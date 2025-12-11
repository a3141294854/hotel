package config

import (
	"fmt"
	"os"
	"time"

	"github.com/goccy/go-yaml"
)

// Config 应用程序配置结构
type Config struct {
	Database     DatabaseConfig     `yaml:"database"`
	Redis        RedisConfig        `yaml:"redis"`
	Server       ServerConfig       `yaml:"server"`
	JWT          JWTConfig          `yaml:"jwt"`
	RateLimiting RateLimitingConfig `yaml:"rate_limiting"`
}

// DatabaseConfig 数据库配置
type DatabaseConfig struct {
	Driver    string `yaml:"driver"`
	Host      string `yaml:"host"`
	Port      int    `yaml:"port"`
	Username  string `yaml:"username"`
	Password  string `yaml:"password"`
	Name      string `yaml:"name"`
	Charset   string `yaml:"charset"`
	ParseTime bool   `yaml:"parse_time"`
	Loc       string `yaml:"loc"`
}

// GetDSN 获取数据库连接字符串
func (dc DatabaseConfig) GetDSN() string {
	return fmt.Sprintf("%s:%s@tcp(%s:%d)/%s?charset=%s&parseTime=%t&loc=%s",
		dc.Username, dc.Password, dc.Host, dc.Port, dc.Name, dc.Charset, dc.ParseTime, dc.Loc)
}

// RedisConfig Redis配置
type RedisConfig struct {
	Host      string               `yaml:"host"`
	Port      int                  `yaml:"port"`
	Password  string               `yaml:"password"`
	Databases RedisDatabasesConfig `yaml:"databases"`
}

// RedisDatabasesConfig Redis数据库配置
type RedisDatabasesConfig struct {
	AccessToken  int `yaml:"access_token"`
	RefreshToken int `yaml:"refresh_token"`
	Cache        int `yaml:"cache"`
	RateLimit    int `yaml:"rate_limit"`
	MessageQueue int `yaml:"message_queue"`
}

// ServerConfig 服务器配置
type ServerConfig struct {
	Host string `yaml:"host"`
	Port int    `yaml:"port"`
	Mode string `yaml:"mode"`
}

// JWTConfig JWT配置
type JWTConfig struct {
	SecretKey            string        `yaml:"secret_key"`
	AccessTokenDuration  time.Duration `yaml:"access_token_duration"`
	RefreshTokenDuration time.Duration `yaml:"refresh_token_duration"`
}

// RateLimitingConfig 限流配置
type RateLimitingConfig struct {
	Default   location   `yaml:"default"`
	Locations []location `yaml:"locations"`
}

type location struct {
	Name     string        `yaml:"name"`
	Capacity int           `yaml:"capacity"`
	FillRate time.Duration `yaml:"fill_rate"`
}

// LoadConfig 加载配置文件
func LoadConfig(configPath string) (*Config, error) {
	// 默认配置路径
	if configPath == "" {
		configPath = "configs/config.yaml"
	}

	// 检查文件是否存在
	if _, err := os.Stat(configPath); os.IsNotExist(err) {
		return nil, fmt.Errorf("配置文件不存在: %s", configPath)
	}

	// 读取文件
	data, err := os.ReadFile(configPath)
	if err != nil {
		return nil, fmt.Errorf("读取配置文件失败: %v", err)
	}

	// 解析YAML
	var config Config
	if err := yaml.Unmarshal(data, &config); err != nil {
		return nil, fmt.Errorf("解析配置文件失败: %v", err)
	}

	return &config, nil
}
