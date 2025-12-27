package util

import (
	"fmt"
	"io"
	"os"
	"path/filepath"
	"strings"

	"github.com/sirupsen/logrus"
	"gopkg.in/natefinch/lumberjack.v2"
)

var Logger *logrus.Logger

// InitLogger 初始化日志
func InitLogger(logLevel, output, filePath string, maxSize, maxBackups, maxAge int) {
	Logger = logrus.New()

	// 设置日志级别
	level, err := logrus.ParseLevel(logLevel)
	if err != nil {
		level = logrus.InfoLevel
	}
	Logger.SetLevel(level)

	// 设置日志格式
	Logger.SetFormatter(&CustomFormatter{})

	// 设置输出，全部转换成小写
	switch strings.ToLower(output) {
	case "file":
		// 确保目录存在，函数的作用是获取路径的目录部分
		//os.Stat() 获取文件或目录的信息
		//os.IsNotExist(err) 检查错误是否表示文件/目录不存在
		//os.MkdirAll() 创建目录，包括所有必要的父目录
		//0755 是权限设置，表示所有者有读写执行权限，组和其他用户有读执行权限
		dir := filepath.Dir(filePath)
		if _, err := os.Stat(dir); os.IsNotExist(err) {
			err = os.MkdirAll(dir, 0755)
			if err != nil {
				fmt.Println("创建日志目录失败:", err.Error())
			}
		}

		// 日志轮转
		//lumberjack.Logger 是一个日志轮转库，用于管理日志文件
		//Filename 日志文件路径
		//MaxSize 单个日志文件的最大大小(MB)
		//MaxBackups 保留的旧日志文件数量
		//MaxAge 保留日志文件的最大天数
		//Compress 是否压缩旧日志文件
		Logger.SetOutput(&lumberjack.Logger{
			Filename:   filePath,
			MaxSize:    maxSize,    // MB
			MaxBackups: maxBackups, // 保留的旧日志文件数量
			MaxAge:     maxAge,     // 保留日志文件的最大天数
			Compress:   true,       // 压缩旧日志文件
		})
	case "both":
		// 同时输出到文件和控制台
		dir := filepath.Dir(filePath)
		if _, err := os.Stat(dir); os.IsNotExist(err) {
			err = os.MkdirAll(dir, 0755)
			if err != nil {
				fmt.Println("创建日志目录失败:", err.Error())
			}
		}

		fileWriter := &lumberjack.Logger{
			Filename:   filePath,
			MaxSize:    maxSize,
			MaxBackups: maxBackups,
			MaxAge:     maxAge,
			Compress:   true,
		}
		//io.MultiWriter() 创建一个写入器，可以将数据同时写入多个写入器
		//os.Stdout 代表标准输出（控制台）
		//fileWriter 是文件写入器
		//这样设置后，日志会同时输出到控制台和文件

		Logger.SetOutput(io.MultiWriter(os.Stdout, fileWriter))
	default: // console
		Logger.SetOutput(os.Stdout)
	}
}

// CustomFormatter 自定义日志格式
type CustomFormatter struct{}

func (f *CustomFormatter) Format(entry *logrus.Entry) ([]byte, error) {
	timestamp := entry.Time.Format("2006-01-02 15:04:05")
	level := strings.ToUpper(entry.Level.String())
	msg := entry.Message

	// 如果有字段信息，按固定顺序输出
	var fieldsStr string
	if len(entry.Data) > 0 {
		// 定义字段输出顺序（添加 "客户端IP"）
		fieldOrder := []string{"请求id", "请求方法", "路径", "状态码", "使用时间", "操作员", "客户端IP"}

		fields := make([]string, 0, len(entry.Data))

		// 按固定顺序添加字段
		for _, key := range fieldOrder {
			if value, exists := entry.Data[key]; exists {
				fields = append(fields, fmt.Sprintf("%s=%v", key, value))
			}
		}

		// 添加其他未在 fieldOrder 中的字段
		for k, v := range entry.Data {
			found := false
			for _, key := range fieldOrder {
				if k == key {
					found = true
					break
				}
			}
			if !found {
				fields = append(fields, fmt.Sprintf("%s=%v", k, v))
			}
		}

		fieldsStr = fmt.Sprintf(" - %s", strings.Join(fields, ", "))
	}

	// 格式: [时间] [级别] 消息 - 字段1=值1, 字段2=值2
	output := fmt.Sprintf("[%s] [%s] %s%s\n", timestamp, level, msg, fieldsStr)
	return []byte(output), nil
}
