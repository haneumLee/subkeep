package models

import (
	"log"

	"gorm.io/gorm"
)

// BillingCycle represents the billing frequency of a subscription.
type BillingCycle string

const (
	BillingCycleWeekly  BillingCycle = "weekly"
	BillingCycleMonthly BillingCycle = "monthly"
	BillingCycleYearly  BillingCycle = "yearly"
)

// SubscriptionStatus represents the current state of a subscription.
type SubscriptionStatus string

const (
	SubscriptionStatusActive    SubscriptionStatus = "active"
	SubscriptionStatusPaused    SubscriptionStatus = "paused"
	SubscriptionStatusCancelled SubscriptionStatus = "cancelled"
)

// SplitType represents how a shared subscription cost is divided.
type SplitType string

const (
	SplitTypeEqual        SplitType = "equal"
	SplitTypeCustomAmount SplitType = "custom_amount"
	SplitTypeCustomRatio  SplitType = "custom_ratio"
)

// AuthProvider represents an OAuth authentication provider.
type AuthProvider string

const (
	AuthProviderGoogle AuthProvider = "google"
	AuthProviderApple  AuthProvider = "apple"
	AuthProviderNaver  AuthProvider = "naver"
	AuthProviderKakao  AuthProvider = "kakao"
)

// AutoMigrateAll runs GORM AutoMigrate for all models in the correct order
// respecting foreign key dependencies.
func AutoMigrateAll(db *gorm.DB) error {
	return db.AutoMigrate(
		&User{},
		&Category{},
		&Subscription{},
		&ShareGroup{},
		&ShareMember{},
		&SubscriptionShare{},
	)
}

// SeedDefaultCategories inserts the predefined system categories if they don't
// already exist.
func SeedDefaultCategories(db *gorm.DB) error {
	type defaultCat struct {
		Name  string
		Color string
		Icon  string
		Order int
	}

	defaults := []defaultCat{
		{Name: "엔터테인먼트", Color: "#FF6B6B", Icon: "entertainment", Order: 1},
		{Name: "생산성", Color: "#4ECDC4", Icon: "productivity", Order: 2},
		{Name: "클라우드/스토리지", Color: "#45B7D1", Icon: "cloud", Order: 3},
		{Name: "AI 서비스", Color: "#96CEB4", Icon: "ai", Order: 4},
		{Name: "쇼핑", Color: "#FFEAA7", Icon: "shopping", Order: 5},
		{Name: "기타", Color: "#DFE6E9", Icon: "etc", Order: 6},
	}

	for _, d := range defaults {
		cat := Category{
			Name:     d.Name,
			Color:    &d.Color,
			Icon:     &d.Icon,
			SortOrder: d.Order,
			IsSystem: true,
		}

		result := db.Where("name = ? AND is_system = ?", d.Name, true).FirstOrCreate(&cat)
		if result.Error != nil {
			log.Printf("failed to seed category %s: %v", d.Name, result.Error)
			return result.Error
		}
	}

	return nil
}
