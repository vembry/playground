package model

import "github.com/go-redsync/redsync/v4"

type Mutex struct {
	Name string
	// A reference to the Redsync mutex
	Locker *redsync.Mutex
}
