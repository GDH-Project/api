package domain

import (
	"context"
	"database/sql"
	"time"

	"github.com/jackc/pgx/v5"
)

type DeviceRepository interface {
	// WithTransaction 트렌잭션 함수
	WithTransaction(ctx context.Context, f func(tx pgx.Tx) error) error
	// CreateDeviceInfoTx deviceInfo를 생성하고 ID 와 오류를 반환
	CreateDeviceInfoTx(ctx context.Context, tx pgx.Tx, in *RawDeviceInfo) (string, error)
	// GetDeviceInfoByID deviceInfo 를 ID 기준으로 찾아서 반환
	GetDeviceInfoByID(ctx context.Context, id string) (*DeviceInfo, error)
	// GetDeviceInfoListByParamAndPage deviceInfo 를 파라미터와 페이지 옵션을 기준으로 찾아서 리스트로 반환한다.
	GetDeviceInfoListByParamAndPage(ctx context.Context, in *DeviceInfo, page *Page) ([]*DeviceInfo, *Page, error)
	// UpdateDeviceInfo 장치 정보 업데이트
	UpdateDeviceInfo(ctx context.Context, in *RawDeviceInfo) error
	// DeleteDeviceInfoByID 장비 제거
	DeleteDeviceInfoByID(ctx context.Context, id string, userID string) error
	// CreateDeviceReqeustSchemaListTx req_to_sensor 에 삽입되는 디바이스 응답(JSON) 키와 센서를 연결하는 부분
	CreateDeviceReqeustSchemaListTx(ctx context.Context, tx pgx.Tx, deviceID string, schemas []*RawDeviceRequestSchema) error
	GetDeviceRequestSchemaListByDeviceID(ctx context.Context, deviceID string) ([]*DeviceRequestSchema, error)
	GetDeviceRequestSchemaByID(ctx context.Context, id int) (*DeviceRequestSchema, error)
	UpdateDeviceRequestSchema(ctx context.Context, in *RawDeviceRequestSchema) error
	DeleteDeviceRequestSchemaByID(ctx context.Context, id int) error
	CreateDeviceApiKey(ctx context.Context, in *RawApiKey) (*RawApiKey, error)
	// GetDeviceApiKeyByApiKey 반환되는 값은 장치 ID 이다.
	GetDeviceApiKeyByApiKey(ctx context.Context, apiKey string) (string, error)
	GetDeviceApiKeyListByUserIDAndDeviceID(ctx context.Context, userID string, deviceID string) ([]*RawApiKey, error)
}

type DeviceService interface {
	CreateDevice(ctx context.Context, deviceInfoData *RawDeviceInfo, deviceSchemaDataList []*RawDeviceRequestSchema) error
	GetDeviceInfoListByParamAndPage(ctx context.Context, in *DeviceInfo, page *Page) ([]*DeviceInfo, *Page, error)
	GetDeviceInfoByID(ctx context.Context, id string, userID string) (*DeviceInfo, error)
	UpdateDeviceInfo(ctx context.Context, in *RawDeviceInfo) error
	GetDeviceReqeustSchemaListByID(ctx context.Context, deviceID string, userID string) ([]*DeviceRequestSchema, error)
	UpdateDeviceReqeustSchemaByID(ctx context.Context, in *RawDeviceRequestSchema, userID string) error
	DeleteDeviceInfoByID(ctx context.Context, deviceID string, userID string) error
	CreateDeviceReqeustSchema(ctx context.Context, in *RawDeviceRequestSchema) error
	CreateDeviceApiKey(ctx context.Context, in *RawApiKey) (*ApiKey, error)
	// GetDeviceApiKeyByApiKey 반횐되는 값은 장치 ID 이다.
	GetDeviceApiKeyByApiKey(ctx context.Context, apiKey string) (string, error)
	GetDeviceApiKeyListByUserIDAndDeviceID(ctx context.Context, userID string, deviceID string) ([]*RawApiKey, error)
}

