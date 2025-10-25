package repository

import (
	"context"

	"github.com/GDH-Project/api/internal/domain"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
	"go.uber.org/zap"
)

type deviceRepository struct {
	log *zap.Logger
	db  *pgxpool.Pool
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

	defer func() {
		if r := recover(); r != nil {
			_ = tx.Rollback(ctx)
			panic(r)
		}
	}()

	if err := f(tx); err != nil {
		_ = tx.Rollback(ctx)
		return err
	}

	return tx.Commit(ctx)
}

func NewDeviceRepository(log *zap.Logger, db *pgxpool.Pool) domain.DeviceRepository {
	return &deviceRepository{
		log: log,
		db:  db,
	}

}
