package handler

import (
	"context"
	"errors"
	"net/http"

	"github.com/GDH-Project/api/internal/domain"
	"github.com/GDH-Project/api/internal/util"
	"github.com/danielgtaylor/huma/v2"
	"github.com/jackc/pgx/v5"
	"go.uber.org/zap"
)

type sensorListResponse struct {
	util.CacheHeader
	Body struct {
		Data []*domain.Sensor `json:"data" doc:"센서 정보 JSON 배열 입니다."`
	}
}
type sensorResponse struct {
	util.CacheHeader
	Body struct {
		Data *domain.Sensor `json:"data" doc:"센서 정보 JSON 입니다."`
	}
}

// addressStateList 도/특별시 리스트 응답 구조체
type addressStateListResponse struct {
	util.CacheHeader
	Body struct {
		Data []*domain.AddressState `json:"data" doc:"도/특별시 주소 정보 JSON 배열 입니다."`
	}
}

// addressCityList 시/군/구 응답 구조체
type addressCityListResponse struct {
	util.CacheHeader
	Body struct {
		Data []*domain.AddressCity `json:"data" doc:"시/군/구 주소 정보 JSON 배열 입니다."`
	}
}

type cropListResponse struct {
	util.CacheHeader
	Body struct {
		Data []*domain.Crop `json:"data" doc:"작물 정보 JSON 배열 입니다."`
	}
}

type cropResponse struct {
	util.CacheHeader
	Body struct {
		Data *domain.Crop `json:"data" doc:"작물 정보 JSON 입니다."`
	}
}

type updateCycleListResponse struct {
	util.CacheHeader
	Body struct {
		Data []*domain.UpdateCycle `json:"data" doc:"업데이트 주기 정보 JSON 배열 입니다."`
	}
}

