//go:build install
// +build install

package main

import (
	_ "github.com/contiamo/openapi-generator-go"
	_ "github.com/go-ozzo/ozzo-validation/v4"
)

//go:generate go install github.com/contiamo/openapi-generator-go
