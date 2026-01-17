package repository

import (
	"errors"
)

var (
	// ErrEmptyDatabaseDSN - DatabaseDSN is empty.
	ErrEmptyDatabaseDSN = errors.New("DatabaseDSN is empty")
	// ErrCollectingDBConf - unable to collect DB conf.
	ErrCollectingDBConf = errors.New("unable to collect DB conf")
	// ErrCreateExistLogin - unable to create exists login.
	ErrCreateExistLogin = errors.New("unable to create exists login")
	// ErrDataNotFound - unable to find data.
	ErrDataNotFound = errors.New("unable to find data")
	// ErrDataAlreadyExist - data already exist.
	ErrDataAlreadyExists = errors.New("data already exist")
	// ErrNoData - no data.
	ErrNoData = errors.New("no data")

	ErrUserAlreadyExists = errors.New("data already exist")
)
