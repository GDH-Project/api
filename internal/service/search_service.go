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
