package cassandra

import (
	"context"
	"path/filepath"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"go.opentelemetry.io/collector/component/componenttest"
	"go.opentelemetry.io/collector/config"
	"go.opentelemetry.io/collector/exporter/exporterhelper"
	"go.opentelemetry.io/collector/service/servicetest"

	"github.com/jaegertracing/jaeger/plugin/storage/cassandra"
)

func TestLoadConfig(t *testing.T) {
	factories, err := componenttest.NopFactories()
	assert.NoError(t, err)

	factory := NewFactory()
	factories.Exporters[typeStr] = factory
	cfg, err := servicetest.LoadConfigAndValidate(filepath.Join("testdata", "config.yaml"), factories)

	require.NoError(t, err)
	require.NotNil(t, cfg)

	e0 := cfg.Exporters[config.NewComponentID(typeStr)]
	assert.Equal(t, e0, factory.CreateDefaultConfig())

	e1 := cfg.Exporters[config.NewComponentIDWithName(typeStr, "2")]
	options := cassandra.NewOptions("cassandra")
	options.Primary.Configuration.Servers = []string{"192.168.1.1"}
	options.Primary.Configuration.Keyspace = "jaeger"
	assert.Equal(t, e1,
		&Config{
			ExporterSettings: config.NewExporterSettings(config.NewComponentIDWithName(typeStr, "2")),
			TimeoutSettings: exporterhelper.TimeoutSettings{
				Timeout: 10 * time.Second,
			},
			RetrySettings: exporterhelper.RetrySettings{
				Enabled:         true,
				InitialInterval: 10 * time.Second,
				MaxInterval:     1 * time.Minute,
				MaxElapsedTime:  10 * time.Minute,
			},
			QueueSettings: exporterhelper.QueueSettings{
				Enabled:      true,
				NumConsumers: 2,
				QueueSize:    10,
			},
			Options: options,
		})

	set := componenttest.NewNopExporterCreateSettings()
	te, err := factory.CreateTracesExporter(context.Background(), set, e1)
	require.NoError(t, err)
	require.NotNil(t, te)
}
