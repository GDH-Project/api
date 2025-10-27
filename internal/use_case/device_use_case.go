package usecase

import (
	"context"
	"errors"

	"github.com/GDH-Project/api/internal/domain"
	"github.com/GDH-Project/api/internal/util"
	"go.uber.org/zap"
)

type deviceUseCase struct {
	log       *zap.Logger
	deviceSvc domain.DeviceService
	metaSvc   domain.MetaService
}

func (uc *deviceUseCase) UpdateDeviceReqeustSchemaByID(ctx context.Context, in *domain.DeviceRequestSchema, deviceID string) error {
	validate := util.ValidateUser{
		Ctx:        ctx,
		TargetRole: domain.UserRoleDevice,
	}
	if !validate.Exec() {
		err := errors.New("권한이 존재하지 않습니다")
		uc.log.Info("device.uc.UpdateDeviceReqeustSchemaByID() 오류 - 권한이 존재하지 않습니다.", zap.Error(err))
		return err
	}

	param := &domain.RawDeviceRequestSchema{
		ID:       in.ID,
		DeviceID: deviceID,
		Key:      in.Key,
	}
	sensorData, err := uc.metaSvc.GetSensorByParam(ctx, &domain.Sensor{Title: in.Target})
	if err != nil {
		uc.log.Info("device.uc.UpdateDeviceReqeustSchemaByID() 오류 - 센서 정보를 받아올 수 없습니다.", zap.Error(err))
		return errors.New("존재하지 않는 센서입니다")
	}
	param.TargetSensorID = sensorData.ID

	if err := uc.deviceSvc.UpdateDeviceReqeustSchemaByID(ctx, param, validate.UserID()); err != nil {
		uc.log.Info("device.uc.UpdateDeviceReqeustSchemaByID() 오류", zap.Error(err))
		return err
	}

	return nil

}

func (uc *deviceUseCase) GetDeviceReqeustSchemaListByID(ctx context.Context, deviceID string) ([]*domain.DeviceRequestSchema, error) {
	validate := util.ValidateUser{
		Ctx:        ctx,
		TargetRole: domain.UserRoleDevice,
	}
	if !validate.Exec() {
		err := errors.New("권한이 존재하지 않습니다")
		uc.log.Info("device.uc.GetDeviceReqeustSchemaListByID() 오류 - 권한이 존재하지 않습니다.", zap.Error(err))
		return nil, err
	}

	list, err := uc.deviceSvc.GetDeviceReqeustSchemaListByID(ctx, deviceID, validate.UserID())
	if err != nil {
		uc.log.Info("device.uc.GetDeviceReqeustSchemaListByID() 오류", zap.Error(err))
		return nil, errors.New("장비 스키마 정보를 불러올 수 없습니다")
	}

	return list, nil
}

func (uc *deviceUseCase) UpdateDeviceInfo(ctx context.Context, in *domain.DeviceInfo) error {
	validate := util.ValidateUser{
		Ctx:        ctx,
		TargetRole: domain.UserRoleDevice,
	}
	if !validate.Exec() {
		err := errors.New("권한이 존재하지 않습니다")
		uc.log.Info("device.uc.UpdateDeviceInfo() 오류 - 권한이 존재하지 않습니다.", zap.Error(err))
		return err
	}

	rawInfo := &domain.RawDeviceInfo{
		ID:     in.ID,
		UserID: validate.UserID(),
		Title:  in.Title,
		Name:   *in.Name,
	}

	if in.UpdateCycle > 0 {
		updateCycleList, err := uc.metaSvc.GetUpdateCycleList(ctx)
		if err != nil {
			uc.log.Info("device.uc.UpdateDeviceInfo() 오류", zap.Error(err))
			return errors.New("갱신 주기를 받아올 수 없습니다")
		}

		for _, item := range updateCycleList {
			if item.Interval == in.UpdateCycle {
				rawInfo.UpdateCycleID = item.ID
				break
			}
		}

		if rawInfo.UpdateCycleID == 0 {
			err := errors.New("갱신 주기가 유효하지 않습니다")
			uc.log.Info("device.uc.UpdateDeviceInfo() 오류", zap.Error(err))
			return err
		}
	}

	if err := uc.deviceSvc.UpdateDeviceInfo(ctx, rawInfo); err != nil {
		uc.log.Info("device.uc.UpdateDeviceInfo() 오류", zap.Error(err))
		return errors.New("유저 정보를 업데이트 하는 도중 오류가 발생했습니다")
	}

	return nil
}

func (uc *deviceUseCase) GetDeviceInfoByID(ctx context.Context, id string) (*domain.DeviceInfo, error) {
	validate := util.ValidateUser{
		Ctx:        ctx,
		TargetRole: domain.UserRoleDevice,
	}
	if !validate.Exec() {
		err := errors.New("권한이 존재하지 않습니다")
		uc.log.Info("device.uc.GetDeviceInfoByID() 오류 - 권한이 존재하지 않습니다.", zap.Error(err))
		return nil, err
	}

	data, err := uc.deviceSvc.GetDeviceInfoByID(ctx, id, validate.UserID())
	if err != nil {
		return nil, errors.New("id를 확인해주세요")
	}

	return data, nil
}

