package message_queue

import (
	"context"
	"encoding/json"
	"errors"
	"time"

	"github.com/redis/go-redis/v9"
	"github.com/sirupsen/logrus"
	"hotel/internal/util/logger"
	"hotel/services"
)

type Task struct {
	ID         string                 `json:"id"`
	Type       string                 `json:"type"`   // 操作类型
	Data       map[string]interface{} `json:"data"`   // 请求数据
	Status     string                 `json:"status"` // pending, processing, completed, failed
	Result     interface{}            `json:"result"` // 操作结果
	Error      string                 `json:"error"`  // 错误信息
	CreateAt   time.Time              `json:"create_at"`
	CompleteAt *time.Time             `json:"complete_at,omitempty"`
}

// AddTask 提交任务到队列
func AddTask(ctx context.Context, s *services.Services, taskType string, data map[string]interface{}) (string, error) {
	// 生成任务ID
	taskID := generateTaskID()

	// 创建任务
	task := Task{
		ID:       taskID,
		Type:     taskType,
		Data:     data,
		Status:   "pending",
		CreateAt: time.Now(),
	}

	// 保存任务到Redis（用于状态查询）
	taskBytes, _ := json.Marshal(task)
	if err := s.RdbMq.Set(ctx, "task:"+taskID, taskBytes, time.Hour*24).Err(); err != nil {
		return "", err
	}

	// 将任务添加到队列
	err := s.RdbMq.XAdd(ctx, &redis.XAddArgs{
		Stream: "task_queue",
		Values: map[string]interface{}{
			"task_id": taskID,
		},
	}).Err()
	if err != nil {
		return "", err
	}

	return taskID, nil
}

// GetTaskStatus 查询任务状态
func GetTaskStatus(ctx context.Context, s *services.Services, taskID string) (*Task, error) {
	taskBytes, err := s.RdbMq.Get(ctx, "task:"+taskID).Result()
	if err != nil {
		return nil, err
	}

	var task Task
	if err := json.Unmarshal([]byte(taskBytes), &task); err != nil {
		return nil, err
	}

	return &task, nil
}

// UpdateTaskStatus 更新任务状态
func UpdateTaskStatus(ctx context.Context, s *services.Services, taskID, status string, result interface{}, errMsg string) {
	task, err := GetTaskStatus(ctx, s, taskID)
	if err != nil {
		logger.Logger.WithFields(logrus.Fields{
			"error":   err.Error(),
			"task_id": taskID,
		}).Error("获取任务失败")
		return
	}

	task.Status = status
	task.Result = result
	if status == "completed" {
		now := time.Now()
		task.CompleteAt = &now
	}
	if errMsg != "" {
		task.Error = errMsg
	}

	taskBytes, _ := json.Marshal(task)
	if err := s.RdbMq.Set(ctx, "task:"+taskID, taskBytes, time.Hour*24).Err(); err != nil {
		logger.Logger.WithFields(logrus.Fields{
			"error":   err.Error(),
			"task_id": taskID,
		}).Error("更新任务状态失败")
	}
}

// GetNextTask 从队列获取下一个任务
func GetNextTask(ctx context.Context, s *services.Services, group, consumer string) (string, error) {
	// 确保消费者组存在
	err := s.RdbMq.XGroupCreate(ctx, "task_queue", group, "0").Err()
	if err != nil {
		return "", err
	}

	// 读取任务
	result, err := s.RdbMq.XReadGroup(ctx, &redis.XReadGroupArgs{
		Group:    group,
		Consumer: consumer,
		Streams:  []string{"task_queue", ">"},
		Count:    1,
		Block:    time.Second * 5,
	}).Result()

	if err != nil || len(result) == 0 || len(result[0].Messages) == 0 {
		return "", nil // 没有任务
	}

	// 获取任务ID
	taskID, ok := result[0].Messages[0].Values["task_id"].(string)
	if !ok {
		return "", errors.New("无效的任务消息")
	}

	// 确认消息
	msgID := result[0].Messages[0].ID
	if err := s.RdbMq.XAck(ctx, "task_queue", group, msgID).Err(); err != nil {
		logger.Logger.WithFields(logrus.Fields{
			"error":  err.Error(),
			"msg_id": msgID,
			"group":  group,
		}).Error("确认任务消息失败")
	}

	return taskID, nil
}

// generateTaskID 简单的任务ID生成器
func generateTaskID() string {
	return time.Now().Format("20060102150405") + "-" + randomString(6)
}

// randomString 简单的随机字符串生成器
func randomString(length int) string {
	const charset = "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789"
	b := make([]byte, length)
	for i := range b {
		b[i] = charset[time.Now().UnixNano()%int64(len(charset))]
	}
	return string(b)
}
