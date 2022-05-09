package components

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"go.opentelemetry.io/collector/config"
)

func TestDefaultComponents(t *testing.T) {
	factories, err := Components()
	assert.NoError(t, err)

	exts := factories.Extensions
	for k, v := range exts {
		assert.Equal(t, k, v.Type())
		assert.Equal(t, config.NewComponentID(k), v.CreateDefaultConfig().ID())
	}

	recvs := factories.Receivers
	for k, v := range recvs {
		assert.Equal(t, k, v.Type())
		assert.Equal(t, config.NewComponentID(k), v.CreateDefaultConfig().ID())
	}

	procs := factories.Processors
	for k, v := range procs {
		assert.Equal(t, k, v.Type())
		assert.Equal(t, config.NewComponentID(k), v.CreateDefaultConfig().ID())
	}

	exps := factories.Exporters
	for k, v := range exps {
		assert.Equal(t, k, v.Type())
		assert.Equal(t, config.NewComponentID(k), v.CreateDefaultConfig().ID())
	}
}
