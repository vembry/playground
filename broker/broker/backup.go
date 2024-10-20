package broker

// // restore is a simple way to restore broker's queue backup
// func (b *broker) restore() {
// 	data, err := os.ReadFile("broker-backup")
// 	if err != nil {
// 		log.Fatal(err)
// 	}

// 	json.Unmarshal(data, &b.idleQueue)
// 	// os.Remove("broker-backup")
// }

// // backupQueue backs up queues to temporary storage
// func (b *broker) backupQueue() {
// 	// move 'active queue' back to 'queue'
// 	for _, value := range b.activeQueue {
// 		val := b.idleQueue[value.QueueName]
// 		val.Items = append(val.Items, value.Payload)
// 		b.idleQueue[value.QueueName] = val
// 	}

// 	rawQueue, _ := json.Marshal(b.idleQueue)

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
