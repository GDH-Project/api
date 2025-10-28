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

func (uc *searchUseCase) GetDeviceInfoByDeviceID(ctx context.Context, deviceID string) (*domain.DeviceInfo, error) {
	data, err := uc.searchService.GetDeviceInfoByDeviceID(ctx, deviceID)
	if err != nil {
		uc.log.Info("search.uc.GetDeviceInfoByDeviceID() 오류 - 장치 ID 조회 실패", zap.Error(err))
		return nil, err
	}

	return data, nil
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
