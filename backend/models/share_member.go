package models

import (
	"time"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

// ShareMember represents a member within a ShareGroup.
type ShareMember struct {
	ID           uuid.UUID `gorm:"type:uuid;primaryKey;default:gen_random_uuid()" json:"id"`
	ShareGroupID uuid.UUID `gorm:"type:uuid;not null;index" json:"shareGroupId" validate:"required"`
	Nickname     string    `gorm:"type:varchar(50);not null" json:"nickname" validate:"required,min=1,max=50"`
	IsOwner      bool      `gorm:"not null;default:false" json:"isOwner"`
	CreatedAt    time.Time `gorm:"not null" json:"createdAt"`

	// Associations
	ShareGroup ShareGroup `gorm:"foreignKey:ShareGroupID;constraint:OnDelete:CASCADE" json:"shareGroup,omitempty"`
}

// TableName overrides the default table name.
func (ShareMember) TableName() string {
	return "share_members"
}

// BeforeCreate sets a new UUID before inserting.
func (sm *ShareMember) BeforeCreate(tx *gorm.DB) error {
	if sm.ID == uuid.Nil {
		sm.ID = uuid.New()
	}
	return nil
}
