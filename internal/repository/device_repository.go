package repository

import (
	"context"
	"errors"
	"fmt"
	"strings"

	"github.com/GDH-Project/api/internal/domain"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
	"go.uber.org/zap"
)

type deviceRepository struct {
	log *zap.Logger
	db  *pgxpool.Pool
}

func (r *deviceRepository) GetDeviceApiKeyByApiKey(ctx context.Context, apiKey string) (string, error) {
	var deviceID string
	q := `SELECT device_info_id FROM device.api_key WHERE api_key = $1 ;`
	if err := r.db.QueryRow(ctx, q, apiKey).Scan(&deviceID); err != nil {
		r.log.Info("device.r.GetDeviceApiKeyByApiKey() 오류", zap.Error(err))
		return "", err
	}

	return deviceID, nil
}

func (r *deviceRepository) CreateDeviceApiKey(ctx context.Context, in *domain.RawApiKey) (*domain.RawApiKey, error) {
	var apiKey string
	var id int
	q := `INSERT INTO device.api_key(api_key, user_id, device_info_id, title, description) VALUES ($1, $2, $3, $4, $5) RETURNING id, api_key;`
	if err := r.db.QueryRow(ctx, q,
		in.APIKey,
		in.UserID,
		in.DeviceID,
		in.Title,
		in.Desc,
	).Scan(&id, &apiKey); err != nil {
		r.log.Info("device.r.CreateDeviceApiKey() 오류", zap.Error(err))
		return nil, err
	}

	// 생성후 반환된 API키가 전달된 키와 일치하지 않는 경우 --> 사실상 인서트 실패
	if apiKey != in.APIKey {
		err := errors.New("API키 생성에 실패했습니다")
		r.log.Info("device.r.CreateDeviceApiKey() 오류",
			zap.Any("apiKey", in.APIKey),
			zap.Any("data", in),
			zap.Error(err),
		)
		return nil, err
	}

	in.ID = id

	return in, nil
}

func (r *deviceRepository) DeleteDeviceRequestSchemaByID(ctx context.Context, id int) error {
	var successID int
	q := `DELETE FROM device.req_to_sensor WHERE id = $1 RETURNING id;`
	if err := r.db.QueryRow(ctx, q, id).Scan(&successID); err != nil {
		r.log.Info("device.r.DeleteDeviceRequestSchemaByID() 오류", zap.Error(err))
		return err
	}

	return nil
}

func (r *deviceRepository) UpdateDeviceRequestSchema(ctx context.Context, in *domain.RawDeviceRequestSchema) error {
	var successID int
	q := `
			UPDATE device.req_to_sensor 
			SET
			    key = COALESCE(NULLIF($2,''), key),
			    sensor_id = COALESCE(NULLIF($3,0), sensor_id)
			WHERE 
			    id = $1
			RETURNING id;
		`
	if err := r.db.QueryRow(ctx, q,
		in.ID,
		in.Key,
		in.TargetSensorID,
	).Scan(&successID); err != nil {
		r.log.Error("device.r.UpdateDeviceRequestSchema() 오류", zap.Error(err))
		return err
	}

	return nil
}

func (r *deviceRepository) GetDeviceRequestSchemaByID(ctx context.Context, id int) (*domain.DeviceRequestSchema, error) {
	var deviceSchema domain.DeviceRequestSchema
	q := `
			SELECT rq.id, rq.key ,s.title
    		FROM device.req_to_sensor rq
    			JOIN device.sensor s ON rq.sensor_id = s.id
    		WHERE rq.id = $1
    	`
	if err := r.db.QueryRow(ctx, q, id).Scan(
		&deviceSchema.ID,
		&deviceSchema.Key,
		&deviceSchema.Target,
	); err != nil {
		r.log.Info("device.r.GetDeviceRequestSchemaByID() 오류", zap.Error(err))
		return nil, err
	}

	return &deviceSchema, nil
}

