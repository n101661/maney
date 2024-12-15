package utils

type Option[T any] func(*T)

func ApplyOptions[T any](o *T, opts []Option[T]) *T {
	for _, opt := range opts {
		opt(o)
	}
	return o
}
