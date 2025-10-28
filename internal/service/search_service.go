package service

import (
	"context"
	"errors"

	"github.com/GDH-Project/api/internal/domain"
	"go.uber.org/zap"
)

type searchService struct {
	device domain.DeviceRepository
	log    *zap.Logger
}

func (svc *searchService) GetRankingDeviceInfoList(ctx context.Context, filter domain.DeviceRanking, limit int) ([]*domain.DeviceInfo, error) {
	data, err := svc.device.GetRankingDeviceInfoList(ctx, filter, limit)
	if err != nil {
		svc.log.Error("search.svc.GetRankingDeviceInfoList", zap.Error(err))
		return nil, errors.New("랭킹 정보 로드중 오류가 발생했습니다")
	}

	return data, nil
}

func (svc *searchService) GetDeviceInfoByDeviceID(ctx context.Context, deviceID string) (*domain.DeviceInfo, error) {
	data, err := svc.device.GetDeviceInfoByID(ctx, deviceID)
	if err != nil {
		svc.log.Info("search.svc.GetDeviceInfoByDeviceID() 오류 - 장치 조회 실패", zap.Error(err))
		return nil, errors.New("장치 ID를 확인해 주세요")
	}

	return data, nil
}

func (svc *searchService) GetDeviceDataListByDeviceID(ctx context.Context, deviceID string) ([]map[string]interface{}, error) {
	dataList, err := svc.device.GetDeviceDataListByDeviceID(ctx, deviceID)
	if err != nil {
		svc.log.Info("search.svc.GetDeviceDataListByDeviceID() 오류 - 장치 조회 실패", zap.Error(err))
		return nil, errors.New("장치 조회 실패")
	}

	var jsonArr []map[string]interface{}

	for _, data := range dataList {
		data.Data["time"] = data.Time
		jsonArr = append(jsonArr, data.Data)
	}

	return jsonArr, nil
}

func NewSearchService(log *zap.Logger, deviceRepository domain.DeviceRepository) domain.SearchService {
	return &searchService{
		log:    log,
		device: deviceRepository,
	}
}
