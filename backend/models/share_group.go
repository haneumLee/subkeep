package models

import (
	"time"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

// ShareGroup represents a group of people sharing subscription costs.
type ShareGroup struct {
	ID          uuid.UUID      `gorm:"type:uuid;primaryKey;default:gen_random_uuid()" json:"id"`
	OwnerUserID uuid.UUID     `gorm:"type:uuid;not null;index" json:"ownerUserId" validate:"required"`
	Name        string         `gorm:"type:varchar(100);not null" json:"name" validate:"required,min=1,max=100"`
	Description *string        `gorm:"type:text" json:"description" validate:"omitempty"`
	CreatedAt   time.Time      `gorm:"not null" json:"createdAt"`
	UpdatedAt   time.Time      `gorm:"not null" json:"updatedAt"`
	DeletedAt   gorm.DeletedAt `gorm:"index" json:"deletedAt,omitempty"`

	// Associations
	Owner   User          `gorm:"foreignKey:OwnerUserID;constraint:OnDelete:CASCADE" json:"owner,omitempty"`
	Members []ShareMember `gorm:"foreignKey:ShareGroupID;constraint:OnDelete:CASCADE" json:"members,omitempty"`
}

// TableName overrides the default table name.
func (ShareGroup) TableName() string {
	return "share_groups"
}

// BeforeCreate sets a new UUID before inserting.
func (sg *ShareGroup) BeforeCreate(tx *gorm.DB) error {
	if sg.ID == uuid.Nil {
		sg.ID = uuid.New()
	}
	return nil
}
