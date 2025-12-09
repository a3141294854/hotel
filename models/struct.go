package models

import (
	"gorm.io/gorm"
	"time"
)

type Employee struct {
	ID       uint   `json:"id" gorm:"primaryKey;autoIncrement"`
	Name     string `json:"name" gorm:"column:name"`
	User     string `json:"user" gorm:"column:user"`
	Password string `json:"password" gorm:"column:password"`

	RoleID uint `json:"role_id" gorm:"column:role_id"`
	Role   Role `json:"role" gorm:"foreignKey:RoleID"`
}
type Role struct {
	ID   uint   `json:"id" gorm:"primaryKey;autoIncrement"`
	Name string `json:"name" gorm:"column:name"`

	Permissions []Permission `json:"permissions" gorm:"many2many:role_permission"`
}

type Permission struct {
	ID   uint   `json:"id" gorm:"primaryKey;autoIncrement"`
	Name string `json:"name" gorm:"column:name"`
}

type Luggage struct {
	ID       uint    `json:"id" gorm:"primaryKey;autoIncrement"`
	GuestID  uint    `json:"guest_id" gorm:"column:guest_id"`
	Tag      string  `json:"tag" gorm:"column:tag"`
	Weight   float32 `json:"weight" gorm:"column:weight"`
	Status   string  `json:"status" gorm:"column:status"`
	Location string  `json:"location" gorm:"column:location"`

	CreatedAt time.Time      `json:"-"`
	UpdatedAt time.Time      `json:"-"`
	DeletedAt gorm.DeletedAt `gorm:"index" json:"-"`

	Guest Guest `json:"guest" gorm:"foreignKey:GuestID"`
}
type Guest struct {
	ID   uint   `json:"id" gorm:"primaryKey;autoIncrement;type:int unsigned"`
	Name string `json:"name" gorm:"column:guest_name"`
}

type TokenBucketLimiter struct {
	Capacity     int           // 桶容量
	FillRate     time.Duration // 添加令牌速率，如每10ms加1个令牌
	Tokens       int           // 当前令牌数
	LastFillTime time.Time     // 上次添加令牌时间
}
