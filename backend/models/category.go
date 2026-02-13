package models

import (
	"time"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

// Category represents a grouping label for subscriptions.
// System categories (IsSystem=true, UserID=nil) are shared across all users.
type Category struct {
	ID        uuid.UUID  `gorm:"type:uuid;primaryKey;default:gen_random_uuid()" json:"id"`
	UserID    *uuid.UUID `gorm:"type:uuid;index" json:"userId"`
	Name      string     `gorm:"type:varchar(50);not null" json:"name" validate:"required,min=1,max=50"`
	Color     *string    `gorm:"type:varchar(7)" json:"color" validate:"omitempty,hexcolor"`
	Icon      *string    `gorm:"type:varchar(50)" json:"icon" validate:"omitempty,max=50"`
	SortOrder int        `gorm:"type:int;not null;default:0" json:"sortOrder"`
	IsSystem  bool       `gorm:"not null;default:false" json:"isSystem"`
	CreatedAt time.Time  `gorm:"not null" json:"createdAt"`
	UpdatedAt time.Time  `gorm:"not null" json:"updatedAt"`

	// Associations
	User          *User          `gorm:"foreignKey:UserID;constraint:OnDelete:CASCADE" json:"user,omitempty"`
	Subscriptions []Subscription `gorm:"foreignKey:CategoryID" json:"subscriptions,omitempty"`
}

// TableName overrides the default table name.
func (Category) TableName() string {
	return "categories"
}

// BeforeCreate sets a new UUID before inserting.
func (c *Category) BeforeCreate(tx *gorm.DB) error {
	if c.ID == uuid.Nil {
		c.ID = uuid.New()
	}
	return nil
}
