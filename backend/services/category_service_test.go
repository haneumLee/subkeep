package services

import (
	"net/http"
	"testing"
	"time"

	"github.com/google/uuid"
	"gorm.io/gorm"

	"github.com/subkeep/backend/models"
)

// ---------------------------------------------------------------------------
// Mock repository
// ---------------------------------------------------------------------------

type mockCategoryRepo struct {
	categories map[string]*models.Category
	createErr  error
	updateErr  error
	deleteErr  error
}

func newMockCategoryRepo() *mockCategoryRepo {
	return &mockCategoryRepo{categories: make(map[string]*models.Category)}
}

func (m *mockCategoryRepo) FindByID(id string) (*models.Category, error) {
	cat, ok := m.categories[id]
	if !ok {
		return nil, gorm.ErrRecordNotFound
	}
	return cat, nil
}

func (m *mockCategoryRepo) FindByUserID(userID string) ([]*models.Category, error) {
	var result []*models.Category
	for _, cat := range m.categories {
		if cat.IsSystem {
			result = append(result, cat)
			continue
		}
		if cat.UserID != nil && cat.UserID.String() == userID {
			result = append(result, cat)
		}
	}
	return result, nil
}

func (m *mockCategoryRepo) FindSystemCategories() ([]*models.Category, error) {
	var result []*models.Category
	for _, cat := range m.categories {
		if cat.IsSystem {
			result = append(result, cat)
		}
	}
	return result, nil
}

func (m *mockCategoryRepo) Create(cat *models.Category) error {
	if m.createErr != nil {
		return m.createErr
	}
	if cat.ID == uuid.Nil {
		cat.ID = uuid.New()
	}
	cat.CreatedAt = time.Now()
	cat.UpdatedAt = time.Now()
	m.categories[cat.ID.String()] = cat
	return nil
}

func (m *mockCategoryRepo) Update(cat *models.Category) error {
	if m.updateErr != nil {
		return m.updateErr
	}
	cat.UpdatedAt = time.Now()
	m.categories[cat.ID.String()] = cat
	return nil
}

func (m *mockCategoryRepo) Delete(id string) error {
	if m.deleteErr != nil {
		return m.deleteErr
	}
	delete(m.categories, id)
	return nil
}

func (m *mockCategoryRepo) ReassignSubscriptions(categoryID, targetCategoryID string) error {
	return nil
}

// seedSystemCategory inserts a system category into the mock repo.
func (m *mockCategoryRepo) seedSystemCategory(name string) *models.Category {
	cat := &models.Category{
		ID:        uuid.New(),
		UserID:    nil,
		Name:      name,
		IsSystem:  true,
		SortOrder: 0,
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
	}
	m.categories[cat.ID.String()] = cat
	return cat
}

// seedUserCategory inserts a user custom category into the mock repo.
func (m *mockCategoryRepo) seedUserCategory(userID uuid.UUID, name string) *models.Category {
	cat := &models.Category{
		ID:        uuid.New(),
		UserID:    &userID,
		Name:      name,
		IsSystem:  false,
		SortOrder: 0,
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
	}
	m.categories[cat.ID.String()] = cat
	return cat
}

// ---------------------------------------------------------------------------
// Tests – GetCategories
// ---------------------------------------------------------------------------

func TestGetCategories(t *testing.T) {
	userID := uuid.New()
	otherUserID := uuid.New()

	t.Run("returns system + user custom categories", func(t *testing.T) {
		repo := newMockCategoryRepo()
		svc := NewCategoryService(repo)

		repo.seedSystemCategory("엔터테인먼트")
		repo.seedSystemCategory("기타")
		repo.seedUserCategory(userID, "내 카테고리")

		cats, err := svc.GetCategories(userID.String())
		assertNil(t, err)
		assertEqual(t, len(cats), 3)
	})

	t.Run("returns only system categories for user with no custom ones", func(t *testing.T) {
		repo := newMockCategoryRepo()
		svc := NewCategoryService(repo)

		repo.seedSystemCategory("엔터테인먼트")
		repo.seedSystemCategory("기타")

		cats, err := svc.GetCategories(userID.String())
		assertNil(t, err)
		assertEqual(t, len(cats), 2)
	})

	t.Run("does not return other users custom categories", func(t *testing.T) {
		repo := newMockCategoryRepo()
		svc := NewCategoryService(repo)

		repo.seedSystemCategory("엔터테인먼트")
		repo.seedUserCategory(userID, "내 카테고리")
		repo.seedUserCategory(otherUserID, "다른 사용자 카테고리")

		cats, err := svc.GetCategories(userID.String())
		assertNil(t, err)
		assertEqual(t, len(cats), 2) // system + own custom only
		for _, c := range cats {
			if !c.IsSystem {
				assertEqual(t, c.UserID.String(), userID.String())
			}
		}
	})
}

// ---------------------------------------------------------------------------
// Tests – CreateCategory
// ---------------------------------------------------------------------------

