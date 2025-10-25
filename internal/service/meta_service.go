package service

import (
	"context"

	"github.com/GDH-Project/api/internal/domain"
	"go.uber.org/zap"
)

type metaService struct {
	log *zap.Logger
	r   domain.MetaRepository
}

func (svc *metaService) GetSensorList(ctx context.Context) ([]*domain.Sensor, error) {
	return svc.r.GetSensorList(ctx)

}

func (svc *metaService) GetSensorByParam(ctx context.Context, in *domain.Sensor) (*domain.Sensor, error) {
	return svc.r.GetSensorByParam(ctx, in)
}

func (svc *metaService) GetCropList(ctx context.Context) ([]*domain.Crop, error) {
	return svc.r.GetCropList(ctx)
}

func (svc *metaService) GetCropByParam(ctx context.Context, in *domain.Crop) (*domain.Crop, error) {
	return svc.r.GetCropByParam(ctx, in)
}

func (svc *metaService) GetUpdateCycleList(ctx context.Context) ([]*domain.UpdateCycle, error) {
	return svc.r.GetUpdateCycleList(ctx)
}

func (svc *metaService) GetAddressStateList(ctx context.Context) ([]*domain.AddressState, error) {
	return svc.r.GetAddressStateList(ctx)
}

func (svc *metaService) GetAddressCityListByState(ctx context.Context, state string) ([]*domain.AddressCity, error) {
	return svc.r.GetAddressCityListByState(ctx, state)
}

func NewMetaService(log *zap.Logger, metaRepository domain.MetaRepository) domain.MetaService {
	return &metaService{
		log: log,
		r:   metaRepository,
	}
}