func (uc *deviceUseCase) GetDeviceInfoListByParamAndPage(ctx context.Context, in *domain.DeviceInfo, page *domain.Page) ([]*domain.DeviceInfo, *domain.Page, error) {
	// 사용자 권한 체크
	validate := util.ValidateUser{
		Ctx:        ctx,
		TargetRole: domain.UserRoleDevice,
	}
	if !validate.Exec() {
		err := errors.New("권한이 존재하지 않습니다")
		uc.log.Info("device.uc.CreateDeviceInfo() 오류 - 권한이 존재하지 않습니다.", zap.Error(err))
		return nil, nil, err
	}

	// 유저 아이디 삽입
	in.UserID = validate.UserID()

	dataList, page, err := uc.deviceSvc.GetDeviceInfoListByParamAndPage(ctx, in, page)
	if err != nil {
		uc.log.Info("device.uc.GetDeviceInfoListByParamAndPage()", zap.Error(err))
		return nil, nil, errors.New("장비 조회 실패")
	}

	return dataList, page, nil
}

func (uc *deviceUseCase) CreateDevice(ctx context.Context, deviceInfoData *domain.DeviceInfo, deviceSchemaDataList []*domain.DeviceRequestSchema) error {
	// 사용자 권한 체크
	validateUser := util.ValidateUser{
		Ctx:        ctx,
		TargetRole: domain.UserRoleDevice,
	}

	if !validateUser.Exec() {
		err := errors.New("권한이 존재하지 않습니다")
		uc.log.Info("device.uc.CreateDeviceInfo() 오류 - 권한이 존재하지 않습니다.", zap.Error(err))
		return err
	}

	var rawDeviceInfo domain.RawDeviceInfo
	var rawDeviceRequestSchemaList []*domain.RawDeviceRequestSchema

	rawDeviceInfo.UserID = validateUser.UserID()
	rawDeviceInfo.Title = deviceInfoData.Title
	rawDeviceInfo.Name = *deviceInfoData.Name

	// 작물 정보 ID
	cropData, err := uc.metaSvc.GetCropByParam(ctx, &domain.Crop{Title: deviceInfoData.Crop})
	if err != nil {
		uc.log.Info("device.uc.CreateDeviceInfo() 오류 - 작물 정보를 받아올 수 없습니다.", zap.Error(err))
		return errors.New("존재하지 않는 작물입니다")
	}
	rawDeviceInfo.CropID = cropData.ID

	// 갱신 주기 ID
	updateCycleData, err := uc.metaSvc.GetUpdateCycleList(ctx)
	if err != nil {
		uc.log.Info("device.uc.CreateDeviceInfo() 오류 - 갱신주기를 받아올 수 없습니다.", zap.Error(err))
		return errors.New("갱신주기 조회중 오류가 발생했습니다")
	}
	for _, data := range updateCycleData {
		if data.Interval == deviceInfoData.UpdateCycle {
			rawDeviceInfo.UpdateCycleID = data.ID
			break
		}
	}
	if rawDeviceInfo.UpdateCycleID == 0 {
		err := errors.New("존재하지 않는 갱신주기 입니다")
		uc.log.Info("device.uc.CreateDeviceInfo() 오류 - 존재하지 않는 갱신 주기 입니다.", zap.Error(err))
		return err
	}

	// 주소 ID
	addressData, err := uc.metaSvc.GetAddressIDByStateTitleAndCityTitle(ctx, &domain.AddressCity{
		StateTitle: deviceInfoData.Address.State,
		Title:      deviceInfoData.Address.City,
	})
	if err != nil {
		uc.log.Info("device.uc.CreateDeviceInfo() 오류 - 존재하지 않는 주소 입니다.", zap.Error(err))
		return errors.New("잘못된 주소 입니다")
	}
	rawDeviceInfo.AddressStateID = addressData.StateID
	rawDeviceInfo.AddressCityID = addressData.CityID

	// 스키마 배열이 존재한다면 센서 정보를 바인딩 해서 Raw스키마 정보로 변환한다.
	if len(deviceSchemaDataList) > 0 {
		// 스키마 삽입 배열 생성

		// 센서 정보
		sensorList, err := uc.metaSvc.GetSensorList(ctx)
		if err != nil {
			uc.log.Info("device.uc.CreateDeviceInfo() 오류 - 센서 정보를 받아올 수 없습니다.", zap.Error(err))
			return err
		}
		for _, item := range deviceSchemaDataList {
			temp := &domain.RawDeviceRequestSchema{
				Key: item.Key,
			}
			for _, sensor := range sensorList {
				if sensor.Title == item.Target {
					temp.TargetSensorID = sensor.ID
					break
				}
			}
			if temp.TargetSensorID == 0 {
				uc.log.Info("device.uc.CreateDeviceInfo() 오류 - 존재하지 않는 센서 이름 입니다.", zap.Error(err))
				return errors.New("존재하지 않는 센서 입니다")
			}

			rawDeviceRequestSchemaList = append(rawDeviceRequestSchemaList, temp)
		}
	}

	if err := uc.deviceSvc.CreateDevice(ctx, &rawDeviceInfo, rawDeviceRequestSchemaList); err != nil {
		uc.log.Info("device.uc.CreateDeviceInfo() 오류 - 장비 생성중 오류가 발생했습니다.",
			zap.Any("deviceData", rawDeviceInfo),
			zap.Any("schemaDataList", rawDeviceRequestSchemaList),
			zap.Error(err),
		)
		return errors.New("장비 생성중 오류가 발생했습니다")
	}

	return nil
}

func NewDeviceUseCase(log *zap.Logger, deviceService domain.DeviceService, metaService domain.MetaService) domain.DeviceUseCase {
	return &deviceUseCase{
		log:       log,
		deviceSvc: deviceService,
		metaSvc:   metaService,
	}
}
