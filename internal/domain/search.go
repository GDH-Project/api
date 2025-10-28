package domain

import "context"

type SearchService interface {
	GetDeviceDataListByDeviceID(ctx context.Context, deviceID string) ([]map[string]interface{}, error)
	GetDeviceInfoByDeviceID(ctx context.Context, deviceID string) (*DeviceInfo, error)
	GetRankingDeviceInfoList(ctx context.Context, filter DeviceRanking, limit int) ([]*DeviceInfo, error)
}

type SearchUseCase interface {
	GetDeviceDataListByDeviceID(ctx context.Context, deviceID string) ([]map[string]interface{}, error)
	GetDeviceInfoByDeviceID(ctx context.Context, deviceID string) (*DeviceInfo, error)
	GetDeviceInfoListByParamAndPageInfo(ctx context.Context, in *DeviceInfo, p *Page) ([]*DeviceInfo, *Page, error)
	GetRankingDeviceInfoList(ctx context.Context, filter DeviceRanking, limit int) ([]*DeviceInfo, error)
}
