package service

import (
	"context"
	"errors"

	"github.com/GDH-Project/api/internal/domain"
	"github.com/jackc/pgx/v5"
	"go.uber.org/zap"
)

type deviceService struct {
	log    *zap.Logger
	device domain.DeviceRepository
	meta   domain.MetaRepository
}

func (svc *deviceService) GetDeviceInfoByID(ctx context.Context, id string, userID string) (*domain.DeviceInfo, error) {
	data, err := svc.device.GetDeviceInfoByID(ctx, id)
	if err != nil {
		return nil, err
	}
	if data.UserID != userID {
		err := errors.New("접근 권한이 없습니다")
		svc.log.Warn("device.svc.GetDeviceInfoByID() 권한이 없는 유저가 조회를 시도했습니다.",
			zap.String("user_id", userID),
			zap.String("id", id),
			zap.Error(err),
		)
		return nil, err
	}

	return data, nil
}

func (svc *deviceService) GetDeviceInfoListByParamAndPage(ctx context.Context, in *domain.DeviceInfo, page *domain.Page) ([]*domain.DeviceInfo, *domain.Page, error) {
	return svc.device.GetDeviceInfoListByParamAndPage(ctx, in, page)
}

func (svc *deviceService) CreateDevice(ctx context.Context, deviceInfoData *domain.RawDeviceInfo, deviceSchemaDataList []*domain.RawDeviceRequestSchema) error {
	// 장비 생성 트랜잭션 시작
	err := svc.device.WithTransaction(ctx, func(tx pgx.Tx) error {
		deviceID, err := svc.device.CreateDeviceInfoTx(ctx, tx, deviceInfoData)
		if err != nil {
			return err
		}

		// 스키마 데이터가 존재하지 않을 경우
		if len(deviceSchemaDataList) == 0 {
			svc.log.Info("device.svc.CreateDeviceInfo() 스키마 생성 스킵",
				zap.String("device_id", deviceID),
			)
			return nil
		}

		// 스키마 데이터가 존재하는 경우
		// DB에 스키마 삽입
		if err := svc.device.CreateDeviceReqeustSchemaListTx(ctx, tx, deviceID, deviceSchemaDataList); err != nil {
			svc.log.Info("device.svc.CreateDeviceInfo() 오류 - 스키마를 생성할 수 없습니다.", zap.Error(err))
			return err
		}

		return nil
	})
	// 장비 생성 트랜잭션 종료

	if err != nil {
		svc.log.Info("device.svc.CreateDeviceInfo() 오류 - 트랜잭션 도중 오류가 발생했습니다..", zap.Error(err))
		return err
	}

	return nil
}

func NewDeviceService(log *zap.Logger, deviceRepository domain.DeviceRepository) domain.DeviceService {
	return &deviceService{
		log:    log,
		device: deviceRepository,
	}
}
