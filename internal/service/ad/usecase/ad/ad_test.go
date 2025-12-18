package ad

import (
	modelad "2025_2_404/internal/service/ad/domain/ad"
	modelfullad "2025_2_404/internal/service/ad/domain/ad_full_info"
	modeluser "2025_2_404/internal/service/ad/domain/user"
	"context"
	"errors"
	"testing"
	"time"

	"github.com/google/uuid"
	"go.uber.org/zap"
)

type mockAdRepository struct {
	findByUserIDFunc      func(ctx context.Context, userID modeluser.ID) ([]modelfullad.AdFullInfo, error)
	createFunc            func(ctx context.Context, ad modelad.Ads) error
	getOneAdFunc          func(ctx context.Context, adID modelad.ID, clientID modeluser.ID) (modelfullad.AdFullInfo, error)
	updateFunc            func(ctx context.Context, ad modelad.Ads) error
	deleteFunc            func(ctx context.Context, adID modelad.ID, clientID modeluser.ID) error
	getAdDetailForSlotFunc func(ctx context.Context, id modelad.ID, click, impression int) (modelfullad.DetailID, error)
	getAdSlotFunc         func(ctx context.Context, minCost uint32) (modelad.Ads, error)
	getAdCountFunc        func(ctx context.Context, clientID modeluser.ID) (int64, error)
}

func (m *mockAdRepository) FindByUserID(ctx context.Context, userID modeluser.ID) ([]modelfullad.AdFullInfo, error) {
	if m.findByUserIDFunc != nil {
		return m.findByUserIDFunc(ctx, userID)
	}
	return nil, nil
}

func (m *mockAdRepository) Create(ctx context.Context, ad modelad.Ads) error {
	if m.createFunc != nil {
		return m.createFunc(ctx, ad)
	}
	return nil
}

func (m *mockAdRepository) GetOneAd(ctx context.Context, adID modelad.ID, clientID modeluser.ID) (modelfullad.AdFullInfo, error) {
	if m.getOneAdFunc != nil {
		return m.getOneAdFunc(ctx, adID, clientID)
	}
	return modelfullad.AdFullInfo{}, nil
}

func (m *mockAdRepository) Update(ctx context.Context, ad modelad.Ads) error {
	if m.updateFunc != nil {
		return m.updateFunc(ctx, ad)
	}
	return nil
}

func (m *mockAdRepository) Delete(ctx context.Context, adID modelad.ID, clientID modeluser.ID) error {
	if m.deleteFunc != nil {
		return m.deleteFunc(ctx, adID, clientID)
	}
	return nil
}

func (m *mockAdRepository) GetAdDetailForSlot(ctx context.Context, id modelad.ID, click, impression int) (modelfullad.DetailID, error) {
	if m.getAdDetailForSlotFunc != nil {
		return m.getAdDetailForSlotFunc(ctx, id, click, impression)
	}
	return modelfullad.DetailID{}, nil
}

func (m *mockAdRepository) GetAdSlot(ctx context.Context, minCost uint32) (modelad.Ads, error) {
	if m.getAdSlotFunc != nil {
		return m.getAdSlotFunc(ctx, minCost)
	}
	return modelad.Ads{}, nil
}

func (m *mockAdRepository) GetAdCount(ctx context.Context, clientID modeluser.ID) (int64, error) {
	if m.getAdCountFunc != nil {
		return m.getAdCountFunc(ctx, clientID)
	}
	return 0, nil
}

func TestCreate_BudgetLessThan100(t *testing.T) {
	mockRepo := &mockAdRepository{
		createFunc: func(ctx context.Context, ad modelad.Ads) error {
			if ad.Status != "non-active" {
				t.Errorf("expected status 'non-active', got %s", ad.Status)
			}
			return nil
		},
	}

	logger, _ := zap.NewDevelopment()
	useCase := New(mockRepo, logger)

	ad := modelad.Ads{
		ID:       uuid.New(),
		ClientID: uuid.New(),
		Title:    "Test Ad",
		Budget:   50,
	}

	err := useCase.Create(context.Background(), ad)
	if err != nil {
		t.Errorf("expected no error, got %v", err)
	}
}

