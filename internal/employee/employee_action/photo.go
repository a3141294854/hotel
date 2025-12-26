package employee_action

import (
	"fmt"
	"hotel/internal/util"
	"hotel/models"
	"hotel/services"
	"net/http"
	"os"
	"path/filepath"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/sirupsen/logrus"
)

// UploadPhoto 上传照片
func UploadPhoto(c *gin.Context, s *services.Services) {
	// 获取上传的文件
	file, err := c.FormFile("photo")
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"success": false,
			"message": "请上传照片",
		})
		util.Logger.WithFields(logrus.Fields{
			"error": err.Error(),
		}).Error("获取上传文件失败")
		return
	}

	// 检查文件大小（限制 5MB）
	const maxSize = 5 * 1024 * 1024
	if file.Size > maxSize {
		c.JSON(http.StatusBadRequest, gin.H{
			"success": false,
			"message": "照片大小不能超过 5MB",
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
			"message": "只支持 jpg、jpeg、png、gif 格式",
		})
		return
	}

	// 创建保存目录
	uploadDir := "./uploads/photos"
	if err := os.MkdirAll(uploadDir, 0755); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"success": false,
			"message": "创建上传目录失败",
		})
		util.Logger.WithFields(logrus.Fields{
			"error": err.Error(),
		}).Error("创建上传目录失败")
		return
	}

	// 生成唯一文件名
	timestamp := time.Now().Format("20060102150405")
	filename := fmt.Sprintf("%s_%s%s", timestamp, generateRandomString(6), ext)
	path := filepath.Join(uploadDir, filename)

	// 保存文件
	if err := c.SaveUploadedFile(file, path); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"success": false,
			"message": "保存照片失败",
		})
		util.Logger.WithFields(logrus.Fields{
			"error": err.Error(),
		}).Error("保存照片失败")
		return
	}

	// 返回文件访问路径
	photoURL := fmt.Sprintf("/photos/%s", filename)
	s.DB.Model(&models.Photo{}).Create(&models.Photo{
		FileName: filename,
		Url:      photoURL,
	})

	c.JSON(http.StatusOK, gin.H{
		"success": true,
		"message": "上传成功",
		"data": gin.H{
			"url":      photoURL,
			"filename": filename,
			"size":     file.Size,
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

func PhotoTouchLuggageStorage(c *gin.Context, s *services.Services) {
	var req struct {
		LuggageStorageID uint     `json:"luggage_storage_id"`
		FileName         []string `json:"file_name"`
	}

	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"success": false,
			"message": "请求格式不正确",
		})
		return
	}

	if len(req.FileName) == 0 || req.LuggageStorageID == 0 {
		c.JSON(http.StatusBadRequest, gin.H{
			"success": false,
			"message": "请提供照片名称和行李寄存ID",
		})
		return
	}
	ok, err := util.ExIf(s.DB, "id", &models.LuggageStorage{}, fmt.Sprintf("%d", req.LuggageStorageID))
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

	for _, fileName := range req.FileName {
		result := s.DB.Model(&models.Photo{}).Where("file_name = ?", fileName).Update("luggage_storage_id", req.LuggageStorageID)
		if result.Error != nil {
			c.JSON(http.StatusInternalServerError, gin.H{
				"success": false,
				"message": "更新照片失败",
			})
			util.Logger.WithFields(logrus.Fields{
				"error": result.Error.Error(),
			}).Error("更新照片失败")
			return
		}
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
		b[i] = charset[i%len(charset)]
	}
	return string(b)
}
