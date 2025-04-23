package slugid

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestNew(t *testing.T) {
	s := New("prefix", 11)
	assert.NotEmpty(t, s)
	assert.Len(t, s, 18)
}

func BenchmarkNew(b *testing.B) {
	for range b.N {
		New("prefix", 11)
	}
}
