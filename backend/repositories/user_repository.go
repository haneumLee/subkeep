package repositories

import (
	"fmt"
	"time"

	"github.com/subkeep/backend/models"
	"gorm.io/gorm"
)

// UserRepository defines the interface for user data access.
type UserRepository interface {
	FindByID(id string) (*models.User, error)
	FindByProviderID(provider, providerUserID string) (*models.User, error)
	FindByEmail(email string) (*models.User, error)
	Create(user *models.User) error
	Update(user *models.User) error
	Delete(id string) error
	UpdateLastLogin(id string) error
}

// userRepository is the GORM implementation of UserRepository.
type userRepository struct {
	db *gorm.DB
}

// NewUserRepository creates a new GORM-backed UserRepository.
func NewUserRepository(db *gorm.DB) UserRepository {
	return &userRepository{db: db}
}

// FindByID retrieves a user by their UUID.
func (r *userRepository) FindByID(id string) (*models.User, error) {
	var user models.User
	if err := r.db.Where("id = ?", id).First(&user).Error; err != nil {
		return nil, fmt.Errorf("find user by id: %w", err)
	}
	return &user, nil
}

// FindByProviderID retrieves a user by their OAuth provider and provider user ID.
func (r *userRepository) FindByProviderID(provider, providerUserID string) (*models.User, error) {
	var user models.User
	if err := r.db.Where("provider = ? AND provider_user_id = ?", provider, providerUserID).First(&user).Error; err != nil {
		return nil, fmt.Errorf("find user by provider id: %w", err)
	}
	return &user, nil
}

// FindByEmail retrieves a user by their email address.
func (r *userRepository) FindByEmail(email string) (*models.User, error) {
	var user models.User
	if err := r.db.Where("email = ?", email).First(&user).Error; err != nil {
		return nil, fmt.Errorf("find user by email: %w", err)
	}
	return &user, nil
}

// Create inserts a new user into the database.
func (r *userRepository) Create(user *models.User) error {
	if err := r.db.Create(user).Error; err != nil {
		return fmt.Errorf("create user: %w", err)
	}
	return nil
}

// Update saves changes to an existing user.
func (r *userRepository) Update(user *models.User) error {
	if err := r.db.Save(user).Error; err != nil {
		return fmt.Errorf("update user: %w", err)
	}
	return nil
}

// Delete performs a soft delete on a user by their UUID.
func (r *userRepository) Delete(id string) error {
	if err := r.db.Where("id = ?", id).Delete(&models.User{}).Error; err != nil {
		return fmt.Errorf("delete user: %w", err)
	}
	return nil
}

// UpdateLastLogin sets the last login timestamp to the current time.
func (r *userRepository) UpdateLastLogin(id string) error {
	now := time.Now()
	if err := r.db.Model(&models.User{}).Where("id = ?", id).Update("last_login_at", now).Error; err != nil {
		return fmt.Errorf("update last login: %w", err)
	}
	return nil
}
