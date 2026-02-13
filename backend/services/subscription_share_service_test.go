package services

import (
	"net/http"
	"testing"
	"time"

	"github.com/google/uuid"
	"gorm.io/gorm"

	"github.com/subkeep/backend/models"
	"github.com/subkeep/backend/repositories"
)

// ---------------------------------------------------------------------------
// Mock SubscriptionShare repository
// ---------------------------------------------------------------------------

type mockSubscriptionShareRepo struct {
	shares    map[string]*models.SubscriptionShare
	createErr error
	updateErr error
	deleteErr error
}

func newMockSubscriptionShareRepo() *mockSubscriptionShareRepo {
	return &mockSubscriptionShareRepo{
		shares: make(map[string]*models.SubscriptionShare),
	}
}

func (m *mockSubscriptionShareRepo) FindByID(id string) (*models.SubscriptionShare, error) {
	s, ok := m.shares[id]
	if !ok {
		return nil, gorm.ErrRecordNotFound
	}
	return s, nil
}

func (m *mockSubscriptionShareRepo) FindBySubscriptionID(subscriptionID string) (*models.SubscriptionShare, error) {
	for _, s := range m.shares {
		if s.SubscriptionID.String() == subscriptionID {
			return s, nil
		}
	}
	return nil, gorm.ErrRecordNotFound
}

func (m *mockSubscriptionShareRepo) FindByUserID(userID string) ([]*models.SubscriptionShare, error) {
	var result []*models.SubscriptionShare
	for _, s := range m.shares {
		result = append(result, s)
	}
	return result, nil
}

func (m *mockSubscriptionShareRepo) Create(share *models.SubscriptionShare) error {
	if m.createErr != nil {
		return m.createErr
	}
	if share.ID == uuid.Nil {
		share.ID = uuid.New()
	}
	share.CreatedAt = time.Now()
	share.UpdatedAt = time.Now()
	m.shares[share.ID.String()] = share
	return nil
}

func (m *mockSubscriptionShareRepo) Update(share *models.SubscriptionShare) error {
	if m.updateErr != nil {
		return m.updateErr
	}
	share.UpdatedAt = time.Now()
	m.shares[share.ID.String()] = share
	return nil
}

func (m *mockSubscriptionShareRepo) Delete(id string) error {
	if m.deleteErr != nil {
		return m.deleteErr
	}
	delete(m.shares, id)
	return nil
}

func (m *mockSubscriptionShareRepo) DeleteBySubscriptionID(subscriptionID string) error {
	for id, s := range m.shares {
		if s.SubscriptionID.String() == subscriptionID {
			delete(m.shares, id)
		}
	}
	return nil
}

// ---------------------------------------------------------------------------
// Mock Subscription repository (for ownership checks)
// ---------------------------------------------------------------------------

type mockSubRepoForShare struct {
	subs map[string]*models.Subscription
}

func newMockSubRepoForShare() *mockSubRepoForShare {
	return &mockSubRepoForShare{subs: make(map[string]*models.Subscription)}
}

func (m *mockSubRepoForShare) FindByID(id string) (*models.Subscription, error) {
	sub, ok := m.subs[id]
	if !ok {
		return nil, gorm.ErrRecordNotFound
	}
	return sub, nil
}

func (m *mockSubRepoForShare) FindByUserID(userID string, filter repositories.SubscriptionFilter) ([]*models.Subscription, int64, error) {
	return nil, 0, nil
}

func (m *mockSubRepoForShare) Create(sub *models.Subscription) error { return nil }
func (m *mockSubRepoForShare) Update(sub *models.Subscription) error { return nil }
func (m *mockSubRepoForShare) Delete(id string) error                { return nil }
func (m *mockSubRepoForShare) Restore(id string) error               { return nil }
func (m *mockSubRepoForShare) CountByUserID(userID string) (int64, error) {
	return 0, nil
}
func (m *mockSubRepoForShare) FindDuplicateName(userID, serviceName string) (bool, error) {
	return false, nil
}

func (m *mockSubRepoForShare) seedSubscription(userID uuid.UUID, name string) *models.Subscription {
	sub := &models.Subscription{
		ID:           uuid.New(),
		UserID:       userID,
		ServiceName:  name,
		Amount:       15000,
		BillingCycle: models.BillingCycleMonthly,
		Currency:     "KRW",
		Status:       models.SubscriptionStatusActive,
		StartDate:    time.Now(),
		NextBillingDate: time.Now().AddDate(0, 1, 0),
		CreatedAt:    time.Now(),
		UpdatedAt:    time.Now(),
	}
	m.subs[sub.ID.String()] = sub
	return sub
}

