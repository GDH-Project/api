package domain

import (
	"context"
	"time"
)

type DeviceRepository interface {
	// GetSensorList 센서 전체 리스트 조회
	GetSensorList(ctx context.Context) ([]*Sensor, error)
	// GetSensorByParam 센서를 파라미터를 통해 조회
	GetSensorByParam(ctx context.Context, in *Sensor) (*Sensor, error)

	// GetAllCropList 모든 작물 정보 조회
	GetAllCropList(ctx context.Context) ([]*Crop, error)
	// GetCropByParam 작물 파라미터를 통해 조회
	GetCropByParam(ctx context.Context, in *Sensor) (*Crop, error)

	// GetAllUpdateCycleList 업데이트 주기 조회
	GetAllUpdateCycleList(ctx context.Context) ([]*UpdateCycle, error)

	// GetAllAddressStateList 도/특별시 전체 리스트 조회
	GetAllAddressStateList(ctx context.Context) ([]*AddressState, error)
	// GetAddressCityListByState 도/특별시 정보를 통해 시/군/구 리스트 반환
	GetAddressCityListByState(ctx context.Context, state string) ([]*AddressCity, error)

	// GetDeviceDataListByUserIDWithPage(ctx context.Context, userID string, page Page) ([]*DeviceData, error)
	//
	// GetDeviceRequestSchemaListByDeviceID(ctx context.Context, deviceID string) ([]*DeviceRequestSchema, error)
	//
	// GetDeviceInfoListByParamWithPage(ctx context.Context, in *DeviceInfo, page Page) ([]*DeviceInfo, error)
	// GetDeviceInfoByParam(ctx context.Context, in *DeviceInfo) (*DeviceInfo, error)
}

type DeviceService interface {
	// GetSensorList 센서 전체 리스트 조회
	GetSensorList(ctx context.Context) ([]*Sensor, error)
	// GetSensorByParam 센서를 파라미터를 통해 조회
	GetSensorByParam(ctx context.Context, in *Sensor) (*Sensor, error)

	// GetAllCropList 모든 작물 정보 조회
	GetAllCropList(ctx context.Context) ([]*Crop, error)
	// GetCropByParam 작물 파라미터를 통해 조회
	GetCropByParam(ctx context.Context, in *Sensor) (*Crop, error)

	// GetAllUpdateCycleList 업데이트 주기 조회
	GetAllUpdateCycleList(ctx context.Context) ([]*UpdateCycle, error)

	// GetAllAddressStateList 도/특별시 전체 리스트 조회
	GetAllAddressStateList(ctx context.Context) ([]*AddressState, error)
	// GetAddressCityListByState 도/특별시 정보를 통해 시/군/구 리스트 반환
	GetAddressCityListByState(ctx context.Context, state string) ([]*AddressCity, error)
}

type DeviceUseCase interface {
	// GetSensorList 센서 전체 리스트 조회
	GetSensorList(ctx context.Context) ([]*Sensor, error)
	// GetSensorByParam 센서를 파라미터를 통해 조회
	GetSensorByParam(ctx context.Context, in *Sensor) (*Sensor, error)

	// GetAllCropList 모든 작물 정보 조회
	GetAllCropList(ctx context.Context) ([]*Crop, error)
	// GetCropByParam 작물 파라미터를 통해 조회
	GetCropByParam(ctx context.Context, in *Sensor) (*Crop, error)

	// GetAllUpdateCycleList 업데이트 주기 조회
	GetAllUpdateCycleList(ctx context.Context) ([]*UpdateCycle, error)

	// GetAllAddressStateList 도/특별시 전체 리스트 조회
	GetAllAddressStateList(ctx context.Context) ([]*AddressState, error)
	// GetAddressCityListByState 도/특별시 정보를 통해 시/군/구 리스트 반환
	GetAddressCityListByState(ctx context.Context, state string) ([]*AddressCity, error)
}

type Page struct {
	Size int
	Page int
}

