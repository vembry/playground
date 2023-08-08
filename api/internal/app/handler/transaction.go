package handler

import (
	"api/internal/model"
	"fmt"

	"github.com/gin-gonic/gin"
)

// CreateTransaction is to create a transaction entry
func (s *server) CreateTransaction(ctx *gin.Context) {
	userId := ctx.GetHeader("x-user-id")

	var payload model.CreateTransaction
	if err := ctx.Bind(&payload); err != nil {
		// when error, return 4xx
		ctx.JSON(400, GenericResponse{
			Status:  false,
			Message: fmt.Errorf("found error on parsing request. err=%w", err).Error(),
		})
		return
	}

	// assign user id to request
	payload.UserId = userId

	// call transaction domain
	err := s.transactionDomain.Create(ctx.Request.Context(), &payload)
	if err != nil {
		ctx.JSON(500, model.CommonResponse{
			Status:  false,
			Message: fmt.Errorf("error on creating transaction entry. err=%w", err).Error(),
		})
	}

	ctx.JSON(200, model.CommonResponse{
		Status:  true,
		Message: "ok",
	})
}
