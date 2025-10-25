package handler

import (
	"context"
	"errors"
	"net/http"

	"github.com/GDH-Project/api/internal/domain"
	"github.com/danielgtaylor/huma/v2"
	"github.com/jackc/pgx/v5"
	"go.uber.org/zap"
)

type sensorListResponse struct {
	Body struct {
		Data []*domain.Sensor `json:"data" doc:"센서 정보를 답은 JSON 배열 입니다."`
	}
}
type sensorResponse struct {
	Body struct {
		Data *domain.Sensor `json:"data" doc:"센서 데이터 JSON 입니다."`
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
		if err != nil {
			log.Error("meta.h.v1MetaGetSensorList 오류", zap.Error(err))
			return nil, huma.Error500InternalServerError("전체 센서 데이터를 불러오는 도중 오류가 발생했습니다.")
		}

		resp.Body.Data = sensorList

		return &resp, nil
	})

	huma.Register(v1, huma.Operation{
		OperationID:   "v1MetaGetSensorByID",
		Method:        http.MethodGet,
		Path:          "/meta/sensor/{id}",
		Summary:       "센서 정보 조회 by ID",
		Description:   "센서 정보 조회 by ID API 입니다.",
		Tags:          []string{"Meta"},
		DefaultStatus: http.StatusOK,
	}, func(ctx context.Context, i *struct {
		ID int `path:"id" doc:"센서 ID 입니다." example:"1"`
	}) (*sensorResponse, error) {
		var resp sensorResponse
		sensor, err := metaUseCase.GetSensorByParam(ctx, &domain.Sensor{ID: i.ID})
		if err != nil {
			if errors.Is(err, pgx.ErrNoRows) {
				log.Info("meta.h.v1MetaGetSensorByID 잘못된 ID 검색 발생", zap.Error(err))
				return nil, huma.Error400BadRequest("존재하지 않는 센서 ID 입니다.")
			}
			log.Error("meta.h.v1MetaGetSensorByID 오류", zap.Error(err))
			return nil, huma.Error500InternalServerError("센서 데이터 불러오는 도중 오류가 발생했습니다.")
		}

		resp.Body.Data = sensor

		return &resp, nil
	})
}
