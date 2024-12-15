package storage

import (
	"errors"
)

var (
	ErrExists   = errors.New("the data exists")
	ErrNotFound = errors.New("the data is not found")
)
