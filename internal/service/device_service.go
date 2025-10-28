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
}

func (svc *deviceService) GetDeviceApiKeyByApiKey(ctx context.Context, apiKey string) (string, error) {
	deviceID, err := svc.device.GetDeviceApiKeyByApiKey(ctx, apiKey)
	if err != nil {
		svc.log.Info("device.svc.GetDeviceApiKeyByID() 오류 - API키가 일치하지 않습니다",
			zap.String("apiKey", apiKey),
			zap.Error(err))
		return "", errors.New("잘못된 접근입니다")
	}

	return deviceID, nil
}

func (svc *deviceService) CreateDeviceApiKey(ctx context.Context, in *domain.RawApiKey) (*domain.ApiKey, error) {

	data, err := svc.device.CreateDeviceApiKey(ctx, in)
	if err != nil {
		return nil, err
	}

	apiKey := &domain.ApiKey{
		ID:       data.ID,
		Key:      data.APIKey,
		DeviceID: data.DeviceID,
		Title:    data.Title,
		Desc:     data.Desc.String,
	}

	return apiKey, nil
}

func (svc *deviceService) CreateDeviceReqeustSchema(ctx context.Context, in *domain.RawDeviceRequestSchema) error {
	if err := svc.device.WithTransaction(ctx, func(tx pgx.Tx) error {
		if err := svc.device.CreateDeviceReqeustSchemaListTx(ctx, tx, in.DeviceID, []*domain.RawDeviceRequestSchema{in}); err != nil {
			return err
		}
		return nil
	}); err != nil {
		svc.log.Info("device.svc.CreateDeviceReqeustSchema() 오류", zap.Error(err))
		return errors.New("장치 요청 스키마 생성에 실패했습니다")
	}

	return nil
}

func (svc *deviceService) DeleteDeviceInfoByID(ctx context.Context, deviceID string, userID string) error {
	if err := svc.device.DeleteDeviceInfoByID(ctx, deviceID, userID); err != nil {
		svc.log.Info("device.svc.DeleteDeviceInfoByID() 오류", zap.Error(err))
		return errors.New("장비 데이터 제거에 실패했습니다")
	}
	return nil
}

// 장비 접근 권한 체크
func (svc *deviceService) checkDeviceAccessState(ctx context.Context, deviceID, userID string) (*domain.DeviceInfo, error) {
	deviceInfo, err := svc.device.GetDeviceInfoByID(ctx, deviceID)
	if err != nil {
		svc.log.Info("device.svc.checkDeviceAccessState() 오류 - 존재 하지 않는 장비 ID 입니다",
			zap.String("id", deviceID),
			zap.Error(err),
		)
		return nil, errors.New("존재하지 않는 장비 ID 입니다")
	}
	if deviceInfo.UserID != userID {
		err := errors.New("장비 데이터 접근 권한이 없습니다")
		svc.log.Info("device.svc.checkDeviceAccessState() 오류",
			zap.String("id", deviceID),
			zap.String("userID", userID),
			zap.Error(err),
		)
		return nil, err
	}
	return deviceInfo, nil
}

func (svc *deviceService) UpdateDeviceReqeustSchemaByID(ctx context.Context, in *domain.RawDeviceRequestSchema, userID string) error {
	// 장비 접근 권한 체크
	_, err := svc.checkDeviceAccessState(ctx, in.DeviceID, userID)
	if err != nil {
		svc.log.Info("device.svc.UpdateDeviceReqeustSchemaByID() 오류 - 장비 접근 권한이 없습니다.")
		return err
	}

	if err := svc.device.UpdateDeviceRequestSchema(ctx, in); err != nil {
		svc.log.Info("device.svc.UpdateDeviceReqeustSchemaByID() 오류 - 장비 요청 스키마 업데이트중 오류가 발생했습니다.",
			zap.Any("data", in),
			zap.Error(err),
		)
		return errors.New("장비 요청 스키마 업데이트중 오류가 발생했습니다")
	}

	return nil
}

func (svc *deviceService) GetDeviceReqeustSchemaListByID(ctx context.Context, deviceID string, userID string) ([]*domain.DeviceRequestSchema, error) {
	// 장비 접근 권한 체크
	_, err := svc.checkDeviceAccessState(ctx, deviceID, userID)
	if err != nil {
		svc.log.Info("device.svc.UpdateDeviceInfo() 오류 - 장비 접근 권한이 없습니다.")
		return nil, err
	}

	// --- 권한 체크 완료 ---

	// --- 데이터 조회 ---
	list, err := svc.device.GetDeviceRequestSchemaListByDeviceID(ctx, deviceID)
	if err != nil {
		svc.log.Info("device.svc.UpdateDeviceInfo() 오류",
			zap.String("id", deviceID),
			zap.Error(err),
		)
		return nil, errors.New("장비 요청 스키마 데이터를 불러올 수 없습니다")
	}

	return list, nil
}

func (svc *deviceService) UpdateDeviceInfo(ctx context.Context, in *domain.RawDeviceInfo) error {

	if err := svc.device.UpdateDeviceInfo(ctx, in); err != nil {
		svc.log.Info("device.svc.UpdateDeviceInfo() 오류",
			zap.Any("data", in),
			zap.Error(err),
		)
		return errors.New("장비 정보를 업데이트 하는 도중 오류가 발생했습니다")
	}

	return nil
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
