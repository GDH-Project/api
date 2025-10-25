package domain

import "context"

type MetaRepository interface {
	// GetSensorList 센서 전체 리스트 조회
	GetSensorList(ctx context.Context) ([]*Sensor, error)
	// GetSensorByParam 센서를 파라미터를 통해 조회
	GetSensorByParam(ctx context.Context, in *Sensor) (*Sensor, error)

	// GetCropList 모든 작물 정보 조회
	GetCropList(ctx context.Context) ([]*Crop, error)
	// GetCropByParam 작물 파라미터를 통해 조회
	GetCropByParam(ctx context.Context, in *Crop) (*Crop, error)

	// GetUpdateCycleList 업데이트 주기 조회
	GetUpdateCycleList(ctx context.Context) ([]*UpdateCycle, error)

	// GetAddressStateList 도/특별시 전체 리스트 조회
	GetAddressStateList(ctx context.Context) ([]*AddressState, error)
	// GetAddressCityListByState 도/특별시 정보를 통해 시/군/구 리스트 반환
	GetAddressCityListByState(ctx context.Context, state string) ([]*AddressCity, error)

	// GetDeviceDataListByUserIDWithPage(ctx context.Context, userID string, page Page) ([]*DeviceData, error)
	//
	// GetDeviceRequestSchemaListByDeviceID(ctx context.Context, deviceID string) ([]*DeviceRequestSchema, error)
	//
	// GetDeviceInfoListByParamWithPage(ctx context.Context, in *DeviceInfo, page Page) ([]*DeviceInfo, error)
	// GetDeviceInfoByParam(ctx context.Context, in *DeviceInfo) (*DeviceInfo, error)
}

type MetaService interface {
	// GetSensorList 센서 전체 리스트 조회
	GetSensorList(ctx context.Context) ([]*Sensor, error)
	// GetSensorByParam 센서를 파라미터를 통해 조회
	GetSensorByParam(ctx context.Context, in *Sensor) (*Sensor, error)

	// GetCropList 모든 작물 정보 조회
	GetCropList(ctx context.Context) ([]*Crop, error)
	// GetCropByParam 작물 파라미터를 통해 조회
	GetCropByParam(ctx context.Context, in *Crop) (*Crop, error)

	// GetUpdateCycleList 업데이트 주기 조회
	GetUpdateCycleList(ctx context.Context) ([]*UpdateCycle, error)

	// GetAddressStateList 도/특별시 전체 리스트 조회
	GetAddressStateList(ctx context.Context) ([]*AddressState, error)
	// GetAddressCityListByState 도/특별시 정보를 통해 시/군/구 리스트 반환
	GetAddressCityListByState(ctx context.Context, state string) ([]*AddressCity, error)
}

type MetaUseCase interface {
	// GetSensorList 센서 전체 리스트 조회
	GetSensorList(ctx context.Context) ([]*Sensor, error)
	// GetSensorByParam 센서를 파라미터를 통해 조회
	GetSensorByParam(ctx context.Context, in *Sensor) (*Sensor, error)

	// GetCropList 모든 작물 정보 조회
	GetCropList(ctx context.Context) ([]*Crop, error)
	// GetCropByParam 작물 파라미터를 통해 조회
	GetCropByParam(ctx context.Context, in *Crop) (*Crop, error)

	// GetUpdateCycleList 업데이트 주기 조회
	GetUpdateCycleList(ctx context.Context) ([]*UpdateCycle, error)

	// GetAddressStateList 도/특별시 전체 리스트 조회
	GetAddressStateList(ctx context.Context) ([]*AddressState, error)
	// GetAddressCityListByState 도/특별시 정보를 통해 시/군/구 리스트 반환
	GetAddressCityListByState(ctx context.Context, state string) ([]*AddressCity, error)
}

// Sensor
//
// 수집되는 센서 데이터 정보 입니다.
type Sensor struct {
	ID       int     `json:"id" doc:"센서의 고유 ID 입니다." example:"1"`
	Title    string  `json:"title" doc:"센서의 한글 명칭 입니다." example:"기온"`
	EngTitle string  `json:"eng_title" doc:"센서의 영어 명칭 입니다." example:"Air Temperature"`
	Desc     string  `json:"desc" doc:"센서 설명 입니다." example:"작물의 광합성, 호흡, 증산 작용에 직접적인 영향을 미치는 대기의 온도"`
	Unit     *string `json:"unit,omitempty" doc:"센서 단위 입니다." example:"°C"`
	UnitDesc *string `json:"unit_desc,omitempty" doc:"센서 단위 설명 입니다." example:"섭씨"`
}

// Crop
//
// 작물 정보 입니다.
type Crop struct {
	ID    int     `json:"-"`
	Title string  `json:"title" doc:"작물명" example:"토마토"`
	Desc  *string `json:"desc,omitempty" doc:"작물의 설명" example:"토마토에 대한 설명입니다."`
}

// UpdateCycle
//
// 통신 주기 입니다.
type UpdateCycle struct {
	ID       int     `json:"-"`
	Interval int     `json:"interval" doc:"통신 주기(분)" example:"60"`
	Desc     *string `json:"desc,omitempty" doc:"설명입니다." example:"1시간 주기 업데이트"`
}

// AddressState
//
// 주소 도/특별시 정보 입니다.
type AddressState struct {
	ID    int    `json:"-"`
	Title string `json:"title" doc:"도/특별시 입니다." example:"서울특별시"`
}

// AddressCity
//
// 시/군/구 정보 입니다.
type AddressCity struct {
	ID         int    `json:"-"`
	StateTitle string `json:"state_title" doc:"도/특별시" example:"서울특별시"`
	Title      string `json:"title" doc:"시/군/구" example:"동대문구"`
}
