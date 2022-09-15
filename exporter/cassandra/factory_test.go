package cassandra

import (
	"context"
	"testing"

	"github.com/stretchr/testify/assert"
	"go.opentelemetry.io/collector/component"
	"go.opentelemetry.io/collector/component/componenttest"
	"go.opentelemetry.io/collector/config/configtest"
)

func TestCreateDefaultConfig(t *testing.T) {
	factory := NewFactory()
	cfg := factory.CreateDefaultConfig()
	assert.NotNil(t, cfg, "failed to create default config")
	assert.NoError(t, configtest.CheckConfigStruct(cfg))
}

func TestCreateMetricsExporter(t *testing.T) {
	factory := NewFactory()
	cfg := factory.CreateDefaultConfig()

	set := componenttest.NewNopExporterCreateSettings()
	_, err := factory.CreateMetricsExporter(context.Background(), set, cfg)
	assert.Error(t, err, component.ErrDataTypeIsNotSupported)
}

func TestCreateInstanceViaFactory(t *testing.T) {
	factory := NewFactory()

	cfg := factory.CreateDefaultConfig()
	set := componenttest.NewNopExporterCreateSettings()
	exp, err := factory.CreateTracesExporter(context.Background(), set, cfg)
	assert.NoError(t, err)
	exp, err = factory.CreateTracesExporter(context.Background(), set, cfg)
	assert.NoError(t, err)
	assert.NotNil(t, exp)
	assert.NoError(t, exp.Shutdown(context.Background()))
}