func (r *deviceRepository) GetDeviceRequestSchemaListByDeviceID(ctx context.Context, deviceID string) ([]*domain.DeviceRequestSchema, error) {
	var deviceSchemaList []*domain.DeviceRequestSchema
	q := `
			SELECT rq.id, rq.key ,s.title
    		FROM device.req_to_sensor rq
    			JOIN device.sensor s ON rq.sensor_id = s.id
    		WHERE rq.device_id = $1
    	`

	rows, err := r.db.Query(ctx, q, deviceID)
	if err != nil {
		r.log.Error("device.r.GetDeviceRequestSchemaListByDeviceID() 오류", zap.Error(err))
		return nil, err
	}
	defer rows.Close()

	for rows.Next() {
		var deviceRequestSchema domain.DeviceRequestSchema
		if err := rows.Scan(
			&deviceRequestSchema.ID,
			&deviceRequestSchema.Key,
			&deviceRequestSchema.Target,
		); err != nil {
			r.log.Error("device.r.GetDeviceRequestSchemaListByDeviceID() 오류", zap.Error(err))
			return nil, err
		}
		deviceSchemaList = append(deviceSchemaList, &deviceRequestSchema)
	}

	if err := rows.Err(); err != nil {
		r.log.Error("device.r.GetDeviceRequestSchemaListByDeviceID() 오류", zap.Error(err))
		return nil, err
	}

	return deviceSchemaList, nil
}

func (r *deviceRepository) CreateDeviceReqeustSchemaListTx(ctx context.Context, tx pgx.Tx, deviceID string, schemas []*domain.RawDeviceRequestSchema) error {
	// deviceID 가 없거나 스키마가 없는 경우 필터링
	if deviceID == "" || len(schemas) == 0 {
		err := errors.New("파라미터를 확인해 주세요")
		r.log.Info("device.r.CreateDeviceReqeustSchemaListTx() 오류",
			zap.String("deviceID", deviceID),
			zap.Any("schemas", schemas),
			zap.Error(err),
		)

		return err
	}

	// 쿼리에 필요한 총 $1 같은 플레이스홀더 갯수
	totalPlaceholders := len(schemas) * 3

	valuePlaceholders := make([]string, len(schemas))
	args := make([]interface{}, 0, totalPlaceholders)

	for i, item := range schemas {
		// $1, $2 행 문자
		ph1 := i*3 + 1
		ph2 := i*3 + 2
		ph3 := i*3 + 3

		valuePlaceholders[i] = fmt.Sprintf("($%d::uuid, $%d, $%d)", ph1, ph2, ph3)

		args = append(args, deviceID, item.Key, item.TargetSensorID)
	}

	sb := strings.Builder{}
	sb.WriteString("INSERT INTO device.req_to_sensor(device_id, key, sensor_id) VALUES ")
	sb.WriteString(strings.Join(valuePlaceholders, ","))
	sb.WriteString(";")

	q := sb.String()

	_, err := tx.Exec(ctx, q, args...)
	if err != nil {
		r.log.Error("device.r.CreateDeviceReqeustSchemaListTx() 오류",
			zap.String("deviceID", deviceID),
			zap.String("q", q),
			zap.Any("args", args),

			zap.Error(err),
		)
		return err
	}

	return nil
}

func (r *deviceRepository) DeleteDeviceInfoByID(ctx context.Context, id string, userID string) error {
	var successID string
	q := `
			UPDATE device.device_info
			SET deleted_at = NOW()
			WHERE 
			    deleted_at IS NULL 
			  AND 
			    id = $1::uuid 
			  AND 
			    user_id = $2::uuid
			RETURNING id;
		`
	if err := r.db.QueryRow(ctx, q, id, userID).Scan(&successID); err != nil {
		r.log.Error("device.r.DeleteDeviceInfoByID() 오류", zap.String("id", id), zap.Error(err))
		return err
	}

	if successID == "" || id != successID {
		err := errors.New("id에 해당하는 필드가 존재하지 않습니다")
		r.log.Info("device.r.DeleteDeviceInfoByID() 오류",
			zap.String("id", id),
			zap.Error(err),
		)
		return err
	}

	return nil
}

func (r *deviceRepository) UpdateDeviceInfo(ctx context.Context, in *domain.RawDeviceInfo) error {
	var successID string
	q := `
			UPDATE device.device_info 
			SET
				title = COALESCE(NULLIF($2, ''), title),
				device_name = COALESCE(NULLIF($3, ''), device_name),
				-- 작물 정보도 변경이 가능하지만 변경 하지 못하게 조치 필요
				crop_id = COALESCE(NULLIF($4, 0), crop_id),
			    update_cycle_id = COALESCE(NULLIF($5, 0), update_cycle_id),
			    address_state_id = COALESCE(NULLIF($6, 0), address_state_id),
			    address_city_id = COALESCE(NULLIF($7, 0), address_city_id)
			WHERE 
			    deleted_at IS NULL 
			  AND
			    id = $1::UUID
			RETURNING id;			    
		`

	if err := r.db.QueryRow(ctx, q,
		in.ID,
		in.Title,
		in.Name,
		in.CropID,
		in.UpdateCycleID,
		in.AddressStateID,
		in.AddressCityID,
	).Scan(&successID); err != nil {
		r.log.Info("device.r.UpdateDeviceInfo() 오류",
			zap.Error(err),
			zap.Any("data", in),
		)
		return err
	}

	if successID == "" {
		err := errors.New("id에 해당하는 필드가 존재하지 않습니다")
		r.log.Info("device.r.UpdateDeviceInfo() 오류",
			zap.Error(err),
			zap.Any("data", in),
		)

		return err
	}
	return nil
}

