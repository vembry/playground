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
