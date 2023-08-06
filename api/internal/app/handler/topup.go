package handler

import (
	"fmt"

	"github.com/gin-gonic/gin"
)

// CreateTopup is to create a topup entry.
func (s *server) CreateTopup(ctx *gin.Context) {

	// CreateTopupRequest contains fields submitted via create-topup api
	type CreateTopupRequest struct {
		Amount float64 `json:"amount"`
	}

	var topupRequest CreateTopupRequest
	if err := ctx.Bind(&topupRequest); err != nil {
		// when error, return 4xx
		ctx.JSON(400, GenericResponse{
			Status:  false,
			Message: fmt.Errorf("found error. err=%w", err).Error(),
		})
		return
	}

	ctx.JSON(200, map[string]interface{}{
		"status":  true,
		"payload": topupRequest,
	})
}