func TestCreate_BudgetGreaterOrEqualTo100(t *testing.T) {
	mockRepo := &mockAdRepository{
		createFunc: func(ctx context.Context, ad modelad.Ads) error {
			if ad.Status != "active" {
				t.Errorf("expected status 'active', got %s", ad.Status)
			}
			return nil
		},
	}

	logger, _ := zap.NewDevelopment()
	useCase := New(mockRepo, logger)

	ad := modelad.Ads{
		ID:       uuid.New(),
		ClientID: uuid.New(),
		Title:    "Test Ad",
		Budget:   100,
	}

	err := useCase.Create(context.Background(), ad)
	if err != nil {
		t.Errorf("expected no error, got %v", err)
	}
}

func TestCreate_RepositoryError(t *testing.T) {
	expectedErr := errors.New("repository error")
	mockRepo := &mockAdRepository{
		createFunc: func(ctx context.Context, ad modelad.Ads) error {
			return expectedErr
		},
	}

	logger, _ := zap.NewDevelopment()
	useCase := New(mockRepo, logger)

	ad := modelad.Ads{
		ID:       uuid.New(),
		ClientID: uuid.New(),
		Title:    "Test Ad",
		Budget:   100,
	}

	err := useCase.Create(context.Background(), ad)
	if err != expectedErr {
		t.Errorf("expected error %v, got %v", expectedErr, err)
	}
}

func TestGetOneAd_ConversionCalculation(t *testing.T) {
	tests := []struct {
		name               string
		impressions        int
		clicks             int
		expectedConversion int
	}{
		{"no impressions", 0, 0, -1},
		{"zero clicks", 100, 0, 0},
		{"50% conversion", 100, 50, 50},
		{"100% conversion", 100, 100, 100},
		{"25% conversion", 200, 50, 25},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mockRepo := &mockAdRepository{
				getOneAdFunc: func(ctx context.Context, adID modelad.ID, clientID modeluser.ID) (modelfullad.AdFullInfo, error) {
					return modelfullad.AdFullInfo{
						ID:          adID,
						Impressions: tt.impressions,
						Clicks:      tt.clicks,
						Budget:      100,
					}, nil
				},
			}

			logger, _ := zap.NewDevelopment()
			useCase := New(mockRepo, logger)

			adID := uuid.New()
			clientID := uuid.New()

			adInfo, conversion, err := useCase.GetOneAd(context.Background(), adID, clientID)
			if err != nil {
				t.Errorf("expected no error, got %v", err)
			}

			if conversion != tt.expectedConversion {
				t.Errorf("expected conversion %d, got %d", tt.expectedConversion, conversion)
			}

			if adInfo.Impressions != tt.impressions {
				t.Errorf("expected impressions %d, got %d", tt.impressions, adInfo.Impressions)
			}

			if adInfo.Clicks != tt.clicks {
				t.Errorf("expected clicks %d, got %d", tt.clicks, adInfo.Clicks)
			}
		})
	}
}

func TestGetOneAd_StatusUpdate(t *testing.T) {
	tests := []struct {
		name           string
		budget         uint32
		expectedStatus string
	}{
		{"budget below 100", 50, "non-active"},
		{"budget at 100", 100, "active"},
		{"budget above 100", 200, "active"},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mockRepo := &mockAdRepository{
				getOneAdFunc: func(ctx context.Context, adID modelad.ID, clientID modeluser.ID) (modelfullad.AdFullInfo, error) {
					return modelfullad.AdFullInfo{
						ID:          adID,
						Budget:      tt.budget,
						Status:      "active",
						Impressions: 100,
						Clicks:      50,
					}, nil
				},
			}

			logger, _ := zap.NewDevelopment()
			useCase := New(mockRepo, logger)

			adID := uuid.New()
			clientID := uuid.New()

			adInfo, _, err := useCase.GetOneAd(context.Background(), adID, clientID)
			if err != nil {
				t.Errorf("expected no error, got %v", err)
			}

			if adInfo.Status != tt.expectedStatus {
				t.Errorf("expected status %s, got %s", tt.expectedStatus, adInfo.Status)
			}
		})
	}
}

func TestGetAdDetailForSlot_Impression(t *testing.T) {
	mockRepo := &mockAdRepository{
		getAdDetailForSlotFunc: func(ctx context.Context, id modelad.ID, click, impression int) (modelfullad.DetailID, error) {
			if impression != 1 || click != 0 {
				t.Errorf("expected impression=1, click=0, got impression=%d, click=%d", impression, click)
			}
			return modelfullad.DetailID{}, nil
		},
	}

	logger, _ := zap.NewDevelopment()
	useCase := New(mockRepo, logger)

	_, err := useCase.GetAdDetailForSlot(context.Background(), uuid.New(), "impression")
	if err != nil {
		t.Errorf("expected no error, got %v", err)
	}
}

