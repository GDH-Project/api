package domain

import (
	"context"
	"time"

	"github.com/jackc/pgx/v5"
)

type DeviceRepository interface {
	// WithTransaction 트렌잭션 함수
	WithTransaction(ctx context.Context, f func(tx pgx.Tx) error) error
	// CreateDeviceInfoTx device를 생성하고 ID 와 오류를 반환
	CreateDeviceInfoTx(ctx context.Context, tx pgx.Tx, in *CreateDeviceInfo) (string, error)
}

// DeviceData
//
// 장비에서 수집된 데이터 JSON 배열 입니다.
type DeviceData struct {
	Time     time.Time              `json:"-"` // Datajson에 추가할 시간 정보
	DeviceID string                 `json:"-"` // 장치 ID
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

type CreateDeviceInfo struct {
	UserID         string // 유저의 UUID 입니다.
	Title          string
	Name           string
	CropID         int
	UpdateCycleID  int
	AddressStateID int
	AddressCityID  int
}
type DeviceInfo struct {
	UserID string // 유저의 UUID 입니다.

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
