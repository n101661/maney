//go:build generate
// +build generate

package main

//go:generate rm -rf ./server/models
//go:generate openapi-generator-go generate models --spec ./openapi.yaml --output ./server/models --package-name models
