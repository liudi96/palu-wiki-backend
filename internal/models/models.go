package models

import (
	"time"
)

// Guide represents a game guide entry.
type Guide struct {
	ID            uint      `gorm:"primaryKey" json:"id"`
	Title         string    `gorm:"type:varchar(255);not null" json:"title"`
	Content       string    `gorm:"type:longtext;not null" json:"content"`
	Tags          string    `gorm:"type:varchar(255)" json:"tags"` // Comma-separated tags
	Category      string    `gorm:"type:varchar(100)" json:"category"`
	SourceURL     string    `gorm:"type:varchar(255)" json:"source_url"`
	CreatedAt     time.Time `json:"created_at"`
	UpdatedAt     time.Time `json:"updated_at"`
	IsAIGenerated bool      `json:"is_ai_generated"`
	Version       string    `gorm:"type:varchar(50)" json:"version"`
}

// UserQuery represents a user's query to the AI bot.
type UserQuery struct {
	ID         uint      `gorm:"primaryKey" json:"id"`
	UserID     string    `gorm:"type:varchar(100)" json:"user_id"` // WeChat OpenID or UnionID
	QueryText  string    `gorm:"type:text;not null" json:"query_text"`
	AIResponse string    `gorm:"type:longtext;not null" json:"ai_response"`
	QueryTime  time.Time `json:"query_time"`
}

// OfficialUpdate represents an official game update log.
type OfficialUpdate struct {
	ID          uint      `gorm:"primaryKey" json:"id"`
	Title       string    `gorm:"type:varchar(255);not null" json:"title"`
	Content     string    `gorm:"type:longtext;not null" json:"content"`
	PublishDate time.Time `json:"publish_date"`
	SourceURL   string    `gorm:"type:varchar(255)" json:"source_url"`
	Processed   bool      `gorm:"default:false" json:"processed"`
}
