package models

import (
	"time"
)

// Client represents a registered client device
type Client struct {
	UUID       string `gorm:"type:uuid;primaryKey;default:gen_random_uuid()"`
	Token      string `gorm:"type:varchar(255);unique;not null"`
	ClientName string `gorm:"type:varchar(100);not null"`
	CreatedAt  time.Time
	UpdatedAt  time.Time
}

// ClientConfig stores client monitoring preferences
type ClientConfig struct {
	ClientUUID  string `gorm:"type:uuid;primaryKey;foreignKey:ClientUUID;references:UUID;constraint:OnDelete:CASCADE"`
	CPU         bool   `gorm:"default:true"`
	GPU         bool   `gorm:"default:true"`
	RAM         bool   `gorm:"default:true"`
	SWAP        bool   `gorm:"default:true"`
	LOAD        bool   `gorm:"default:true"`
	UPTIME      bool   `gorm:"default:true"`
	TEMP        bool   `gorm:"default:true"`
	OS          bool   `gorm:"default:true"`
	DISK        bool   `gorm:"default:true"`
	NET         bool   `gorm:"default:true"`
	PROCESS     bool   `gorm:"default:true"`
	Connections bool   `gorm:"default:true"`
	Interval    int    `gorm:"default:3;check:Interval >= 1"` // Ensure positive interval
	CreatedAt   time.Time
	UpdatedAt   time.Time
}

// User represents an authenticated user
type User struct {
	UUID      string `gorm:"type:uuid;primaryKey;default:gen_random_uuid()"`
	Username  string `gorm:"type:varchar(50);unique;not null"`
	Passwd    string `gorm:"type:varchar(255);not null"` // Hashed password
	CreatedAt time.Time
	UpdatedAt time.Time
}

// Session manages user sessions
type Session struct {
	UUID      string    `gorm:"type:uuid;primaryKey;default:gen_random_uuid()"`
	UserUUID  string    `gorm:"type:uuid;foreignKey:UserUUID;references:UUID;constraint:OnDelete:CASCADE"`
	Session   string    `gorm:"type:varchar(255);unique;not null"`
	Expires   time.Time `gorm:"not null"`
	CreatedAt time.Time
}

// History logs client metrics over time
type History struct {
	ID             uint64    `gorm:"primaryKey;autoIncrement"`
	ClientUUID     string    `gorm:"type:uuid;index;foreignKey:ClientUUID;references:UUID;constraint:OnDelete:CASCADE"`
	Time           time.Time `gorm:"index;default:CURRENT_TIMESTAMP"`
	CPU            float32   `gorm:"type:decimal(5,2)"` // e.g., 75.50%
	GPU            float32   `gorm:"type:decimal(5,2)"`
	RAM            int64     `gorm:"type:bigint"`
	RAMTotal       int64     `gorm:"type:bigint"`
	SWAP           int64     `gorm:"type:bigint"`
	SWAPTotal      int64     `gorm:"type:bigint"`
	LOAD           float32   `gorm:"type:decimal(5,2)"`
	TEMP           float32   `gorm:"type:decimal(5,2)"`
	DISK           int64     `gorm:"type:bigint"`
	DISKTotal      int64     `gorm:"type:bigint"`
	NETIn          int64     `gorm:"type:bigint"`
	NETOut         int64     `gorm:"type:bigint"`
	NETTotalUp     int64     `gorm:"type:bigint"`
	NETTotalDown   int64     `gorm:"type:bigint"`
	PROCESS        int
	Connections    int
	ConnectionsUDP int
}

// ClientInfo stores static client information
type ClientInfo struct {
	ClientUUID string `gorm:"type:uuid;primaryKey;foreignKey:ClientUUID;references:UUID;constraint:OnDelete:CASCADE"`
	CPUNAME    string `gorm:"type:varchar(100)"`
	CPUARCH    string `gorm:"type:varchar(50)"`
	CPUCORES   int
	OS         string `gorm:"type:varchar(100)"`
	GPUNAME    string `gorm:"type:varchar(100)"`
	CreatedAt  time.Time
	UpdatedAt  time.Time
}

// Config stores site-wide settings
type Config struct {
	ID        uint   `gorm:"primaryKey;autoIncrement"`
	Sitename  string `gorm:"type:varchar(100);not null"`
	CreatedAt time.Time
	UpdatedAt time.Time
}

// Custom stores custom UI configurations
type Custom struct {
	ID        uint   `gorm:"primaryKey;autoIncrement"`
	CustomCSS string `gorm:"type:text"`
	CustomJS  string `gorm:"type:text"`
	CreatedAt time.Time
	UpdatedAt time.Time
}
