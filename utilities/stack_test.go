package utilities

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestStack(t *testing.T) {
	stack := NewStack()
	for i := int16(0); i < 20; i++ {
		stack.Push(i)
	}
	for i := int16(19); i >= 0; i-- {
		v := stack.Pop()
		j := v.(int16)
		assert.True(t, j == i, "The stack is not popping properly")
	}
	assert.Nil(t, stack.Pop(), "Default element should be nil")
}
