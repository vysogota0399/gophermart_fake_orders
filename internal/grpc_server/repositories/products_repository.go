package repositories

import (
	"context"
	"database/sql"
	"fmt"
	"math/rand/v2"

	"github.com/vysogota0399/gophermart/internal/grpc_server/models"
	"github.com/vysogota0399/gophermart/internal/logging"
	"github.com/vysogota0399/gophermart/internal/storage"
	"go.uber.org/zap"
)

type ProductsRepository struct {
	strg ProductsStorage
	lg   *logging.ZapLogger
}

func NewProductsRepository(strg *storage.Storage, lg *logging.ZapLogger) *ProductsRepository {
	return &ProductsRepository{strg: strg.DB, lg: lg}
}

type ProductsStorage interface {
	QueryRowContext(ctx context.Context, query string, args ...any) *sql.Row
	QueryContext(ctx context.Context, query string, args ...any) (*sql.Rows, error)
}

func (rep *ProductsRepository) GenerateRandomOrder(ctx context.Context, number string) (*models.Order, error) {
	rowsCount := rand.Int64N(20)
	rows, err := rep.strg.QueryContext(
		ctx,
		`
			WITH random_ids AS (
					SELECT id
					FROM generate_series((SELECT MIN(id) FROM products), (SELECT MAX(id) FROM products)) AS id
					ORDER BY RANDOM()
					LIMIT $1
			)

			SELECT t.match, t.reward, t.reward_type
			FROM products t
			WHERE t.id IN (select id from random_ids);
		`,
		rowsCount,
	)
	if err != nil {
		return nil, fmt.Errorf("product_repository: generate order error %w", err)
	}
	defer rows.Close()

	order := &models.Order{Number: number}
	products := []*models.Product{}

	rep.lg.DebugCtx(ctx, "start generate report", zap.Int64("products count", rowsCount))
	for rows.Next() {
		product := models.Product{}
		if err := rows.Scan(&product.Match, &product.Reward, &product.RewardType); err != nil {
			return nil, fmt.Errorf("product_repository: generate order error %w", err)
		}

		rep.lg.DebugCtx(ctx, "add product to report", zap.Any("product", product))
		products = append(products, &product)
	}

	rep.lg.DebugCtx(ctx, "generate report finished", zap.Any("result", order))
	order.Products = products
	return order, nil
}
