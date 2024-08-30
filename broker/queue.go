package main

import (
	"bufio"
	"encoding/json"
	"log"
	"net/http"
	"os"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/segmentio/ksuid"
)

type ActiveQueue struct {
	Id         ksuid.KSUID `json:"id"`
	QueueName  string      `json:"queue_name"`
	PollExpiry time.Time   `json:"poll_expiry"`
	Payload    string      `json:"payload"`
}

type enqueuePayload struct {
	Name    string `json:"name"`
	Payload string `json:"payload"`
}

type queue struct {
	activeQueue map[ksuid.KSUID]ActiveQueue
	queue       map[string][]string
}

func newQueue() *queue {

	return &queue{
		activeQueue: map[ksuid.KSUID]ActiveQueue{},
		queue:       map[string][]string{},
	}
}

func (q *queue) get(c *gin.Context) {
	c.JSON(http.StatusOK, map[string]interface{}{
		"message": "ok",
		"data": map[string]interface{}{
			"activeQueue": q.activeQueue,
			"queue":       q.queue,
		},
	})
}

// enqueue is to enqueues queue
func (q *queue) enqueue(c *gin.Context) {
	var payload enqueuePayload

	// retrieve queue payload
	c.BindJSON(&payload) // need to handle error

	// retrieve queue maps
	val, ok := q.queue[payload.Name]
	if !ok {
		// when not exists, create list
		val = []string{}
	}

	// add enqueued payload to queue maps
	val = append(val, payload.Payload)
	q.queue[payload.Name] = val

	c.Status(http.StatusOK)
}

// poll is to get entry from queue head
func (q *queue) poll(c *gin.Context) {
	queueName := c.Param("queue_name")

	// attempt to get queue
	val, ok := q.queue[queueName]
	if !ok {
		val = []string{}
		q.queue[queueName] = val
	}

	// break away when queue has no entry
	if len(val) == 0 {
		c.JSON(http.StatusOK, map[string]interface{}{
			"message": "no active queue",
			"data":    nil,
		})
		return
	}

	// extract value from "q.queue" head
	tempQueue := val[0]

	// remove it from "q.queue"
	val = val[1:]
	q.queue[queueName] = val

	queueId := ksuid.New()

	// construct active queue entry
	activeQueue := ActiveQueue{
		Id:         queueId,
		QueueName:  queueName,
		PollExpiry: time.Now().UTC().Add(1 * time.Minute), // this is for sweeping purposes
		Payload:    tempQueue,
	}

	q.activeQueue[queueId] = activeQueue

	// return the polled queue
	c.JSON(http.StatusOK, map[string]interface{}{
		"message": "ok",
		"data":    activeQueue,
	})

}

// completePoll is to ack-ed out poll-ed queue so it wont get poll-ed anymore
func (q *queue) completePoll(c *gin.Context) {
	queueIdRaw := c.Param("queue_id")

	queueId, _ := ksuid.Parse(queueIdRaw)

	// attempt to get queue
	_, ok := q.activeQueue[queueId]
	if !ok {
		c.Status(http.StatusNotFound)
		return
	}

	delete(q.activeQueue, queueId)
	c.Status(http.StatusOK)
}

// shutdown is a simple way to backup broker's queues
func (q *queue) shutdown() {
	// move 'active queue' back to 'queue'
	for _, value := range q.activeQueue {
		val := q.queue[value.QueueName]
		val = append(val, value.Payload)
		q.queue[value.QueueName] = val
	}

	rawQueue, _ := json.Marshal(q.queue)

	f, _ := os.Create("broker-backup")
	defer func() {
		f.Close()
	}()

	// make a write buffer
	w := bufio.NewWriter(f)

	// write a chunk
	if _, err := w.Write(rawQueue); err != nil {
		panic(err)
	}

	w.Flush()
}

// restore is a simple way to restore broker's queue backup
func (q *queue) restore() {
	data, err := os.ReadFile("broker-backup")
	if err != nil {
		log.Fatal(err)
	}

	json.Unmarshal(data, &q.queue)
	os.Remove("broker-backup")
}
