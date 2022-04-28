package elasticsearch

import (
	"go.opentelemetry.io/collector/component"
	"go.opentelemetry.io/collector/exporter/exporterhelper"
)

// newTracesExporter returns a new Jaeger gRPC exporter.
// The exporter name is the name to be used in the observability of the exporter.
// The collectorEndpoint should be of the form "hostname:14250" (a gRPC target).
func newTracesExporter(cfg *Config, set component.ExporterCreateSettings) (component.TracesExporter, error) {
	return exporterhelper.NewTracesExporter(cfg, set, nil)
}
