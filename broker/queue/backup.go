package queue

// // restore is a simple way to restore broker's queue backup
// func (q *queue) restore() {
// 	data, err := os.ReadFile("broker-backup")
// 	if err != nil {
// 		log.Fatal(err)
// 	}

// 	json.Unmarshal(data, &q.idleQueue)
// 	// os.Remove("broker-backup")
// }

// // backupQueue backs up queues to temporary storage
// func (q *queue) backupQueue() {
// 	// move 'active queue' back to 'queue'
// 	for _, value := range q.activeQueue {
// 		val := q.idleQueue[value.QueueName]
// 		val.Items = append(val.Items, value.Payload)
// 		q.idleQueue[value.QueueName] = val
// 	}

// 	rawQueue, _ := json.Marshal(q.idleQueue)

// 	f, _ := os.Create("broker-backup")
// 	defer func() {
// 		f.Close()
// 	}()

// 	// make a write buffer
// 	w := bufio.NewWriter(f)

// 	// write a chunk
// 	if _, err := w.Write(rawQueue); err != nil {
// 		panic(err)
// 	}

// 	w.Flush()
// }
