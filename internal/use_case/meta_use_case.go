package usecase

import (
	"context"

	"github.com/GDH-Project/api/internal/domain"
	"go.uber.org/zap"
)

type metaUseCase struct {
	log *zap.Logger
	svc domain.MetaService
}

func (uc *metaUseCase) GetSensorList(ctx context.Context) ([]*domain.Sensor, error) {
	return uc.svc.GetSensorList(ctx)
}

func (uc *metaUseCase) GetSensorByParam(ctx context.Context, in *domain.Sensor) (*domain.Sensor, error) {
	return uc.svc.GetSensorByParam(ctx, in)
}

func (uc *metaUseCase) GetCropList(ctx context.Context) ([]*domain.Crop, error) {
	return uc.svc.GetCropList(ctx)
}

func (uc *metaUseCase) GetCropByParam(ctx context.Context, in *domain.Crop) (*domain.Crop, error) {
	return uc.svc.GetCropByParam(ctx, in)
}

func (uc *metaUseCase) GetUpdateCycleList(ctx context.Context) ([]*domain.UpdateCycle, error) {
	return uc.svc.GetUpdateCycleList(ctx)
}

func (uc *metaUseCase) GetAddressStateList(ctx context.Context) ([]*domain.AddressState, error) {
	return uc.svc.GetAddressStateList(ctx)
}

func (uc *metaUseCase) GetAddressCityListByState(ctx context.Context, state string) ([]*domain.AddressCity, error) {
	return uc.svc.GetAddressCityListByState(ctx, state)
}

func NewMetaUseCase(log *zap.Logger, metaService domain.MetaService) domain.MetaUseCase {
	return &metaUseCase{
		log: log,
		svc: metaService,
	}
}