// ---------------------------------------------------------------------------
// Mock ShareGroup repository (for ownership checks)
// ---------------------------------------------------------------------------

type mockShareGroupRepoForShare struct {
	groups  map[string]*models.ShareGroup
	members map[string]*models.ShareMember
}

func newMockShareGroupRepoForShare() *mockShareGroupRepoForShare {
	return &mockShareGroupRepoForShare{
		groups:  make(map[string]*models.ShareGroup),
		members: make(map[string]*models.ShareMember),
	}
}

func (m *mockShareGroupRepoForShare) FindByID(id string) (*models.ShareGroup, error) {
	g, ok := m.groups[id]
	if !ok {
		return nil, gorm.ErrRecordNotFound
	}
	// Rebuild members slice.
	var members []models.ShareMember
	for _, mem := range m.members {
		if mem.ShareGroupID.String() == id {
			members = append(members, *mem)
		}
	}
	g.Members = members
	return g, nil
}

func (m *mockShareGroupRepoForShare) FindByOwnerID(ownerID string) ([]*models.ShareGroup, error) {
	return nil, nil
}

func (m *mockShareGroupRepoForShare) Create(group *models.ShareGroup) error  { return nil }
func (m *mockShareGroupRepoForShare) Update(group *models.ShareGroup) error  { return nil }
func (m *mockShareGroupRepoForShare) Delete(id string) error                 { return nil }
func (m *mockShareGroupRepoForShare) AddMember(member *models.ShareMember) error { return nil }
func (m *mockShareGroupRepoForShare) RemoveMember(memberID string) error     { return nil }
func (m *mockShareGroupRepoForShare) RemoveAllSubscriptionShares(groupID string) error {
	return nil
}