func TestGetAdDetailForSlot_Click(t *testing.T) {
	mockRepo := &mockAdRepository{
		getAdDetailForSlotFunc: func(ctx context.Context, id modelad.ID, click, impression int) (modelfullad.DetailID, error) {
			if click != 1 || impression != 0 {
				t.Errorf("expected click=1, impression=0, got click=%d, impression=%d", click, impression)
			}
			return modelfullad.DetailID{}, nil
		},
	}

	logger, _ := zap.NewDevelopment()
	useCase := New(mockRepo, logger)

	_, err := useCase.GetAdDetailForSlot(context.Background(), uuid.New(), "click")
	if err != nil {
		t.Errorf("expected no error, got %v", err)
	}
}

func TestUpdate_Success(t *testing.T) {
	mockRepo := &mockAdRepository{
		updateFunc: func(ctx context.Context, ad modelad.Ads) error {
			return nil
		},
	}

	logger, _ := zap.NewDevelopment()
	useCase := New(mockRepo, logger)

	ad := modelad.Ads{
		ID:       uuid.New(),
		ClientID: uuid.New(),
		Title:    "Updated Ad",
		Budget:   150,
		Status:   "active",
		StartAt:  time.Now(),
		EndAt:    time.Now().Add(24 * time.Hour),
	}

	err := useCase.Update(context.Background(), ad)
	if err != nil {
		t.Errorf("expected no error, got %v", err)
	}
}

func TestDelete_Success(t *testing.T) {
	mockRepo := &mockAdRepository{
		deleteFunc: func(ctx context.Context, adID modelad.ID, clientID modeluser.ID) error {
			return nil
		},
	}

	logger, _ := zap.NewDevelopment()
	useCase := New(mockRepo, logger)

	err := useCase.Delete(context.Background(), uuid.New(), uuid.New())
	if err != nil {
		t.Errorf("expected no error, got %v", err)
	}
}

func TestFindByUserID_Success(t *testing.T) {
	expectedAds := []modelfullad.AdFullInfo{
		{ID: uuid.New(), Title: "Ad 1"},
		{ID: uuid.New(), Title: "Ad 2"},
	}

	mockRepo := &mockAdRepository{
		findByUserIDFunc: func(ctx context.Context, userID modeluser.ID) ([]modelfullad.AdFullInfo, error) {
			return expectedAds, nil
		},
	}

	logger, _ := zap.NewDevelopment()
	useCase := New(mockRepo, logger)

	ads, err := useCase.FindByUserID(context.Background(), uuid.New())
	if err != nil {
		t.Errorf("expected no error, got %v", err)
	}

	if len(ads) != len(expectedAds) {
		t.Errorf("expected %d ads, got %d", len(expectedAds), len(ads))
	}
}

func TestGetAdSlot_Success(t *testing.T) {
	expectedAd := modelad.Ads{
		ID:     uuid.New(),
		Title:  "Slot Ad",
		Budget: 150,
	}

	mockRepo := &mockAdRepository{
		getAdSlotFunc: func(ctx context.Context, minCost uint32) (modelad.Ads, error) {
			if minCost != 100 {
				t.Errorf("expected minCost 100, got %d", minCost)
			}
			return expectedAd, nil
		},
	}

	logger, _ := zap.NewDevelopment()
	useCase := New(mockRepo, logger)

	ad, err := useCase.GetAdSlot(context.Background(), 100)
	if err != nil {
		t.Errorf("expected no error, got %v", err)
	}

	if ad.ID != expectedAd.ID {
		t.Errorf("expected ad ID %v, got %v", expectedAd.ID, ad.ID)
	}
}

func TestGetAdCount_Success(t *testing.T) {
	expectedCount := int64(42)

	mockRepo := &mockAdRepository{
		getAdCountFunc: func(ctx context.Context, clientID modeluser.ID) (int64, error) {
			return expectedCount, nil
		},
	}

	logger, _ := zap.NewDevelopment()
	useCase := New(mockRepo, logger)

	count, err := useCase.GetAdCount(context.Background(), uuid.New())
	if err != nil {
		t.Errorf("expected no error, got %v", err)
	}

	if count != expectedCount {
		t.Errorf("expected count %d, got %d", expectedCount, count)
	}
}
