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

type mockShareGroupRepo struct {
	groups    map[string]*models.ShareGroup
	members   map[string]*models.ShareMember
	createErr error
	updateErr error
	deleteErr error
}

func newMockShareGroupRepo() *mockShareGroupRepo {
	return &mockShareGroupRepo{
		groups:  make(map[string]*models.ShareGroup),
		members: make(map[string]*models.ShareMember),
	}
}

func (m *mockShareGroupRepo) FindByID(id string) (*models.ShareGroup, error) {
	g, ok := m.groups[id]
	if !ok {
		return nil, gorm.ErrRecordNotFound
	}
	// Rebuild members slice from members map.
	var members []models.ShareMember
	for _, mem := range m.members {
		if mem.ShareGroupID.String() == id {
			members = append(members, *mem)
		}
	}
	g.Members = members
	return g, nil
}

func (m *mockShareGroupRepo) FindByOwnerID(ownerID string) ([]*models.ShareGroup, error) {
	var result []*models.ShareGroup
	for _, g := range m.groups {
		if g.OwnerUserID.String() == ownerID {
			// Attach members.
			var members []models.ShareMember
			for _, mem := range m.members {
				if mem.ShareGroupID.String() == g.ID.String() {
					members = append(members, *mem)
				}
			}
			g.Members = members
			result = append(result, g)
		}
	}
	return result, nil
}

func (m *mockShareGroupRepo) Create(group *models.ShareGroup) error {
	if m.createErr != nil {
		return m.createErr
	}
	if group.ID == uuid.Nil {
		group.ID = uuid.New()
	}
	group.CreatedAt = time.Now()
	group.UpdatedAt = time.Now()

	// Persist members.
	for i := range group.Members {
		mem := &group.Members[i]
		if mem.ID == uuid.Nil {
			mem.ID = uuid.New()
		}
		mem.ShareGroupID = group.ID
		mem.CreatedAt = time.Now()
		m.members[mem.ID.String()] = mem
	}

	m.groups[group.ID.String()] = group
	return nil
}

func (m *mockShareGroupRepo) Update(group *models.ShareGroup) error {
	if m.updateErr != nil {
		return m.updateErr
	}
	group.UpdatedAt = time.Now()
	m.groups[group.ID.String()] = group
	return nil
}

func (m *mockShareGroupRepo) Delete(id string) error {
	if m.deleteErr != nil {
		return m.deleteErr
	}
	delete(m.groups, id)
	return nil
}

func (m *mockShareGroupRepo) AddMember(member *models.ShareMember) error {
	if member.ID == uuid.Nil {
		member.ID = uuid.New()
	}
	member.CreatedAt = time.Now()
	m.members[member.ID.String()] = member
	return nil
}

func (m *mockShareGroupRepo) RemoveMember(memberID string) error {
	delete(m.members, memberID)
	return nil
}

func (m *mockShareGroupRepo) RemoveAllSubscriptionShares(groupID string) error {
	return nil
}

// seedGroup inserts a share group with an owner member and optional extra members.
func (m *mockShareGroupRepo) seedGroup(ownerID uuid.UUID, name string, extraNicknames ...string) *models.ShareGroup {
	group := &models.ShareGroup{
		ID:          uuid.New(),
		OwnerUserID: ownerID,
		Name:        name,
		CreatedAt:   time.Now(),
		UpdatedAt:   time.Now(),
	}

	ownerMember := &models.ShareMember{
		ID:           uuid.New(),
		ShareGroupID: group.ID,
		Nickname:     "나",
		IsOwner:      true,
		CreatedAt:    time.Now(),
	}
	m.members[ownerMember.ID.String()] = ownerMember

	for _, nick := range extraNicknames {
		mem := &models.ShareMember{
			ID:           uuid.New(),
			ShareGroupID: group.ID,
			Nickname:     nick,
			IsOwner:      false,
			CreatedAt:    time.Now(),
		}
		m.members[mem.ID.String()] = mem
	}

	m.groups[group.ID.String()] = group
	return group
}

// ---------------------------------------------------------------------------
// Tests – GetShareGroups
// ---------------------------------------------------------------------------