func (r *deviceRepository) GetDeviceInfoListByParamAndPage(ctx context.Context, in *domain.DeviceInfo, page *domain.Page) ([]*domain.DeviceInfo, *domain.Page, error) {
	var deviceInfoList []*domain.DeviceInfo
	var count int
	// info 수 카운트
	q := `
			SELECT COUNT(info.id)
			FROM device.device_info info
				JOIN device.crop crop ON info.crop_id = crop.id
				JOIN device.update_cycle uc ON info.update_cycle_id = uc.id
			    JOIN device.address_state state ON info.address_state_id = state.id
				JOIN device.address_city city ON info.address_city_id = city.id
			WHERE info.deleted_at IS NULL 
			    AND (NULLIF($1, '') IS NULL OR info.user_id = NULLIF($1, '')::uuid)
				AND (NULLIF($2, '') IS NULL OR info.title LIKE '%' || $2 || '%')
				AND (NULLIF($3, '') IS NULL OR crop.title = NULLIF($3, ''))
				AND (NULLIF($4, '') IS NULL OR state.title = NULLIF($4, ''))
				AND (NULLIF($5, '') IS NULL OR city.title = NULLIF($5, ''))			          
		`
	if err := r.db.QueryRow(ctx, q,
		in.UserID,
		in.Title,
		in.Crop,
		in.Address.State,
		in.Address.City,
	).Scan(&count); err != nil {
		r.log.Error("device.r.GetDeviceInfoListByParamAndPage() 오류", zap.Error(err))
		return nil, nil, err
	}

	q = `
			SELECT 
			    info.user_id,
			    info.id,
			    info.title,
			    info.device_name,
			    crop.title,
			    uc.interval,
			    state.title,
			    city.title,
			    info.created_at,
			    info.updated_at
			FROM device.device_info info
				JOIN device.crop crop ON info.crop_id = crop.id
				JOIN device.update_cycle uc ON info.update_cycle_id = uc.id
			    JOIN device.address_state state ON info.address_state_id = state.id
				JOIN device.address_city city ON info.address_city_id = city.id
			WHERE info.deleted_at IS NULL 
			    -- $1 (user_id): 값이 있을 때만 user_id 필터링
			    AND (NULLIF($1, '') IS NULL OR info.user_id = NULLIF($1, '')::uuid)
				-- $2 (title LIKE): 값이 있을 때만 LIKE 검색
				AND (NULLIF($2, '') IS NULL OR info.title LIKE '%' || $2 || '%')
				-- $3 (crop.title): 값이 있을 때만 crop.title 필터링
				AND (NULLIF($3, '') IS NULL OR crop.title = NULLIF($3, ''))
				-- $4 (state.title): 값이 있을 때만 state.title 필터링
				AND (NULLIF($4, '') IS NULL OR state.title = NULLIF($4, ''))
				-- $5 (city.title): 값이 있을 때만 city.title 필터링
				AND (NULLIF($5, '') IS NULL OR city.title = NULLIF($5, ''))			          
			ORDER BY INFO.created_at DESC
			LIMIT $6 OFFSET $7;
		`
	rows, err := r.db.Query(ctx, q,
		in.UserID,
		in.Title,
		in.Crop,
		in.Address.State,
		in.Address.City,
		page.Size,
		page.Offset(),
	)
	if err != nil {
		r.log.Error("device.r.GetDeviceInfoListByParamAndPage() 오류", zap.Error(err))
		return nil, nil, err
	}
	defer rows.Close()

	for rows.Next() {
		var deviceInfo domain.DeviceInfo
		if err := rows.Scan(
			&deviceInfo.UserID,
			&deviceInfo.ID,
			&deviceInfo.Title,
			&deviceInfo.Name,
			&deviceInfo.Crop,
			&deviceInfo.UpdateCycle,
			&deviceInfo.Address.State,
			&deviceInfo.Address.City,
			&deviceInfo.CreatedAt,
			&deviceInfo.UpdatedAt,
		); err != nil {
			r.log.Error("device.r.GetDeviceInfoListByParamAndPage() 오류", zap.Error(err))
			return nil, nil, err
		}

		deviceInfoList = append(deviceInfoList, &deviceInfo)
	}

	if err := rows.Err(); err != nil {
		r.log.Error("device.r.GetDeviceInfoListByParamAndPage() 오류", zap.Error(err))
		return nil, nil, err
	}

	// 페이징 정보
	p := &domain.Page{
		Size:        page.Size,
		Page:        page.Page,
		RecordCount: count,
	}
	return deviceInfoList, p, nil
}

