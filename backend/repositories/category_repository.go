package repositories

import (
	"fmt"

	"github.com/subkeep/backend/models"
	"gorm.io/gorm"
)

// CategoryRepository defines the interface for category data access.
type CategoryRepository interface {
	FindByID(id string) (*models.Category, error)
	FindByUserID(userID string) ([]*models.Category, error)
	FindSystemCategories() ([]*models.Category, error)
	Create(cat *models.Category) error
	Update(cat *models.Category) error
	Delete(id string) error
	ReassignSubscriptions(categoryID, targetCategoryID string) error
}

// categoryRepository is the GORM implementation of CategoryRepository.
type categoryRepository struct {
	db *gorm.DB
}

// NewCategoryRepository creates a new GORM-backed CategoryRepository.
func NewCategoryRepository(db *gorm.DB) CategoryRepository {
	return &categoryRepository{db: db}
}

// FindByID retrieves a category by its UUID.
func (r *categoryRepository) FindByID(id string) (*models.Category, error) {
	var cat models.Category
	if err := r.db.Where("id = ?", id).First(&cat).Error; err != nil {
		return nil, fmt.Errorf("find category by id: %w", err)
	}
	return &cat, nil
}

// FindByUserID retrieves system categories and the user's custom categories,
// ordered by sort_order ASC.
func (r *categoryRepository) FindByUserID(userID string) ([]*models.Category, error) {
	var cats []*models.Category
	if err := r.db.
		Where("user_id = ? OR is_system = true", userID).
		Order("sort_order ASC").
		Find(&cats).Error; err != nil {
		return nil, fmt.Errorf("find categories by user id: %w", err)
	}
	return cats, nil
}

// FindSystemCategories retrieves only system categories.
func (r *categoryRepository) FindSystemCategories() ([]*models.Category, error) {
	var cats []*models.Category
	if err := r.db.Where("is_system = true").Order("sort_order ASC").Find(&cats).Error; err != nil {
		return nil, fmt.Errorf("find system categories: %w", err)
	}
	return cats, nil
}

// Create inserts a new category into the database.
func (r *categoryRepository) Create(cat *models.Category) error {
	if err := r.db.Create(cat).Error; err != nil {
		return fmt.Errorf("create category: %w", err)
	}
	return nil
}

// Update saves changes to an existing category.
func (r *categoryRepository) Update(cat *models.Category) error {
	if err := r.db.Save(cat).Error; err != nil {
		return fmt.Errorf("update category: %w", err)
	}
	return nil
}

// Delete hard-deletes a category by its UUID.
func (r *categoryRepository) Delete(id string) error {
	if err := r.db.Where("id = ?", id).Delete(&models.Category{}).Error; err != nil {
		return fmt.Errorf("delete category: %w", err)
	}
	return nil
}

// ReassignSubscriptions moves all subscriptions from one category to another.
func (r *categoryRepository) ReassignSubscriptions(categoryID, targetCategoryID string) error {
	if err := r.db.Model(&models.Subscription{}).
		Where("category_id = ?", categoryID).
		Update("category_id", targetCategoryID).Error; err != nil {
		return fmt.Errorf("reassign subscriptions: %w", err)
	}
	return nil
}
