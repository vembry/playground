package model

import "github.com/go-redsync/redsync/v4"

// Mutex is a wrapper for redsync.Mutex
type Mutex struct {
	Name string
	// A reference to the Redsync mutex
	Locker *redsync.Mutex
}
