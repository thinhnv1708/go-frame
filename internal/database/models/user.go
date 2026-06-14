package models

import "time"

type UserModel struct {
	ID        string     `gorm:"column:id;type:varchar(36);primaryKey" json:"id"`
	Name      string     `gorm:"column:name;type:varchar(255);not null" json:"name"`
	Username  string     `gorm:"column:username;type:varchar(50);uniqueIndex;not null" json:"username"`
	Password  string     `gorm:"column:password;type:varchar(255);not null" json:"-"`
	Dob       time.Time  `gorm:"column:dob;type:date;not null" json:"dob"`
	CreatedAt time.Time  `gorm:"column:created_at" json:"created_at"`
	UpdatedAt time.Time  `gorm:"column:updated_at" json:"updated_at"`
	DeletedAt *time.Time `gorm:"column:deleted_at;index" json:"-"`
}

func (UserModel) TableName() string {
	return "users"
}
