package api

import (
	"context"
	"fmt"
	"math/rand/v2"
	"net/http"
	"strconv"

	"github.com/theplant/luhn"
	"github.com/vysogota0399/gophermart/internal/logging"
	"go.uber.org/zap"
)

type Handler struct {
	lg *logging.ZapLogger
}

func NewHandler(lg *logging.ZapLogger) *Handler {
	return &Handler{lg: lg}
}

func (h *Handler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	n := 16
	bnumber := make([]byte, n-1)
	for i := range bnumber {
		bnumber[i] = byte(rand.Int32N(10) + '0')
	}

	number, err := strconv.Atoi(string(bnumber))
	if err != nil {
		h.lg.ErrorCtx(context.Background(), "generate number error", zap.Error(err))
	}

	w.Write([]byte(fmt.Sprintf("%d", Luhner(number))))
}

func Luhner(numb int) int {
	if luhn.Valid(numb) {
		return numb
	}
	return 10*numb + luhn.CalculateLuhn(numb)
}
