package app

import (
	"testing"

	"github.com/open-telemetry/opentelemetry-collector/config"
	"github.com/stretchr/testify/assert"
	"go.uber.org/zap"

	"github.com/jaegertracing/jaeger-opentelemetry-collector/pkg/exporter/cassandra"
	"github.com/jaegertracing/jaeger-opentelemetry-collector/pkg/exporter/elasticsearch"

	"github.com/open-telemetry/opentelemetry-collector/config/configmodels"
	"github.com/spf13/viper"
	"github.com/stretchr/testify/require"

	"github.com/jaegertracing/jaeger-opentelemetry-collector/pkg/defaults"
)

func TestDefaultConfig(t *testing.T) {
	factories, err := defaults.Components(viper.New())
	require.NoError(t, err)
	tests := []struct {
		storageType   string
		exporterTypes []string
		pipeline      map[string]*configmodels.Pipeline
	}{
		{
			storageType:   "elasticsearch",
			exporterTypes: []string{elasticsearch.TypeStr},
			pipeline: map[string]*configmodels.Pipeline{
				"traces": {
					InputType:  configmodels.TracesDataType,
					Receivers:  []string{"jaeger"},
					Exporters:  []string{elasticsearch.TypeStr},
					Processors: []string{"batch"},
				},
			},
		},
		{
			storageType:   "cassandra",
			exporterTypes: []string{cassandra.TypeStr},
			pipeline: map[string]*configmodels.Pipeline{
				"traces": {
					InputType:  configmodels.TracesDataType,
					Receivers:  []string{"jaeger"},
					Exporters:  []string{cassandra.TypeStr},
					Processors: []string{"batch"},
				},
			},
		},
		{
			storageType:   "cassandra,elasticsearch",
			exporterTypes: []string{cassandra.TypeStr, elasticsearch.TypeStr},
			pipeline: map[string]*configmodels.Pipeline{
				"traces": {
					InputType:  configmodels.TracesDataType,
					Receivers:  []string{"jaeger"},
					Exporters:  []string{cassandra.TypeStr, elasticsearch.TypeStr},
					Processors: []string{"batch"},
				},
			},
		},
	}
	for _, test := range tests {
		t.Run(test.storageType, func(t *testing.T) {
			cfg := DefaultConfig(test.storageType, factories)
			require.NoError(t, config.ValidateConfig(cfg, zap.NewNop()))

			assert.Equal(t, 1, len(cfg.Receivers))
			assert.Equal(t, "jaeger", cfg.Receivers["jaeger"].Name())
			assert.Equal(t, 1, len(cfg.Processors))
			assert.Equal(t, "batch", cfg.Processors["batch"].Name())
			assert.Equal(t, len(test.exporterTypes), len(cfg.Exporters))

			types := []string{}
			for _, v := range cfg.Exporters {
				types = append(types, v.Type())
			}
			assert.Equal(t, test.exporterTypes, types)
			assert.EqualValues(t, test.pipeline, cfg.Service.Pipelines)
		})
	}

}
