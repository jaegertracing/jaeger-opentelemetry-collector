package elasticsearch

import (
	"github.com/stretchr/testify/assert"

	"testing"
)

func TestType(t *testing.T) {
	assert.Equal(t, "elasticsearch", typeStr)
}
