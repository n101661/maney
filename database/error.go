package database

import "errors"

var (
	ErrResourceExisted = errors.New("database: the resource has existed")
)
