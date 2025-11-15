package models

type Employee struct {
	ID       uint   `gorm:"primaryKey;autoIncrement"`
	User     string `json:"user" gorm:"column:user_name"`
	Password string `json:"password" gorm:"column:password"`
}

type Luggage struct {
	ID        uint    `gorm:"primaryKey;autoIncrement"`
	GuestID   uint    `json:"guest_id" gorm:"column:guest_id"`
	GuestName string  `json:"guest_name" gorm:"column:guest_name"`
	Tag       string  `json:"tag" gorm:"column:tag;size:50"`
	Weight    float32 `json:"weight" gorm:"column:weight"`
	Status    string  `json:"status" gorm:"column:status;size:20"`
	Location  string  `json:"location" gorm:"column:location;size:100"`
}
type Guest struct {
	ID   uint   `gorm:"primaryKey;autoIncrement"`
	Name string `json:"name" gorm:"column:guest_name"`
}
