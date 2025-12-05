package message_queue

import (
	"context"
	"errors"
	"hotel/models"
	"hotel/services"
	"log"
	"time"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

// StartTaskProcessor 启动任务处理器
func StartTaskProcessor(ctx context.Context, s *services.Services) {
	log.Println("启动任务处理器...")
	go func() {
		for {
			select {
			case <-ctx.Done():
				log.Println("任务处理器停止")
				return
			default:
				// 获取下一个任务
				taskID, err := GetNextTask(ctx, s, "processor", "processor_1")
				if err != nil {
					log.Println("获取任务失败:", err)
					continue
				}

				if taskID == "" {
					// 没有任务，稍后重试
					time.Sleep(time.Millisecond * 500)
					continue
				}

				// 处理任务
				log.Println("开始处理任务:", taskID)
				ProcessTask(ctx, s, taskID)
			}
		}
	}()
}

// ProcessTask 处理任务
func ProcessTask(ctx context.Context, s *services.Services, taskID string) {
	// 更新任务状态为处理中
	UpdateTaskStatus(ctx, s, taskID, "processing", nil, "")

	// 获取任务详情
	task, err := GetTaskStatus(ctx, s, taskID)
	if err != nil {
		UpdateTaskStatus(ctx, s, taskID, "failed", nil, "获取任务失败: "+err.Error())
		return
	}

	// 根据任务类型执行相应操作
	var result interface{}
	var processErr error

	switch task.Type {
	case "add_luggage":
		result, processErr = processAddLuggage(s, task.Data)
	case "delete_luggage":
		result, processErr = processDeleteLuggage(s, task.Data)
	case "update_luggage":
		result, processErr = processUpdateLuggage(s, task.Data)
	default:
		processErr = errors.New("未知任务类型: " + task.Type)
	}

	// 更新任务状态
	if processErr != nil {
		UpdateTaskStatus(ctx, s, taskID, "failed", nil, processErr.Error())
		log.Printf("任务 %s 处理失败: %v", taskID, processErr)
	} else {
		UpdateTaskStatus(ctx, s, taskID, "completed", result, "")
		log.Printf("任务 %s 处理成功", taskID)
	}
}

// processAddLuggage 处理添加行李任务
func processAddLuggage(s *services.Services, data map[string]interface{}) (interface{}, error) {
	// 解析请求数据
	guestName, ok := data["guest_name"].(string)
	if !ok {
		return nil, errors.New("客户姓名不能为空")
	}

	tag, _ := data["tag"].(string)
	weight, _ := data["weight"].(float32)
	status, _ := data["status"].(string)
	location, _ := data["location"].(string)

	// 先创建或获取客户记录
	var guest models.Guest
	result := s.DB.Where("guest_name = ?", guestName).First(&guest)
	if result.Error != nil {
		if errors.Is(result.Error, gorm.ErrRecordNotFound) {
			// 客户不存在，创建新客户
			guest = models.Guest{
				Name: guestName,
			}
			if err := s.DB.Create(&guest).Error; err != nil {
				return nil, err
			}
		} else {
			// 其他查询错误
			return nil, result.Error
		}
	}

	// 创建行李记录
	luggage := models.Luggage{
		GuestID:  guest.ID,
		Tag:      tag,
		Weight:   weight,
		Status:   status,
		Location: location,
	}

	result = s.DB.Create(&luggage)
	if result.Error != nil {
		return nil, result.Error
	}

	// 返回成功结果
	return gin.H{
		"id":         luggage.ID,
		"guest_id":   luggage.GuestID,
		"guest_name": guestName,
		"tag":        luggage.Tag,
		"weight":     luggage.Weight,
		"status":     luggage.Status,
		"location":   luggage.Location,
	}, nil
}

// processDeleteLuggage 处理删除行李任务
func processDeleteLuggage(s *services.Services, data map[string]interface{}) (interface{}, error) {
	// 解析请求数据
	luggageID, ok := data["id"].(float64)
	if !ok {
		return nil, errors.New("行李ID不能为空")
	}

	// 检查行李是否存在
	var existingLuggage models.Luggage
	if err := s.DB.Where("status = ?", "寄存中").First(&existingLuggage, uint(luggageID)).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, errors.New("行李记录不存在或已取出")
		}
		return nil, err
	}

	// 更新行李状态
	existingLuggage.Status = "已取出"
	existingLuggage.Location = "已取出"
	result := s.DB.Model(&models.Luggage{}).Where("id = ?", luggageID).Updates(existingLuggage)
	if result.Error != nil {
		return nil, result.Error
	}

	// 软删除记录
	if err := s.DB.Where("id = ?", luggageID).Delete(&models.Luggage{}).Error; err != nil {
		return nil, err
	}

	return gin.H{
		"success": true,
		"message": "行李删除成功",
	}, nil
}

// processUpdateLuggage 处理更新行李任务
func processUpdateLuggage(s *services.Services, data map[string]interface{}) (interface{}, error) {
	// 解析请求数据
	luggageID, ok := data["id"].(float64)
	if !ok {
		return nil, errors.New("行李ID不能为空")
	}

	// 先检查记录是否存在
	var existingLuggage models.Luggage
	if err := s.DB.First(&existingLuggage, uint(luggageID)).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, errors.New("行李记录不存在")
		}
		return nil, err
	}

	// 更新字段
	if status, ok := data["status"].(string); ok {
		existingLuggage.Status = status
	}
	if location, ok := data["location"].(string); ok {
		existingLuggage.Location = location
	}

	// 执行更新
	result := s.DB.Model(&existingLuggage).Updates(existingLuggage)
	if result.Error != nil {
		return nil, result.Error
	}

	// 获取客户信息
	var guest models.Guest
	s.DB.First(&guest, existingLuggage.GuestID)

	return gin.H{
		"id":         existingLuggage.ID,
		"guest_id":   existingLuggage.GuestID,
		"guest_name": guest.Name,
		"tag":        existingLuggage.Tag,
		"weight":     existingLuggage.Weight,
		"status":     existingLuggage.Status,
		"location":   existingLuggage.Location,
	}, nil
}
