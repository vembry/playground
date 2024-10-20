package broker

import (
	"broker/model"
	"bufio"
	"encoding/json"
	"log"
	"os"
)

// fileDumper is backup/restore handler for the broker's data in file
type fileDumper struct {
	path string
}

func NewFileDumper(path string) *fileDumper {
	return &fileDumper{
		path: path,
	}
}

func (fd *fileDumper) Restore() map[string]*model.IdleQueue {
	out := map[string]*model.IdleQueue{}

	data, err := os.ReadFile(fd.path)
	if err != nil {
		log.Printf("err=%v", err)
		return out
	}

	json.Unmarshal(data, &out)
	os.Remove(fd.path)

	return out
}

func (fd *fileDumper) Backup(queue map[string]*model.IdleQueue) {
	rawQueue, _ := json.Marshal(queue)

	f, _ := os.Create(fd.path)
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
