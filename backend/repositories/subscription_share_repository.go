package repositories

import (
	"fmt"

	"gorm.io/gorm"

	"github.com/subkeep/backend/models"
)

// SubscriptionShareRepository defines the interface for subscription share data access.
type SubscriptionShareRepository interface {
	FindByID(id string) (*models.SubscriptionShare, error)
	FindBySubscriptionID(subscriptionID string) (*models.SubscriptionShare, error)
	FindByUserID(userID string) ([]*models.SubscriptionShare, error)
	Create(share *models.SubscriptionShare) error
	Update(share *models.SubscriptionShare) error
	Delete(id string) error
	DeleteBySubscriptionID(subscriptionID string) error
}

// subscriptionShareRepository is the GORM implementation of SubscriptionShareRepository.
type subscriptionShareRepository struct {
	db *gorm.DB
}

// NewSubscriptionShareRepository creates a new GORM-backed SubscriptionShareRepository.
func NewSubscriptionShareRepository(db *gorm.DB) SubscriptionShareRepository {
	return &subscriptionShareRepository{db: db}
}

// FindByID retrieves a subscription share by its UUID, preloading associations.
func (r *subscriptionShareRepository) FindByID(id string) (*models.SubscriptionShare, error) {
	var share models.SubscriptionShare
	if err := r.db.Preload("Subscription").Preload("ShareGroup.Members").
		Where("id = ?", id).First(&share).Error; err != nil {
		return nil, fmt.Errorf("find subscription share by id: %w", err)
	}
	return &share, nil
}

// FindBySubscriptionID retrieves a subscription share by subscription UUID.
func (r *subscriptionShareRepository) FindBySubscriptionID(subscriptionID string) (*models.SubscriptionShare, error) {
	var share models.SubscriptionShare
	if err := r.db.Preload("Subscription").Preload("ShareGroup.Members").
		Where("subscription_id = ?", subscriptionID).First(&share).Error; err != nil {
		return nil, fmt.Errorf("find subscription share by subscription id: %w", err)
	}
	return &share, nil
}

// FindByUserID retrieves all subscription shares for subscriptions owned by the given user.
func (r *subscriptionShareRepository) FindByUserID(userID string) ([]*models.SubscriptionShare, error) {
	var shares []*models.SubscriptionShare
	if err := r.db.Preload("Subscription").Preload("ShareGroup.Members").
		Joins("JOIN subscriptions ON subscriptions.id = subscription_shares.subscription_id").
		Where("subscriptions.user_id = ?", userID).
		Order("subscription_shares.created_at DESC").
		Find(&shares).Error; err != nil {
		return nil, fmt.Errorf("find subscription shares by user id: %w", err)
	}
	return shares, nil
}

// Create inserts a new subscription share into the database.
func (r *subscriptionShareRepository) Create(share *models.SubscriptionShare) error {
	if err := r.db.Create(share).Error; err != nil {
		return fmt.Errorf("create subscription share: %w", err)
	}
	return nil
}

// Update saves changes to an existing subscription share.
func (r *subscriptionShareRepository) Update(share *models.SubscriptionShare) error {
	if err := r.db.Save(share).Error; err != nil {
		return fmt.Errorf("update subscription share: %w", err)
	}
	return nil
}

// Delete removes a subscription share by its UUID.
func (r *subscriptionShareRepository) Delete(id string) error {
	if err := r.db.Where("id = ?", id).Delete(&models.SubscriptionShare{}).Error; err != nil {
		return fmt.Errorf("delete subscription share: %w", err)
	}
	return nil
}

// DeleteBySubscriptionID removes all subscription shares for a given subscription.
func (r *subscriptionShareRepository) DeleteBySubscriptionID(subscriptionID string) error {
	if err := r.db.Where("subscription_id = ?", subscriptionID).
		Delete(&models.SubscriptionShare{}).Error; err != nil {
		return fmt.Errorf("delete subscription shares by subscription id: %w", err)
	}
	return nil
}
