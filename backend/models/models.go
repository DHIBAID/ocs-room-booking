package models

import (
	"time"
)

type User struct {
	ID           uint      `json:"id" gorm:"primaryKey"`
	Username     string    `json:"username" gorm:"unique;not null"`
	PasswordHash string    `json:"-" gorm:"not null"`
	Role         string    `json:"role" gorm:"type:varchar(20);not null"` // admin, core, viewer
	IsActive     bool      `json:"is_active" gorm:"default:true"`
	CreatedAt    time.Time `json:"created_at"`
	UpdatedAt    time.Time `json:"updated_at"`
}

type Room struct {
	ID              uint   `json:"id" gorm:"primaryKey"`
	Block           string `json:"block" gorm:"not null"`
	Name            string `json:"name" gorm:"unique;not null"`
	Capacity        int    `json:"capacity" gorm:"not null"`
	Status          string `json:"status" gorm:"type:varchar(20);default:'available'"`
	AllowedPurposes string `json:"allowed_purposes" gorm:"type:text"`
	Notes           string `json:"notes"`
}

type Booking struct {
	ID           uint      `json:"id" gorm:"primaryKey"`
	RoomID       uint      `json:"room_id" gorm:"not null"`
	Room         Room      `json:"room" gorm:"foreignKey:RoomID"`
	UserID       uint      `json:"user_id" gorm:"not null"`
	User         User      `json:"user" gorm:"foreignKey:UserID"`
	Purpose      string    `json:"purpose" gorm:"not null"` // OA, Interview, PPT
	Participants int       `json:"participants" gorm:"not null"`
	StartTime    time.Time `json:"start_time" gorm:"not null"`
	EndTime      time.Time `json:"end_time" gorm:"not null"`
	Status       string    `json:"status" gorm:"type:varchar(20);default:'confirmed'"` // confirmed, cancelled
	CreatedAt    time.Time `json:"created_at"`
	UpdatedAt    time.Time `json:"updated_at"`
}
