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
		Data     []*domain.DeviceInfo `json:"data" doc:"장치 정보 JSON 배열 입니다."`
		PageInfo pageInfo             `json:"page_info" doc:"페이지 정보 입니다."`
	}
}

type deviceInfoResponse struct {
	Body struct {
		Data *domain.DeviceInfo `json:"data" doc:"장치 정보 JSON 입니다."`
	}
}

type deviceRequestSchemaListResponse struct {
	Body struct {
		Data []*domain.DeviceRequestSchema `json:"data" doc:"장치 요청 스키마 JSON 배열 입니다."`
	}
}

type firstDeviceApiKeyResponse struct {
	Body struct {
		domain.ApiKey
		Key string `json:"key" doc:"Api 키 입니다. 최초 한번만 확인 할 수 있습니다."`
	}
}

func RegisterDeviceHandler(api huma.API, log *zap.Logger, deviceUseCase domain.DeviceUseCase, m middleware.Middleware) {
	v1 := huma.NewGroup(api, "/api/v1")

	// 장치 생성 API
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
				Key    string `json:"key" doc:"장치에서 보내는 데이터 JSON의 키값 입니다." example:"soil_temp"`
				Target string `json:"target" doc:"지정한 키를 바인딩 할 센서 명칭 입니다." example:"토양 온도"`
			} `json:"schema,omitempty" doc:"장치의 요청과 센서 값을 바인딩 하는 스키마 입니다."`
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

	// 나의 장치 리스트 조회 API
	huma.Register(v1, m.WithAuth(
		huma.Operation{
			OperationID:   "v1DeviceGetMyDeviceInfoList",
			Method:        http.MethodGet,
			Path:          "/devices",
			Summary:       "장치 정보 리스트 조회",
			Description:   "장치 정보 리스트 조회 API 입니다.",
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

	// 장치 정보 조회 By ID
	huma.Register(v1, m.WithAuth(
		huma.Operation{
			OperationID:   "v1DeviceGetMyDeviceInfoByID",
			Method:        http.MethodGet,
			Path:          "/device/{device_id}",
			Summary:       "장치 정보 조회 By ID",
			Description:   "장치 정보 조회 By ID API 입니다.",
			Tags:          []string{"Device"},
			DefaultStatus: http.StatusOK,
		},
		domain.UserRoleDevice,
	), func(ctx context.Context, i *struct {
		DeviceID string `path:"device_id" doc:"장치 정보 고유 ID 입니다." format:"uuid"`
	}) (*deviceInfoResponse, error) {
		var resp deviceInfoResponse
		data, err := deviceUseCase.GetDeviceInfoByID(ctx, i.DeviceID)
		if err != nil {
			return nil, huma.Error400BadRequest(err.Error())
		}

		resp.Body.Data = data

		return &resp, nil
	})

	// 장치 정보 업데이트 By ID API
	huma.Register(v1, m.WithAuth(
		huma.Operation{
			OperationID:   "v1DeviceUpdateMyDeviceInfoByID",
			Method:        http.MethodPut,
			Path:          "/device/{device_id}",
			Summary:       "장치 정보 업데이트 By ID",
			Description:   "장치 정보 업데이트 By ID API 입니다.",
			Tags:          []string{"Device"},
			DefaultStatus: http.StatusOK,
		},
		domain.UserRoleDevice,
	), func(ctx context.Context, i *struct {
		DeviceID string `path:"device_id" doc:"장치 정보 고유 ID 입니다." format:"uuid"`
		Body     struct {
			Title       string `json:"title,omitempty" minLength:"5" doc:"장치의 이름 입니다. 검색 시 노출되는 이름 입니다." example:"경기도 안양시 토마토 스마트팜"`
			Name        string `json:"name,omitempty" doc:"장치 등록자만 확인 가능한 값입니다. 개인의 장치 식별에 사용하면 됩니다." example:"A-B1 섹터 3구역"`
			UpdateCycle int    `json:"interval,omitempty" doc:"데이터의 업데이트 주기 입니다.(분)" example:"30"`
		}
	}) (*struct{}, error) {
		var param domain.DeviceInfo
		param.ID = i.DeviceID
		param.Title = i.Body.Title
		param.Name = &i.Body.Name
		param.UpdateCycle = i.Body.UpdateCycle

		if err := deviceUseCase.UpdateDeviceInfo(ctx, &param); err != nil {
			return nil, huma.Error400BadRequest(err.Error())
		}

		return nil, nil
	})

	// 장치 정보 제거 By ID API
	huma.Register(v1, m.WithAuth(
		huma.Operation{
			OperationID:   "v1DeviceDeleteDeviceInfoByID",
			Method:        http.MethodDelete,
			Path:          "/device/{device_id}",
			Summary:       "장치 정보 제거 By ID",
			Description:   "장치 정보 제거 By ID API 입니다.",
			Tags:          []string{"Device"},
			DefaultStatus: http.StatusOK,
		},
		domain.UserRoleDevice,
	), func(ctx context.Context, i *struct {
		DeviceID string `path:"device_id" doc:"장치 정보 고유 ID 입니다." format:"uuid"`
	}) (*struct{}, error) {
		if err := deviceUseCase.DeleteDeviceInfoByID(ctx, i.DeviceID); err != nil {
			return nil, huma.Error400BadRequest(err.Error())
		}
		return nil, nil
	})

	// 장치 요청 스키마 생성 API
	huma.Register(v1, m.WithAuth(
		huma.Operation{
			OperationID:   "v1DeviceCreateDeviceReqeustSchemaByDeviceID",
			Method:        http.MethodPost,
			Path:          "/device/{device_id}/schema",
			Summary:       "장치 요청 스키마 생성",
			Description:   "장치 요청 스키마 생성 API 입니다.",
			Tags:          []string{"Device"},
			DefaultStatus: http.StatusCreated,
		},
		domain.UserRoleDevice,
	), func(ctx context.Context, i *struct {
		DeviceID string `path:"device_id" doc:"장치 정보 고유 ID 입니다." format:"uuid"`
		Body     struct {
			Key    string `json:"key" doc:"장치에서 보내는 데이터 JSON의 키값 입니다." example:"soil_temp"`
			Target string `json:"target" doc:"지정한 키를 바인딩 할 센서 명칭 입니다." example:"토양 온도"`
		}
	}) (*struct{}, error) {
		var param domain.DeviceRequestSchema

		param.Key = i.Body.Key
		param.Target = i.Body.Target

		if err := deviceUseCase.CreateDeviceReqeustSchema(ctx, i.DeviceID, &param); err != nil {
			return nil, huma.Error400BadRequest(err.Error())
		}

		return nil, nil
	})

	// 장치 요청 스키마 리스트 조회 API
	huma.Register(v1, m.WithAuth(
		huma.Operation{
			OperationID:   "v1DeviceGetDeviceReqeustSchemaListByDeviceID",
			Method:        http.MethodGet,
			Path:          "/device/{device_id}/schema",
			Summary:       "장치 요청 스키마 리스트 조회 By DeviceID",
			Description:   "장치 요청 스키마 리스트 조회 By DeviceID API 입니다.",
			Tags:          []string{"Device"},
			DefaultStatus: http.StatusOK,
		},
		domain.UserRoleDevice,
	), func(ctx context.Context, i *struct {
		DeviceID string `path:"device_id" doc:"장치 정보 고유 ID 입니다." format:"uuid"`
	}) (*deviceRequestSchemaListResponse, error) {
		var resp deviceRequestSchemaListResponse

		schemaList, err := deviceUseCase.GetDeviceReqeustSchemaListByID(ctx, i.DeviceID)
		if err != nil {
			return nil, huma.Error400BadRequest(err.Error())
		}

		resp.Body.Data = schemaList
		return &resp, nil
	})

	// 장치 요청 스키마 수정 By ID API
	huma.Register(v1, m.WithAuth(
		huma.Operation{
			OperationID:   "v1DeviceUpdateDeviceReqeustSchemaByDeviceIDAndSchemaID",
			Method:        http.MethodPut,
			Path:          "/device/{device_id}/schema/{schema_id}",
			Summary:       "장치 요청 스키마 수정 By ID",
			Description:   "장치 요청 스키마 수정 By ID API 입니다.",
			Tags:          []string{"Device"},
			DefaultStatus: http.StatusOK,
		},
		domain.UserRoleDevice,
	), func(ctx context.Context, i *struct {
		DeviceID string `path:"device_id" doc:"장치 정보 고유 ID 입니다." format:"uuid"`
		SchemaID int    `path:"schema_id" doc:"요청 스키마 ID 입니다."`
		Body     struct {
			Key    string `json:"key" doc:"바인딩할 요청시 JSON 키 입니다." example:"temp"`
			Target string `json:"target" doc:"바인딩할 key -> sensor title 입니다." example:"기온"`
		}
	}) (*struct{}, error) {
		var param domain.DeviceRequestSchema
		param.ID = i.SchemaID
		param.Key = i.Body.Key
		param.Target = i.Body.Target

		if err := deviceUseCase.UpdateDeviceReqeustSchemaByID(ctx, &param, i.DeviceID); err != nil {
			return nil, huma.Error400BadRequest(err.Error())
		}

		return nil, nil
	})

	// 장치 API 키 생성 API
	huma.Register(v1, m.WithAuth(
		huma.Operation{
			OperationID:   "v1DeviceCreateApiKeyByDeviceID",
			Method:        http.MethodPost,
			Path:          "/device/{device_id}/api-key",
			Summary:       "장치 Api Key 생성 By ID",
			Description:   "장치 Api Key 생성 By ID API 입니다.",
			Tags:          []string{"Device"},
			DefaultStatus: http.StatusCreated,
		},
		domain.UserRoleDevice,
	), func(ctx context.Context, i *struct {
		DeviceID string `path:"device_id" doc:"장치 정보 고유 ID 입니다." format:"uuid"`
		Body     struct {
			Title string `json:"title" maxLength:"40" doc:"사용자가 API키 식별에 사용되는 값입니다." example:"안양시 토마토 농장"`
			Desc  string `json:"desc,omitempty" doc:"사용자가 추가로 남길 API키에 대한 설명입니다." example:"A-B1 섹터 2층 1번 센서"`
		}
	}) (*firstDeviceApiKeyResponse, error) {
		var resp firstDeviceApiKeyResponse

		param := &domain.ApiKey{
			DeviceID: i.DeviceID,
			Title:    i.Body.Title,
			Desc:     i.Body.Desc,
		}
		data, err := deviceUseCase.CreateDeviceApiKey(ctx, param)
		if err != nil {
			return nil, huma.Error400BadRequest(err.Error())
		}

		resp.Body.ID = data.ID
		resp.Body.Key = data.Key
		resp.Body.DeviceID = i.DeviceID
		resp.Body.Title = i.Body.Title
		resp.Body.Desc = i.Body.Desc

		return &resp, nil
	})

	log.Info("Device Handler 등록")
}
