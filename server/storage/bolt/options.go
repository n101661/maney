package bolt

import (
	"encoding/json"
)

type options struct {
	marshalValue   func(interface{}) ([]byte, error)
	unmarshalValue func([]byte, interface{}) error
}

func defaultOptions() *options {
	return &options{
		marshalValue:   json.Marshal,
		unmarshalValue: json.Unmarshal,
	}
}

func WithMarshaller(marshaler func(interface{}) ([]byte, error)) func(*options) {
	return func(o *options) {
		o.marshalValue = marshaler
	}
}

func WithUnmarshaller(unmarshaler func([]byte, interface{}) error) func(*options) {
	return func(o *options) {
		o.unmarshalValue = unmarshaler
	}
}
