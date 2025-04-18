package repository

import (
	"errors"
)

// Repository errors.
var (
	ErrDataExists   = errors.New("the data exists")
	ErrDataNotFound = errors.New("the data is not found")
)
