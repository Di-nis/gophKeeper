package repository

import (
	"errors"
)

var (
	// ErrEmptyDatabasePath - Database path is empty.
	ErrEmptyDatabasePath = errors.New("Database path is empty")
)
