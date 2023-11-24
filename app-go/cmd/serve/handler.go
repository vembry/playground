package serve

import (
	"app/internal/domain"

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

func (h *handler) SetupBalance(c *gin.Context) {
	c.JSON(200, baseResponse[string]{
		Object: "ok",
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
