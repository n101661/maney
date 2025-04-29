package postgres

import (
	"github.com/lib/pq"
)

func UniqueViolationError(err error) bool {
	if e, ok := err.(*pq.Error); ok {
		return e.Code == "23505"
	}
	return false
}