func (m *mockShareGroupRepoForShare) seedGroup(ownerID uuid.UUID, name string, extraNicknames ...string) *models.ShareGroup {
	group := &models.ShareGroup{
		ID:          uuid.New(),
		OwnerUserID: ownerID,
		Name:        name,
		CreatedAt:   time.Now(),
		UpdatedAt:   time.Now(),
	}

	// Owner member.
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
// Helper to build service with mocks
// ---------------------------------------------------------------------------

func float64Ptr(f float64) *float64 { return &f }

func newTestShareService() (*SubscriptionShareService, *mockSubscriptionShareRepo, *mockSubRepoForShare, *mockShareGroupRepoForShare) {
	shareRepo := newMockSubscriptionShareRepo()
	subRepo := newMockSubRepoForShare()
	groupRepo := newMockShareGroupRepoForShare()
	svc := NewSubscriptionShareService(shareRepo, subRepo, groupRepo)
	return svc, shareRepo, subRepo, groupRepo
}

// ---------------------------------------------------------------------------
// Tests – LinkSubscriptionToShareGroup
// ---------------------------------------------------------------------------

func TestLinkSubscriptionToShareGroup_EqualSplit(t *testing.T) {
	userID := uuid.New()
	svc, _, subRepo, groupRepo := newTestShareService()

	sub := subRepo.seedSubscription(userID, "넷플릭스")
	group := groupRepo.seedGroup(userID, "넷플릭스 공유", "친구1")

	req := &LinkShareRequest{
		SubscriptionID: sub.ID.String(),
		ShareGroupID:   group.ID.String(),
		SplitType:      "equal",
	}

	share, err := svc.LinkSubscriptionToShareGroup(userID.String(), req)
	assertNil(t, err)
	assertNotNil(t, share)
	assertEqual(t, share.SplitType, models.SplitTypeEqual)
	assertEqual(t, share.SubscriptionID.String(), sub.ID.String())
	assertEqual(t, share.ShareGroupID.String(), group.ID.String())
	assertEqual(t, share.TotalMembersSnapshot, 2) // owner + 1 friend
}

func TestLinkSubscriptionToShareGroup_CustomAmount(t *testing.T) {
	userID := uuid.New()
	svc, _, subRepo, groupRepo := newTestShareService()

	sub := subRepo.seedSubscription(userID, "유튜브 프리미엄")
	group := groupRepo.seedGroup(userID, "유튜브 공유", "친구1", "친구2")

	req := &LinkShareRequest{
		SubscriptionID: sub.ID.String(),
		ShareGroupID:   group.ID.String(),
		SplitType:      "custom_amount",
		MyShareAmount:  intPtr(5000),
	}

	share, err := svc.LinkSubscriptionToShareGroup(userID.String(), req)
	assertNil(t, err)
	assertNotNil(t, share)
	assertEqual(t, share.SplitType, models.SplitTypeCustomAmount)
	assertEqual(t, *share.MyShareAmount, 5000)
	assertEqual(t, share.TotalMembersSnapshot, 3) // owner + 2 friends
}

func TestLinkSubscriptionToShareGroup_CustomRatio(t *testing.T) {
	userID := uuid.New()
	svc, _, subRepo, groupRepo := newTestShareService()

	sub := subRepo.seedSubscription(userID, "스포티파이")
	group := groupRepo.seedGroup(userID, "스포티파이 공유", "친구1")

	ratio := 0.6
	req := &LinkShareRequest{
		SubscriptionID: sub.ID.String(),
		ShareGroupID:   group.ID.String(),
		SplitType:      "custom_ratio",
		MyShareRatio:   &ratio,
	}

	share, err := svc.LinkSubscriptionToShareGroup(userID.String(), req)
	assertNil(t, err)
	assertNotNil(t, share)
	assertEqual(t, share.SplitType, models.SplitTypeCustomRatio)
	assertEqual(t, *share.MyShareRatio, 0.6)
}

func TestLinkSubscriptionToShareGroup_ForbiddenSubscription(t *testing.T) {
	userID := uuid.New()
	otherUserID := uuid.New()
	svc, _, subRepo, groupRepo := newTestShareService()

	sub := subRepo.seedSubscription(otherUserID, "남의 구독")
	group := groupRepo.seedGroup(userID, "내 그룹", "친구1")

	req := &LinkShareRequest{
		SubscriptionID: sub.ID.String(),
		ShareGroupID:   group.ID.String(),
		SplitType:      "equal",
	}

	_, err := svc.LinkSubscriptionToShareGroup(userID.String(), req)
	assertAppErrorCode(t, err, http.StatusForbidden)
}

func TestLinkSubscriptionToShareGroup_ForbiddenShareGroup(t *testing.T) {
	userID := uuid.New()
	otherUserID := uuid.New()
	svc, _, subRepo, groupRepo := newTestShareService()

	sub := subRepo.seedSubscription(userID, "내 구독")
	group := groupRepo.seedGroup(otherUserID, "남의 그룹", "친구1")

	req := &LinkShareRequest{
		SubscriptionID: sub.ID.String(),
		ShareGroupID:   group.ID.String(),
		SplitType:      "equal",
	}

	_, err := svc.LinkSubscriptionToShareGroup(userID.String(), req)
	assertAppErrorCode(t, err, http.StatusForbidden)
}

func TestLinkSubscriptionToShareGroup_DuplicateLink(t *testing.T) {
	userID := uuid.New()
	svc, _, subRepo, groupRepo := newTestShareService()

	sub := subRepo.seedSubscription(userID, "넷플릭스")
	group := groupRepo.seedGroup(userID, "넷플릭스 공유", "친구1")

	req := &LinkShareRequest{
		SubscriptionID: sub.ID.String(),
		ShareGroupID:   group.ID.String(),
		SplitType:      "equal",
	}

	// First link should succeed.
	_, err := svc.LinkSubscriptionToShareGroup(userID.String(), req)
	assertNil(t, err)

	// Second link should fail.
	_, err = svc.LinkSubscriptionToShareGroup(userID.String(), req)
	assertAppErrorCode(t, err, http.StatusBadRequest)
}

func TestLinkSubscriptionToShareGroup_CustomAmountMissing(t *testing.T) {
	userID := uuid.New()
	svc, _, subRepo, groupRepo := newTestShareService()

	sub := subRepo.seedSubscription(userID, "유튜브")
	group := groupRepo.seedGroup(userID, "공유 그룹", "친구1")

	req := &LinkShareRequest{
		SubscriptionID: sub.ID.String(),
		ShareGroupID:   group.ID.String(),
		SplitType:      "custom_amount",
		// MyShareAmount is intentionally missing.
	}

	_, err := svc.LinkSubscriptionToShareGroup(userID.String(), req)
	assertAppErrorCode(t, err, http.StatusUnprocessableEntity)
}

func TestLinkSubscriptionToShareGroup_CustomRatioMissing(t *testing.T) {
	userID := uuid.New()
	svc, _, subRepo, groupRepo := newTestShareService()

	sub := subRepo.seedSubscription(userID, "스포티파이")
	group := groupRepo.seedGroup(userID, "공유 그룹", "친구1")

	req := &LinkShareRequest{
		SubscriptionID: sub.ID.String(),
		ShareGroupID:   group.ID.String(),
		SplitType:      "custom_ratio",
		// MyShareRatio is intentionally missing.
	}

	_, err := svc.LinkSubscriptionToShareGroup(userID.String(), req)
	assertAppErrorCode(t, err, http.StatusUnprocessableEntity)
}

func TestLinkSubscriptionToShareGroup_SubscriptionNotFound(t *testing.T) {
	userID := uuid.New()
	svc, _, _, groupRepo := newTestShareService()

	group := groupRepo.seedGroup(userID, "공유 그룹", "친구1")

	req := &LinkShareRequest{
		SubscriptionID: uuid.New().String(),
		ShareGroupID:   group.ID.String(),
		SplitType:      "equal",
	}

	_, err := svc.LinkSubscriptionToShareGroup(userID.String(), req)
	assertAppErrorCode(t, err, http.StatusNotFound)
}

// ---------------------------------------------------------------------------
// Tests – UnlinkSubscriptionShare
// ---------------------------------------------------------------------------

func TestUnlinkSubscriptionShare_Success(t *testing.T) {
	userID := uuid.New()
	svc, shareRepo, subRepo, groupRepo := newTestShareService()

	sub := subRepo.seedSubscription(userID, "넷플릭스")
	group := groupRepo.seedGroup(userID, "넷플릭스 공유", "친구1")

	// Link first.
	req := &LinkShareRequest{
		SubscriptionID: sub.ID.String(),
		ShareGroupID:   group.ID.String(),
		SplitType:      "equal",
	}
	share, err := svc.LinkSubscriptionToShareGroup(userID.String(), req)
	assertNil(t, err)
	assertNotNil(t, share)

	// Unlink.
	err = svc.UnlinkSubscriptionShare(userID.String(), share.ID.String())
	assertNil(t, err)

	// Verify removed.
	_, lookupErr := shareRepo.FindByID(share.ID.String())
	assertError(t, lookupErr)
}

// ---------------------------------------------------------------------------
// Tests – UpdateSubscriptionShare
// ---------------------------------------------------------------------------

func TestUpdateSubscriptionShare_Success(t *testing.T) {
	userID := uuid.New()
	svc, _, subRepo, groupRepo := newTestShareService()

	sub := subRepo.seedSubscription(userID, "넷플릭스")
	group := groupRepo.seedGroup(userID, "넷플릭스 공유", "친구1")

	// Link with equal split.
	linkReq := &LinkShareRequest{
		SubscriptionID: sub.ID.String(),
		ShareGroupID:   group.ID.String(),
		SplitType:      "equal",
	}
	share, err := svc.LinkSubscriptionToShareGroup(userID.String(), linkReq)
	assertNil(t, err)
	assertNotNil(t, share)

	// Update to custom_amount.
	newSplitType := "custom_amount"
	updateReq := &UpdateShareRequest{
		SplitType:     &newSplitType,
		MyShareAmount: intPtr(7000),
	}
	updated, err := svc.UpdateSubscriptionShare(userID.String(), share.ID.String(), updateReq)
	assertNil(t, err)
	assertNotNil(t, updated)
	assertEqual(t, updated.SplitType, models.SplitTypeCustomAmount)
	assertEqual(t, *updated.MyShareAmount, 7000)
}

// ---------------------------------------------------------------------------
// Tests – GetSubscriptionShare
// ---------------------------------------------------------------------------

func TestGetSubscriptionShare_Success(t *testing.T) {
	userID := uuid.New()
	svc, _, subRepo, groupRepo := newTestShareService()

	sub := subRepo.seedSubscription(userID, "넷플릭스")
	group := groupRepo.seedGroup(userID, "넷플릭스 공유", "친구1")

	linkReq := &LinkShareRequest{
		SubscriptionID: sub.ID.String(),
		ShareGroupID:   group.ID.String(),
		SplitType:      "equal",
	}
	_, err := svc.LinkSubscriptionToShareGroup(userID.String(), linkReq)
	assertNil(t, err)

	got, err := svc.GetSubscriptionShare(userID.String(), sub.ID.String())
	assertNil(t, err)
	assertNotNil(t, got)
	assertEqual(t, got.SubscriptionID.String(), sub.ID.String())
}

func TestGetSubscriptionShare_SubscriptionNotFound(t *testing.T) {
	userID := uuid.New()
	svc, _, _, _ := newTestShareService()

	_, err := svc.GetSubscriptionShare(userID.String(), uuid.New().String())
	assertAppErrorCode(t, err, http.StatusNotFound)
}
