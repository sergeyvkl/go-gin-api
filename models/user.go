package models

import (
	"time"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

type User struct {
	ID        string         `gorm:"primaryKey;type:varchar(36);not null" json:"id"`
	Name      string         `gorm:"type:varchar(100);not null;index" json:"name" binding:"required,min=2,max=100"`
	Email     string         `gorm:"type:varchar(100);uniqueIndex;not null" json:"email" binding:"required,email"`
	Age       int            `gorm:"type:int;not null" json:"age" binding:"required,min=1,max=150"`
	CreatedAt time.Time      `gorm:"autoCreateTime" json:"created_at"`
	UpdatedAt time.Time      `gorm:"autoUpdateTime" json:"updated_at"`
	DeletedAt gorm.DeletedAt `gorm:"index" json:"-"` // Soft delete
}

// TableName задает имя таблицы в базе данных
func (User) TableName() string {
	return "users"
}

// BeforeCreate - хук GORM, выполняемый перед созданием записи
func (u *User) BeforeCreate(tx *gorm.DB) error {
	if u.ID == "" {
		u.ID = generateUUID()
	}
	return nil
}

// generateUUID генерирует UUID v4
func generateUUID() string {
	// Используем стандартный uuid пакет
	// Для простоты здесь используется заглушка, в реальном коде импортируйте github.com/google/uuid
	return uuid.New().String() // В реальном коде: uuid.New().String()
}