func TestCreateCategory(t *testing.T) {
	userID := uuid.New()

	t.Run("creates custom category with valid data", func(t *testing.T) {
		repo := newMockCategoryRepo()
		svc := NewCategoryService(repo)

		req := &CreateCategoryRequest{
			Name:  "게임",
			Color: strPtr("#FF5733"),
			Icon:  strPtr("game-icon"),
		}

		cat, err := svc.CreateCategory(userID.String(), req)
		assertNil(t, err)
		assertNotNil(t, cat)
		assertEqual(t, cat.Name, "게임")
		assertEqual(t, *cat.Color, "#FF5733")
		assertEqual(t, *cat.Icon, "game-icon")
		assertEqual(t, cat.IsSystem, false)
		assertEqual(t, cat.UserID.String(), userID.String())
	})

	t.Run("rejects empty name", func(t *testing.T) {
		repo := newMockCategoryRepo()
		svc := NewCategoryService(repo)

		req := &CreateCategoryRequest{Name: ""}
		_, err := svc.CreateCategory(userID.String(), req)
		assertError(t, err)
		assertAppErrorCode(t, err, http.StatusUnprocessableEntity)
	})

	t.Run("sets IsSystem to false automatically", func(t *testing.T) {
		repo := newMockCategoryRepo()
		svc := NewCategoryService(repo)

		req := &CreateCategoryRequest{Name: "직접 만든 카테고리"}
		cat, err := svc.CreateCategory(userID.String(), req)
		assertNil(t, err)
		assertNotNil(t, cat)
		assertEqual(t, cat.IsSystem, false)
	})

	t.Run("handles optional fields", func(t *testing.T) {
		repo := newMockCategoryRepo()
		svc := NewCategoryService(repo)

		req := &CreateCategoryRequest{
			Name:      "옵션 테스트",
			SortOrder: intPtr(5),
		}

		cat, err := svc.CreateCategory(userID.String(), req)
		assertNil(t, err)
		assertNotNil(t, cat)
		assertEqual(t, cat.SortOrder, 5)
		// Color and Icon are nil when not provided.
		assertNil(t, cat.Color)
		assertNil(t, cat.Icon)
	})
}

// ---------------------------------------------------------------------------
// Tests – UpdateCategory
// ---------------------------------------------------------------------------

func TestUpdateCategory(t *testing.T) {
	userID := uuid.New()
	otherUserID := uuid.New()

	t.Run("updates category name successfully", func(t *testing.T) {
		repo := newMockCategoryRepo()
		svc := NewCategoryService(repo)
		cat := repo.seedUserCategory(userID, "이전 이름")

		req := &UpdateCategoryRequest{Name: strPtr("새 이름")}
		updated, err := svc.UpdateCategory(userID.String(), cat.ID.String(), req)
		assertNil(t, err)
		assertNotNil(t, updated)
		assertEqual(t, updated.Name, "새 이름")
	})

	t.Run("rejects updating system category", func(t *testing.T) {
		repo := newMockCategoryRepo()
		svc := NewCategoryService(repo)
		sysCat := repo.seedSystemCategory("엔터테인먼트")

		req := &UpdateCategoryRequest{Name: strPtr("변경")}
		_, err := svc.UpdateCategory(userID.String(), sysCat.ID.String(), req)
		assertAppErrorCode(t, err, http.StatusForbidden)
	})

	t.Run("rejects updating another users category", func(t *testing.T) {
		repo := newMockCategoryRepo()
		svc := NewCategoryService(repo)
		cat := repo.seedUserCategory(otherUserID, "남의 카테고리")

		req := &UpdateCategoryRequest{Name: strPtr("변경 시도")}
		_, err := svc.UpdateCategory(userID.String(), cat.ID.String(), req)
		assertAppErrorCode(t, err, http.StatusForbidden)
	})

	t.Run("returns error for non-existent category", func(t *testing.T) {
		repo := newMockCategoryRepo()
		svc := NewCategoryService(repo)

		req := &UpdateCategoryRequest{Name: strPtr("이름")}
		_, err := svc.UpdateCategory(userID.String(), uuid.New().String(), req)
		assertAppErrorCode(t, err, http.StatusNotFound)
	})

	t.Run("handles partial update only color", func(t *testing.T) {
		repo := newMockCategoryRepo()
		svc := NewCategoryService(repo)
		cat := repo.seedUserCategory(userID, "원래 이름")

		req := &UpdateCategoryRequest{Color: strPtr("#AABBCC")}
		updated, err := svc.UpdateCategory(userID.String(), cat.ID.String(), req)
		assertNil(t, err)
		assertNotNil(t, updated)
		assertEqual(t, updated.Name, "원래 이름") // unchanged
		assertEqual(t, *updated.Color, "#AABBCC")
	})
}

// ---------------------------------------------------------------------------
// Tests – DeleteCategory
// ---------------------------------------------------------------------------

func TestDeleteCategory(t *testing.T) {
	userID := uuid.New()
	otherUserID := uuid.New()

	t.Run("deletes custom category successfully", func(t *testing.T) {
		repo := newMockCategoryRepo()
		svc := NewCategoryService(repo)

		// Seed the system "기타" category for reassignment.
		repo.seedSystemCategory("기타")

		cat := repo.seedUserCategory(userID, "삭제할 카테고리")

		err := svc.DeleteCategory(userID.String(), cat.ID.String())
		assertNil(t, err)

		// Verify it was removed from the mock store.
		_, lookupErr := repo.FindByID(cat.ID.String())
		assertError(t, lookupErr)
	})

	t.Run("rejects deleting system category", func(t *testing.T) {
		repo := newMockCategoryRepo()
		svc := NewCategoryService(repo)
		sysCat := repo.seedSystemCategory("엔터테인먼트")

		err := svc.DeleteCategory(userID.String(), sysCat.ID.String())
		assertAppErrorCode(t, err, http.StatusForbidden)
	})

	t.Run("rejects deleting another users category", func(t *testing.T) {
		repo := newMockCategoryRepo()
		svc := NewCategoryService(repo)
		cat := repo.seedUserCategory(otherUserID, "남의 카테고리")

		err := svc.DeleteCategory(userID.String(), cat.ID.String())
		assertAppErrorCode(t, err, http.StatusForbidden)
	})

	t.Run("returns error for non-existent category", func(t *testing.T) {
		repo := newMockCategoryRepo()
		svc := NewCategoryService(repo)

		err := svc.DeleteCategory(userID.String(), uuid.New().String())
		assertAppErrorCode(t, err, http.StatusNotFound)
	})
}
