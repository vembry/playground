package handler

import (
	"app/internal/model"

	"github.com/gin-gonic/gin"
)

// HealthCheck is handle incoming health-check request.
func (s *handler) HealthCheck(c *gin.Context) {
	c.JSON(200, model.CommonResponse{
		Status:  true,
		Message: "ok",
	})
}
