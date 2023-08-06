package handler

import (
	"fmt"

	"github.com/gin-gonic/gin"
)

// GetBalance is to get user's latest balance
func (s *server) GetBalance(ctx *gin.Context) {

	userId := ctx.GetHeader("x-user-id")

	ctx.JSON(200, map[string]interface{}{
		"status":  true,
		"payload": fmt.Sprintf("userId=%s", userId),
	})
}
