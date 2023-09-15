package elasticsearch

import (
	"context"

	"go.opentelemetry.io/collector/component"
	"go.opentelemetry.io/collector/exporter"
	"go.opentelemetry.io/collector/exporter/exporterhelper"
	"go.opentelemetry.io/collector/pdata/ptrace"
)

// newTracesExporter returns a new Jaeger gRPC exporter.
// The exporter name is the name to be used in the observability of the exporter.
// The collectorEndpoint should be of the form "hostname:14250" (a gRPC target).
func newTracesExporter(cfg *Config, set exporter.CreateSettings) (exporter.Traces, error) {
	s := newSender(cfg, set.TelemetrySettings)
	return exporterhelper.NewTracesExporter(context.TODO(), set, cfg, s.pushTraces)
}

type sender struct {
	name     string
	settings component.TelemetrySettings
}

func newSender(cfg *Config, settings component.TelemetrySettings) *sender {
	return &sender{
		// name:     cfg.ID().String(),
		settings: settings,
	}
}

func (s *sender) pushTraces(_ context.Context, _ ptrace.Traces) error {
	return nil
}
