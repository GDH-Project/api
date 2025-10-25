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

func (r *deviceRepository) CreateDeviceInfoTx(ctx context.Context, tx pgx.Tx, in *domain.CreateDeviceInfo) (string, error) {
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
