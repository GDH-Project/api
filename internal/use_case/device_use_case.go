package usecase

import (
	"context"
	"crypto/rand"
	"database/sql"
	"encoding/base64"
	"errors"
	"fmt"
	"strings"

	"github.com/GDH-Project/api/internal/domain"
	"github.com/GDH-Project/api/internal/util"
	"go.uber.org/zap"
)

type deviceUseCase struct {
	log       *zap.Logger
	deviceSvc domain.DeviceService
	metaSvc   domain.MetaService
}

func (uc *deviceUseCase) CreateDeviceDataWithApiKey(ctx context.Context, apiKey string, in map[string]interface{}) error {
	// API KEY에 해당하는 장비 ID 추출
	deviceId, err := uc.deviceSvc.GetDeviceApiKeyByApiKey(ctx, apiKey)
	if err != nil {
		return err
	}

	// 장비 ID에 해당하는 스키마 추출
	schemaList, err := uc.deviceSvc.GetDeviceRequestSchemaListByID(ctx, deviceId)
	if err != nil {
		uc.log.Info("device.uc.CreateDeviceDataWithApiKey()오류 - 스키마 정보를 받아올 수 없습니다.", zap.Error(err))
		return err
	}

	// 스키마에 해당하는 json string 파싱
	var jsonKeyValue []string
	for _, schema := range schemaList {
		data, ok := in[schema.Key]
		if ok {
			parseData, idParsed := data.(float64)
			if idParsed {
				temp := fmt.Sprintf("\"%s\":%s", schema.Target, util.FormatFloat(parseData, 2))
				jsonKeyValue = append(jsonKeyValue, temp)
			}
		}
	}
	if len(jsonKeyValue) == 0 {
		uc.log.Info("device.uc.CreateDeviceDataWithApiKey() 오류 - 삽입할 데이터가 없음")
		return errors.New("스키마와 연결된 데이터가 존재하지 않습니다")
	}
	jsonStr := fmt.Sprintf("{%s}", strings.Join(jsonKeyValue, ","))

	uc.log.Debug("debug api",
		zap.Any("schemaList", schemaList),
		zap.Any("deviceId", deviceId),
		zap.Any("in", in),
		zap.Any("jsonKeyValue", jsonKeyValue),
		zap.Any("jsonStr", jsonStr),
	)

	if err := uc.deviceSvc.CreateDeviceDataWithDeviceID(ctx, deviceId, jsonStr); err != nil {
		uc.log.Info("device.uc.CreateDeviceDataWithApiKey() 오류", zap.Error(err))
		return err
	}
	return nil
}

func (uc *deviceUseCase) DeleteDeviceApiKeyByUserIDAndDeviceID(ctx context.Context, deviceID, id string) error {
	// 유저 권한 체크
	validate, err := uc.validateUser(ctx, domain.UserRoleDevice)
	if err != nil {
		return err
	}
	userID := validate.UserID()

	// 장치와 유저 연결성 체크는 불필요
	return uc.deviceSvc.DeleteDeviceApiKeyByUserIDAndDeviceID(ctx, userID, deviceID, id)
}

func (uc *deviceUseCase) GetDeviceApiKeyListByUserIDAndDeviceID(ctx context.Context, deviceID string) ([]*domain.ApiKey, error) {
	// 유저 권한 체크
	validate, err := uc.validateUser(ctx, domain.UserRoleDevice)
	if err != nil {
		return nil, err
	}
	userID := validate.UserID()
	deviceInfo, err := uc.deviceSvc.GetDeviceInfoByID(ctx, deviceID, userID)
	if err != nil {
		uc.log.Info("device.uc.GetDeviceApiKeyListByUserIDAndDeviceID() 오류 - 장치 정보를 불러올 수 없습니다.",
			zap.String("userId", userID),
			zap.String("deviceId", deviceID),
			zap.Error(err))
		return nil, err
	}

	dataList, err := uc.deviceSvc.GetDeviceApiKeyListByUserIDAndDeviceID(ctx, userID, deviceInfo.ID)
	if err != nil {
		uc.log.Info("device.uc.GetDeviceApiKeyListByUserIDAndDeviceID() 오류", zap.Error(err))
		return nil, err
	}

	var apiKeys []*domain.ApiKey
	for _, item := range dataList {
		apiKeys = append(apiKeys, &domain.ApiKey{
			ID:        item.ID,
			DeviceID:  item.DeviceID,
			Title:     item.Title,
			Desc:      item.Desc.String,
			CreatedAt: item.CreatedAt,
		})
	}

	return apiKeys, nil
}

