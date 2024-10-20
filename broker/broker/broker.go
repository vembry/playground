package broker

import (
	"broker/model"
	"fmt"
	"log"
	"sync"
	"time"

	"github.com/segmentio/ksuid"
)

type broker struct {
	idleQueue   sync.Map
	activeQueue sync.Map

	mutexMap sync.Map     // for locking purposes
	ticker   *time.Ticker // for sweeping purposes

	backupProvider IBackup // to handle broker backups
}

func New(backupProvider IBackup) *broker {
	// force
	if backupProvider == nil {
		backupProvider = NewFileDumper("broker-backup")
	}

	return &broker{
		idleQueue:      sync.Map{},
		activeQueue:    sync.Map{},
		ticker:         time.NewTicker(1 * time.Second),
		mutexMap:       sync.Map{},
		backupProvider: backupProvider,
	}
}

// Get is to retrieve all available queues
func (b *broker) Get() model.QueueData {

	i, j := 0, 0
	b.idleQueue.Range(func(key, value any) bool {
		i += len(value.(*model.IdleQueue).Items)
		return true
	})
	b.activeQueue.Range(func(key, value any) bool {
		j++
		return true
	})

	return model.QueueData{
		IdleQueueCount:   int64(i),
		ActiveQueueCount: int64(j),
	}
}

// Enqueue is to enqueues queue
func (b *broker) Enqueue(request model.EnqueuePayload) error {
	// retrieve idle queue
	idleQueue, unlocker := b.retrieveIdle(request.Name)
	defer unlocker()

	// add enqueued payload to queue maps
	idleQueue.Items = append(idleQueue.Items, &model.Queue{Payload: request.Payload})

	return nil
}

// poll is to get entry from queue head
func (b *broker) Poll(queueName string) (*model.ActiveQueue, error) {
	// retrieve idle queue
	idleQueue, unlocker := b.retrieveIdle(queueName)
	defer unlocker()

	// break away when queue has no entry
	if len(idleQueue.Items) == 0 {
		return nil, nil
	}

	// extract value from idleQueue's head
	queue := idleQueue.Items[0]

	// slice extracted-queue from idleQueue
	idleQueue.Items = idleQueue.Items[1:]

	queueId := ksuid.New()

	// construct active queue entry
	activeQueue := &model.ActiveQueue{
		Id:         queueId,
		QueueName:  queueName,
		PollExpiry: time.Now().UTC().Add(20 * time.Second), // this is for sweeping purposes
		Queue:      queue,
	}

	b.activeQueue.Store(queueId, activeQueue)

	// return the polled queue
	return activeQueue, nil
}

// CompletePoll is to ack-ed out poll-ed queue so it wont get poll-ed anymore
func (b *broker) CompletePoll(queueId ksuid.KSUID) error {
	// attempt to get queue
	_, ok := b.activeQueue.Load(queueId)
	if !ok {
		return fmt.Errorf("queue not found")
	}

	// remove queue from active queue
	b.activeQueue.Delete(queueId)
	return nil
}

// Stop handler to shutdown broker
func (b *broker) Stop() {
	b.deactivateQueues()

	maps := map[string]*model.IdleQueue{}
	b.idleQueue.Range(func(key, value any) bool {
		maps[key.(string)] = value.(*model.IdleQueue)
		return true
	})

	b.backupProvider.Backup(maps)
}

// Start handler to start broker
func (b *broker) Start() {
	b.restore()
	go b.sweepWorker()
}

func (b *broker) restore() {
	maps := b.backupProvider.Restore()

	for key, val := range maps {
		b.idleQueue.Store(key, val)
	}
}

// retrieveIdle loads and lock targeted queue
func (b *broker) retrieveIdle(queueName string) (*model.IdleQueue, func()) {
	// Get or create a mutex for the specific queueName
	mutex, _ := b.mutexMap.LoadOrStore(queueName, &sync.Mutex{})

	// Lock the mutex for this specific queue
	mutex.(*sync.Mutex).Lock()

	// retrieve queue from map
	val, _ := b.idleQueue.LoadOrStore(queueName, &model.IdleQueue{})

	return val.(*model.IdleQueue), func() {
		mutex.(*sync.Mutex).Unlock()
	}
}

// sweepWorker is to sweep expiring active queues
func (b *broker) sweepWorker() {
	for range b.ticker.C {
		b.activeQueue.Range(b.sweepActual)
	}
}

// sweepActual is to check and remove if an active-queue entry has expired
func (b *broker) sweepActual(key, value any) bool {
	val := value.(*model.ActiveQueue)
	if time.Now().After(val.PollExpiry) {
		log.Printf("sweeping out %s...", val.Id)

		b.deactivateQueue(val)
	}

	return true
}

// deactivateQueue deactivate queue and put it back to idle queue
func (b *broker) deactivateQueue(queue *model.ActiveQueue) {
	// remove queue from active queue
	b.activeQueue.Delete(queue.Id)

	idleQueue, unlocker := b.retrieveIdle(queue.QueueName)
	defer unlocker()

	// add active queue back to idle queue
	idleQueue.Items = append(idleQueue.Items, queue.Queue)
}

// deactivateQueue deactivate queue and put it back to idle queue
func (b *broker) deactivateQueues() {
	b.activeQueue.Range(func(key, value any) bool {
		b.deactivateQueue(value.(*model.ActiveQueue))
		return true
	})
}
