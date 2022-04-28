package elasticsearch

import (
	"context"

	"go.opentelemetry.io/collector/component"
	"go.opentelemetry.io/collector/config"
)

const (
	// The value of "type" key in configuration.
	typeStr = "elasticsearch"
)

// NewFactory creates a factory for Jaeger exporter
func NewFactory() component.ExporterFactory {
	return component.NewExporterFactory(
		typeStr,
		createDefaultConfig,
		component.WithTracesExporter(createTracesExporter))
}

func createDefaultConfig() config.Exporter {
	return &Config{}
}

func createTracesExporter(
	_ context.Context,
	set component.ExporterCreateSettings,
	config config.Exporter,
) (component.TracesExporter, error) {
	expCfg := config.(*Config)
	return newTracesExporter(expCfg, set)
}