// GetDeviceApiKeyByApiKey 권한 검사가 존재하지 않기 때문에 조심할 것
func (uc *deviceUseCase) GetDeviceApiKeyByApiKey(ctx context.Context, apiKey string) (string, error) {
	return uc.deviceSvc.GetDeviceApiKeyByApiKey(ctx, apiKey)
}

func (uc *deviceUseCase) CreateDeviceApiKey(ctx context.Context, in *domain.ApiKey) (*domain.ApiKey, error) {
	// 유저 권한 확인
	validate, err := uc.validateUser(ctx, domain.UserRoleDevice)
	if err != nil {
		return nil, err
	}

	// 장비 접근 권한 확인
	_, err = uc.deviceSvc.GetDeviceInfoByID(ctx, in.DeviceID, validate.UserID())
	if err != nil {
		uc.log.Info("device.uc.GetDeviceInfoByID() 오류 - 장치 정보를 불러올 수 없습니다.", zap.Error(err))
		return nil, errors.New("장치 정보를 불러올 수 없습니다")
	}

	// 32바이트 문자열 생성
	b := make([]byte, 24) // 24 -> base64 인코딩시 32자리
	if _, err := rand.Read(b); err != nil {
		uc.log.Info("device.uc.GetDeviceInfoByID() 오류 - 32바이트 문자열 생성 실패", zap.Error(err))
		return nil, err
	}
	key := base64.URLEncoding.EncodeToString(b)

	param := &domain.RawApiKey{
		APIKey:   key,
		UserID:   validate.UserID(),
		DeviceID: in.DeviceID,
		Title:    in.Title,
		Desc: sql.NullString{
			String: in.Desc,
			Valid:  in.Desc != "",
		},
	}

	data, err := uc.deviceSvc.CreateDeviceApiKey(ctx, param)
	if err != nil {
		uc.log.Info("device.uc.GetDeviceInfoByID() 오류 - API키 생성 실패",
			zap.Any("data", param),
			zap.Error(err),
		)
		return nil, err
	}

	return data, nil
}

func (uc *deviceUseCase) CreateDeviceRequestSchema(ctx context.Context, deviceID string, in *domain.DeviceRequestSchema) error {
	validate, err := uc.validateUser(ctx, domain.UserRoleDevice)
	if err != nil {
		return err
	}

	_, err = uc.deviceSvc.GetDeviceInfoByID(ctx, deviceID, validate.UserID())
	if err != nil {
		uc.log.Info("device.uc.CreateDeviceRequestSchema() 오류 - 장치 정보를 불러올 수 없습니다.", zap.Error(err))
		return errors.New("장치 정보를 불러올 수 없습니다")
	}

	sensorInfo, err := uc.metaSvc.GetSensorByParam(ctx, &domain.Sensor{Title: in.Target})
	if err != nil {
		uc.log.Info("device.uc.CreateDeviceRequestSchema() 오류 - 센서 정보를 불러올 수 없습니다.", zap.Error(err))
		return errors.New("센서 정보를 불러올 수 없습니다")
	}

	param := &domain.RawDeviceRequestSchema{
		DeviceID:       deviceID,
		Key:            in.Key,
		TargetSensorID: sensorInfo.ID,
	}

	if err := uc.deviceSvc.CreateDeviceRequestSchema(ctx, param); err != nil {
		uc.log.Info("device.uc.CreateDeviceRequestSchema() 오류 - 장치 요청 데이터를 생성할 수 없습니다..", zap.Error(err))
		return err
	}

	return nil
}