func TestGetShareGroups(t *testing.T) {
	userID := uuid.New()
	otherUserID := uuid.New()

	t.Run("returns all groups owned by user", func(t *testing.T) {
		repo := newMockShareGroupRepo()
		svc := NewShareGroupService(repo)

		repo.seedGroup(userID, "넷플릭스 공유", "친구1")
		repo.seedGroup(userID, "유튜브 공유", "친구2")

		groups, err := svc.GetShareGroups(userID.String())
		assertNil(t, err)
		assertEqual(t, len(groups), 2)
	})

	t.Run("returns empty list for user with no groups", func(t *testing.T) {
		repo := newMockShareGroupRepo()
		svc := NewShareGroupService(repo)

		groups, err := svc.GetShareGroups(uuid.New().String())
		assertNil(t, err)
		assertEqual(t, len(groups), 0)
	})

	t.Run("does not return other users groups", func(t *testing.T) {
		repo := newMockShareGroupRepo()
		svc := NewShareGroupService(repo)

		repo.seedGroup(userID, "내 그룹", "친구1")
		repo.seedGroup(otherUserID, "남의 그룹", "친구2")

		groups, err := svc.GetShareGroups(userID.String())
		assertNil(t, err)
		assertEqual(t, len(groups), 1)
		assertEqual(t, groups[0].Name, "내 그룹")
	})
}

// ---------------------------------------------------------------------------
// Tests – GetShareGroup
// ---------------------------------------------------------------------------

func TestGetShareGroup(t *testing.T) {
	userID := uuid.New()
	otherUserID := uuid.New()

	t.Run("returns group when user is owner", func(t *testing.T) {
		repo := newMockShareGroupRepo()
		svc := NewShareGroupService(repo)
		group := repo.seedGroup(userID, "넷플릭스 공유", "친구1")

		got, err := svc.GetShareGroup(userID.String(), group.ID.String())
		assertNil(t, err)
		assertNotNil(t, got)
		assertEqual(t, got.Name, "넷플릭스 공유")
	})

	t.Run("rejects when user is not owner", func(t *testing.T) {
		repo := newMockShareGroupRepo()
		svc := NewShareGroupService(repo)
		group := repo.seedGroup(otherUserID, "남의 그룹", "친구1")

		_, err := svc.GetShareGroup(userID.String(), group.ID.String())
		assertAppErrorCode(t, err, http.StatusForbidden)
	})

	t.Run("returns error for non-existent group", func(t *testing.T) {
		repo := newMockShareGroupRepo()
		svc := NewShareGroupService(repo)

		_, err := svc.GetShareGroup(userID.String(), uuid.New().String())
		assertAppErrorCode(t, err, http.StatusNotFound)
	})
}

// ---------------------------------------------------------------------------
// Tests – CreateShareGroup
// ---------------------------------------------------------------------------

func TestCreateShareGroup(t *testing.T) {
	userID := uuid.New()

	t.Run("creates group with valid data including owner as member", func(t *testing.T) {
		repo := newMockShareGroupRepo()
		svc := NewShareGroupService(repo)

		req := &CreateShareGroupRequest{
			Name:        "넷플릭스 공유",
			Description: strPtr("넷플릭스 가족 공유"),
			Members: []CreateShareMemberRequest{
				{Nickname: "친구1"},
			},
		}

		group, err := svc.CreateShareGroup(userID.String(), req)
		assertNil(t, err)
		assertNotNil(t, group)
		assertEqual(t, group.Name, "넷플릭스 공유")
		assertEqual(t, group.OwnerUserID.String(), userID.String())
		// Should have 2 members: owner + 1 friend.
		assertEqual(t, len(group.Members), 2)
	})

	t.Run("rejects empty name", func(t *testing.T) {
		repo := newMockShareGroupRepo()
		svc := NewShareGroupService(repo)

		req := &CreateShareGroupRequest{
			Name: "",
			Members: []CreateShareMemberRequest{
				{Nickname: "친구1"},
			},
		}
		_, err := svc.CreateShareGroup(userID.String(), req)
		assertError(t, err)
		assertAppErrorCode(t, err, http.StatusUnprocessableEntity)
	})

	t.Run("rejects fewer than 1 additional member", func(t *testing.T) {
		repo := newMockShareGroupRepo()
		svc := NewShareGroupService(repo)

		req := &CreateShareGroupRequest{
			Name:    "혼자 그룹",
			Members: []CreateShareMemberRequest{},
		}
		_, err := svc.CreateShareGroup(userID.String(), req)
		assertError(t, err)
		assertAppErrorCode(t, err, http.StatusUnprocessableEntity)
	})

	t.Run("owner member is auto-created with isOwner true", func(t *testing.T) {
		repo := newMockShareGroupRepo()
		svc := NewShareGroupService(repo)

		req := &CreateShareGroupRequest{
			Name: "공유 테스트",
			Members: []CreateShareMemberRequest{
				{Nickname: "친구1"},
			},
		}

		group, err := svc.CreateShareGroup(userID.String(), req)
		assertNil(t, err)
		assertNotNil(t, group)

		var ownerFound bool
		for _, mem := range group.Members {
			if mem.IsOwner {
				ownerFound = true
				assertEqual(t, mem.Nickname, "나")
			}
		}
		assertEqual(t, ownerFound, true)
	})
}

