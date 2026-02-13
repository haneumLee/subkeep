package repositories

import (
	"fmt"

	"github.com/subkeep/backend/models"
	"gorm.io/gorm"
)

// FolderRepository defines the interface for folder data access.
type FolderRepository interface {
	FindByID(id string) (*models.Folder, error)
	FindByUserID(userID string) ([]*models.Folder, error)
	Create(folder *models.Folder) error
	Update(folder *models.Folder) error
	Delete(id string) error
}

// folderRepository is the GORM implementation of FolderRepository.
type folderRepository struct {
	db *gorm.DB
}

// NewFolderRepository creates a new GORM-backed FolderRepository.
func NewFolderRepository(db *gorm.DB) FolderRepository {
	return &folderRepository{db: db}
}

// FindByID retrieves a folder by its UUID.
func (r *folderRepository) FindByID(id string) (*models.Folder, error) {
	var folder models.Folder
	if err := r.db.Where("id = ?", id).First(&folder).Error; err != nil {
		return nil, fmt.Errorf("find folder by id: %w", err)
	}
	return &folder, nil
}

// FindByUserID retrieves all folders for a given user, ordered by sort_order ASC.
func (r *folderRepository) FindByUserID(userID string) ([]*models.Folder, error) {
	var folders []*models.Folder
	if err := r.db.
		Where("user_id = ?", userID).
		Order("sort_order ASC, created_at ASC").
		Find(&folders).Error; err != nil {
		return nil, fmt.Errorf("find folders by user id: %w", err)
	}
	return folders, nil
}

// Create inserts a new folder into the database.
func (r *folderRepository) Create(folder *models.Folder) error {
	if err := r.db.Create(folder).Error; err != nil {
		return fmt.Errorf("create folder: %w", err)
	}
	return nil
}

// Update saves changes to an existing folder.
func (r *folderRepository) Update(folder *models.Folder) error {
	if err := r.db.Save(folder).Error; err != nil {
		return fmt.Errorf("update folder: %w", err)
	}
	return nil
}

// Delete soft-deletes a folder by its UUID and clears folderId from subscriptions.
func (r *folderRepository) Delete(id string) error {
	// Clear folderId from subscriptions with this folder.
	if err := r.db.Model(&models.Subscription{}).
		Where("folder_id = ?", id).
		Update("folder_id", nil).Error; err != nil {
		return fmt.Errorf("clear folder from subscriptions: %w", err)
	}

	if err := r.db.Where("id = ?", id).Delete(&models.Folder{}).Error; err != nil {
		return fmt.Errorf("delete folder: %w", err)
	}
	return nil
}
