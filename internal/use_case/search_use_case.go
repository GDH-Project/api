package usecase

import (
	"context"
	"errors"

	"github.com/GDH-Project/api/internal/domain"
	"go.uber.org/zap"
)

type searchUseCase struct {
	log           *zap.Logger
	searchService domain.SearchService
	deviceService domain.DeviceService
}

func (uc *searchUseCase) GetDeviceInfoListByParamAndPageInfo(ctx context.Context, in *domain.DeviceInfo, p *domain.Page) ([]*domain.DeviceInfo, *domain.Page, error) {
	infoList, page, err := uc.deviceService.GetDeviceInfoListByParamAndPage(ctx, in, p)
	if err != nil {
		uc.log.Info("search.uc.GetDeviceInfoListByParamAndPageInfo() 오류",
			zap.Any("param", in),
			zap.Any("page", page),
			zap.Error(err),
		)
		return nil, nil, errors.New("검색중 오류가 발생했습니다")
	}

	return infoList, page, nil
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

func NewSearchUseCase(log *zap.Logger, searchService domain.SearchService, deviceService domain.DeviceService) domain.SearchUseCase {
	return &searchUseCase{
		log:           log,
		searchService: searchService,
		deviceService: deviceService,
	}
}
