package handler

import (
	"api/internal/model"
	"fmt"
	"log"

	"github.com/gin-gonic/gin"
	"github.com/segmentio/ksuid"
)

// CreateTopup is to create a topup entry.
func (s *handler) CreateTopup(ctx *gin.Context) {
	userIdStr := ctx.GetHeader("x-user-id")

	// CreateTopupRequest contains fields submitted via create-topup api
	type CreateTopupRequest struct {
		Amount float64 `json:"amount"`
	}

	var topupRequest CreateTopupRequest
	if err := ctx.Bind(&topupRequest); err != nil {
		log.Printf("found error on parsing request. err=%v", err)
		// when error, return 4xx
		ctx.JSON(400, GenericResponse{
			Status:  false,
			Message: "found error on parsing request",
		})
		return
	}

	// parse user id
	userId, err := ksuid.Parse(userIdStr)
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
		UserId: userId,
		Amount: topupRequest.Amount,
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
		"payload": topupRequest,
	})
}
