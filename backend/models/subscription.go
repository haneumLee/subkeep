package models

import (
	"math"
	"time"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

// Subscription represents a user's subscription to a service.
type Subscription struct {
	ID              uuid.UUID          `gorm:"type:uuid;primaryKey;default:gen_random_uuid()" json:"id"`
	UserID          uuid.UUID          `gorm:"type:uuid;not null;index" json:"userId" validate:"required"`
	ServiceName     string             `gorm:"type:varchar(100);not null" json:"serviceName" validate:"required,min=1,max=100"`
	CategoryID      *uuid.UUID         `gorm:"type:uuid;index" json:"categoryId" validate:"omitempty"`
	Amount          int                `gorm:"type:int;not null" json:"amount" validate:"required,gte=0"`
	BillingCycle    BillingCycle       `gorm:"type:varchar(20);not null" json:"billingCycle" validate:"required,oneof=weekly monthly yearly"`
	Currency        string             `gorm:"type:varchar(3);not null;default:'KRW'" json:"currency" validate:"required,len=3"`
	NextBillingDate time.Time          `gorm:"type:date;not null" json:"nextBillingDate" validate:"required"`
	AutoRenew       bool               `gorm:"not null;default:true" json:"autoRenew"`
	Status          SubscriptionStatus `gorm:"type:varchar(20);not null;default:'active'" json:"status" validate:"required,oneof=active paused cancelled"`
	SatisfactionScore *int             `gorm:"type:int" json:"satisfactionScore" validate:"omitempty,min=1,max=5"`
	Note            *string            `gorm:"type:text" json:"note" validate:"omitempty,max=500"`
	ServiceURL      *string            `gorm:"type:varchar(255)" json:"serviceUrl" validate:"omitempty,url,max=255"`
	StartDate       time.Time          `gorm:"type:date;not null" json:"startDate" validate:"required"`
	CreatedAt       time.Time          `gorm:"not null" json:"createdAt"`
	UpdatedAt       time.Time          `gorm:"not null" json:"updatedAt"`
	DeletedAt       gorm.DeletedAt     `gorm:"index" json:"deletedAt,omitempty"`

	// Associations
	User     User      `gorm:"foreignKey:UserID;constraint:OnDelete:CASCADE" json:"user,omitempty"`
	Category *Category `gorm:"foreignKey:CategoryID;constraint:OnDelete:SET NULL" json:"category,omitempty"`
}

// TableName overrides the default table name.
func (Subscription) TableName() string {
	return "subscriptions"
}

// BeforeCreate sets a new UUID before inserting.
func (s *Subscription) BeforeCreate(tx *gorm.DB) error {
	if s.ID == uuid.Nil {
		s.ID = uuid.New()
	}
	return nil
}

// MonthlyAmount returns the monthly-equivalent cost of this subscription.
//   - monthly: as-is
//   - yearly:  amount / 12 (rounded)
//   - weekly:  amount * 52 / 12 (rounded)
func (s *Subscription) MonthlyAmount() int {
	switch s.BillingCycle {
	case BillingCycleMonthly:
		return s.Amount
	case BillingCycleYearly:
		return int(math.Round(float64(s.Amount) / 12.0))
	case BillingCycleWeekly:
		return int(math.Round(float64(s.Amount) * 52.0 / 12.0))
	default:
		return s.Amount
	}
}

// AnnualAmount returns the annual-equivalent cost of this subscription.
func (s *Subscription) AnnualAmount() int {
	return s.MonthlyAmount() * 12
}
