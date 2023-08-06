package handler

import (
	"api/internal/model"

	"github.com/gin-gonic/gin"
)

// HealthCheck is handle incoming health-check request.
func (s *server) HealthCheck(c *gin.Context) {
	c.JSON(200, model.CommonResponse{
		Status:  true,
		Message: "ok",
	})
}
