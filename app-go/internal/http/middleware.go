package http

import (
	"app/internal/app"
	"strconv"
	"time"

	"github.com/gin-gonic/gin"
)

func inboundMetric(metric *app.Metric) func(*gin.Context) {
	return func(c *gin.Context) {
		// initiate time
		start := time.Now()

		// continue request
		c.Next()

		// construct values to be passed onto histogram observation for latency
		duration := time.Since(start)
		route := c.FullPath()
		method := c.Request.Method
		statusCode := strconv.Itoa(c.Writer.Status())

		// save latency observation
		metric.RecordInbound(route, method, statusCode, duration)
	}
}
