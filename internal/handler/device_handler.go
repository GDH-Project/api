package handler

import (
	"context"
	"net/http"

	"github.com/GDH-Project/api/internal/domain"
	"github.com/GDH-Project/api/internal/middleware"
	"github.com/danielgtaylor/huma/v2"
	"go.uber.org/zap"
)

type pageInfo struct {
	Size        int  `json:"size" doc:"현재 페이지 크기 입니다." example:"10"`
	Page        int  `json:"page" doc:"현재 페이지 입니다." example:"1"`
	NextPage    int  `json:"next_page,omitempty" doc:"다음 페이지 번호 입니다." example:"2"`
	RecordCount int  `json:"record_count" doc:"총 데이터 수 입니다." example:"20"`
	HasNextPage bool `json:"has_next_page" doc:"다음 페이지 존재 여부 입니다." example:"true"`
}
type deviceInfoListResponse struct {
	Body struct {
		Data     []*domain.DeviceInfo `json:"data" doc:"장비 정보 배열 입니다."`
		PageInfo pageInfo             `json:"page_info" doc:"페이지 정보 입니다."`
	}
}

func RegisterDeviceHandler(api huma.API, log *zap.Logger, deviceUseCase domain.DeviceUseCase, m middleware.Middleware) {
	v1 := huma.NewGroup(api, "/api/v1")

	// 장비 생성 API
	huma.Register(v1, m.WithAuth(
		huma.Operation{
			OperationID:   "v1DeviceCreateDeviceInfoWithSchema",
			Method:        http.MethodPost,
			Path:          "/device",
			Summary:       "장치 생성",
			Description:   "장치 생성 API 입니다. 장치 데이터와 스키마 데이터를 받아 새로운 장치를 생성합니다. 스키마 데이터는 생략 가능합니다.",
			Tags:          []string{"Device"},
			DefaultStatus: http.StatusCreated,
		},
		domain.UserRoleDevice,
	), func(ctx context.Context, i *struct {
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

		// 파라미터 변환
		var deviceInfo domain.DeviceInfo
		deviceInfo.Title = i.Body.Title
		deviceInfo.Name = &i.Body.Name
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

	// 나의 장비 리스트 조회 API
	huma.Register(v1, m.WithAuth(
		huma.Operation{
			OperationID:   "v1DeviceGetMyDeviceInfoList",
			Method:        http.MethodGet,
			Path:          "/devices",
			Summary:       "나의 장치 정보 리스트 조회",
			Description:   "나의 장치 정보 리스트 조회 API 입니다.",
			Tags:          []string{"Device"},
			DefaultStatus: http.StatusOK,
		},
		domain.UserRoleDevice,
	), func(ctx context.Context, i *struct {
		Size     int    `query:"size" doc:"조회할 페이지 크기 입니다." default:"10"`
		Page     int    `query:"page" doc:"조회할 페이지 입니다." default:"1"`
		Title    string `query:"title" doc:"title 명칭 입니다." example:"스마트"`
		Crop     string `query:"crop" doc:"작물명입니다." example:"토마토"`
		Interval int    `query:"interval" doc:"데이터 갱신 주기 입니다." example:"10"`
		State    string `query:"state" doc:"도/특별시 명칭 입니다." example:"경기도"`
		City     string `query:"city" doc:"시/군/구 명칭 입니다." example:"안양시"`
	}) (*deviceInfoListResponse, error) {
		var resp deviceInfoListResponse

		var param domain.DeviceInfo
		param.Title = i.Title
		param.Crop = i.Crop
		param.UpdateCycle = i.Interval
		param.Address.State = i.State
		param.Address.City = i.City

		var p domain.Page
		p.Size = i.Size
		p.Page = i.Page
		dataList, page, err := deviceUseCase.GetDeviceInfoListByParamAndPage(ctx, &param, &p)
		if err != nil {
			return nil, huma.Error400BadRequest(err.Error())
		}

		resp.Body.Data = dataList

		pageInfo := pageInfo{
			Size:        page.Size,
			Page:        page.Page,
			RecordCount: page.RecordCount,
		}
		if page.HasNext() {
			pageInfo.NextPage = page.Page + 1
			pageInfo.HasNextPage = true
		}

		resp.Body.PageInfo = pageInfo

		return &resp, nil
	})

}
