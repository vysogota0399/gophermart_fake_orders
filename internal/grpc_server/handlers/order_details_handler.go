package handlers

import (
	"context"
	"math/rand"

	"github.com/vysogota0399/gophermart/internal/grpc_server/models"
	"github.com/vysogota0399/gophermart/internal/logging"
	grpc "github.com/vysogota0399/gophermart_protos/gen/services/denormalized_order"
	"go.uber.org/zap"
)

type OrderDetailsHandler struct {
	grpc.UnimplementedDenormalizedOrderServiceServer

	lg  *logging.ZapLogger
	rep ProductsRepository
}

func NewOrderDetailsHandler(rep ProductsRepository, lg *logging.ZapLogger) *OrderDetailsHandler {
	return &OrderDetailsHandler{lg: lg, rep: rep}
}

type ProductsRepository interface {
	GenerateRandomOrder(ctx context.Context, number string) (*models.Order, error)
}

func (h *OrderDetailsHandler) OrderDetails(ctx context.Context, req *grpc.OrderDetailsRequest) (*grpc.DenormalizedOrder, error) {
	order, err := h.rep.GenerateRandomOrder(ctx, req.OrderNumber)
	if err != nil {
		h.lg.ErrorCtx(ctx, "handle order details failed", zap.Error(err))
		return nil, err
	}
	h.lg.DebugCtx(ctx, "generated order", zap.Any("order", order), zap.Any("products", order.Products))

	goods := []*grpc.DenormalizedOrder_Product{}
	for _, p := range order.Products {
		goods = append(
			goods,
			&grpc.DenormalizedOrder_Product{
				Price: rand.Int63n(9999999),
				Name:  p.Match,
			},
		)
	}

	response := &grpc.DenormalizedOrder{
		Goods:  goods,
		Number: order.Number,
	}

	return response, nil
}
