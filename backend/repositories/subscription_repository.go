package repositories

import (
	"fmt"

	"github.com/subkeep/backend/models"
	"gorm.io/gorm"
)

// SubscriptionFilter holds query parameters for filtering, sorting, and
// paginating subscription lists.
type SubscriptionFilter struct {
	Status     string // "active", "paused", "cancelled", "" (all)
	CategoryID string // filter by category UUID
	SortBy     string // "amount", "satisfaction", "next_billing_date", "created_at"
	SortOrder  string // "asc", "desc"
	Page       int
	PerPage    int
}

// Defaults applies default values for missing filter fields.
func (f *SubscriptionFilter) Defaults() {
	if f.Page <= 0 {
		f.Page = 1
	}
	if f.PerPage <= 0 {
		f.PerPage = 20
	}
	if f.PerPage > 100 {
		f.PerPage = 100
	}
	if f.SortOrder != "asc" && f.SortOrder != "desc" {
		f.SortOrder = "desc"
	}
}

// SubscriptionRepository defines the interface for subscription data access.
type SubscriptionRepository interface {
	FindByID(id string) (*models.Subscription, error)
	FindByUserID(userID string, filter SubscriptionFilter) ([]*models.Subscription, int64, error)
	Create(sub *models.Subscription) error
	Update(sub *models.Subscription) error
	Delete(id string) error // soft delete
	Restore(id string) error // 소프트 삭제 복원 (deleted_at = NULL)
	CountByUserID(userID string) (int64, error)
	FindDuplicateName(userID, serviceName string) (bool, error)
	FindSimilarInCategory(userID string, categoryID string, excludeSubID string) ([]*models.Subscription, error)
}

// subscriptionRepository is the GORM implementation of SubscriptionRepository.
type subscriptionRepository struct {
	db *gorm.DB
}

// NewSubscriptionRepository creates a new GORM-backed SubscriptionRepository.
func NewSubscriptionRepository(db *gorm.DB) SubscriptionRepository {
	return &subscriptionRepository{db: db}
}

// FindByID retrieves a subscription by its UUID, preloading the Category.
func (r *subscriptionRepository) FindByID(id string) (*models.Subscription, error) {
	var sub models.Subscription
	if err := r.db.Preload("Category").Where("id = ?", id).First(&sub).Error; err != nil {
		return nil, fmt.Errorf("find subscription by id: %w", err)
	}
	return &sub, nil
}

// FindByUserID retrieves subscriptions for a given user with filtering,
// sorting, and pagination.
func (r *subscriptionRepository) FindByUserID(userID string, filter SubscriptionFilter) ([]*models.Subscription, int64, error) {
	filter.Defaults()

	query := r.db.Model(&models.Subscription{}).Where("user_id = ?", userID)

	// Apply status filter.
	if filter.Status != "" {
		query = query.Where("status = ?", filter.Status)
	}

	// Apply category filter.
	if filter.CategoryID != "" {
		query = query.Where("category_id = ?", filter.CategoryID)
	}

	// Count total before pagination.
	var total int64
	if err := query.Count(&total).Error; err != nil {
		return nil, 0, fmt.Errorf("count subscriptions: %w", err)
	}

	// Determine sort column.
	sortColumn := "created_at"
	switch filter.SortBy {
	case "amount":
		sortColumn = "amount"
	case "satisfaction":
		sortColumn = "satisfaction_score"
	case "next_billing_date":
		sortColumn = "next_billing_date"
	case "created_at":
		sortColumn = "created_at"
	}

	orderClause := fmt.Sprintf("%s %s", sortColumn, filter.SortOrder)

	offset := (filter.Page - 1) * filter.PerPage

	var subs []*models.Subscription
	if err := query.
		Preload("Category").
		Order(orderClause).
		Offset(offset).
		Limit(filter.PerPage).
		Find(&subs).Error; err != nil {
		return nil, 0, fmt.Errorf("find subscriptions by user id: %w", err)
	}

	return subs, total, nil
}

// Create inserts a new subscription into the database.
func (r *subscriptionRepository) Create(sub *models.Subscription) error {
	if err := r.db.Create(sub).Error; err != nil {
		return fmt.Errorf("create subscription: %w", err)
	}
	return nil
}

// Update saves changes to an existing subscription.
func (r *subscriptionRepository) Update(sub *models.Subscription) error {
	if err := r.db.Save(sub).Error; err != nil {
		return fmt.Errorf("update subscription: %w", err)
	}
	return nil
}

// Delete performs a soft delete on a subscription by its UUID.
func (r *subscriptionRepository) Delete(id string) error {
	if err := r.db.Where("id = ?", id).Delete(&models.Subscription{}).Error; err != nil {
		return fmt.Errorf("delete subscription: %w", err)
	}
	return nil
}

// Restore reverses a soft delete by setting deleted_at back to NULL.
func (r *subscriptionRepository) Restore(id string) error {
	if err := r.db.Model(&models.Subscription{}).
		Unscoped().
		Where("id = ?", id).
		Update("deleted_at", nil).Error; err != nil {
		return fmt.Errorf("restore subscription: %w", err)
	}
	return nil
}

// CountByUserID returns the number of non-deleted subscriptions for a user.
func (r *subscriptionRepository) CountByUserID(userID string) (int64, error) {
	var count int64
	if err := r.db.Model(&models.Subscription{}).Where("user_id = ?", userID).Count(&count).Error; err != nil {
		return 0, fmt.Errorf("count subscriptions by user id: %w", err)
	}
	return count, nil
}

// FindDuplicateName checks whether a subscription with the same service name
// already exists for the given user.
func (r *subscriptionRepository) FindDuplicateName(userID, serviceName string) (bool, error) {
	var count int64
	if err := r.db.Model(&models.Subscription{}).
		Where("user_id = ? AND service_name = ?", userID, serviceName).
		Count(&count).Error; err != nil {
		return false, fmt.Errorf("find duplicate subscription name: %w", err)
	}
	return count > 0, nil
}

// FindSimilarInCategory retrieves subscriptions in the same category for a user,
// optionally excluding a specific subscription.
func (r *subscriptionRepository) FindSimilarInCategory(userID string, categoryID string, excludeSubID string) ([]*models.Subscription, error) {
	query := r.db.Model(&models.Subscription{}).
		Preload("Category").
		Where("user_id = ? AND category_id = ?", userID, categoryID)

	if excludeSubID != "" {
		query = query.Where("id != ?", excludeSubID)
	}

	var subs []*models.Subscription
	if err := query.Find(&subs).Error; err != nil {
		return nil, fmt.Errorf("find similar subscriptions in category: %w", err)
	}
	return subs, nil
}
