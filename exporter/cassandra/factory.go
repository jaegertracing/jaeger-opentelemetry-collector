package cassandra

import (
	"context"

	"go.opentelemetry.io/collector/component"
	"go.opentelemetry.io/collector/config"
	"go.opentelemetry.io/collector/exporter/exporterhelper"

	"github.com/jaegertracing/jaeger/plugin/storage/cassandra"
)

const (
	// The value of "type" key in configuration.
	typeStr   = "cassandra"
	stability = component.StabilityLevelInDevelopment
)

// NewFactory creates a factory for the Jaeger cassandra exporter
func NewFactory() component.ExporterFactory {
	return component.NewExporterFactory(
		typeStr,
		createDefaultConfig,
		component.WithTracesExporter(createTracesExporter, stability))
}

func createDefaultConfig() config.Exporter {
	return &Config{
		ExporterSettings: config.NewExporterSettings(config.NewComponentID(typeStr)),
		TimeoutSettings:  exporterhelper.NewDefaultTimeoutSettings(),
		RetrySettings:    exporterhelper.NewDefaultRetrySettings(),
		QueueSettings:    exporterhelper.NewDefaultQueueSettings(),
		Options:          cassandra.NewOptions(typeStr),
	}
}

func createTracesExporter(
	_ context.Context,
	set component.ExporterCreateSettings,
	config config.Exporter,
) (component.TracesExporter, error) {
	expCfg := config.(*Config)
	return newTracesExporter(expCfg, set)
}
