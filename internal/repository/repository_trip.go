package repository

import (
	"context"
	"errors"
	"fmt"
	"github.com/google/uuid"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
	"job4j.ru/share-trip/internal/dto"
)

type RepoPg struct {
	pool *pgxpool.Pool
}

func NewRepoPg(pool *pgxpool.Pool) *RepoPg {
	return &RepoPg{pool: pool}
}

func (r *RepoPg) Create(ctx context.Context, it dto.Trip) (*dto.Trip, error) {
	// запись в основную таблицу
	_, err := r.pool.Exec(
		ctx,
		`insert into trips(id, driver_id, from_point, to_point, departure_time, seats, status) values($1, $2, $3, $4, $5, $6, $7)`,
		it.ID, it.DriverId, it.FromPoint, it.ToPoint, it.DepartureTime, it.AvailableSeats, it.Status,
	)
	if err != nil {
		return nil, fmt.Errorf("r.pool.Exec: %w", err)
	}
	// запись в историческую таблицу
	id := uuid.New().String()
	_, err = r.pool.Exec(
		ctx,
		`insert into trip_history(id, trip_id, to_status) values($1, $2, $3)`,
		id, it.ID, it.Status,
	)
	if err != nil {
		return nil, fmt.Errorf("r.pool.Exec: %w", err)
	}

	return &it, nil
}

func (r *RepoPg) List(ctx context.Context) ([]dto.Trip, error) {
	rows, err := r.pool.Query(ctx, `select id, name from trips`)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var trips []dto.Trip
	for rows.Next() {
		var item dto.Trip
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

func (r *RepoPg) Get(ctx context.Context, tripId string) (dto.Trip, error) {
	var it dto.Trip
	err := r.pool.QueryRow(
		ctx,
		`select id, driver_id, from_point, to_point, COALESCE(to_char(departure_time, 'MM-DD-YYYY HH24:MI'), ''), seats, status from trips where id = $1`,
		tripId,
	).Scan(&it.ID, &it.DriverId, &it.FromPoint, &it.ToPoint, &it.DepartureTime, &it.AvailableSeats, &it.Status)

	return it, err
}

func (r *RepoPg) GetByID(
	ctx context.Context,
	tx pgx.Tx,
	id string,
) (*dto.Trip, error) {
	trip := &dto.Trip{}

	err := tx.QueryRow(
		ctx,
		`select id, driver_id, from_point, to_point, COALESCE(to_char(departure_time, 'MM-DD-YYYY HH24:MI'), ''), seats, status, created_at from trips where id = $1 `,
		id).Scan(
		&trip.ID,
		&trip.DriverId,
		&trip.FromPoint,
		&trip.ToPoint,
		&trip.DepartureTime,
		&trip.AvailableSeats,
		&trip.Status,
		&trip.CreatedAt,
	)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, ErrTripNotFound
		}
		return nil, fmt.Errorf("query trip by id %s: %w", id, err)
	}

	return trip, nil
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

func (r *RepoPg) UpdateStatus(ctx context.Context, tx pgx.Tx, id string, oldStatus string, newStatus string) error {
	_, err := tx.Exec(ctx, "UPDATE trips SET status = $2 WHERE id = $1", id, newStatus)
	if err != nil {
		return fmt.Errorf("r.pool.Exec: %w", err)
	}
	// отразить смену статуса в исторической таблице
	idHist := uuid.New().String()
	_, err = tx.Exec(ctx, "INSERT INTO trip_history(id, trip_id, from_status, to_status) values($1, $2, $3, $4)",
		idHist, id, oldStatus, newStatus)

	if err != nil {
		return fmt.Errorf("r.pool.Exec: %w", err)
	}

	// добавить сообщение в таблицу уведомлений
	idEvent := uuid.New().String()
	_, err = tx.Exec(ctx, "INSERT INTO outbox_event(id, event_name, aggregate_id, payload) values($1, $2, $3, $4)",
		idEvent, dto.TripEventPublished, id, dto.SentNotificationTripPublishRequest{
			TripID: id,
		})

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

func (r *RepoPg) GetForUpdateByID(
	ctx context.Context,
	tx pgx.Tx,
	id string,
) (*dto.Trip, error) {
	trip := &dto.Trip{}
	err := tx.QueryRow(ctx, "SELECT "+
		"id, "+
		"driver_id, "+
		"from_point, "+
		"to_point, "+
		"COALESCE(to_char(departure_time, 'MM-DD-YYYY HH24:MI'), '') AS departure_time, "+
		"seats, "+
		"status, "+
		"COALESCE(to_char(created_at, 'MM-DD-YYYY HH24:MI'), '') AS created_at "+
		"FROM trips WHERE id = $1 FOR UPDATE", id).Scan(
		&trip.ID,
		&trip.DriverId,
		&trip.FromPoint,
		&trip.ToPoint,
		&trip.DepartureTime,
		&trip.AvailableSeats,
		&trip.Status,
		&trip.CreatedAt,
	)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, ErrTripNotFound
		}
		return nil, fmt.Errorf("query trip by id %s: %w", id, err)
	}

	return trip, nil
}

func (r *RepoPg) EventList(ctx context.Context, tripId string) ([]dto.TripEvent, error) {
	rows, err := r.pool.Query(ctx, `select id, event_name from outbox_event  WHERE aggregate_id = $1`, tripId)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var trips []dto.TripEvent
	for rows.Next() {
		var item dto.TripEvent
		if err := rows.Scan(&item.ID, &item.Name); err != nil {
			return nil, err
		}
		trips = append(trips, item)
	}

	if err := rows.Err(); err != nil {
		return nil, err
	}

	return trips, nil
}
