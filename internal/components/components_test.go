package components

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestDefaultComponents(t *testing.T) {
	_, err := Components()
	assert.NoError(t, err)
}
