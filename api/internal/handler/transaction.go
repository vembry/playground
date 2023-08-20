package handler

import (
	"api/internal/model"
	"fmt"
	"log"

	"github.com/gin-gonic/gin"
	"github.com/segmentio/ksuid"
)

// CreateTransaction is to create a transaction entry
func (s *handler) CreateTransaction(ctx *gin.Context) {
	userIdStr := ctx.GetHeader("x-user-id")

	var payload model.CreateTransaction
	if err := ctx.Bind(&payload); err != nil {
		log.Printf("found error on parsing request. err=%v", err)

		// when error, return 4xx
		ctx.JSON(400, GenericResponse{
			Status:  false,
			Message: "found error on parsing request",
		})
		return
	}

	var err error

	// parse user id
	payload.UserId, err = ksuid.Parse(userIdStr)
	if err != nil {
		log.Printf("found error on parsing user id. err=%v", err)
		// when error, return 4xx
		ctx.JSON(400, GenericResponse{
			Status:  false,
			Message: fmt.Errorf("found error on parsing user id. err=%w", err).Error(),
		})
		return
	}

	// call transaction domain
	err = s.transactionDomain.Create(ctx.Request.Context(), &payload)
	if err != nil {
		log.Printf("found error on creating transaction entry. err=%v", err)
		ctx.JSON(500, model.CommonResponse{
			Status:  false,
			Message: "found error on creating transaction entry",
		})
		return
	}

	ctx.JSON(200, model.CommonResponse{
		Status:  true,
		Message: "ok",
	})
}