type DeviceUseCase interface {
	CreateDevice(ctx context.Context, deviceInfoData *DeviceInfo, deviceSchemaDataList []*DeviceRequestSchema) error
	GetDeviceInfoListByParamAndPage(ctx context.Context, in *DeviceInfo, page *Page) ([]*DeviceInfo, *Page, error)
	GetDeviceInfoByID(ctx context.Context, id string) (*DeviceInfo, error)
	UpdateDeviceInfo(ctx context.Context, in *DeviceInfo) error

	GetDeviceReqeustSchemaListByID(ctx context.Context, deviceID string) ([]*DeviceRequestSchema, error)
	UpdateDeviceReqeustSchemaByID(ctx context.Context, in *DeviceRequestSchema, deviceID string) error
	DeleteDeviceInfoByID(ctx context.Context, deviceID string) error

	CreateDeviceReqeustSchema(ctx context.Context, deviceID string, in *DeviceRequestSchema) error
	CreateDeviceApiKey(ctx context.Context, in *ApiKey) (*ApiKey, error)
	// GetDeviceApiKeyByApiKey 반횐되는 값은 장치 ID 이다.
	GetDeviceApiKeyByApiKey(ctx context.Context, apiKey string) (string, error)
	GetDeviceApiKeyListByUserIDAndDeviceID(ctx context.Context, deviceID string) ([]*ApiKey, error)
}

// DeviceData
//
// 장비에서 수집된 데이터 JSON 배열 입니다.
type DeviceData struct {
	Time     time.Time              `json:"-"` // Datajson에 추가할 시간 정보
	DeviceID string                 `json:"-"` // 장치 ID
	Data     map[string]interface{} `json:"data" doc:"장비에서 수집된 데이터 JSON문자열 + time 정보"`
}

type RawDeviceRequestSchema struct {
	ID             int
	DeviceID       string
	Key            string
	TargetSensorID int
}

// DeviceRequestSchema
//
// 장비의 요청과 센서 정보를 바인딩 하는 스키마 입니다.
type DeviceRequestSchema struct {
	ID     int    `json:"id" doc:"고유 ID 입니다."`
	Key    string `json:"key" doc:"장비에서 보내는 데이터의 json key 입니다." example:"degree"`
	Target string `json:"target" doc:"센서 데이터 리스트의 title 명칭 입니다." example:"기온"`
}

type RawDeviceInfo struct {
	ID             string // 장치 고유 ID
	UserID         string // 유저의 UUID 입니다.
	Title          string
	Name           string
	CropID         int
	UpdateCycleID  int
	AddressStateID int
	AddressCityID  int
}
type DeviceInfo struct {
	UserID string `json:"-"` // 유저의 UUID 입니다.

	ID          string  `json:"id" doc:"장치 고유 ID 입니다." format:"uuid"`
	Title       string  `json:"title" doc:"검색에 노출되는 명칭입니다." example:"안양시 자동 재배 시설 토마토 데이터"`
	Name        *string `json:"name,omitempty" doc:"장치관리자에게 보이는 고유 명칭 입니다." example:"안양시 스마트 펙토리 토마토 A-B1 섹터"`
	Crop        string  `json:"crop" doc:"작물 정보 입니다." example:"토마토"`
	UpdateCycle int     `json:"update_cycle" doc:"데이터 업데이트 주기 입니다." example:"60"`
	Address     struct {
		State string `json:"state" doc:"도/특별시 명칭 입니다." example:"경기도"`
		City  string `json:"city" doc:"시/군/구 명칭 입니다." example:"안양시"`
	} `json:"address" doc:"주소지"`

	CreatedAt time.Time `json:"created_at" doc:"최초 장치 등록 시간 입니다." example:"2025-10-24 22:54:52.874221 +09:00"`
	UpdatedAt time.Time `json:"updated_at" doc:"장치 정보 업데이트 시간 입니다." example:"2025-10-24 22:54:52.874221 +09:00"`
}

type RawApiKey struct {
	ID        int            `json:"-"` // 32자리 문자열로 직접 생성
	APIKey    string         `json:"-"`
	DeviceID  string         `json:"device_id"`
	UserID    string         `json:"user_id"` // uuid
	Title     string         `json:"title"`
	Desc      sql.NullString `json:"desc"`
	CreatedAt time.Time      `json:"created_at"`
}

type ApiKey struct {
	ID        int       `json:"id" doc:"API Key의 고유 ID 입니다." example:"1"`
	Key       string    `json:"-"`
	DeviceID  string    `json:"device_id" doc:"장치 고유 ID 입니다."`
	Title     string    `json:"title" maxLength:"50" doc:"API 키에 대한 이름 입니다." example:"경기도 안양시 토마토 농장 A-B1 섹터 센서"`
	Desc      string    `json:"desc,omitempty" doc:"API 키에 대한 셜명입니다." example:"2층 토마토 센서"`
	CreatedAt time.Time `json:"created_at" doc:"API 키 생성 시간 입니다."`
}
