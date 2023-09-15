package elasticsearch

import (
	"context"

	"go.opentelemetry.io/collector/component"
	"go.opentelemetry.io/collector/exporter"
)

const (
	// The value of "type" key in configuration.
	typeStr   = "elasticsearch"
	stability = component.StabilityLevelAlpha
)

// NewFactory creates a factory for Jaeger exporter
func NewFactory() exporter.Factory {
	return exporter.NewFactory(
		typeStr,
		createDefaultConfig,
		exporter.WithTraces(createTracesExporter, stability))
}

func createDefaultConfig() component.Config {
	return &Config{
		// ExporterSettings: component.NewExporterSettings(component.NewID(typeStr)),
	}
}

func createTracesExporter(
	ctx context.Context,
	set exporter.CreateSettings,
	cfg component.Config,
) (exporter.Traces, error) {
	expCfg := cfg.(*Config)
	return newTracesExporter(expCfg, set)
}
