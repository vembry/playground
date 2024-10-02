package queue

import (
	"broker/model"
	"bufio"
	"encoding/json"
	"fmt"
	"log"
	"os"
	"sync"
	"time"

	"github.com/segmentio/ksuid"
)

type queue struct {
	activeQueue map[ksuid.KSUID]*model.ActiveQueue
	idleQueue   map[string]*model.IdleQueue

	ticker *time.Ticker
	mutex  sync.Mutex // Add mutex for thread safety
}

func New() *queue {
	return &queue{
		activeQueue: map[ksuid.KSUID]*model.ActiveQueue{},
		idleQueue:   map[string]*model.IdleQueue{},
		ticker:      time.NewTicker(1 * time.Second),
	}
}

// Get is to retrieve all available queues
func (q *queue) Get() model.QueueData {
	return model.QueueData{
		ActiveQueue: q.activeQueue,
		IdleQueue:   q.idleQueue,
	}
}

// Enqueue is to enqueues queue
func (q *queue) Enqueue(payload model.EnqueuePayload) error {
	// retrieve idle queue
	idleQueue, unlocker := q.retrieveIdle(payload.Name)
	defer unlocker()

	// add enqueued payload to queue maps
	idleQueue.Items = append(idleQueue.Items, payload.Payload)
	q.idleQueue[payload.Name] = idleQueue

	return nil
}

// poll is to get entry from queue head
func (q *queue) Poll(queueName string) (*model.ActiveQueue, error) {
	// retrieve idle queue
	idleQueue, unlocker := q.retrieveIdle(queueName)
	defer unlocker()

	// break away when queue has no entry
	if len(idleQueue.Items) == 0 {
		return nil, nil
	}

	// extract value from "q.queue" head
	tempQueue := idleQueue.Items[0]

	// remove it from "q.queue"
	idleQueue.Items = idleQueue.Items[1:]
	q.idleQueue[queueName] = idleQueue

	queueId := ksuid.New()

	// construct active queue entry
	activeQueue := &model.ActiveQueue{
		Id:         queueId,
		QueueName:  queueName,
		PollExpiry: time.Now().UTC().Add(20 * time.Second), // this is for sweeping purposes
		Payload:    tempQueue,
	}

	q.activeQueue[queueId] = activeQueue

	// return the polled queue
	return activeQueue, nil
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

// Shutdown shutdown broker gracefully
func (q *queue) Shutdown() {
	// q.backupQueue()
}

// backupQueue backs up queues to temporary storage
func (q *queue) backupQueue() {
	// move 'active queue' back to 'queue'
	for _, value := range q.activeQueue {
		val := q.idleQueue[value.QueueName]
		val.Items = append(val.Items, value.Payload)
		q.idleQueue[value.QueueName] = val
	}

	rawQueue, _ := json.Marshal(q.idleQueue)

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

func (q *queue) Start() {
	// q.restore()
	go q.sweep()
}

// restore is a simple way to restore broker's queue backup
func (q *queue) restore() {
	data, err := os.ReadFile("broker-backup")
	if err != nil {
		log.Fatal(err)
	}

	json.Unmarshal(data, &q.idleQueue)
	// os.Remove("broker-backup")
}

// retrieveIdle loads and lock targeted queue
func (q *queue) retrieveIdle(queueName string) (*model.IdleQueue, func()) {
	q.mutex.Lock()

	val, ok := q.idleQueue[queueName]
	if !ok {
		val = &model.IdleQueue{}
	}
	return val, func() {
		q.mutex.Unlock()
	}
}

// sweep is to sweep active queues for expiring polled queues
func (q *queue) sweep() {
	for range q.ticker.C {
		// execute sweep
		for key, val := range q.activeQueue {
			if time.Now().After(val.PollExpiry) {
				log.Printf("executing sweep...")

				// remove queue from active queue
				delete(q.activeQueue, key)

				// load/lock idle queue
				idleQueue, unlocker := q.retrieveIdle(val.QueueName)

				// add it back to queue
				idleQueue.Items = append(idleQueue.Items, val.Payload)
				q.idleQueue[val.QueueName] = idleQueue

				unlocker()
			}
		}
	}
}
