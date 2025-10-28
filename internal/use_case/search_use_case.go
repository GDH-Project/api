package usecase

import (
	"context"

	"github.com/GDH-Project/api/internal/domain"
	"go.uber.org/zap"
)

type searchUseCase struct {
	log           *zap.Logger
	searchService domain.SearchService
}

func (uc *searchUseCase) GetDeviceDataListByDeviceID(ctx context.Context, deviceID string) ([]map[string]interface{}, error) {
	return uc.searchService.GetDeviceDataListByDeviceID(ctx, deviceID)
}

func NewSearchUseCase(log *zap.Logger, searchService domain.SearchService) domain.SearchUseCase {
	return &searchUseCase{
		log:           log,
		searchService: searchService,
	}
}
