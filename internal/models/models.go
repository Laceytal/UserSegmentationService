package models

import "time"

// Таблица segments: список сегментов (MAIL_GPT, CLOUD_DISCOUNT_30 и т.д.)
type Segment struct {
	ID          uint   `gorm:"primaryKey"`
	Key         string `gorm:"unique;not null"` // уникальный ключ сегмента
	Description string // описание
	CreatedAt   time.Time
	UpdatedAt   time.Time
}

// Таблица user_segments: явные назначения сегмента пользователю
type UserSegment struct {
	UserID     uint `gorm:"primaryKey"`
	SegmentID  uint `gorm:"primaryKey"`
	AssignedAt time.Time
}

// Таблица segment_rules: правила автоматического распределения
type SegmentRule struct {
	SegmentID uint   `gorm:"primaryKey"`
	Pct       int    `gorm:"not null;check:pct>=1 AND pct<=100"` // процент от 1 до 100
	Seed      string `gorm:"not null"`                           // seed для детерминированного хеширования
	CreatedAt time.Time
}
