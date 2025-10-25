package repository

import (
	"context"
	"errors"

	"github.com/GDH-Project/api/internal/domain"
	"github.com/jackc/pgx/v5/pgxpool"
	"go.uber.org/zap"
)

type metaRepository struct {
	log *zap.Logger
	db  *pgxpool.Pool
}

func (r *metaRepository) GetSensorList(ctx context.Context) ([]*domain.Sensor, error) {
	q := `SELECT id, title, eng_title, description, unit, unit_description FROM device.sensor;`
	rows, err := r.db.Query(ctx, q)
	if err != nil {
		r.log.Error("device.r.GetSensorList() 오류", zap.Error(err))
		return nil, err
	}
	defer rows.Close()

	var sensorList []*domain.Sensor

	for rows.Next() {
		var s domain.Sensor
		if err := rows.Scan(
			&s.ID,
			&s.Title,
			&s.EngTitle,
			&s.Desc,
			&s.Unit,
			&s.UnitDesc,
		); err != nil {
			r.log.Error("device.r.GetSensorList() 오류", zap.Error(err))
			return nil, err
		}

		sensorList = append(sensorList, &s)
	}

	if err = rows.Err(); err != nil {
		r.log.Error("device.r.GetSensorList() 오류", zap.Error(err))
		return nil, err
	}

	return sensorList, nil
}

func (r *metaRepository) GetSensorByParam(ctx context.Context, in *domain.Sensor) (*domain.Sensor, error) {

	if in.ID == 0 && in.Title == "" {
		r.log.Error("ID 혹은 Title은 필수 입니다")
		return nil, errors.New("ID 혹은 Title은 필수 입니다")
	}
	var sensor domain.Sensor
	q := `
			SELECT id, title, eng_title, description, unit, unit_description 
				FROM device.sensor 
				WHERE 
				    id = NULLIF($1, 0) 
					OR title = NULLIF($2, '')::TEXT;
		`
	if err := r.db.QueryRow(ctx, q,
		in.ID,
		in.Title,
	).Scan(
		&sensor.ID,
		&sensor.Title,
		&sensor.EngTitle,
		&sensor.Desc,
		&sensor.Unit,
		&sensor.UnitDesc,
	); err != nil {
		r.log.Error("device.r.GetSensorByParam() 오류", zap.Error(err))
		return nil, err
	}

	return &sensor, nil
}

func (r *metaRepository) GetCropList(ctx context.Context) ([]*domain.Crop, error) {
	var cropList []*domain.Crop

	q := `SELECT id, title, description FROM device.crop`
	rows, err := r.db.Query(ctx, q)
	if err != nil {
		r.log.Error("device.r.GetCropList() 오류", zap.Error(err))
		return nil, err
	}
	defer rows.Close()

	for rows.Next() {
		var crop domain.Crop
		if err := rows.Scan(
			&crop.ID,
			&crop.Title,
			&crop.Desc,
		); err != nil {
			r.log.Error("device.r.GetCropList() 오류", zap.Error(err))
			return nil, err
		}

		cropList = append(cropList, &crop)
	}

	if err := rows.Err(); err != nil {
		r.log.Error("device.r.GetCropList() 오류", zap.Error(err))
		return nil, err
	}

	return cropList, nil
}

func (r *metaRepository) GetCropByParam(ctx context.Context, in *domain.Crop) (*domain.Crop, error) {
	var crop domain.Crop
	q := `
		SELECT id, title, description 
			FROM device.crop 
			WHERE
			    id = NULLIF($1, 0)
				OR title = NULLIF($2, '')::TEXT;           
	   `
	if err := r.db.QueryRow(ctx, q,
		in.ID,
		in.Title,
	).Scan(
		&crop.ID,
		&crop.Title,
		&crop.Desc,
	); err != nil {
		r.log.Error("device.r.GetCropByParam() 오류", zap.Error(err))
		return nil, err
	}

	return &crop, nil
}

func (r *metaRepository) GetUpdateCycleList(ctx context.Context) ([]*domain.UpdateCycle, error) {
	var updateCycleList []*domain.UpdateCycle
	q := `SELECT id, interval, description FROM device.update_cycle`
	rows, err := r.db.Query(ctx, q)
	if err != nil {
		r.log.Error("device.r.GetUpdateCycleList() 오류", zap.Error(err))
		return nil, err
	}
	defer rows.Close()

	for rows.Next() {
		var updateCycle domain.UpdateCycle
		if err := rows.Scan(
			&updateCycle.ID,
			&updateCycle.Interval,
			&updateCycle.Desc,
		); err != nil {
			r.log.Error("device.r.GetUpdateCycleList() 오류", zap.Error(err))
			return nil, err
		}

		updateCycleList = append(updateCycleList, &updateCycle)
	}

	if err := rows.Err(); err != nil {
		r.log.Error("device.r.GetUpdateCycleList() 오류", zap.Error(err))
		return nil, err
	}

	return updateCycleList, nil
}

func (r *metaRepository) GetAddressStateList(ctx context.Context) ([]*domain.AddressState, error) {
	var addressStateList []*domain.AddressState
	q := `SELECT id,title FROM device.address_state`
	rows, err := r.db.Query(ctx, q)
	if err != nil {
		r.log.Error("device.r.GetAddressStateList() 오류", zap.Error(err))
		return nil, err
	}
	defer rows.Close()

	for rows.Next() {
		var addressState domain.AddressState
		if err := rows.Scan(
			&addressState.ID,
			&addressState.Title,
		); err != nil {
			r.log.Error("device.r.GetAddressStateList() 오류", zap.Error(err))
			return nil, err
		}
		addressStateList = append(addressStateList, &addressState)
	}

	if err := rows.Err(); err != nil {
		r.log.Error("device.r.GetAddressStateList() 오류", zap.Error(err))
		return nil, err
	}

	return addressStateList, nil
}

func (r *metaRepository) GetAddressCityListByState(ctx context.Context, state string) ([]*domain.AddressCity, error) {
	var addressCityList []*domain.AddressCity
	q := `
		SELECT 
		    c.id,
		    s.title AS stateTitle,
		    c.title
		FROM 
		    device.address_city  c
		JOIN  device.address_state s  ON c.address_state_id = s.id
		WHERE s.title = $1;
		`
	rows, err := r.db.Query(ctx, q,
		state,
	)

	if err != nil {
		r.log.Error("device.r.GetAddressCityListByState() 오류", zap.Error(err))
		return nil, err
	}
	defer rows.Close()

	for rows.Next() {
		var addressCity domain.AddressCity
		if err := rows.Scan(
			&addressCity.ID,
			&addressCity.StateTitle,
			&addressCity.Title,
		); err != nil {
			r.log.Error("device.r.GetAddressCityListByState() 오류", zap.Error(err))
			return nil, err
		}

		addressCityList = append(addressCityList, &addressCity)
	}

	if err := rows.Err(); err != nil {
		r.log.Error("device.r.GetAddressCityListByState() 오류", zap.Error(err))
		return nil, err
	}

	return addressCityList, nil
}

func MetaRepository(logger *zap.Logger, db *pgxpool.Pool) domain.MetaRepository {
	return &metaRepository{
		log: logger,
		db:  db,
	}
}
