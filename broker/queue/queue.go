package queue

import (
	"broker/model"
	"bufio"
	"encoding/json"
	"fmt"
	"log"
	"os"
	"time"

	"github.com/segmentio/ksuid"
)

type queue struct {
	activeQueue map[ksuid.KSUID]model.ActiveQueue
	queue       map[string]model.IdleQueue
}

func New() *queue {
	return &queue{
		activeQueue: map[ksuid.KSUID]model.ActiveQueue{},
		queue:       map[string]model.IdleQueue{},
	}
}

// Get is to retrieve all available queues
func (q *queue) Get() model.QueueData {
	return model.QueueData{
		ActiveQueue: q.activeQueue,
		Queue:       q.queue,
	}
}

// Enqueue is to enqueues queue
func (q *queue) Enqueue(payload model.EnqueuePayload) error {
	// retrieve queue maps
	val, ok := q.queue[payload.Name]
	if !ok {
		// when not exists, create list
		val = model.IdleQueue{}
	}

	// add enqueued payload to queue maps
	val.Items = append(val.Items, payload.Payload)
	q.queue[payload.Name] = val

	return nil
}

// poll is to get entry from queue head
func (q *queue) Poll(queueName string) (*model.ActiveQueue, error) {
	// queueName := c.Param("queue_name")

	// attempt to get queue
	val, ok := q.queue[queueName]
	if !ok {
		val = model.IdleQueue{}
		q.queue[queueName] = val
	}

	// break away when queue has no entry
	if len(val.Items) == 0 {
		return nil, fmt.Errorf("no active queue")
	}

	// extract value from "q.queue" head
	tempQueue := val.Items[0]

	// remove it from "q.queue"
	val.Items = val.Items[1:]
	q.queue[queueName] = val

	queueId := ksuid.New()

	// construct active queue entry
	activeQueue := model.ActiveQueue{
		Id:         queueId,
		QueueName:  queueName,
		PollExpiry: time.Now().UTC().Add(1 * time.Minute), // this is for sweeping purposes
		Payload:    tempQueue,
	}

	q.activeQueue[queueId] = activeQueue

	// return the polled queue
	return &activeQueue, nil
}

// CompletePoll is to ack-ed out poll-ed queue so it wont get poll-ed anymore
func (q *queue) CompletePoll(queueId ksuid.KSUID) error {
	// attempt to get queue
	_, ok := q.activeQueue[queueId]
	if !ok {
		return fmt.Errorf("queue not found")
	}

	delete(q.activeQueue, queueId)
	return nil
}

// shutdown is a simple way to backup broker's queues
func (q *queue) Shutdown() {
	// move 'active queue' back to 'queue'
	for _, value := range q.activeQueue {
		val := q.queue[value.QueueName]
		val.Items = append(val.Items, value.Payload)
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
func (q *queue) Restore() {
	data, err := os.ReadFile("broker-backup")
	if err != nil {
		log.Fatal(err)
	}

	json.Unmarshal(data, &q.queue)
	// os.Remove("broker-backup")
}
