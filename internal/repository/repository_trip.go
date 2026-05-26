package repository

import (
	"context"
	"fmt"
	"github.com/google/uuid"
	"github.com/jackc/pgx/v5/pgxpool"
	shareTrip "job4j.ru/share-trip/internal/domain"
)

type RepoPg struct {
	pool *pgxpool.Pool
}

func NewRepoPg(pool *pgxpool.Pool) *RepoPg {
	return &RepoPg{pool: pool}
}

func (r *RepoPg) Create(ctx context.Context, it shareTrip.Trip) error {
	// запись в основную таблицу
	_, err := r.pool.Exec(
		ctx,
		`insert into trips(id, driver_id, from_point, to_point, departure_time, seats) values($1, $2, $3, $4, $5, $6)`,
		it.ID, it.DriverId, it.FromPoint, it.ToPoint, it.DepartureTime, it.AvailableSeats,
	)
	if err != nil {
		return fmt.Errorf("r.pool.Exec: %w", err)
	}
	// запись в историческую таблицу
	id := uuid.New().String()
	_, err = r.pool.Exec(
		ctx,
		`insert into trip_history(id, trip_id, to_status) values($1, $2, $3)`,
		id, it.ID, "draft",
	)
	if err != nil {
		return fmt.Errorf("r.pool.Exec: %w", err)
	}

	return nil
}

func (r *RepoPg) List(ctx context.Context) ([]shareTrip.Trip, error) {
	rows, err := r.pool.Query(ctx, `select id, name from trips`)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var trips []shareTrip.Trip
	for rows.Next() {
		var item shareTrip.Trip
		if err := rows.Scan(&item.ID, &item.DriverId, &item.FromPoint, &item.ToPoint, &item.DepartureTime, &item.AvailableSeats); err != nil {
			return nil, err
		}
		trips = append(trips, item)
	}

	if err := rows.Err(); err != nil {
		return nil, err
	}

	return trips, nil
}

func (r *RepoPg) Get(ctx context.Context, tripId string) (shareTrip.Trip, error) {
	var it shareTrip.Trip
	err := r.pool.QueryRow(
		ctx,
		`select id, driver_id, from_point, to_point, COALESCE(to_char(departure_time, 'MM-DD-YYYY HH24:MI'), ''), seats from trips where id = $1`,
		tripId,
	).Scan(&it.ID, &it.DriverId, &it.FromPoint, &it.ToPoint, &it.DepartureTime, &it.AvailableSeats)

	return it, err
}

func (r *RepoPg) Update(ctx context.Context, name string, newName string) error {
	_, err := r.pool.Exec(
		ctx,
		"UPDATE trips SET name = $2 WHERE name = $1",
		name, newName,
	)
	if err != nil {
		return fmt.Errorf("r.pool.Exec: %w", err)
	}

	return nil
}

func (r *RepoPg) Delete(ctx context.Context, name string) error {
	_, err := r.pool.Exec(
		ctx,
		"DELETE trips WHERE name = $1",
		name,
	)
	if err != nil {
		return fmt.Errorf("r.pool.Exec: %w", err)
	}

	return nil
}

func (r *RepoPg) GetCount(ctx context.Context) (string, error) {
	var count string
	err := r.pool.QueryRow(
		ctx,
		`select count(*) from trips`,
	).Scan(&count)

	return count, err
}

func (r *RepoPg) DoPing(ctx context.Context) error {
	err := r.pool.Ping(ctx)
	return err
}