// ---------------------------------------------------------------------------
// Tests – UpdateShareGroup
// ---------------------------------------------------------------------------

func TestUpdateShareGroup(t *testing.T) {
	userID := uuid.New()
	otherUserID := uuid.New()

	t.Run("updates name successfully", func(t *testing.T) {
		repo := newMockShareGroupRepo()
		svc := NewShareGroupService(repo)
		group := repo.seedGroup(userID, "이전 이름", "친구1")

		req := &UpdateShareGroupRequest{Name: strPtr("새 이름")}
		updated, err := svc.UpdateShareGroup(userID.String(), group.ID.String(), req)
		assertNil(t, err)
		assertNotNil(t, updated)
		assertEqual(t, updated.Name, "새 이름")
	})

	t.Run("rejects when user is not owner", func(t *testing.T) {
		repo := newMockShareGroupRepo()
		svc := NewShareGroupService(repo)
		group := repo.seedGroup(otherUserID, "남의 그룹", "친구1")

		req := &UpdateShareGroupRequest{Name: strPtr("변경")}
		_, err := svc.UpdateShareGroup(userID.String(), group.ID.String(), req)
		assertAppErrorCode(t, err, http.StatusForbidden)
	})

	t.Run("returns error for non-existent group", func(t *testing.T) {
		repo := newMockShareGroupRepo()
		svc := NewShareGroupService(repo)

		req := &UpdateShareGroupRequest{Name: strPtr("이름")}
		_, err := svc.UpdateShareGroup(userID.String(), uuid.New().String(), req)
		assertAppErrorCode(t, err, http.StatusNotFound)
	})

	t.Run("handles partial update description only", func(t *testing.T) {
		repo := newMockShareGroupRepo()
		svc := NewShareGroupService(repo)
		group := repo.seedGroup(userID, "원래 이름", "친구1")

		req := &UpdateShareGroupRequest{Description: strPtr("새 설명")}
		updated, err := svc.UpdateShareGroup(userID.String(), group.ID.String(), req)
		assertNil(t, err)
		assertNotNil(t, updated)
		assertEqual(t, updated.Name, "원래 이름") // unchanged
		assertEqual(t, *updated.Description, "새 설명")
	})
}

// ---------------------------------------------------------------------------
// Tests – DeleteShareGroup
// ---------------------------------------------------------------------------

func TestDeleteShareGroup(t *testing.T) {
	userID := uuid.New()
	otherUserID := uuid.New()

	t.Run("deletes group successfully", func(t *testing.T) {
		repo := newMockShareGroupRepo()
		svc := NewShareGroupService(repo)
		group := repo.seedGroup(userID, "삭제할 그룹", "친구1")

		err := svc.DeleteShareGroup(userID.String(), group.ID.String())
		assertNil(t, err)

		// Verify it was removed from the mock store.
		_, lookupErr := repo.FindByID(group.ID.String())
		assertError(t, lookupErr)
	})

	t.Run("rejects when user is not owner", func(t *testing.T) {
		repo := newMockShareGroupRepo()
		svc := NewShareGroupService(repo)
		group := repo.seedGroup(otherUserID, "남의 그룹", "친구1")

		err := svc.DeleteShareGroup(userID.String(), group.ID.String())
		assertAppErrorCode(t, err, http.StatusForbidden)
	})

	t.Run("returns error for non-existent group", func(t *testing.T) {
		repo := newMockShareGroupRepo()
		svc := NewShareGroupService(repo)

		err := svc.DeleteShareGroup(userID.String(), uuid.New().String())
		assertAppErrorCode(t, err, http.StatusNotFound)
	})
}
