package files

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

const fn = "files.go"

func TestExists(t *testing.T) {
	assert.True(t, Exists(fn))
	assert.False(t, Exists("garbage"))
}
