package domain

import "context"

type SearchService interface {
	GetDeviceDataListByDeviceID(ctx context.Context, deviceID string) ([]map[string]interface{}, error)
}

type SearchUseCase interface {
	GetDeviceDataListByDeviceID(ctx context.Context, deviceID string) ([]map[string]interface{}, error)
}
