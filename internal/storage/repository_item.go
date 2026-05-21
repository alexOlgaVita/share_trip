package storage

import (
	"context"
	"fmt"
	"github.com/jackc/pgx/v5/pgxpool"
	"job4j.ru/share-trip/internal/domain"
)

type RepoPg struct {
	pool *pgxpool.Pool
}

func NewRepoPg(pool *pgxpool.Pool) *RepoPg {
	return &RepoPg{pool: pool}
}

func (r *RepoPg) Create(ctx context.Context, it tracker.Item) error {
	_, err := r.pool.Exec(
		ctx,
		`insert into items(id, name) values($1, $2)`,
		it.ID, it.Name,
	)
	if err != nil {
		return fmt.Errorf("r.pool.Exec: %w", err)
	}
	return nil
}

func (r *RepoPg) List(ctx context.Context) ([]tracker.Item, error) {
	rows, err := r.pool.Query(ctx, `select id, name from items`)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var items []tracker.Item
	for rows.Next() {
		var item tracker.Item
		if err := rows.Scan(&item.ID, &item.Name); err != nil {
			return nil, err
		}
		items = append(items, item)
	}

	if err := rows.Err(); err != nil {
		return nil, err
	}

	return items, nil
}

func (r *RepoPg) Get(ctx context.Context, name string) (tracker.Item, error) {
	var it tracker.Item
	err := r.pool.QueryRow(
		ctx,
		`select id, name from items where name = $1`,
		name,
	).Scan(&it.ID, &it.Name)

	return it, err
}

func (r *RepoPg) Update(ctx context.Context, name string, newName string) error {
	_, err := r.pool.Exec(
		ctx,
		"UPDATE items SET name = $2 WHERE name = $1",
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
		"DELETE items WHERE name = $1",
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
		`select count(*) from items`,
	).Scan(&count)

	return count, err
}

func (r *RepoPg) DoPing(ctx context.Context) error {
	err := r.pool.Ping(ctx)
	return err
}
