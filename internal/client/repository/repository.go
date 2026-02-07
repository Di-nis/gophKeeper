package repository

import (
	"errors"
)

var (
	// ErrEmptyDatabasePath - Database path is empty.
	ErrEmptyDatabasePath = errors.New("database path is empty")
	// ErrDataAlreadyExist - data already exist.
	ErrDataAlreadyExists = errors.New("data already exist")
)
