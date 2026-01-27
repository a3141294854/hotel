package employee_action

import (
	"fmt"
	"github.com/google/uuid"
	"github.com/tencentyun/cos-go-sdk-v5"

	"context"

	"github.com/gin-gonic/gin"
	"github.com/sirupsen/logrus"
	"hotel/internal/util"
	"hotel/models"
	"hotel/services"
	"math/rand"
	"net/http"
	"os"
	"path/filepath"
	"strconv"
	"strings"
)

// UploadPhoto 上传照片（支持多张）
func UploadPhoto(c *gin.Context, s *services.Services, client *cos.Client, cfg *util.Config) {
	// 获取表单数据
	form, err := c.MultipartForm()
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"success": false,
			"message": "获取上传文件失败",
		})
		util.Logger.WithFields(logrus.Fields{
			"error": err.Error(),
			"请求id":  c.GetUint("request_id"),
		}).Error("获取上传文件失败")
		return
	}

	// 获取上传的文件列表
	files := form.File["photos"] // ← 前端需要使用 "photos[]" 字段名
	if len(files) == 0 {
		c.JSON(http.StatusBadRequest, gin.H{
			"success": false,
			"message": "请上传照片",
		})
		return
	}

	// 检查文件数量限制（可选）
	if len(files) > 10 {
		c.JSON(http.StatusBadRequest, gin.H{
			"success": false,
			"message": "一次最多上传 10 张照片",
		})
		return
	}

	/* 创建保存目录
	uploadDir := "./uploads/photos"
	if err := os.MkdirAll(uploadDir, 0755); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"success": false,
			"message": "创建上传目录失败",
		})
		util.Logger.WithFields(logrus.Fields{
			"error": err.Error(),
			"请求id":  c.GetUint("request_id"),
		}).Error("创建上传目录失败")
		return
	}*/

	// 处理每个文件
	var uploadedPhotos []gin.H

	for _, file := range files {
		// 检查文件大小（限制 5MB）
		const maxSize = 5 * 1024 * 1024
		if file.Size > maxSize {
			c.JSON(http.StatusBadRequest, gin.H{
				"success": false,
				"message": fmt.Sprintf("照片 %s 大小超过 5MB", file.Filename),
			})
			return
		}

		// 检查文件类型
		ext := filepath.Ext(file.Filename)
		allowedExts := map[string]bool{
			".jpg":  true,
			".jpeg": true,
			".png":  true,
			".gif":  true,
		}
		if !allowedExts[ext] {
			c.JSON(http.StatusBadRequest, gin.H{
				"success": false,
				"message": fmt.Sprintf("照片 %s 格式不支持，只支持 jpg、jpeg、png、gif", file.Filename),
			})
			return
		}

		// 生成唯一文件名
		uploadDir := "uploads/photos"
		filename := uuid.New().String() + ext
		cosKey := strings.Join([]string{uploadDir, filename}, "/")

		// 打开文件
		fileReader, err := file.Open()
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{
				"success": false,
				"message": "打开文件失败",
			})
			util.Logger.WithFields(logrus.Fields{
				"error": err.Error(),
				"请求id":  c.GetString("request_id"),
			}).Error("打开文件失败")
			return
		}
		defer fileReader.Close()

		//直接上传到 COS
		_, err = client.Object.Put(
			context.Background(),
			cosKey,
			fileReader,
			nil,
		)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{
				"success": false,
				"message": "上传到 COS 失败",
			})
			util.Logger.WithFields(logrus.Fields{
				"error":   err.Error(),
				"cos_key": cosKey,
				"请求id":    c.GetString("request_id"),
			}).Error("上传到 COS 失败")
			return
		}

		/* 保存文件
		if err := c.SaveUploadedFile(file, path); err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{
				"success": false,
				"message": fmt.Sprintf("保存照片 %s 失败", file.Filename),
			})
			util.Logger.WithFields(logrus.Fields{
				"error": err.Error(),
				"请求id":  c.GetUint("request_id"),
			}).Error("保存照片失败")
			return
		}*/

		urlPath := fmt.Sprintf("%s/%s", cfg.Cos.Website, cosKey)

		// 保存到数据库
		photoURL := fmt.Sprintf("%s", urlPath)
		result := s.DB.Model(&models.Photo{}).Create(&models.Photo{
			Url: photoURL,
		})
		if result.Error != nil {
			c.JSON(http.StatusInternalServerError, gin.H{
				"success": false,
				"message": "保存照片失败",
			})
			util.Logger.WithFields(logrus.Fields{
				"error": result.Error.Error(),
				"请求id":  c.GetString("request_id"),
			}).Error("保存照片失败")
			return
		}

		// 添加到结果列表
		uploadedPhotos = append(uploadedPhotos, gin.H{
			"url": photoURL,
		})
	}

	c.JSON(http.StatusOK, gin.H{
		"success": true,
		"message": fmt.Sprintf("成功上传 %d 张照片", len(uploadedPhotos)),
		"data": gin.H{
			"count":  len(uploadedPhotos),
			"photos": uploadedPhotos,
		},
	})
}