func RegisterMetaHandler(api huma.API, log *zap.Logger, metaUseCase domain.MetaUseCase) {
	v1 := huma.NewGroup(api, "/api/v1")

	// 센서 정보 전체 조회 API
	huma.Register(v1, huma.Operation{
		OperationID:   "v1MetaGetSensorList",
		Method:        http.MethodGet,
		Path:          "/meta/sensors",
		Summary:       "전체 센서 정보 조회",
		Description:   "전체 센서 정보 조회 API 입니다.",
		Tags:          []string{"Meta"},
		DefaultStatus: http.StatusOK,
	}, func(ctx context.Context, i *struct{}) (*sensorListResponse, error) {
		var resp sensorListResponse
		sensorList, err := metaUseCase.GetSensorList(ctx)
		if err != nil || len(sensorList) == 0 {
			log.Error("meta.h.v1MetaGetSensorList 오류", zap.Error(err))
			return nil, huma.Error500InternalServerError("전체 센서 데이터를 불러오는 도중 오류가 발생했습니다.")
		}

		resp.Body.Data = sensorList

		cacheHeader := util.CacheHeaderBuilder{
			CacheType: util.CacheTypePublic,
			TTL:       60,
		}
		resp.CacheControl = cacheHeader.String()

		return &resp, nil
	})

	// 센서 ID로 센서 조회 API
	huma.Register(v1, huma.Operation{
		OperationID:   "v1MetaGetSensorByID",
		Method:        http.MethodGet,
		Path:          "/meta/sensor/{id}",
		Summary:       "센서 정보 조회 by ID",
		Description:   "센서 정보 조회 by ID API 입니다.",
		Tags:          []string{"Meta"},
		DefaultStatus: http.StatusOK,
	}, func(ctx context.Context, i *struct {
		ID int `path:"id" required:"true" doc:"센서 ID 입니다." example:"1"`
	}) (*sensorResponse, error) {
		var resp sensorResponse
		sensor, err := metaUseCase.GetSensorByParam(ctx, &domain.Sensor{ID: i.ID})
		if err != nil {
			if errors.Is(err, pgx.ErrNoRows) {
				log.Info("meta.h.v1MetaGetSensorByID 잘못된 ID 검색 발생",
					zap.Int("id", i.ID),
					zap.Error(err))
				return nil, huma.Error400BadRequest("존재하지 않는 센서 ID 입니다.")
			}
			log.Error("meta.h.v1MetaGetSensorByID 오류", zap.Error(err))
			return nil, huma.Error500InternalServerError("센서 데이터 불러오는 도중 오류가 발생했습니다.")
		}

		resp.Body.Data = sensor

		cacheHeader := util.CacheHeaderBuilder{
			CacheType: util.CacheTypePublic,
			TTL:       60,
		}
		resp.CacheControl = cacheHeader.String()

		return &resp, nil
	})

	// 주소 도/특별시 조회
	huma.Register(v1, huma.Operation{
		OperationID:   "v1MetaGetAddressState",
		Method:        http.MethodGet,
		Path:          "/meta/address/state",
		Summary:       "주소 도/특별시 조회",
		Description:   "주소 도/특별시 조회 API 입니다.",
		Tags:          []string{"Meta"},
		DefaultStatus: http.StatusOK,
	}, func(ctx context.Context, i *struct{}) (*addressStateListResponse, error) {
		var resp addressStateListResponse
		addressStateList, err := metaUseCase.GetAddressStateList(ctx)
		if err != nil || len(addressStateList) == 0 {
			log.Error("meta.h.v1MetaGetAddressState 오류", zap.Error(err))
			return nil, huma.Error500InternalServerError("주소 데이터를 불러오는 도중 오류가 발생했습니다.")
		}

		resp.Body.Data = addressStateList

		cacheHeader := util.CacheHeaderBuilder{
			CacheType: util.CacheTypePublic,
			TTL:       1800,
		}
		resp.CacheControl = cacheHeader.String()

		return &resp, nil
	})

	// 주소 시/군/구 조회
	huma.Register(v1, huma.Operation{
		OperationID:   "v1MetaGetAddressCityByStateID",
		Method:        http.MethodGet,
		Path:          "/meta/address/city",
		Summary:       "주소 시/군/구 조회",
		Description:   "주소 시/군/구 조회 API 입니다.",
		Tags:          []string{"Meta"},
		DefaultStatus: http.StatusOK,
	}, func(ctx context.Context, i *struct {
		State string `query:"state" required:"true" doc:"도/특별시 Title 입니다." example:"서울특별시"`
	}) (*addressCityListResponse, error) {
		var resp addressCityListResponse
		addressCityList, err := metaUseCase.GetAddressCityListByState(ctx, i.State)
		if err != nil {
			log.Error("meta.h.v1MetaGetAddressCityByStateID 오류", zap.Error(err))
			return nil, huma.Error500InternalServerError("주소 데이터를 불러오는 도중 오류가 발생했습니다.")
		}

		// 값이 존재하지 않는 경우
		if len(addressCityList) == 0 {
			log.Info("meta.h.v1MetaGetAddressCityByStateID 데이터 조회 실패",
				zap.String("state", i.State),
				zap.Error(err),
			)
			return nil, huma.Error400BadRequest("state에 해당하는 데이터가 존재하지 않습니다.")
		}

		resp.Body.Data = addressCityList

		cacheHeader := util.CacheHeaderBuilder{
			CacheType: util.CacheTypePublic,
			TTL:       1800,
		}
		resp.CacheControl = cacheHeader.String()

		return &resp, nil
	})

	// 전체 작물 조회
	huma.Register(v1, huma.Operation{
		OperationID:   "v1MetaGetCropList",
		Method:        http.MethodGet,
		Path:          "/meta/crops",
		Summary:       "전체 작물 조회",
		Description:   "전체 작물 조회 API 입니다.",
		Tags:          []string{"Meta"},
		DefaultStatus: http.StatusOK,
	}, func(ctx context.Context, i *struct{}) (*cropListResponse, error) {
		var resp cropListResponse
		cropList, err := metaUseCase.GetCropList(ctx)
		if err != nil || len(cropList) == 0 {
			log.Error("meta.h.v1MetaGetCropList 오류", zap.Error(err))
			return nil, huma.Error500InternalServerError("전체 작물 데이터를 불러오는 도중 오류가 발생했습니다.")
		}

		resp.Body.Data = cropList

		cacheHeader := util.CacheHeaderBuilder{
			CacheType: util.CacheTypePublic,
			TTL:       60,
		}
		resp.CacheControl = cacheHeader.String()

		return &resp, nil
	})

	// 작물 조회 By 작물명
	huma.Register(v1, huma.Operation{
		OperationID:   "v1MetaGetCropByTitle",
		Method:        http.MethodGet,
		Path:          "/meta/crop/{title}",
		Summary:       "작물 조회 By 작물명",
		Description:   "작물 조회 By 작물명 API 입니다.",
		Tags:          []string{"Meta"},
		DefaultStatus: http.StatusOK,
	}, func(ctx context.Context, i *struct {
		Title string `path:"title" required:"true" doc:"작물 명칭 입니다." example:"토마토"`
	}) (*cropResponse, error) {
		var resp cropResponse
		crop, err := metaUseCase.GetCropByParam(ctx, &domain.Crop{Title: i.Title})
		if err != nil {
			log.Error("meta.h.v1MetaGetCropList 오류", zap.Error(err))
			return nil, huma.Error500InternalServerError("전체 작물 데이터를 불러오는 도중 오류가 발생했습니다.")
		}

		resp.Body.Data = crop

		cacheHeader := util.CacheHeaderBuilder{
			CacheType: util.CacheTypePublic,
			TTL:       60,
		}
		resp.CacheControl = cacheHeader.String()

		return &resp, nil
	})

	// 갱신 주기 조회
	huma.Register(v1, huma.Operation{
		OperationID:   "v1MetaGetUpdateCycleList",
		Method:        http.MethodGet,
		Path:          "/meta/update-cycle",
		Summary:       "전체 업데이트 주기 조회",
		Description:   "전체 업데이트 주기 조회 API 입니다. 장비의 업데이트 주기에 사용되는 데이터 입니다.",
		Tags:          []string{"Meta"},
		DefaultStatus: http.StatusOK,
	}, func(ctx context.Context, i *struct{}) (*updateCycleListResponse, error) {
		var resp updateCycleListResponse
		updateCycleList, err := metaUseCase.GetUpdateCycleList(ctx)
		if err != nil || len(updateCycleList) == 0 {
			log.Error("meta.h.v1MetaGetUpdateCycleList 오류", zap.Error(err))
			return nil, huma.Error500InternalServerError("업데이트 주기 정보를 불러오는 중 오류가 발생했습니다.")
		}

		resp.Body.Data = updateCycleList

		cacheHeader := util.CacheHeaderBuilder{
			CacheType: util.CacheTypePublic,
			TTL:       1800,
		}
		resp.CacheControl = cacheHeader.String()

		return &resp, nil
	})

	log.Info("Meta Handler 등록")

}
