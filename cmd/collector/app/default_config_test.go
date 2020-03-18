package app

import (
	"testing"

	"github.com/open-telemetry/opentelemetry-collector/config"
	"github.com/open-telemetry/opentelemetry-collector/config/configmodels"
	"github.com/spf13/viper"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"go.uber.org/zap"

	"github.com/jaegertracing/jaeger-opentelemetry-collector/pkg/defaults"
)

func TestDefaultConfig(t *testing.T) {
	factories, err := defaults.Components(viper.New())
	require.NoError(t, err)
	cfg := DefaultConfig(factories)
	require.NoError(t, config.ValidateConfig(cfg, zap.NewNop()))
	assert.Equal(t, 1, len(cfg.Receivers))
	assert.Equal(t, "jaeger", cfg.Receivers["jaeger"].Name())
	assert.Equal(t, 1, len(cfg.Processors))
	assert.Equal(t, "batch", cfg.Processors["batch"].Name())
	assert.Equal(t, 1, len(cfg.Exporters))
	assert.Equal(t, "jaeger_elasticsearch", cfg.Exporters["jaeger_elasticsearch"].Name())
	assert.EqualValues(t, map[string]*configmodels.Pipeline{
		"traces": {
			InputType:  configmodels.TracesDataType,
			Receivers:  []string{"jaeger"},
			Exporters:  []string{"jaeger_elasticsearch"},
			Processors: []string{"batch"},
		},
	}, cfg.Service.Pipelines)
}
