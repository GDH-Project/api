package handler

import (
	"context"
	"net/http"

	"github.com/GDH-Project/api/internal/domain"
	"github.com/GDH-Project/api/internal/middleware"
	"github.com/GDH-Project/api/internal/util"
	"github.com/danielgtaylor/huma/v2"
	"go.uber.org/zap"
)

func RegisterDeviceHandler(api huma.API, log *zap.Logger, deviceUseCase domain.DeviceUseCase, m middleware.Middleware) {
	v1 := huma.NewGroup(api, "/api/v1")

	// 장비 생성 API
	huma.Register(v1, m.WithAuth(huma.Operation{
		OperationID:   "v1DeviceCreateDeviceInfoWithSchema",
		Method:        http.MethodPost,
		Path:          "/device",
		Summary:       "장치 생성",
		Description:   "장치 생성 API 입니다. 장치 데이터와 스키마 데이터를 받아 새로운 장치를 생성합니다. 스키마 데이터는 생략 가능합니다.",
		Tags:          []string{"Device"},
		DefaultStatus: http.StatusCreated,
	}), func(ctx context.Context, i *struct {
		Body struct {
			Title       string `json:"title" minLength:"5" doc:"장치의 이름 입니다. 검색 시 노출되는 이름 입니다." example:"경기도 안양시 토마토 스마트팜"`
			Name        string `json:"name,omitempty" doc:"장치 등록자만 확인 가능한 값입니다. 개인의 장치 식별에 사용하면 됩니다." example:"A-B1 섹터 3구역"`
			Crop        string `json:"crop" doc:"작물명 입니다." example:"토마토"`
			UpdateCycle int    `json:"interval" doc:"데이터의 업데이트 주기 입니다.(분)" example:"30"`
			Address     struct {
				State string `json:"state" doc:"도/특별시 명칭 입니다." example:"경기도"`
				City  string `json:"city" doc:"시/군/구 명칭 입니다." example:"안양시"`
			} `json:"address" doc:"주소 정보 입니다."`
			Schema []struct {
				Key    string `json:"key" doc:"장비에서 보내는 데이터 JSON의 키값 입니다." example:"soil_temp"`
				Target string `json:"target" doc:"지정한 키를 바인딩 할 센서 명칭 입니다." example:"토양 온도"`
			} `json:"schema,omitempty" doc:"장비의 요청과 센서 값을 바인딩 하는 스키마 입니다."`
		}
	}) (*struct{}, error) {
		// 권한 체크
		validate := util.ValidateUser{
			Ctx:        ctx,
			TargetRole: domain.UserRoleDevice,
		}
		if !validate.Exec() {
			return nil, huma.Error403Forbidden("권한이 없습니다.")
		}

		// 파라미터 변환
		var deviceInfo domain.DeviceInfo
		deviceInfo.Title = i.Body.Title
		deviceInfo.Name = i.Body.Name
		deviceInfo.Crop = i.Body.Crop
		deviceInfo.UpdateCycle = i.Body.UpdateCycle
		deviceInfo.Address.State = i.Body.Address.State
		deviceInfo.Address.City = i.Body.Address.City

		var schemas []*domain.DeviceRequestSchema

		for _, item := range i.Body.Schema {
			schemas = append(schemas, &domain.DeviceRequestSchema{
				Key:    item.Key,
				Target: item.Target,
			})
		}

		if err := deviceUseCase.CreateDevice(ctx, &deviceInfo, schemas); err != nil {
			return nil, huma.Error400BadRequest(err.Error())
		}

		return nil, nil
	})

}
