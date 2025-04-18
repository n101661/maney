package bolt

import (
	"encoding/json"

	"github.com/n101661/maney/pkg/utils"
)

type Options struct {
	MarshalValue   func(interface{}) ([]byte, error)
	UnmarshalValue func([]byte, interface{}) error
}

func DefaultOptions() *Options {
	return &Options{
		MarshalValue:   json.Marshal,
		UnmarshalValue: json.Unmarshal,
	}
}

func WithMarshaller(marshaler func(interface{}) ([]byte, error)) utils.Option[Options] {
	return func(o *Options) {
		o.MarshalValue = marshaler
	}
}

func WithUnmarshaler(unmarshaler func([]byte, interface{}) error) utils.Option[Options] {
	return func(o *Options) {
		o.UnmarshalValue = unmarshaler
	}
}
