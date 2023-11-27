package serve

import (
	"app/internal/domain"
	"app/internal/model"

	"github.com/gin-gonic/gin"
	"github.com/segmentio/ksuid"
)

type handler struct {
	balanceDomain domain.IBalance
}

func newHandler(balanceDomain domain.IBalance) *handler {
	return &handler{
		balanceDomain: balanceDomain,
	}
}

// Open opens new balance. Basically creates new balance entry
func (h *handler) Open(c *gin.Context) {
	balance, err := h.balanceDomain.Open(c)
	if err != nil {
		c.JSON(400, model.BaseResponse[struct{}]{
			Error: err.Error(),
		})
	}
	c.JSON(200, model.BaseResponse[*model.Balance]{
		Object: balance,
	})
}

// Get gets balance by balance id
func (h *handler) Get(c *gin.Context) {
	// get param
	balanceIdRaw := c.Param("balance_id")
	balanceId, _ := ksuid.Parse(balanceIdRaw)

	// call service
	balance, err := h.balanceDomain.Get(c, balanceId)
	if err != nil {
		c.JSON(400, model.BaseResponse[struct{}]{
			Error: err.Error(),
		})
	}

	// return
	c.JSON(200, model.BaseResponse[*model.Balance]{
		Object: balance,
	})
}

// Withdraw attempts to withdraw balance
func (h *handler) Withdraw(c *gin.Context) {
	c.JSON(200, model.BaseResponse[string]{
		Object: "ok",
	})
}

// Deposit attempts to deposit balance
func (h *handler) Deposit(c *gin.Context) {
	c.JSON(200, model.BaseResponse[string]{
		Object: "ok",
	})
}

// Transfer attempts to send balance from a balance id to another balance id
func (h *handler) Transfer(c *gin.Context) {
	c.JSON(200, model.BaseResponse[string]{
		Object: "ok",
	})
}
