package repository

import (
	"context"
	"errors"

	"github.com/GDH-Project/api/internal/domain"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
	"go.uber.org/zap"
)

type deviceRepository struct {
	log *zap.Logger
	db  *pgxpool.Pool
}

func (r *deviceRepository) UpdateDeviceInfo(ctx context.Context, in *domain.RawDeviceInfo) error {
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
			    id = $1::UUID;			    
		`

	if _, err := r.db.Exec(ctx, q,
		in.ID,
		in.Title,
		in.Name,
		in.CropID,
		in.UpdateCycleID,
		in.AddressStateID,
		in.AddressCityID,
	); err != nil {
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

	r.log.Info("총 레코드 수",
		zap.Int("count", count),
		zap.Int("pageSize", page.Size),
		zap.Float64("계산된 값", float64(count)/float64(page.Size)),
	)

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
			&deviceInfo.Address.City,
			&deviceInfo.Address.State,
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
		&deviceInfo.Address.City,
		&deviceInfo.Address.State,
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
	tx, _ := r.db.Begin(ctx)

	// 트랜잭션중 패닉 오류 발생시 트랜잭션을 롤백 시키고 패닉을 다시 발생시킨다.
	defer func() {
		if r := recover(); r != nil {
			_ = tx.Rollback(ctx)
			panic(r)
		}
	}()

	// 만약 f(tx)를 실행후 오류 발생시 트랜잭션을 롤백 시키고 오류를 반환한다.
	if err := f(tx); err != nil {
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
