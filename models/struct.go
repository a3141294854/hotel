package models

import (
	"gorm.io/gorm"
	"time"
)

type Employee struct {
	ID       uint   `json:"id" gorm:"primaryKey;autoIncrement"`
	Name     string `json:"name"`
	User     string `json:"user"`
	Password string `json:"password"`
	Phone    string `json:"phone"`
	Status   string `json:"status"`

	HotelID uint   `json:"hotel_id"`
	Hotel   *Hotel `json:"hotel,omitempty" gorm:"foreignKey:HotelId"`

	RoleID uint  `json:"role_id"`
	Role   *Role `json:"role,omitempty" gorm:"foreignKey:RoleID"`

	CreatedAt time.Time
	UpdatedAt time.Time
	DeletedAt gorm.DeletedAt `gorm:"index"`
}
type Role struct {
	ID     uint   `json:"id" gorm:"primaryKey;autoIncrement"`
	Name   string `json:"name"`
	Status string `json:"status"`

	Permissions *[]Permission `json:"permissions,omitempty" gorm:"many2many:role_permission"`

	CreatedAt time.Time
	UpdatedAt time.Time
	DeletedAt gorm.DeletedAt `gorm:"index"`
}

type Permission struct {
	ID     uint   `json:"id" gorm:"primaryKey;autoIncrement"`
	Name   string `json:"name"`
	Status string `json:"status"`

	Roles *[]Role `json:"roles,omitempty" gorm:"many2many:role_permission"`

	CreatedAt time.Time
	UpdatedAt time.Time
	DeletedAt gorm.DeletedAt `gorm:"index"`
}

type Luggage struct {
	ID     uint   `json:"id" gorm:"primaryKey;autoIncrement"`
	Tag    string `json:"tag"`
	Status string `json:"status"`

	LocationID uint      `json:"location_id"`
	Location   *Location `json:"location,omitempty" gorm:"foreignKey:LocationID"`

	HotelID uint   `json:"hotel_id"`
	Hotel   *Hotel `json:"hotel,omitempty" gorm:"foreignKey:HotelId"`

	GuestID uint   `json:"guest_id"`
	Guest   *Guest `json:"guest,omitempty" gorm:"foreignKey:GuestID"`

	LuggageStorageID uint            `json:"luggage_storage_id"`
	LuggageStorage   *LuggageStorage `json:"luggage_storage,omitempty" gorm:"foreignKey:LuggageStorageID"`

	CreatedAt time.Time
	UpdatedAt time.Time
	DeletedAt gorm.DeletedAt `gorm:"index"`
}

type LuggageStorage struct {
	ID      uint   `json:"id" gorm:"primaryKey;autoIncrement"`
	HotelID uint   `json:"hotel_id"`
	Hotel   *Hotel `json:"hotel,omitempty" gorm:"foreignKey:HotelId"`

	OperatorID   uint   `json:"operator_id"`
	OperatorName string `json:"operator_name"`

	BagCount      int `json:"bag_count"`
	BackpackCount int `json:"backpack_count"`
	BoxCount      int `json:"box_count"`

	GuestName  string `json:"guest_name"`
	GuestPhone string `json:"guest_phone"`
	GuestRoom  string `json:"guest_room"`

	PickUpCode string `json:"pick_up_code"`
	Status     string `json:"status"`
	Remark     string `json:"remark"`

	Photos *[]Photo `json:"photos,omitempty" gorm:"foreignKey:LuggageStorageID"`

	CreatedAt time.Time
	UpdatedAt time.Time
	DeletedAt gorm.DeletedAt `gorm:"index"`
}

type Photo struct {
	ID               uint            `json:"id" gorm:"primaryKey;autoIncrement"`
	Url              string          `json:"url"`
	LuggageStorageID uint            `json:"luggage_storage_id"`
	LuggageStorage   *LuggageStorage `json:"luggage_storage,omitempty" gorm:"foreignKey:LuggageStorageID"`

	CreatedAt time.Time
	UpdatedAt time.Time
	DeletedAt gorm.DeletedAt `gorm:"index"`
}

type Location struct {
	ID     uint   `json:"id" gorm:"primaryKey;autoIncrement"`
	Name   string `json:"name"`
	Status string `json:"status"`

	HotelID uint   `json:"hotel_id"`
	Hotel   *Hotel `json:"hotel,omitempty" gorm:"foreignKey:HotelId"`

	Luggage *[]Luggage `json:"luggage,omitempty" gorm:"foreignKey:LocationID"`

	CreatedAt time.Time
	UpdatedAt time.Time
	DeletedAt gorm.DeletedAt `gorm:"index"`
}

type Guest struct {
	ID    uint   `json:"id" gorm:"primaryKey;autoIncrement"`
	Name  string `json:"name"`
	Phone string `json:"phone"`
	Room  string `json:"room"`

	CreatedAt time.Time      ``
	UpdatedAt time.Time      ``
	DeletedAt gorm.DeletedAt `gorm:"index" `
}

type Hotel struct {
	ID     uint   `json:"id" gorm:"primaryKey;autoIncrement"`
	Name   string `json:"name"`
	Status string `json:"status"`
	Place  string `json:"place"`

	CreatedAt time.Time
	UpdatedAt time.Time
	DeletedAt gorm.DeletedAt `gorm:"index"`
}
type TokenBucketLimiter struct {
	Capacity     int           // 桶容量
	FillRate     time.Duration // 添加令牌速率，如每10ms加1个令牌
	Tokens       int           // 当前令牌数
	LastFillTime int64         // 上次添加令牌时间
}
