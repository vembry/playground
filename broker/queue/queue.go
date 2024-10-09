package queue

import (
	"broker/model"
	"fmt"
	"log"
	"sync"
	"time"

	"github.com/segmentio/ksuid"
)

type queue struct {
	idleQueue   sync.Map
	activeQueue sync.Map

	ticker   *time.Ticker
	mutexMap sync.Map
}

func New() *queue {
	return &queue{
		idleQueue:   sync.Map{},
		activeQueue: sync.Map{},
		ticker:      time.NewTicker(1 * time.Second),
		mutexMap:    sync.Map{},
	}
}

// Get is to retrieve all available queues
func (q *queue) Get() model.QueueData {

	i, j := 0, 0
	q.idleQueue.Range(func(key, value any) bool {
		i += len(value.(*model.IdleQueue).Items)
		return true
	})
	q.activeQueue.Range(func(key, value any) bool {
		j++
		return true
	})

	return model.QueueData{
		IdleQueueCount:   int64(i),
		ActiveQueueCount: int64(j),
	}
}

// Enqueue is to enqueues queue
func (q *queue) Enqueue(payload model.EnqueuePayload) error {
	// retrieve idle queue
	idleQueue, unlocker := q.retrieveIdle(payload.Name)
	defer unlocker()

	// add enqueued payload to queue maps
	idleQueue.Items = append(idleQueue.Items, payload.Payload)

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

	// extract value from idleQueue's head
	tempQueue := idleQueue.Items[0]

	// slice extracted-queue from idleQueue
	idleQueue.Items = idleQueue.Items[1:]

	queueId := ksuid.New()

	// construct active queue entry
	activeQueue := &model.ActiveQueue{
		Id:         queueId,
		QueueName:  queueName,
		PollExpiry: time.Now().UTC().Add(20 * time.Second), // this is for sweeping purposes
		Payload:    tempQueue,
	}

	q.activeQueue.Store(queueId, activeQueue)

	// return the polled queue
	return activeQueue, nil
}

// CompletePoll is to ack-ed out poll-ed queue so it wont get poll-ed anymore
func (q *queue) CompletePoll(queueId ksuid.KSUID) error {
	// attempt to get queue
	_, ok := q.activeQueue.Load(queueId)
	if !ok {
		return fmt.Errorf("queue not found")
	}

	q.activeQueue.Delete(queueId)
	return nil
}

// Shutdown shutdown broker gracefully
func (q *queue) Shutdown() {
	// q.backupQueue()
}

func (q *queue) Start() {
	// q.restore()
	go q.sweep()
}

// retrieveIdle loads and lock targeted queue
func (q *queue) retrieveIdle(queueName string) (*model.IdleQueue, func()) {
	// Get or create a mutex for the specific queueName
	mutex, _ := q.mutexMap.LoadOrStore(queueName, &sync.Mutex{})

	// Lock the mutex for this specific queue
	mutex.(*sync.Mutex).Lock()

	// retrieve queue from map
	val, _ := q.idleQueue.LoadOrStore(queueName, &model.IdleQueue{})

	return val.(*model.IdleQueue), func() {
		mutex.(*sync.Mutex).Unlock()
	}
}

// sweep is to sweep active queues for expiring polled queues
func (q *queue) sweep() {
	for range q.ticker.C {
		// execute sweep
		q.activeQueue.Range(func(key, value any) bool {
			val := value.(*model.ActiveQueue)
			if time.Now().After(val.PollExpiry) {
				log.Printf("executing sweep...")

				// remove queue from active queue
				q.activeQueue.Delete(key)

				// load/lock idle queue
				idleQueue, unlocker := q.retrieveIdle(val.QueueName)

				// add it back to queue
				idleQueue.Items = append(idleQueue.Items, val.Payload)

				unlocker()
			}

			return true
		})
	}
}
