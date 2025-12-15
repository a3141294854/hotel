package models

import (
	"gorm.io/gorm"
	"time"
)

type Employee struct {
	ID             uint      `json:"id" gorm:"primaryKey;autoIncrement"`
	Name           string    `json:"name"`
	User           string    `json:"user"`
	Password       string    `json:"password"`
	Phone          string    `json:"phone"`
	Status         string    `json:"status"`
	LastActiveTime time.Time `json:"last_active_time"`

	HotelID uint   `json:"hotel_id" gorm:"index"`
	Hotel   *Hotel `json:"hotel,omitempty" gorm:"foreignKey:HotelID"`

	RoleID uint  `json:"role_id" gorm:"index"`
	Role   *Role `json:"role,omitempty" gorm:"foreignKey:RoleID"`

	CreatedAt time.Time
	UpdatedAt time.Time
	DeletedAt gorm.DeletedAt `gorm:"index"`
}
type Role struct {
	ID     uint   `json:"id" gorm:"primaryKey;autoIncrement"`
	Name   string `json:"name"`
	Status string `json:"status"`

	Permissions []Permission `json:"permissions,omitempty" gorm:"many2many:role_permission"`

	CreatedAt time.Time
	UpdatedAt time.Time
	DeletedAt gorm.DeletedAt `gorm:"index"`
}

type Permission struct {
	ID     uint   `json:"id" gorm:"primaryKey;autoIncrement"`
	Name   string `json:"name"`
	Status string `json:"status"`

	Roles []Role `json:"roles,omitempty" gorm:"many2many:role_permission"`

	CreatedAt time.Time
	UpdatedAt time.Time
	DeletedAt gorm.DeletedAt `gorm:"index"`
}

type Luggage struct {
	ID uint `json:"id" gorm:"primaryKey;autoIncrement"`

	/*TagID uint `json:"tag_id"`
	Tag   *Tag `json:"tag,omitempty" gorm:"foreignKey:TagID"`*/

	BagCount      int `json:"bag_count"`
	BackpackCount int `json:"backpack_count"`
	BoxCount      int `json:"box_count"`
	OtherCount    int `json:"other_count"`

	OperatorID   uint   `json:"operator_id" gorm:"index"`
	OperatorName string `json:"operator_name"`

	LocationID uint      `json:"location_id" gorm:"index"`
	Location   *Location `json:"location,omitempty" gorm:"foreignKey:LocationID"`

	HotelID uint   `json:"hotel_id" gorm:"index"`
	Hotel   *Hotel `json:"hotel,omitempty" gorm:"foreignKey:HotelID"`

	GuestID uint   `json:"guest_id" gorm:"index"`
	Guest   *Guest `json:"guest,omitempty" gorm:"foreignKey:GuestID"`

	//Photos []Photo `json:"photos,omitempty" gorm:"foreignKey:LuggageStorageID"`

	PickUpCode string `json:"pick_up_code"`
	Status     string `json:"status" gorm:"index"`
	Remark     string `json:"remark"`

	CreatedAt time.Time
	UpdatedAt time.Time
	DeletedAt gorm.DeletedAt `gorm:"index"`
}

/*type LuggageStorage struct {
	ID uint `json:"id" gorm:"primaryKey;autoIncrement"`

	HotelID uint   `json:"hotel_id"`
	Hotel   *Hotel `json:"hotel,omitempty" gorm:"foreignKey:HotelID"`

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

	Photos []Photo `json:"photos,omitempty" gorm:"foreignKey:LuggageStorageID"`

	//TagID uint `json:"tag_id"`
	//Tag   *Tag `json:"tag,omitempty" gorm:"foreignKey:TagID"`

	CreatedAt time.Time
	UpdatedAt time.Time
	DeletedAt gorm.DeletedAt `gorm:"index"`
}*/

/*type Tag struct {
	ID     uint   `json:"id" gorm:"primaryKey;autoIncrement"`
	Name   string `json:"name"`
	Status string `json:"status"`
}*/

/*type Photo struct {
	ID  uint   `json:"id" gorm:"primaryKey;autoIncrement"`
	Url string `json:"url"`

	LuggageID uint     `json:"luggage_id"`
	Luggage   *Luggage `json:"luggage,omitempty" gorm:"foreignKey:LuggageID"`

	CreatedAt time.Time
	UpdatedAt time.Time
	DeletedAt gorm.DeletedAt `gorm:"index"`
}*/

type Location struct {
	ID     uint   `json:"id" gorm:"primaryKey;autoIncrement"`
	Name   string `json:"name"`
	Status string `json:"status"`

	HotelID uint   `json:"hotel_id" gorm:"index"`
	Hotel   *Hotel `json:"hotel,omitempty" gorm:"foreignKey:HotelID"`

	Luggage []Luggage `json:"luggage,omitempty" gorm:"foreignKey:LocationID"`

	CreatedAt time.Time
	UpdatedAt time.Time
	DeletedAt gorm.DeletedAt `gorm:"index"`
}

type Guest struct {
	ID    uint   `json:"id" gorm:"primaryKey;autoIncrement"`
	Name  string `json:"name"`
	Phone string `json:"phone"`
	Room  string `json:"room"`

	Luggage []Luggage `json:"luggage,omitempty" gorm:"foreignKey:GuestID"`

	CreatedAt time.Time      ``
	UpdatedAt time.Time      ``
	DeletedAt gorm.DeletedAt `gorm:"index" `
}

type Hotel struct {
	ID     uint   `json:"id" gorm:"primaryKey;autoIncrement"`
	Name   string `json:"name"`
	Status string `json:"status"`
	Place  string `json:"place"`
	Remark string `json:"remark"`

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
