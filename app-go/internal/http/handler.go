package http

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

// HealthCheck is handle incoming health-check request.
func (s *handler) HealthCheck(c *gin.Context) {
	c.JSON(200, BaseResponse[struct{}]{
		Object: struct{}{},
	})
}

// Open opens new balance. Basically creates new balance entry
func (h *handler) Open(c *gin.Context) {
	balance, err := h.balanceDomain.Open(c)
	if err != nil {
		c.JSON(400, BaseResponse[struct{}]{
			Error: err.Error(),
		})
		return
	}
	c.JSON(200, BaseResponse[*model.Balance]{
		Object: balance,
	})
}

// Get gets balance by balance id
func (h *handler) Get(c *gin.Context) {
	// get param
	balanceIdRaw := c.Param("balance_id")
	balanceId, _ := ksuid.Parse(balanceIdRaw)

	// call service
	balance, err := h.balanceDomain.Get(c.Request.Context(), balanceId)
	if err != nil {
		c.JSON(400, BaseResponse[struct{}]{
			Error: err.Error(),
		})
		return
	}

	// return
	c.JSON(200, BaseResponse[*model.Balance]{
		Object: balance,
	})
}

// Withdraw attempts to withdraw balance
func (h *handler) Withdraw(c *gin.Context) {
	var in model.WithdrawParam
	if err := c.ShouldBind(&in); err != nil {
		c.JSON(400, BaseResponse[struct{}]{
			Error: err.Error(),
		})
		return
	}

	balanceIdRaw := c.Param("balance_id")
	in.BalanceId, _ = ksuid.Parse(balanceIdRaw)

	withdrawal, err := h.balanceDomain.Withdraw(c.Request.Context(), &in)
	if err != nil {
		c.JSON(400, BaseResponse[struct{}]{
			Error: err.Error(),
		})
		return
	}
	c.JSON(200, BaseResponse[*model.Withdrawal]{
		Object: withdrawal,
	})
}

// Deposit attempts to deposit balance
func (h *handler) Deposit(c *gin.Context) {
	var in model.DepositParam
	if err := c.ShouldBind(&in); err != nil {
		c.JSON(400, BaseResponse[struct{}]{
			Error: err.Error(),
		})
		return
	}

	balanceIdRaw := c.Param("balance_id")
	in.BalanceId, _ = ksuid.Parse(balanceIdRaw)

	deposit, err := h.balanceDomain.Deposit(c.Request.Context(), &in)
	if err != nil {
		c.JSON(400, BaseResponse[struct{}]{
			Error: err.Error(),
		})
		return
	}
	c.JSON(200, BaseResponse[*model.Deposit]{
		Object: deposit,
	})
}

// Transfer attempts to send balance from a balance id to another balance id
func (h *handler) Transfer(c *gin.Context) {
	var in model.TransferParam
	if err := c.ShouldBind(&in); err != nil {
		c.JSON(400, BaseResponse[struct{}]{
			Error: err.Error(),
		})
		return
	}

	balanceIdRaw := c.Param("balance_id")
	in.BalanceIdFrom, _ = ksuid.Parse(balanceIdRaw)

	transfer, err := h.balanceDomain.Transfer(c.Request.Context(), &in)
	if err != nil {
		c.JSON(400, BaseResponse[struct{}]{
			Error: err.Error(),
		})
		return
	}
	c.JSON(200, BaseResponse[*model.Transfer]{
		Object: transfer,
	})
}
