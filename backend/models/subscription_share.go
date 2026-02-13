package models

import (
	"math"
	"time"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

// SubscriptionShare links a Subscription to a ShareGroup and defines
// how the cost is split among members.
type SubscriptionShare struct {
	ID                   uuid.UUID  `gorm:"type:uuid;primaryKey;default:gen_random_uuid()" json:"id"`
	SubscriptionID       uuid.UUID  `gorm:"type:uuid;not null;index" json:"subscriptionId" validate:"required"`
	ShareGroupID         uuid.UUID  `gorm:"type:uuid;not null;index" json:"shareGroupId" validate:"required"`
	SplitType            SplitType  `gorm:"type:varchar(20);not null" json:"splitType" validate:"required,oneof=equal custom_amount custom_ratio"`
	MyShareAmount        *int       `gorm:"type:int" json:"myShareAmount" validate:"omitempty,gte=0"`
	MyShareRatio         *float64   `gorm:"type:numeric(5,4)" json:"myShareRatio" validate:"omitempty,gte=0,lte=1"`
	TotalMembersSnapshot int        `gorm:"type:int;not null" json:"totalMembersSnapshot" validate:"required,gte=1"`
	CreatedAt            time.Time  `gorm:"not null" json:"createdAt"`
	UpdatedAt            time.Time  `gorm:"not null" json:"updatedAt"`

	// Associations
	Subscription Subscription `gorm:"foreignKey:SubscriptionID;constraint:OnDelete:CASCADE" json:"subscription,omitempty"`
	ShareGroup   ShareGroup   `gorm:"foreignKey:ShareGroupID;constraint:OnDelete:CASCADE" json:"shareGroup,omitempty"`
}

// TableName overrides the default table name.
func (SubscriptionShare) TableName() string {
	return "subscription_shares"
}

// BeforeCreate sets a new UUID before inserting.
func (ss *SubscriptionShare) BeforeCreate(tx *gorm.DB) error {
	if ss.ID == uuid.Nil {
		ss.ID = uuid.New()
	}
	return nil
}

// PersonalAmount calculates how much the current user pays based on the
// split configuration.
//   - equal:         monthlyAmount / totalMembersSnapshot (rounded)
//   - custom_amount: myShareAmount as-is
//   - custom_ratio:  monthlyAmount * myShareRatio (rounded)
func (ss *SubscriptionShare) PersonalAmount(monthlyAmount int) int {
	switch ss.SplitType {
	case SplitTypeEqual:
		if ss.TotalMembersSnapshot == 0 {
			return monthlyAmount
		}
		return int(math.Round(float64(monthlyAmount) / float64(ss.TotalMembersSnapshot)))
	case SplitTypeCustomAmount:
		if ss.MyShareAmount != nil {
			return *ss.MyShareAmount
		}
		return 0
	case SplitTypeCustomRatio:
		if ss.MyShareRatio != nil {
			return int(math.Round(float64(monthlyAmount) * *ss.MyShareRatio))
		}
		return 0
	default:
		return monthlyAmount
	}
}
