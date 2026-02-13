package models

import (
	"time"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

// User represents an authenticated user of the application.
type User struct {
	ID             uuid.UUID      `gorm:"type:uuid;primaryKey;default:gen_random_uuid()" json:"id"`
	Provider       AuthProvider   `gorm:"type:varchar(20);not null" json:"provider" validate:"required,oneof=google apple naver kakao"`
	ProviderUserID string         `gorm:"type:varchar(255);not null;uniqueIndex" json:"providerUserId" validate:"required,max=255"`
	Email          *string        `gorm:"type:varchar(255)" json:"email" validate:"omitempty,email,max=255"`
	Nickname       *string        `gorm:"type:varchar(100)" json:"nickname" validate:"omitempty,max=100"`
	AvatarURL      *string        `gorm:"type:text" json:"avatarUrl" validate:"omitempty,url"`
	CreatedAt      time.Time      `gorm:"not null" json:"createdAt"`
	UpdatedAt      time.Time      `gorm:"not null" json:"updatedAt"`
	LastLoginAt    *time.Time     `json:"lastLoginAt"`
	DeletedAt      gorm.DeletedAt `gorm:"index" json:"deletedAt,omitempty"`

	// Associations
	Subscriptions []Subscription `gorm:"foreignKey:UserID;constraint:OnDelete:CASCADE" json:"subscriptions,omitempty"`
	Categories    []Category     `gorm:"foreignKey:UserID;constraint:OnDelete:CASCADE" json:"categories,omitempty"`
	ShareGroups   []ShareGroup   `gorm:"foreignKey:OwnerUserID;constraint:OnDelete:CASCADE" json:"shareGroups,omitempty"`
}

// TableName overrides the default table name.
func (User) TableName() string {
	return "users"
}

// BeforeCreate sets a new UUID before inserting.
func (u *User) BeforeCreate(tx *gorm.DB) error {
	if u.ID == uuid.Nil {
		u.ID = uuid.New()
	}
	return nil
}
