package handlers

import (
	"context"
	"math/rand"

	"github.com/vysogota0399/gophermart/internal/grpc_server/models"
	"github.com/vysogota0399/gophermart/internal/logging"
	"github.com/vysogota0399/gophermart_protos/gen/queries/order_details"
	"go.uber.org/zap"
	"google.golang.org/genproto/googleapis/type/money"
)

type OrderDetailsHandler struct {
	order_details.UnsafeQueryOrderDetailsServer

	lg  *logging.ZapLogger
	rep ProductsRepository
}

func NewOrderDetailsHandler(rep ProductsRepository, lg *logging.ZapLogger) *OrderDetailsHandler {
	return &OrderDetailsHandler{lg: lg, rep: rep}
}

type ProductsRepository interface {
	GenerateRandomOrder(ctx context.Context, number string) (*models.Order, error)
}

func (h *OrderDetailsHandler) OrderDetails(ctx context.Context, req *order_details.OrderDetailsRequest) (*order_details.OrderDetailsResponse, error) {
	order, err := h.rep.GenerateRandomOrder(ctx, req.OrderNumber)
	if err != nil {
		h.lg.ErrorCtx(ctx, "handle order details failed", zap.Error(err))
		return nil, err
	}
	h.lg.DebugCtx(ctx, "generated order", zap.Any("order", order), zap.Any("products", order.Products))

	goods := []*order_details.OrderDetailsResponse_Product{}
	for _, p := range order.Products {
		price := &money.Money{
			Units: rand.Int63n(1_000_000),
			Nanos: 0,
		}

		goods = append(
			goods,
			&order_details.OrderDetailsResponse_Product{
				Price: price,
				Name:  p.Match,
			},
		)
	}

	response := &order_details.OrderDetailsResponse{
		Goods:  goods,
		Number: order.Number,
	}

	return response, nil
}
