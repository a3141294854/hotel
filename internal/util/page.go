package util

import (
	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
	"strconv"
)

// Paginator 分页器
type Paginator struct {
	Page      int
	PageSize  int
	Total     int64
	TotalPage int
	Offset    int
}

// NewPaginator 从Context创建分页器,默认page=1,page_size=10
func NewPaginator(c *gin.Context) *Paginator {
	page, _ := strconv.Atoi(c.DefaultQuery("page", "1"))
	pageSize, _ := strconv.Atoi(c.DefaultQuery("page_size", "10"))

	if page < 1 {
		page = 1
	}
	if pageSize < 1 {
		pageSize = 10
	}
	if pageSize > 100 {
		pageSize = 100
	}

	return &Paginator{
		Page:     page,
		PageSize: pageSize,
		Offset:   (page - 1) * pageSize, //偏移量，用来跳过前面的数据
	}
}

// Paginate 执行分页查询
func (p *Paginator) Paginate(db *gorm.DB, data interface{}) error {
	// 查询总数
	db.Model(data).Count(&p.Total)

	// 计算总页数
	p.TotalPage = int(p.Total) / p.PageSize
	if int(p.Total)%p.PageSize > 0 {
		p.TotalPage++
	}

	// 执行分页查询
	return db.Offset(p.Offset).Limit(p.PageSize).Find(data).Error
}

// ToMap 转换为Map（用于返回JSON）
func (p *Paginator) ToMap() map[string]interface{} {
	return map[string]interface{}{
		"page":       p.Page,
		"page_size":  p.PageSize,
		"total":      p.Total,
		"total_page": p.TotalPage,
		"has_more":   p.Page < p.TotalPage,
	}
}