// Sensor
//
// 수집되는 센서 데이터 정보 입니다.
type Sensor struct {
	ID       int
	Title    string `json:"title" doc:"센서의 한글 명칭 입니다." example:"기온"`
	EngTitle string `json:"eng_title" doc:"센서의 영어 명칭 입니다." example:"Air Temperature"`
	Desc     string `json:"desc" doc:"센서 설명 입니다." example:"작물의 광합성, 호흡, 증산 작용에 직접적인 영향을 미치는 대기의 온도"`
	Unit     string `json:"unit,omitempty" doc:"센서 단위 입니다." example:"°C"`
	UnitDesc string `json:"unit_desc,omitempty" doc:"센서 단위 설명 입니다." example:"섭씨"`
}

// Crop
//
// 작물 정보 입니다.
type Crop struct {
	ID    int
	Title string `json:"title" doc:"작물명" example:"토마토"`
	Desc  string `json:"desc,omitempty" doc:"작물의 설명"`
}

// UpdateCycle
//
// 통신 주기 입니다.
type UpdateCycle struct {
	ID       int
	Interval int `json:"interval" doc:"통신 주기(분)" example:"60"`
	Desc     int `json:"desc,omitempty" doc:"설명입니다." example:"1시간 주기 업데이트"`
}

// AddressState
//
// 주소 도/특별시 정보 입니다.
type AddressState struct {
	ID    int
	Title string `json:"title" doc:"도/특별시 입니다." example:"서울특별시"`
}

// AddressCity
//
// 시/군/구 정보 입니다.
type AddressCity struct {
	ID         int
	StateTitle string // 도/특별시 title ex)"서울특별시"
	Title      string `json:"title" doc:"시/군/구" example:"동대문구"`
}

// DeviceData
//
// 장비에서 수집된 데이터 JSON 배열 입니다.
type DeviceData struct {
	Time     time.Time              // Datajson에 추가할 시간 정보
	DeviceID string                 // 장치 ID
	Data     map[string]interface{} `json:"data" doc:"장비에서 수집된 데이터 JSON문자열 + time 정보"`
}

// DeviceRequestSchema
//
// 장비의 요청과 센서 정보를 바인딩 하는 스키마 입니다.
type DeviceRequestSchema struct {
	ID       int    `json:"id" doc:"고유 ID 입니다."`
	DeviceID string // 장치 ID
	Key      string `json:"key" doc:"장비에서 보내는 데이터의 json key 입니다." example:"degree"`
	Target   string `json:"target" doc:"센서 데이터 리스트의 title 명칭 입니다." example:"기온"`
}

type DeviceInfo struct {
	UserID string // 유저의 UUID 입니다.

	ID          string `json:"id" doc:"장치 고유 ID 입니다." format:"uuid"`
	Title       string `json:"title" doc:"검색에 노출되는 명칭입니다." example:"경기도 안양시 자동 재배 시설 토마토 데이터"`
	Name        string `json:"name,omitempty" doc:"장치관리자에게 보이는 고유 명칭 입니다." example:"안양시 스마트 펙토리 토마토 A-B1 섹터"`
	Crop        string `json:"crop" doc:"작물 정보 입니다." example:"토마토"`
	UpdateCycle int    `json:"update_cycle" doc:"데이터 업데이트 주기 입니다." example:"60"`
	Address     struct {
		State string `json:"state" doc:"도/특별시 명칭 입니다." example:"경기도"`
		City  string `json:"city" doc:"시/군/구 명칭 입니다." example:"안양시"`
	} `json:"address" doc:"주소지"`

	CreatedAt time.Time `json:"created_at" doc:"최초 장치 등록 시간 입니다." example:"2025-10-24 22:54:52.874221 +09:00"`
	UpdatedAt time.Time `json:"updated_at" doc:"장치 정보 업데이트 시간 입니다." example:"2025-10-24 22:54:52.874221 +09:00"`
}