func (r *deviceRepository) GetDeviceInfoByID(ctx context.Context, id string) (*domain.DeviceInfo, error) {
	var deviceInfo domain.DeviceInfo
	q := `
			SELECT 
			    info.user_id,
			    info.id,
			    info.title,
			    info.device_name,
			    crop.title,
			    uc.interval,
			    state.title,
			    city.title,
			    info.created_at,
			    info.updated_at
			FROM device.device_info info
				JOIN device.crop crop ON info.crop_id = crop.id
				JOIN device.update_cycle uc ON info.update_cycle_id = uc.id
			    JOIN device.address_state state ON info.address_state_id = state.id
				JOIN device.address_city city ON info.address_city_id = city.id
			WHERE info.deleted_at IS NULL 
				AND info.id = NULLIF($1,'')::UUID
			;
		`
	if err := r.db.QueryRow(ctx, q, id).Scan(
		&deviceInfo.UserID,
		&deviceInfo.ID,
		&deviceInfo.Title,
		&deviceInfo.Name,
		&deviceInfo.Crop,
		&deviceInfo.UpdateCycle,
		&deviceInfo.Address.State,
		&deviceInfo.Address.City,
		&deviceInfo.CreatedAt,
		&deviceInfo.UpdatedAt,
	); err != nil {
		// 만약 ID가 존재하지 않아 반환값이 없는경우
		if errors.Is(err, pgx.ErrNoRows) {
			r.log.Info("device.r.GetDeviceInfoByID() id에 해당하는 장비가 존재하지 않습니다", zap.String("id", id))
			return nil, err
		}

		r.log.Error("device.r.GetDeviceInfoByID() 오류", zap.Error(err))
		return nil, err
	}

	return &deviceInfo, nil
}

func (r *deviceRepository) CreateDeviceInfoTx(ctx context.Context, tx pgx.Tx, in *domain.RawDeviceInfo) (string, error) {
	var id string
	q := `
			INSERT INTO 
			    device.device_info(
			                       title,
			                       device_name,
			                       user_id,
			                       crop_id,
			                       update_cycle_id,
			                       address_state_id,
			                       address_city_id
			   ) VALUES ( $1, NULLIF($2,''),$3::uuid , $4, $5, $6, $7)
			RETURNING id;
		`
	if err := tx.QueryRow(ctx, q,
		in.Title,
		in.Name,
		in.UserID,
		in.CropID,
		in.UpdateCycleID,
		in.AddressStateID,
		in.AddressCityID,
	).Scan(&id); err != nil {
		r.log.Error("device.r.CreateDeviceInfoTx() 오류", zap.Error(err))
		return "", err
	}

	return id, nil
}

func (r *deviceRepository) WithTransaction(ctx context.Context, f func(tx pgx.Tx) error) error {
	tx, err := r.db.Begin(ctx)
	if err != nil {
		r.log.Error("device.r.WithTransaction() 트랜잭션 시작 오류", zap.Error(err))
		return err
	}

	// 트랜잭션중 패닉 오류 발생시 트랜잭션을 롤백 시키고 패닉을 다시 발생시킨다.
	defer func() {
		if rc := recover(); rc != nil {
			r.log.Info("device.r.WithTransaction() 트랜잭션 오류 발생으로 인한 롤백", zap.Error(errors.New("트랜잭션 도중 패닉 발생")))
			_ = tx.Rollback(ctx)
			panic(rc)
		}
	}()

	// 만약 f(tx)를 실행후 오류 발생시 트랜잭션을 롤백 시키고 오류를 반환한다.
	if err := f(tx); err != nil {
		r.log.Info("device.r.WithTransaction() 트랜잭션 오류 발생으로 인한 롤백", zap.Error(err))
		_ = tx.Rollback(ctx)
		return err
	}

	// 트랜잭션을 커밋하고 커밋시 오류를 반환한다.
	return tx.Commit(ctx)
}

func NewDeviceRepository(log *zap.Logger, db *pgxpool.Pool) domain.DeviceRepository {
	return &deviceRepository{
		log: log,
		db:  db,
	}

}
