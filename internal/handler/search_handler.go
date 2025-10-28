package handler

import (
	"context"
	"net/http"

	"github.com/GDH-Project/api/internal/domain"
	"github.com/danielgtaylor/huma/v2"
	"go.uber.org/zap"
)

type searchDeviceDataListResponse struct {
	Body struct {
		Data []map[string]interface{} `json:"data" doc:"저장된 데이터 JSON 배열 입니다."`
	}
}

type searchDeviceInfo struct {
	domain.DeviceInfo
	Name string `json:"-"` // 리턴시 생략을 위해 추가
}
type searchDeviceInfoResponse struct {
	Body struct {
		Data searchDeviceInfo `json:"data" doc:"장비 정보 JSON 입니다."`
	}
}

type searchDeviceInfoListResponse struct {
	Body struct {
		Data     []*searchDeviceInfo `json:"data" doc:"장비 정보 JSON 리스트 입니다."`
		PageInfo pageInfo            `json:"page_info" doc:"페이지 정보 입니다."`
	}
}

func RegisterSearchHandler(api huma.API, log *zap.Logger, searchUseCase domain.SearchUseCase) {
	v1 := huma.NewGroup(api, "/api/v1")

	// 장치 데이터 검색 By ID API
	huma.Register(v1, huma.Operation{
		OperationID:   "v1SearchGetDeviceDataListByQuery",
		Method:        http.MethodGet,
		Path:          "/search/{device_id}/data",
		Summary:       "장비 데이터 리스트 검색",
		Description:   "장비 데이터 리스트 검색 API 입니다.",
		Tags:          []string{"Search"},
		DefaultStatus: http.StatusOK,
	}, func(ctx context.Context, i *struct {
		DeviceID string `path:"device_id" doc:"장치 고유 ID 입니다." format:"uuid"`
	}) (*searchDeviceDataListResponse, error) {
		var resp searchDeviceDataListResponse
		data, err := searchUseCase.GetDeviceDataListByDeviceID(ctx, i.DeviceID)
		if err != nil {
			return nil, huma.Error400BadRequest(err.Error())
		}

		resp.Body.Data = data
		return &resp, nil
	})

	// 장치 정보 조회 By ID
	huma.Register(v1, huma.Operation{
		OperationID:   "v1SearchGetDeviceInfoByID",
		Method:        http.MethodGet,
		Path:          "/search/{device_id}",
		Summary:       "장비 검색 By ID",
		Description:   "장비 검색 By ID API 입니다.",
		Tags:          []string{"Search"},
		DefaultStatus: http.StatusOK,
	}, func(ctx context.Context, i *struct {
		DeviceID string `path:"device_id" doc:"장치 고유 ID 입니다." format:"uuid"`
	}) (*searchDeviceInfoResponse, error) {
		var resp searchDeviceInfoResponse
		data, err := searchUseCase.GetDeviceInfoByDeviceID(ctx, i.DeviceID)
		if err != nil {
			return nil, huma.Error400BadRequest(err.Error())
		}

		resp.Body.Data.ID = data.ID
		resp.Body.Data.Title = data.Title
		resp.Body.Data.Crop = data.Crop
		resp.Body.Data.UpdateCycle = data.UpdateCycle
		resp.Body.Data.Address = data.Address
		resp.Body.Data.CreatedAt = data.CreatedAt
		resp.Body.Data.UpdatedAt = data.UpdatedAt

		return &resp, nil
	})

	// 장치 정보 리스트 조회 By Param
	huma.Register(v1, huma.Operation{
		OperationID:   "v1SearchGetDeviceInfoListByParamAndPageInfo",
		Method:        http.MethodGet,
		Path:          "/search",
		Summary:       "장비 검색 By Param",
		Description:   "장비 검색 By Param API 입니다.",
		Tags:          []string{"Search"},
		DefaultStatus: http.StatusOK,
	}, func(ctx context.Context, i *struct {
		Size     int    `query:"size" doc:"조회할 페이지 크기 입니다." default:"10"`
		Page     int    `query:"page" doc:"조회할 페이지 입니다." default:"1"`
		Title    string `query:"title" doc:"title 명칭 입니다." example:"스마트"`
		Crop     string `query:"crop" doc:"작물명입니다." example:"토마토"`
		Interval int    `query:"interval" doc:"데이터 갱신 주기 입니다." example:"10"`
		State    string `query:"state" doc:"도/특별시 명칭 입니다." example:"경기도"`
		City     string `query:"city" doc:"시/군/구 명칭 입니다." example:"안양시"`
	}) (*searchDeviceInfoListResponse, error) {
		var resp searchDeviceInfoListResponse
		param := &domain.DeviceInfo{
			Title:       i.Title,
			Crop:        i.Crop,
			UpdateCycle: i.Interval,
		}
		param.Address.State = i.State
		param.Address.City = i.City

		p := &domain.Page{
			Size: i.Size,
			Page: i.Page,
		}

		data, page, err := searchUseCase.GetDeviceInfoListByParamAndPageInfo(ctx, param, p)
		if err != nil {
			return nil, huma.Error400BadRequest(err.Error())
		}

		var searchDeviceInfoListData []*searchDeviceInfo
		for _, item := range data {
			var t searchDeviceInfo
			t.DeviceInfo = *item

			searchDeviceInfoListData = append(searchDeviceInfoListData, &t)
		}
		resp.Body.Data = searchDeviceInfoListData

		pageInfo := &pageInfo{
			Size:        page.Size,
			Page:        page.Page,
			RecordCount: page.RecordCount,
		}
		if page.HasNext() {
			pageInfo.NextPage = page.Size + 1
			pageInfo.HasNextPage = true
		}
		resp.Body.PageInfo = *pageInfo

		return &resp, nil
	})

	log.Info("Search Handler 등록")
}
