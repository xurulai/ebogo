package model

import (
	"time"
)

// User 用户模型
type User struct {
	ID        int64     `gorm:"primaryKey" json:"id"`                 // 用户ID，主键
	Username  string    `gorm:"uniqueIndex;not null" json:"username"` // 用户名，唯一索引，非空
	Password  string    `gorm:"not null" json:"password"`             // 密码，非空
	Email     string    `gorm:"uniqueIndex;not null" json:"email"`    // 邮箱，唯一索引，非空
	Phone     string    `gorm:"uniqueIndex" json:"phone"`             // 手机号，唯一索引
	CreatedAt time.Time `json:"created_at"`                           // 创建时间
	UpdatedAt time.Time `json:"updated_at"`                           // 更新时间
}