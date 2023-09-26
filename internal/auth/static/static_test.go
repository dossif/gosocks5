package static

import (
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestCredentials(t *testing.T) {
	scs1, _ := NewStatic("user1", "pass1")
	assert.True(t, scs1.Valid("user1", "pass1"))
	assert.False(t, scs1.Valid("user1", "xxx"))
	assert.False(t, scs1.Valid("xxx", "pass1"))
}