// DownloadPhoto 获取照片
func DownloadPhoto(c *gin.Context) {
	filename := c.Param("filename")
	if filename == "" {
		c.JSON(http.StatusBadRequest, gin.H{
			"success": false,
			"message": "请提供文件名",
		})
		return
	}

	if strings.Contains(filename, "..") {
		c.AbortWithStatus(http.StatusForbidden)
		return
	}

	path := fmt.Sprintf("./uploads/photos/%s", filename)

	// 检查文件是否存在
	if _, err := os.Stat(path); os.IsNotExist(err) {
		c.JSON(http.StatusNotFound, gin.H{
			"success": false,
			"message": "照片不存在",
		})
		return
	}

	// 返回文件
	c.File(path)
}

func StaticDownloadPhoto(c *gin.Context) {
	filename := c.Param("filename")
	if filename == "" {
		c.JSON(http.StatusBadRequest, gin.H{
			"success": false,
			"message": "请提供文件地址",
		})
		return
	}
	//规范化路径，防止.和/
	cleanPath := filepath.Clean(filename)

	if strings.Contains(cleanPath, "..") {
		c.AbortWithStatus(http.StatusForbidden)
		return
	}

	localPath := fmt.Sprintf("./uploads/photos/%s", filename)

	// 检查文件是否存在
	if _, err := os.Stat(localPath); os.IsNotExist(err) {
		c.JSON(http.StatusNotFound, gin.H{
			"success": false,
			"message": "照片不存在",
		})
		return
	}

	// 返回文件
	c.File(localPath)

}

func PhotoTouchLuggageStorage(c *gin.Context, s *services.Services) {

	LuggageStorageId := c.Query("luggage_storage_id")
	FileName := c.Query("file_name")
	LuggageStorageID, err := strconv.ParseUint(LuggageStorageId, 10, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"success": false,
			"message": "请提供行李寄存ID",
			"请求id":    c.GetUint("request_id"),
		})
		return
	}

	// 检查参数
	if FileName == "" || LuggageStorageID == 0 {
		c.JSON(http.StatusBadRequest, gin.H{
			"success": false,
			"message": "请提供照片名称和行李寄存ID",
		})
		return
	}
	// 检查行李寄存是否存在
	ok, err := util.ExIf(s.DB, "id", &models.LuggageStorage{}, fmt.Sprintf("%d", LuggageStorageID))
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"success": false,
			"message": "查询行李寄存失败",
		})
		return
	}
	if !ok {
		c.JSON(http.StatusNotFound, gin.H{
			"success": false,
			"message": "行李寄存不存在",
		})
		return
	}
	// 更新照片

	result := s.DB.Model(&models.Photo{}).Where("file_name = ?", FileName).Update("luggage_storage_id", LuggageStorageID)
	if result.Error != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"success": false,
			"message": "更新照片失败",
		})
		util.Logger.WithFields(logrus.Fields{
			"error": result.Error.Error(),
			"请求id":  c.GetUint("request_id"),
		}).Error("更新照片失败")
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"success": true,
		"message": "更新照片成功",
	})

}

func GetAllPhoto(c *gin.Context, s *services.Services) {
	var photos []models.Photo
	result := s.DB.Model(&models.Photo{}).Find(&photos)
	if result.Error != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"success": false,
			"message": "查询照片失败",
		})
		util.Logger.WithFields(logrus.Fields{
			"error": result.Error.Error(),
			"请求id":  c.GetUint("request_id"),
		}).Error("查询照片失败")
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"success": true,
		"message": "查询照片成功",
		"data":    photos,
	})
}

func generateRandomString(length int) string {
	const charset = "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789"
	b := make([]byte, length)
	for i := range b {
		b[i] = charset[rand.Intn(len(charset))] // ✅ 随机选择
	}
	return string(b)
}
