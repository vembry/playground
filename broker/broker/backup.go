package broker

import "broker/model"

// IBackup is the specification to persist queue datas for backup/restore purposes.
type IBackup interface {
	Restore() map[string]*model.IdleQueue
	Backup(map[string]*model.IdleQueue)
}

// TODO: should we consider batched backup/restore?
