package seeds

import (
	"context"
	"database/sql"
	"fmt"
	"log"
	"math/rand/v2"
	"net/http"
	"time"

	"github.com/brianvoe/gofakeit/v7"
	"github.com/go-resty/resty/v2"
	"github.com/vysogota0399/gophermart/internal/grpc_server/models"
	"github.com/vysogota0399/gophermart/internal/storage"
	"golang.org/x/sync/errgroup"
)

type ProductsCreator struct {
	Strg                  *storage.Storage
	AccrualCreateGoodsURL string
}

func truncateProducts(ctx context.Context, tx *sql.Tx) error {
	_, err := tx.ExecContext(ctx, "truncate products")
	if err != nil {
		return fmt.Errorf("truncate products failed")
	}
	return nil
}

func (s *ProductsCreator) Call() error {
	ctx := context.Background()
	tx, err := s.Strg.DB.BeginTx(ctx, &sql.TxOptions{})
	if err != nil {
		log.Println(err)
		return err
	}

	if err := truncateProducts(ctx, tx); err != nil {
		log.Println(err)
		tx.Commit()
		return err
	}

	g, ctx := errgroup.WithContext(ctx)

	s.CreateAccrual(ctx, g, s.CreateDBrecord(ctx, tx, g, productsGenerator()))

	if err := g.Wait(); err != nil {
		tx.Rollback()
		return fmt.Errorf("create products failed error %w", err)
	}

	tx.Commit()
	return nil
}

func (s *ProductsCreator) CreateDBrecord(ctx context.Context, tx *sql.Tx, g *errgroup.Group, products chan *models.Product) chan *models.Product {
	out := make(chan *models.Product)
	g.Go(func() error {
		defer close(out)

		for p := range products {
			if _, err := tx.ExecContext(
				ctx,
				`
					insert into products (match, reward, reward_type)
					values ($1, $2, $3)
				`,
				p.Match, p.Reward, p.RewardType,
			); err != nil {
				return fmt.Errorf("create product %+v failed error %w", p, err)
			}

			log.Printf("[DB] Created %+v", p)
			out <- p
		}

		return nil
	})

	return out
}

func (s *ProductsCreator) CreateAccrual(ctx context.Context, g *errgroup.Group, products chan *models.Product) {
	g.Go(func() error {
		for p := range products {

			req := resty.New().R()
			req.Method = http.MethodPost
			req.SetBody(p)

			req.URL = s.AccrualCreateGoodsURL
			req.Header.Add("Content-Type", "application/json")

			_, err := req.Send()
			if err != nil {
				return fmt.Errorf("accrual create product %+v failed error %w", p, err)
			}

			log.Printf("[ACCRUAL] Created %+v", p)
		}

		return nil
	})
}

func productsGenerator() chan *models.Product {
	out := make(chan *models.Product)

	go func() {
		defer close(out)
		for i := 0; i < 100; i++ {
			time.Sleep(time.Millisecond * 50)
			out <- &models.Product{
				Reward:     rand.Int64N(100),
				Match:      gofakeit.Fruit(),
				RewardType: "%",
			}
		}
	}()

	return out
}
