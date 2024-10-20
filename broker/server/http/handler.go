package http

import (
	"broker/model"
	"broker/server"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/segmentio/ksuid"
)

type handler struct {
	queue server.IBroker
}

func (h *handler) get(c *gin.Context) {
	c.JSON(http.StatusOK, map[string]interface{}{
		"message": "ok",
		"data":    h.queue.Get(),
	})
}

// enqueue is to enqueues queue
func (h *handler) enqueue(c *gin.Context) {
	var payload model.EnqueuePayload

	// retrieve queue payload
	c.BindJSON(&payload) // need to handle error

	// call queue
	h.queue.Enqueue(payload)

	c.Status(http.StatusOK)
}

// poll is to get entry from queue head
func (h *handler) poll(c *gin.Context) {
	queueName := c.Param("queue_name")

	// call queue
	activeQueue, err := h.queue.Poll(queueName)
	if err != nil { // need to make specific error handler for queue
		c.JSON(http.StatusNotFound, map[string]interface{}{
			"message": err.Error(),
			"data":    nil,
		})
		return
	}

	// return the polled queue
	c.JSON(http.StatusOK, map[string]interface{}{
		"message": "ok",
		"data":    activeQueue,
	})

}

// completePoll is to ack-ed out poll-ed queue so it wont get poll-ed anymore
func (h *handler) completePoll(c *gin.Context) {
	queueIdRaw := c.Param("queue_id")
	queueId, _ := ksuid.Parse(queueIdRaw)

	// call queue
	err := h.queue.CompletePoll(queueId)
	if err != nil { // need to make specific error handler for queue
		c.JSON(http.StatusBadRequest, map[string]interface{}{
			"message": err.Error(),
			"data":    nil,
		})
		return
	}

	c.Status(http.StatusOK)
}
