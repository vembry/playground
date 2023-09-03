package model

import "gorm.io/gorm"

// DB is a structure consisting master and slave db connection
type DB struct {
	Master *gorm.DB
	Slave  *gorm.DB
}
