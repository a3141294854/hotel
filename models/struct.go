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
	ID   uint   `json:"id" gorm:"primaryKey;autoIncrement"`
	Name string `json:"name"`

	Permissions []Permission `json:"permissions,omitempty" gorm:"many2many:role_permission"`

	CreatedAt time.Time
	UpdatedAt time.Time
	DeletedAt gorm.DeletedAt `gorm:"index"`
}

type Permission struct {
	ID   uint   `json:"id" gorm:"primaryKey;autoIncrement"`
	Name string `json:"name"`

	Roles []Role `json:"roles,omitempty" gorm:"many2many:role_permission"`

	CreatedAt time.Time
	UpdatedAt time.Time
	DeletedAt gorm.DeletedAt `gorm:"index"`
}

type LuggageStorage struct {
	ID uint `json:"id" gorm:"primaryKey;autoIncrement"`

	/*TagID uint `json:"tag_id"`
	Tag   *Tag `json:"tag,omitempty" gorm:"foreignKey:TagID"`*/

	BagCount      int `json:"bag_count"`
	BackpackCount int `json:"backpack_count"`
	BoxCount      int `json:"box_count"`
	OtherCount    int `json:"other_count"`

	OperatorID   uint   `json:"operator_id" gorm:"index"`
	OperatorName string `json:"operator_name"`

	HotelID uint   `json:"hotel_id" gorm:"index"`
	Hotel   *Hotel `json:"hotel,omitempty" gorm:"foreignKey:HotelID"`

	GuestID uint   `json:"guest_id" gorm:"index"`
	Guest   *Guest `json:"guest,omitempty" gorm:"foreignKey:GuestID;constraint:OnDelete:CASCADE"`

	Luggage []Luggage `json:"luggage,omitempty" gorm:"foreignKey:LuggageStorageID;constraint:OnDelete:CASCADE"`

	//Photos []Photo `json:"photos,omitempty" gorm:"foreignKey:LuggageStorageID"`

	PickUpCode string `json:"pick_up_code"`
	Status     string `json:"status" gorm:"index"`
	Remark     string `json:"remark"`

	CreatedAt time.Time
	UpdatedAt time.Time
	DeletedAt gorm.DeletedAt `gorm:"index"`
}

type Luggage struct {
	ID     uint   `json:"id" gorm:"primaryKey;autoIncrement"`
	Status string `json:"status"`

	LuggageStorageID uint            `json:"luggage_storage_id"`
	LuggageStorage   *LuggageStorage `json:"luggage_storage,omitempty" gorm:"foreignKey:LuggageStorageID"`

	LocationID uint      `json:"location_id" gorm:"index"`
	Location   *Location `json:"location,omitempty" gorm:"foreignKey:LocationID"`

	TagID uint `json:"tag_id" gorm:"index"`
	Tag   *Tag `json:"tag,omitempty" gorm:"foreignKey:TagID"`

	CreatedAt time.Time
	UpdatedAt time.Time
	DeletedAt gorm.DeletedAt `gorm:"index"`
}

type Tag struct {
	ID   uint   `json:"id" gorm:"primaryKey;autoIncrement"`
	Name string `json:"name"`
	Mac  string `json:"mac"`

	CreatedAt time.Time
	UpdatedAt time.Time
	DeletedAt gorm.DeletedAt `gorm:"index"`
}

/*type Photo struct {
	ID  uint   `json:"id" gorm:"primaryKey;autoIncrement"`
	Url string `json:"url"`

	LuggageID uint     `json:"luggage_id"`
	LuggageStorage   *LuggageStorage `json:"luggage,omitempty" gorm:"foreignKey:LuggageID"`

	CreatedAt time.Time
	UpdatedAt time.Time
	DeletedAt gorm.DeletedAt `gorm:"index"`
}*/

type Location struct {
	ID   uint   `json:"id" gorm:"primaryKey;autoIncrement"`
	Name string `json:"name"`

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

	LuggageStorage []LuggageStorage `json:"luggage,omitempty" gorm:"foreignKey:GuestID"`

	CreatedAt time.Time      ``
	UpdatedAt time.Time      ``
	DeletedAt gorm.DeletedAt `gorm:"index" `
}

type Hotel struct {
	ID     uint   `json:"id" gorm:"primaryKey;autoIncrement"`
	Name   string `json:"name"`
	Place  string `json:"place"`
	Remark string `json:"remark"`

	CreatedAt time.Time
	UpdatedAt time.Time
	DeletedAt gorm.DeletedAt `gorm:"index"`
}
