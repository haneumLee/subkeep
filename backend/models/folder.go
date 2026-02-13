package models

import (
	"time"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

// Folder represents a user-defined organizational folder for subscriptions.
type Folder struct {
	ID        uuid.UUID      `gorm:"type:uuid;primaryKey;default:gen_random_uuid()" json:"id"`
	UserID    uuid.UUID      `gorm:"type:uuid;not null;index" json:"userId" validate:"required"`
	Name      string         `gorm:"type:varchar(50);not null" json:"name" validate:"required,min=1,max=50"`
	SortOrder int            `gorm:"type:int;not null;default:0" json:"sortOrder"`
	CreatedAt time.Time      `gorm:"not null" json:"createdAt"`
	UpdatedAt time.Time      `gorm:"not null" json:"updatedAt"`
	DeletedAt gorm.DeletedAt `gorm:"index" json:"deletedAt,omitempty"`

	// Associations
	User          User           `gorm:"foreignKey:UserID;constraint:OnDelete:CASCADE" json:"user,omitempty"`
	Subscriptions []Subscription `gorm:"foreignKey:FolderID" json:"subscriptions,omitempty"`
}

// TableName overrides the default table name.
func (Folder) TableName() string {
	return "folders"
}

// BeforeCreate sets a new UUID before inserting.
func (f *Folder) BeforeCreate(tx *gorm.DB) error {
	if f.ID == uuid.Nil {
		f.ID = uuid.New()
	}
	return nil
}
