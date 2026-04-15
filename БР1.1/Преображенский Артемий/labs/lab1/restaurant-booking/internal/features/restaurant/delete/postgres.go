package delete

import (
	"context"
	"restaurant-booking/internal/adapter/postgres"
	"restaurant-booking/internal/domain"

	"github.com/google/uuid"
)

type postgresRepository struct {
	pool *postgres.Pool
}

func NewPostgres(pool *postgres.Pool) *postgresRepository {
	return &postgresRepository{pool: pool}
}

func (r *postgresRepository) Delete(ctx context.Context, id uuid.UUID) (domain.Restaurant, error) {
	const query = `
		DELETE FROM restaurants WHERE id = $1 RETURNING id
	`
	var restaurantID uuid.UUID
	err := r.pool.Pgx().QueryRow(ctx, query, id).Scan(&restaurantID)
	if err != nil {
		return domain.Restaurant{}, err
	}
	return domain.Restaurant{ID: restaurantID}, nil
}
