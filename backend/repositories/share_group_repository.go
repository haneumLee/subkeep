package repositories

import (
	"fmt"

	"github.com/subkeep/backend/models"
	"gorm.io/gorm"
)

// ShareGroupRepository defines the interface for share group data access.
type ShareGroupRepository interface {
	FindByID(id string) (*models.ShareGroup, error)
	FindByOwnerID(ownerID string) ([]*models.ShareGroup, error)
	Create(group *models.ShareGroup) error
	Update(group *models.ShareGroup) error
	Delete(id string) error
	AddMember(member *models.ShareMember) error
	RemoveMember(memberID string) error
	RemoveAllSubscriptionShares(groupID string) error
}

// shareGroupRepository is the GORM implementation of ShareGroupRepository.
type shareGroupRepository struct {
	db *gorm.DB
}

// NewShareGroupRepository creates a new GORM-backed ShareGroupRepository.
func NewShareGroupRepository(db *gorm.DB) ShareGroupRepository {
	return &shareGroupRepository{db: db}
}

// FindByID retrieves a share group by its UUID, preloading Members.
func (r *shareGroupRepository) FindByID(id string) (*models.ShareGroup, error) {
	var group models.ShareGroup
	if err := r.db.Preload("Members").Where("id = ?", id).First(&group).Error; err != nil {
		return nil, fmt.Errorf("find share group by id: %w", err)
	}
	return &group, nil
}

// FindByOwnerID retrieves all share groups owned by the given user,
// preloading Members.
func (r *shareGroupRepository) FindByOwnerID(ownerID string) ([]*models.ShareGroup, error) {
	var groups []*models.ShareGroup
	if err := r.db.Preload("Members").
		Where("owner_user_id = ?", ownerID).
		Order("created_at DESC").
		Find(&groups).Error; err != nil {
		return nil, fmt.Errorf("find share groups by owner id: %w", err)
	}
	return groups, nil
}

// Create inserts a new share group into the database.
func (r *shareGroupRepository) Create(group *models.ShareGroup) error {
	if err := r.db.Create(group).Error; err != nil {
		return fmt.Errorf("create share group: %w", err)
	}
	return nil
}

// Update saves changes to an existing share group.
func (r *shareGroupRepository) Update(group *models.ShareGroup) error {
	if err := r.db.Save(group).Error; err != nil {
		return fmt.Errorf("update share group: %w", err)
	}
	return nil
}

// Delete performs a soft delete on a share group by its UUID.
func (r *shareGroupRepository) Delete(id string) error {
	if err := r.db.Where("id = ?", id).Delete(&models.ShareGroup{}).Error; err != nil {
		return fmt.Errorf("delete share group: %w", err)
	}
	return nil
}

// AddMember inserts a new share member into the database.
func (r *shareGroupRepository) AddMember(member *models.ShareMember) error {
	if err := r.db.Create(member).Error; err != nil {
		return fmt.Errorf("add share member: %w", err)
	}
	return nil
}

// RemoveMember deletes a share member by its UUID.
func (r *shareGroupRepository) RemoveMember(memberID string) error {
	if err := r.db.Where("id = ?", memberID).Delete(&models.ShareMember{}).Error; err != nil {
		return fmt.Errorf("remove share member: %w", err)
	}
	return nil
}

// RemoveAllSubscriptionShares deletes all SubscriptionShare records
// referencing the given share group.
func (r *shareGroupRepository) RemoveAllSubscriptionShares(groupID string) error {
	if err := r.db.Where("share_group_id = ?", groupID).
		Delete(&models.SubscriptionShare{}).Error; err != nil {
		return fmt.Errorf("remove subscription shares by group id: %w", err)
	}
	return nil
}