func (uc *deviceUseCase) validateUser(ctx context.Context, role domain.UserRole) (*util.ValidateUser, error) {
	validate := &util.ValidateUser{
		Ctx:        ctx,
		TargetRole: domain.UserRoleDevice,
	}
	if !validate.Exec() {
		err := errors.New("권한이 존재하지 않습니다")
		uc.log.Info("device.uc.validateUser() 오류 - 권한이 존재하지 않습니다.", zap.Error(err))
		return nil, err
	}
	return validate, nil
}
func (uc *deviceUseCase) DeleteDeviceInfoByID(ctx context.Context, deviceID string) error {
	validate, err := uc.validateUser(ctx, domain.UserRoleDevice)
	if err != nil {
		return err
	}

	userID := validate.UserID()

	_, err = uc.deviceSvc.GetDeviceInfoByID(ctx, deviceID, userID)
	if err != nil {
		uc.log.Info("device.uc.DeleteDeviceInfoByID() 오류", zap.Error(err))
		return errors.New("장치 정보를 불러올 수 없습니다")
	}

	if err := uc.deviceSvc.DeleteDeviceInfoByID(ctx, deviceID, userID); err != nil {
		uc.log.Info("device.uc.DeleteDeviceInfoByID() 오류", zap.Error(err))
		return err
	}

	return nil
}

func (uc *deviceUseCase) UpdateDeviceRequestSchemaByID(ctx context.Context, in *domain.DeviceRequestSchema, deviceID string) error {
	validate, err := uc.validateUser(ctx, domain.UserRoleDevice)
	if err != nil {
		return err
	}

	param := &domain.RawDeviceRequestSchema{
		ID:       in.ID,
		DeviceID: deviceID,
		Key:      in.Key,
	}
	sensorData, err := uc.metaSvc.GetSensorByParam(ctx, &domain.Sensor{Title: in.Target})
	if err != nil {
		uc.log.Info("device.uc.UpdateDeviceRequestSchemaByID() 오류 - 센서 정보를 받아올 수 없습니다.", zap.Error(err))
		return errors.New("존재하지 않는 센서입니다")
	}
	param.TargetSensorID = sensorData.ID

	if err := uc.deviceSvc.UpdateDeviceRequestSchemaByID(ctx, param, validate.UserID()); err != nil {
		uc.log.Info("device.uc.UpdateDeviceRequestSchemaByID() 오류", zap.Error(err))
		return err
	}

	return nil

}

func (uc *deviceUseCase) GetDeviceRequestSchemaListByID(ctx context.Context, deviceID string) ([]*domain.DeviceRequestSchema, error) {
	validate, err := uc.validateUser(ctx, domain.UserRoleDevice)
	if err != nil {
		return nil, err
	}

	// 장비 접근 권한 체크
	deviceInfo, err := uc.deviceSvc.GetDeviceInfoByID(ctx, deviceID, validate.UserID())
	if err != nil {
		return nil, err
	}

	list, err := uc.deviceSvc.GetDeviceRequestSchemaListByID(ctx, deviceInfo.ID)
	if err != nil {
		uc.log.Info("device.uc.GetDeviceRequestSchemaListByID() 오류", zap.Error(err))
		return nil, errors.New("장비 스키마 정보를 불러올 수 없습니다")
	}

	return list, nil
}

func (uc *deviceUseCase) UpdateDeviceInfo(ctx context.Context, in *domain.DeviceInfo) error {
	validate, err := uc.validateUser(ctx, domain.UserRoleDevice)
	if err != nil {
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
	validate, err := uc.validateUser(ctx, domain.UserRoleDevice)
	if err != nil {
		return nil, err
	}

	data, err := uc.deviceSvc.GetDeviceInfoByID(ctx, id, validate.UserID())
	if err != nil {
		return nil, errors.New("id를 확인해주세요")
	}

	return data, nil
}

func (uc *deviceUseCase) GetDeviceInfoListByParamAndPage(ctx context.Context, in *domain.DeviceInfo, page *domain.Page) ([]*domain.DeviceInfo, *domain.Page, error) {
	validate, err := uc.validateUser(ctx, domain.UserRoleDevice)
	if err != nil {
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
	validate, err := uc.validateUser(ctx, domain.UserRoleDevice)
	if err != nil {
		return err
	}

	var rawDeviceInfo domain.RawDeviceInfo
	var rawDeviceRequestSchemaList []*domain.RawDeviceRequestSchema

	rawDeviceInfo.UserID = validate.UserID()
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
