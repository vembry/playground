package handler

import (
	"api/internal/model"
	"fmt"
	"log"

	"github.com/gin-gonic/gin"
	"github.com/segmentio/ksuid"
)

// GetBalance is to get user's latest balance
func (s *handler) GetBalance(c *gin.Context) {
	userIdStr := c.GetHeader("x-user-id")

	// parse user id
	userId, err := ksuid.Parse(userIdStr)
	if err != nil {
		log.Printf("found error on parsing user id. err=%v", err)
		// when error, return 4xx
		c.JSON(400, GenericResponse{
			Status:  false,
			Message: fmt.Errorf("found error on parsing user id. err=%w", err).Error(),
		})
		return
	}

	// get balance
	res, err := s.balanceDomain.Get(c.Request.Context(), userId)
	if err != nil {
		log.Printf("found error on getting active balance. err=%v", err)
		c.JSON(500, model.CommonResponse{
			Status:  false,
			Message: "found error on getting active balance",
		})
		return
	}

	c.JSON(200, map[string]interface{}{
		"status": true,
		"payload": model.BalanceResponse{
			Amount:    res.Amount,
			CreatedAt: res.CreatedAt,
			UpdatedAt: res.UpdatedAt,
		},
	})
}

// AddBalance is to add user balances
func (s *handler) AddBalance(ctx *gin.Context) {
	// AddBalanceRequest contains fields submitted via add-balance api
	type AddBalanceRequest struct {
		Amount      float64 `json:"amount"`
		UserId      string  `json:"user_id"`
		Description string  `json:"description"`
	}

	var addBalanceRequest AddBalanceRequest
	if err := ctx.Bind(&addBalanceRequest); err != nil {
		log.Printf("found error on parsing request. err=%v", err)
		// when error, return 4xx
		ctx.JSON(400, GenericResponse{
			Status:  false,
			Message: "found error on parsing request",
		})
		return
	}

	// parse user id
	userId, err := ksuid.Parse(addBalanceRequest.UserId)
	if err != nil {
		log.Printf("found error on parsing user id. err=%v", err)
		// when error, return 4xx
		ctx.JSON(400, GenericResponse{
			Status:  false,
			Message: fmt.Errorf("found error on parsing user id. err=%w", err).Error(),
		})
		return
	}

	// add queue task add-balance worker
	err = s.addBalanceHandler.Enqueue(ctx, &model.AddBalanceParam{
		UserId:      userId,
		Description: addBalanceRequest.Description,
		Amount:      addBalanceRequest.Amount,
	})
	if err != nil {
		log.Printf("found error on queuing to balance-worker. err=%v", err)

		// when error, return 4xx
		ctx.JSON(400, GenericResponse{
			Status:  false,
			Message: "found error on doing topup",
		})
		return
	}

	ctx.JSON(200, map[string]interface{}{
		"status":  true,
		"payload": addBalanceRequest,
	})
}
