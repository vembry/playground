package serve

import (
	"app/internal/domain"
	"app/internal/model"

	"github.com/gin-gonic/gin"
)

type handler struct {
	balanceDomain domain.IBalance
}

func newHandler(balanceDomain domain.IBalance) *handler {
	return &handler{
		balanceDomain: balanceDomain,
	}
}

type baseResponse[C any] struct {
	Error  string `json:"error"`
	Object C      `json:"object"`
}

func (h *handler) Open(c *gin.Context) {
	balance, err := h.balanceDomain.Open(c.Request.Context())
	if err != nil {
		c.JSON(400, baseResponse[struct{}]{
			Error: err.Error(),
		})
	}
	c.JSON(200, baseResponse[*model.Balance]{
		Object: balance,
	})
}

func (h *handler) Withdraw(c *gin.Context) {
	c.JSON(200, baseResponse[string]{
		Object: "ok",
	})
}
func (h *handler) Deposit(c *gin.Context) {
	c.JSON(200, baseResponse[string]{
		Object: "ok",
	})
}
func (h *handler) Transfer(c *gin.Context) {
	c.JSON(200, baseResponse[string]{
		Object: "ok",
	})
}
